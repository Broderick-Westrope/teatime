# Tea Time

A secure TUI messaging application.

**What is a TUI?**

A TUI (Textual/Terminal User Interface) is similar to a CLI (Command-Line Interface) but differs in that a TUI provides interactive, text-based graphical elements like menus and windows within the terminal, whereas a CLI relies solely on textual input and output without graphical enhancements.

**Why did you build this?**

I set out to build this in an effort to either learn or understand a range of technologies. Tea Time is built using:
- [Go](https://go.dev/)
- [BubbleTea](https://github.com/charmbracelet/bubbletea) (and [related libraries](https://github.com/charmbracelet))
- [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API) (using [gorilla/websocket](https://github.com/gorilla/websocket))
- [Argon2](https://en.wikipedia.org/wiki/Argon2) (using [alexedwards/argon2id](https://github.com/alexedwards/argon2id))
- [Session-based Authentication](https://roadmap.sh/guides/session-based-authentication)
- [Redis](https://en.wikipedia.org/wiki/Redis)
- [Postgres](https://en.wikipedia.org/wiki/PostgreSQL)
- [SQLite](https://www.sqlite.org/)

A big focus was on exploring security best practices in a hands-on way. As such, all data is encrypted at rest. I also plan to utilise a [Signal Protocol](https://signal.org/docs/) implementation for end-to-end encryption in the near future. 

## Development

This project makes use of a [taskfile](https://taskfile.dev/) to do various things. You can see available task using:
```shell
task --list-all
```

Alternatively, look at [taskfile.yml](taskfile.yml) for more info.

A simple flow is to test, lint, and then run locally:

```sh
# runs test & lint tasks by default
task
# start the server and its dependencies with docker-compose
task dc:up
# start a client locally
go run ./client
```

You can start several clients using the same command in different terminals. You can also run the client in debug mode which will produce logs relative to the working directory at `logs/`:

```sh
DEBUG=t go run ./client
```

Client-side data is stored using SQLite. The location of the database is system-dependent and uses [this XDG package](https://github.com/adrg/xdg) to determine the location. Here's a quick summary for popular OSs:

- MacOS: `~/Library/Application Support`
- Unix: `~/.local/share`
- Windows: `LocalAppData`

Once you locate the `XDG_DATA_HOME` direcotry for your system, you will find `TeaTime/client.db` within which is a SQLite database containing all user data. The only data stored on the server is for authentication and session management.