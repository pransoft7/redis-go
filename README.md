# Redis-Go 🚀

A lightweight, high-performance Redis clone built from scratch in Go. This project demonstrates understanding of network programming, concurrent systems design, and the Redis Serialization Protocol (RESP).

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)

## ✨ Features

- **Custom TCP Server** - Built using Go's `net` package for TCP connection handling
- **RESP Protocol Support** - Compatible with Redis clients via RESP parsing (using `tidwall/resp`)
- **Thread-Safe Key-Value Store** - Concurrent-safe in-memory data structure using Go's sync primitives
- **Multi-Client Support** - Handle concurrent client connections with goroutines
- **Core Commands** - Support for fundamental Redis commands:
  - `SET key value` - Store a key-value pair
  - `GET key` - Retrieve value by key
- **Native Go Client** - Included client library for seamless integration

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      TCP Server                         │
│  ┌─────────────┐     ┌─────────────┐     ┌───────────┐  │
│  │ Accept Loop │────▶│ Peer Handler│────▶│ Message   │  │
│  │             │     │ (goroutine) │     │ Channel   │  │
│  └─────────────┘     └─────────────┘     └─────┬─────┘  │
│                                                │        │
│  ┌─────────────────────────────────────────────▼─────┐  │
│  │                   Event Loop                      │  │
│  │  • Process commands  • Manage peer connections    │  │
│  └─────────────────────────────────────────────┬─────┘  │
│                                                │        │
│  ┌─────────────────────────────────────────────▼─────┐  │
│  │              Key-Value Store (sync.RWMutex)       │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Key Components

| File | Description |
|------|-------------|
| `main.go` | Server initialization, event loop, and command handling |
| `peer.go` | Client connection management and RESP message parsing |
| `proto.go` | Command definitions and RESP protocol utilities |
| `keyval.go` | Thread-safe key-value store implementation |
| `client/client.go` | Go client library for connecting to the server |

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher

### Installation

```bash
git clone https://github.com/pransoft7/redis-go.git
cd redis-go
go mod download
```

### Running the Server

```bash
# Build and run with default port (5432)
make run

# Or run with custom port
go build -o bin/redis-go
./bin/redis-go --listenAddr :6379
```

### Using the Client

```go
package main

import (
    "context"
    "fmt"
    "redis-go/client"
)

func main() {
    // Connect to the server
    c, err := client.New("localhost:5432")
    if err != nil {
        panic(err)
    }
    defer c.Close()

    // Set a value
    err = c.Set(context.Background(), "greeting", "Hello, Redis-Go!")
    if err != nil {
        panic(err)
    }

    // Get the value
    val, err := c.Get(context.Background(), "greeting")
    if err != nil {
        panic(err)
    }
    fmt.Println(val) // Output: Hello, Redis-Go!
}
```

### Using with redis-cli

The server uses RESP protocol and can accept connections from redis-cli:

```bash
redis-cli -p 5432
127.0.0.1:5432> SET mykey myvalue
127.0.0.1:5432> GET mykey
"myvalue"
```

> **Note:** SET currently doesn't return an OK response. Full redis-cli compatibility is a work in progress.

## 🧪 Testing

Run the test suite to verify multi-client concurrent access:

```bash
go test -v ./...
```

The tests validate:
- Server initialization and startup
- Multiple concurrent client connections
- SET/GET operations under load
- Proper connection cleanup

## 🛠️ Technical Highlights

### Networking
- **TCP Server** - Uses Go's `net.Listen()` and `net.Accept()` for connection handling
- **Connection Lifecycle Management** - Custom `Peer` abstraction to track client state and handle disconnections
- **Event Loop in Goroutine** - Central event loop runs in a separate goroutine while accept loop runs on main
- **Persistent Connections** - Reading client requests and writing responses over long-lived TCP connections

### Concurrency Model
- **Channel-based Architecture** - Decoupled components communicate via typed channels (`msgCh`, `addPeerCh`, `delPeerCh`)
- **Goroutine-per-Connection** - Each client gets a dedicated goroutine for reading, enabling true parallelism
- **RWMutex** - Read-write lock on the key-value store allows concurrent reads while serializing writes
- **Select-based Event Loop** - Central dispatcher multiplexes events from multiple channels

### Event Loop Pattern
The server uses a single-threaded event loop (similar to Redis) that processes all state mutations:
```go
func (s *Server) loop() {
    for {
        select {
        case msg := <-s.msgCh:      // Handle incoming commands
        case peer := <-s.addPeerCh: // New connection
        case peer := <-s.delPeerCh: // Disconnection
        case <-s.quitCh:            // Graceful shutdown
        }
    }
}
```

### RESP Protocol
- Uses `tidwall/resp` library for Redis Serialization Protocol parsing
- Command parsing layer built on top to interpret SET/GET operations
- Extensible command structure using Go interfaces

## 📦 Dependencies

| Package | Purpose |
|---------|---------|
| [tidwall/resp](https://github.com/tidwall/resp) | RESP protocol parsing and serialization |

## 🗺️ Roadmap

- [ ] Add support for more Redis commands (DEL, EXISTS, EXPIRE, TTL)
- [ ] Implement data persistence (RDB/AOF)
- [ ] Add pub/sub functionality
- [ ] Implement Redis transactions (MULTI/EXEC)
- [ ] Add cluster mode support

