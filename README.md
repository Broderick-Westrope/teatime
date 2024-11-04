# Tea Time

I want to create a TUI messaging application using Golang and Bubble Tea. It will have a client-server architecture, make use of WebSockets, and persist state.  Eventually, it will leverage the Signal protocol, and have typing indicators and read receipts, but I want to start simple and add these features incrementally.

## Usage

Start server:
```sh
go run ./server/main.go
```

Start client:
```sh
go run ./client/main.go 'Robby.Receiver' 'pa$$word'
```

Write debug logs to `./logs`:
```sh
DEBUG=t go run ./client/main.go 'Robby.Receiver' 'pa$$word'
```
