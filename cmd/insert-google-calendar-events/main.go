package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Worker
// キューから予定を取り出してGoogleカレンダーに登録する
type Worker struct {
}

func (w *Worker) getService(config *oauth2.Config, userName string) (*calendar.Service, error) {
	ctx := context.Background()
	tokFile := filepath.Join("internal", "app", fmt.Sprintf("%s-google-token.json", userName))
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, tok)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

var config *oauth2.Config

func (w *Worker) Start() {
	for event := range queue {
		go func(evt *Event) {
			srv, err := w.getService(config, evt.userId)
			if err != nil {
				log.Fatalf("Unable to retrieve calendar Client %v", err)
			}
			time.Sleep(1 * time.Second)

			calendarId := "primary"
			e, err := srv.Events.Insert(calendarId, evt.Event).Do()
			if err != nil {
				log.Fatalf("Unable to create event. %v\n", err)
			}

			fmt.Printf("Event created: %s\n", e.HtmlLink)
		}(event)
	}
}

type Event struct {
	Event  *calendar.Event
	userId string
}

var queue = make(chan *Event, 10)
var events = []*Event{
	{
		Event: &calendar.Event{
			Summary: "User1-くまさん",
			Start: &calendar.EventDateTime{
				DateTime: "2024-03-15T09:00:00+09:00",
				TimeZone: "Asia/Tokyo",
			},
			End: &calendar.EventDateTime{
				DateTime: "2024-03-15T10:00:00+09:00",
				TimeZone: "Asia/Tokyo",
			},
			ColorId: "1",
		},
		userId: "user1",
	},
	{
		Event: &calendar.Event{
			Summary: "User2-うさぎさん",
			Start: &calendar.EventDateTime{
				DateTime: "2024-03-15T09:00:00+09:00",
				TimeZone: "Asia/Tokyo",
			},
			End: &calendar.EventDateTime{
				DateTime: "2024-03-15T10:00:00+09:00",
				TimeZone: "Asia/Tokyo",
			},
			ColorId: "2",
		},
		userId: "user2",
	},
}

func main() {
	// OAuthクライアントの設定
	b, err := os.ReadFile(filepath.Join("configs", "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err = google.ConfigFromJSON(b,
		calendar.CalendarScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// キューにイベントを登録
	for _, event := range events {
		queue <- event
	}
	close(queue)

	// Workerを作成
	workerCount := 1
	for i := 0; i < workerCount; i++ {
		w := &Worker{}
		go w.Start()
	}

	time.Sleep(10 * time.Second)
}
