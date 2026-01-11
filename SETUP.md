# Setup Instructions

## Before Publishing

Before publishing this library to GitHub, you need to:

1. **Update Module Path**: ✅ Already updated to `github.com/QWQ123321123/go-guacamole-client`
   - `go.mod`
   - `examples/basic_ssh/main.go`
   - `examples/websocket_proxy/main.go`
   - `README.md`
   - Any other documentation files

2. **Update License**: Update the copyright holder in `LICENSE` file

3. **Update Author Information**: Update author name and email in:
   - `README.md`
   - `LICENSE`

## Quick Setup

1. Update module path:
```bash
# ✅ Module path already updated to github.com/QWQ123321123/go-guacamole-client
```

2. Update dependencies:
```bash
go mod tidy
```

3. Build and test:
```bash
make build
make test
```

4. Start guacd:
```bash
docker-compose up -d
```

5. Run examples:
```bash
make run-example
```

## Publishing to GitHub

1. Create a new repository on GitHub (e.g., `go-guacamole-client`)

2. Initialize git (if not already):
```bash
git init
git add .
git commit -m "Initial commit"
```

3. Add remote and push:
```bash
git remote add origin https://github.com/QWQ123321123/go-guacamole-client.git
git branch -M main
git push -u origin main
```

4. Create a release tag:
```bash
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0
```

## Using as a Go Module

After publishing, others can use it as:

```bash
go get github.com/QWQ123321123/go-guacamole-client
```
