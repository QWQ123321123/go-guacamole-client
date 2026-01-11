# Contributing

Thank you for your interest in contributing to go-guacamole-client!

## How to Contribute

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/QWQ123321123/go-guacamole-client.git
cd go-guacamole-client
```

2. Install dependencies:
```bash
go mod download
```

3. Start guacd:
```bash
docker-compose up -d
```

4. Run tests:
```bash
go test ./...
```

## Code Style

- Follow Go standard formatting (`gofmt`)
- Add comments for exported functions and types
- Write tests for new features
- Keep code simple and readable

## Testing

- Write unit tests for new features
- Test with real guacd instances
- Ensure backward compatibility

## Pull Request Guidelines

- Keep PRs focused on a single feature or bug fix
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass

## Reporting Issues

When reporting issues, please include:

- Go version
- guacd version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Relevant logs

Thank you for contributing!
