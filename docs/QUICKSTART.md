# Quick Start Guide

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for running guacd)

## Installation

```bash
go get github.com/QWQ123321123/go-guacamole-client
```

## Setup guacd

Start guacd using Docker Compose:

```bash
docker-compose up -d
```

This will start guacd on port 4822.

## Basic Usage

### Simple SSH Connection

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/QWQ123321123/go-guacamole-client/client"
)

func main() {
    // Create tunnel with default config
    tunnel := client.NewTunnel(client.DefaultConfig())
    
    // Connect to SSH server
    err := tunnel.Connect(
        "your-server.com",
        "username",
        "password",
        22,
        "ssh",
        1024,
        768,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer tunnel.Close()
    
    // Read data
    data, err := tunnel.ReadSome()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(data))
    
    // Write data
    tunnel.Write([]byte("ls\n"))
}
```

### WebSocket Proxy

For a complete WebSocket proxy example, see [examples/websocket_proxy/main.go](../examples/websocket_proxy/main.go).

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

## Running Examples

Make sure guacd is running:

```bash
docker-compose up -d
```

Then run the example:

```bash
go run examples/basic_ssh/main.go
```

Or use the Makefile:

```bash
make run-example
```

## Troubleshooting

### Connection Refused

If you get "connection refused" errors, make sure guacd is running:

```bash
docker-compose ps
```

### Protocol Errors

If you encounter protocol errors, make sure:
- guacd version is compatible (tested with guacd 1.5.0+)
- Connection parameters are correct
- Network connectivity is available
