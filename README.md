# GoRedis Clone

A lightweight Redis-like key-value store server implemented in Go. It supports basic Redis commands and can be used with official Redis clients.

## Features
- RESP protocol support (compatible with Redis clients)
- Basic commands: `SET`, `GET`, `HELLO`, `CLIENT`
- Optional TTL (Time-To-Live) for keys
- Automatic cleanup of expired keys
- Concurrent client handling

## Getting Started

### Prerequisites
- Go 1.18 or newer

### Installation
1. Clone the repository:
   ```sh
   git clone <your-repo-url>
   cd Redis-clone
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```

### Running the Server
Start the server with the default listen address (`:5001`):
```sh
go run main.go
```
Or specify a custom address:
```sh
go run main.go -listenAddr=:6379
```

## Usage
You can connect to the server using any Redis client. Example using the official Go Redis client:
```go
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:5001",
})
rdb.Set(ctx, "foo", "bar", 0)
val, _ := rdb.Get(ctx, "foo").Result()
```

### Supported Commands
- `SET key value [ttl]` — Set a key-value pair with optional TTL
- `GET key` — Get the value for a key
- `HELLO` — Get server info
- `CLIENT` — Client-related command

## Testing
Run unit tests:
```sh
go test -v
```

## Project Structure
- `main.go` — Server setup and event loop
- `keyval.go` — Key-value store logic
- `proto.go` — Command definitions and RESP helpers
- `peer.go` — Client connection handling
- `server_test.go` — Integration tests

## License
MIT
