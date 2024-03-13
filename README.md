
```bash
.
|-- cmd
|   |-- insert-google-calendar-events
|   |   `-- main.go
|   `-- save-google-account-token
|       `-- main.go
|-- configs
|   `-- credentials.json
|-- go.mod
|-- go.sum
`-- internal
    `-- app
        |-- user1-google-token.json
        `-- user2-google-token.json
```

- **cmd/insert-google-calendar-events/main.go**

  Googleカレンダーに予定を登録するワーカーです.
  
  ※ main.go 内のダミーEvent.userIdはトークン保存時に指定したユーザ名に書き換えます.
  
  ```bash
  $ go run ./cmd/insert-google-calendar-events/main.go
  ```

- **cmd/save-google-account-token/main.go**

  Googleカレンダーへ予定を登録するためのtokenを取得し保存します.

  実行するとユーザへGoogleカレンダーの書き込み権限を求める認可画面URLがコンソールに表示されます.

  手順に従い、得た認可コードをコンソールに入力しEnter押下でtokenが保存されます.

  ```
  $ go run ./cmd/save-google-account-token/main.go -user=user1
  ```

- **configs/credentials.json**

  OAuthクライアントアプリのclient_idとclient_secretです.


- **internal/app/_user_name_-google-token.json**

  リソースオーナー（user）のトークンです.

