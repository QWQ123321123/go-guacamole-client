# go-guacamole-client

A Go client library for Apache Guacamole protocol, providing a simple and reliable way to connect to remote servers via SSH, RDP, VNC, etc. through guacd daemon.

## Features

- ✅ **Manual Handshake**: Implements Guacamole protocol handshake manually for better compatibility
- ✅ **Instruction Building**: Uses `seknox/guacamole` library for correct instruction format
- ✅ **Raw Byte Reading**: Reads raw bytes after handshake to avoid parser errors
- ✅ **Simple API**: Clean and easy-to-use interface
- ✅ **Production Ready**: Optimized for real-world usage

## Architecture

```
Your Go Application     go-guacamole-client     guacd              Target Server
      │                          │                │                     │
      │  Connect()               │                │                     │
      ├─────────────────────────>│                │                     │
      │                          │  TCP (Guacamole Protocol)            │
      │                          ├────────────────>│                     │
      │                          │                │  SSH/RDP/VNC        │
      │                          │                ├─────────────────────>│
      │                          │                │                     │
      │  ReadSome() / Write()    │                │                     │
      │<─────────────────────────┤                │                     │
```

## Installation

```bash
go get github.com/QWQ123321123/go-guacamole-client
```

## Quick Start

### 1. Start guacd

```bash
docker-compose up -d
```

### 2. Use the Client

```go
package main

import (
    "fmt"
    "github.com/QWQ123321123/go-guacamole-client/client"
)

func main() {
    // Create tunnel with default config (localhost:4822)
    config := client.DefaultConfig()
    tunnel := client.NewTunnel(config)

    // Connect to target server via SSH
    err := tunnel.Connect(
        "your-server.com", // hostname
        "username",        // username
        "password",        // password
        22,                // port
        "ssh",             // protocol
        1024,              // width
        768,               // height
    )
    if err != nil {
        panic(err)
    }
    defer tunnel.Close()

    // Read data from tunnel
    data, err := tunnel.ReadSome()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))

    // Write data to tunnel
    _, err = tunnel.Write([]byte("ls\n"))
    if err != nil {
        panic(err)
    }
}
```

## Configuration

### Default Config

```go
config := client.DefaultConfig()
// GuacdHost: "localhost"
// GuacdPort: 4822
```

### Custom Config

```go
config := &client.Config{
    GuacdHost: "192.168.1.100",
    GuacdPort: 4822,
}
tunnel := client.NewTunnel(config)
```

## API Reference

### Tunnel

#### `NewTunnel(config *Config) *Tunnel`

Creates a new Guacamole tunnel. If config is nil, uses default config.

#### `Connect(hostname, username, password string, port int, protocol string, width, height int) error`

Connects to guacd and establishes connection to target server.

**Parameters:**
- `hostname`: Target server address
- `username`: Username for authentication
- `password`: Password for authentication
- `port`: Target server port
- `protocol`: Protocol type (`ssh`, `rdp`, `vnc`, etc.)
- `width`: Display width in pixels
- `height`: Display height in pixels

**Returns:** Error if connection fails

#### `ReadSome() ([]byte, error)`

Reads some data from the tunnel. Returns nil, nil on timeout (no data available).

#### `Write(p []byte) (n int, err error)`

Writes data to the tunnel. Implements `io.Writer` interface.

#### `Read(p []byte) (n int, err error)`

Reads data from the tunnel. Implements `io.Reader` interface.

#### `Available() bool`

Checks if there's more data available to read.

#### `Close() error`

Closes the tunnel connection.

### Utilities

#### `FindLastCompleteInstruction(data []byte) int`

Finds the position of the last complete Guacamole instruction in the data buffer.
Useful for ensuring instruction integrity when buffering data.

## Supported Protocols

- SSH (`ssh`)
- RDP (`rdp`)
- VNC (`vnc`)
- Telnet (`telnet`)
- Kubernetes (`kubernetes`)

## Examples

See the [examples](examples/) directory for complete working examples:

- [Basic SSH Connection](examples/basic_ssh/main.go)
- [WebSocket Proxy](examples/websocket_proxy/main.go)

## Requirements

- Go 1.21+
- guacd daemon (Apache Guacamole daemon)
- Docker (optional, for running guacd)

## Docker Compose

A `docker-compose.yml` file is provided for easy guacd setup:

```bash
docker-compose up -d
```

This starts guacd on port 4822.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [Apache Guacamole](https://guacamole.apache.org/) - The underlying protocol
- [seknox/guacamole](https://github.com/seknox/guacamole) - Go library for instruction building

## Author

Your Name - your.email@example.com
