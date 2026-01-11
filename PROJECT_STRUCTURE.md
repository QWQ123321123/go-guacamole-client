# Project Structure

```
go-guacamole-client/
├── client/                  # Core client library
│   ├── tunnel.go           # Main Tunnel implementation
│   └── instruction.go      # Instruction utility functions
├── examples/               # Example applications
│   ├── basic_ssh/         # Basic SSH connection example
│   │   └── main.go
│   └── websocket_proxy/   # WebSocket proxy example
│       └── main.go
├── docs/                   # Documentation
│   ├── ARCHITECTURE.md    # Architecture documentation
│   └── QUICKSTART.md      # Quick start guide
├── docker-compose.yml     # Docker Compose config for guacd
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Makefile               # Build automation
├── README.md              # Main documentation
├── LICENSE                # MIT License
├── CONTRIBUTING.md        # Contribution guidelines
└── .gitignore            # Git ignore rules
```

## Directory Descriptions

### client/

Core client library containing:

- `tunnel.go`: Main `Tunnel` struct and methods for connecting to guacd
- `instruction.go`: Utility functions for Guacamole protocol instructions

### examples/

Example applications demonstrating library usage:

- `basic_ssh/`: Simple SSH connection example
- `websocket_proxy/`: Complete WebSocket proxy implementation

### docs/

Documentation files:

- `ARCHITECTURE.md`: Detailed architecture documentation
- `QUICKSTART.md`: Quick start guide for new users

## Key Files

### go.mod

Defines the module path and dependencies. ✅ Already updated to `github.com/QWQ123321123/go-guacamole-client`.

### docker-compose.yml

Docker Compose configuration for running guacd daemon. Uses the official `guacamole/guacd:latest` image.

### Makefile

Common build tasks:

- `make build`: Build the library
- `make test`: Run tests
- `make clean`: Clean build artifacts
- `make run-example`: Run the basic SSH example
- `make docker-up`: Start guacd
- `make docker-down`: Stop guacd

## Usage

After cloning, update the module path in `go.mod` and example files:

1. ✅ Module path already updated to `github.com/QWQ123321123/go-guacamole-client`
2. Run `go mod tidy` to update dependencies
3. Build and test: `make build && make test`
