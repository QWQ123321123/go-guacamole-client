# Architecture

## Overview

`go-guacamole-client` is a Go client library for Apache Guacamole protocol. It provides a clean interface to connect to remote servers via SSH, RDP, VNC, etc. through the guacd daemon.

## Components

### Client Package

The `client` package contains the core functionality:

- **Tunnel**: Main client struct that manages the connection to guacd
- **Config**: Configuration for guacd connection
- **Instruction utilities**: Helper functions for Guacamole protocol

### Handshake Process

The library implements the Guacamole protocol handshake manually:

1. **Select Protocol**: Send `select` instruction to choose protocol (SSH, RDP, etc.)
2. **Parse Args**: Wait for and parse `args`/`required` response from guacd
3. **Send Capabilities**: Send client capabilities (size, audio, video, image, timezone)
4. **Connect**: Send `connect` instruction with connection parameters
5. **Ready**: Wait for `ready` response from guacd

### Data Flow

```
Application → Tunnel.Connect() → guacd → Target Server
Application ← Tunnel.ReadSome() ← guacd ← Target Server
Application → Tunnel.Write()    → guacd → Target Server
```

## Design Decisions

### Manual Handshake

Instead of using the `Handshake()` method from `seknox/guacamole` library, we implement the handshake manually. This provides:

- Better compatibility with different guacd versions
- More control over the handshake process
- Better error handling

### Instruction Building

We use `seknox/guacamole` library's `NewInstruction()` to build instructions. This ensures:

- Correct instruction format
- Automatic length calculation
- Protocol compliance

### Raw Byte Reading

After handshake, we read raw bytes directly from the connection instead of using the library's parser. This avoids:

- Parser errors with certain instructions
- Compatibility issues
- Performance overhead

## Usage Patterns

### Basic Usage

```go
config := client.DefaultConfig()
tunnel := client.NewTunnel(config)
err := tunnel.Connect(hostname, username, password, port, "ssh", 1024, 768)
// ... use tunnel
tunnel.Close()
```

### WebSocket Proxy

The library is designed to be used with WebSocket proxies. The `ReadSome()` method is optimized for this use case:

- Returns `nil, nil` on timeout (not an error)
- Reads in chunks for efficient buffering
- Supports instruction boundary detection

## Dependencies

- `github.com/seknox/guacamole`: For instruction building
- `go.uber.org/zap`: For logging (optional, can be removed for minimal dependencies)

## Performance Considerations

- **Buffering**: The library uses `bufio.Reader` for efficient reading
- **Instruction Boundaries**: The `FindLastCompleteInstruction()` function ensures data integrity
- **Timeouts**: Read operations have timeouts to prevent blocking
- **Connection Pooling**: Not included - create new tunnels as needed
