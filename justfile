# List available commands
default:
    @just --list

# Build the CLI tool
build:
    go build -o bin/kubelog main.go

# Run tests
test:
    go test ./...

# Build with version information
build-version:
    #!/usr/bin/env bash
    VERSION=$(git describe --tags --always --dirty)
    COMMIT_HASH=$(git rev-parse --short HEAD)
    BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    go build -ldflags "-X github.com/dantech2000/kubelog/lib.commitHash=${COMMIT_HASH} -X github.com/dantech2000/kubelog/lib.buildDate=${BUILD_DATE}" -o bin/kubelog main.go

# Cross-compile for multiple platforms
cross-compile:
    #!/usr/bin/env bash
    VERSION=$(git describe --tags --always --dirty)
    COMMIT_HASH=$(git rev-parse --short HEAD)
    BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LDFLAGS="-X github.com/dantech2000/kubelog/lib.commitHash=${COMMIT_HASH} -X github.com/dantech2000/kubelog/lib.buildDate=${BUILD_DATE}"
    
    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/kubelog-linux-amd64 main.go
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/kubelog-darwin-amd64 main.go
    GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/kubelog-windows-amd64.exe main.go

# Run linter
lint:
    golangci-lint run

# Format code
fmt:
    go fmt ./...

# Clean build artifacts
clean:
    rm -rf bin/*

# Install dependencies
deps:
    go mod tidy
    go mod verify

# Run the CLI tool
run *ARGS:
    go run main.go {{ ARGS }}

# Create a new release tag
release VERSION:
    git tag -a v{{ VERSION }} -m "Release v{{ VERSION }}"
    git push origin v{{ VERSION }}

# Generate and update changelog
changelog:
    git-chglog -o CHANGELOG.md

# Build and run in one step
build-and-run *ARGS: build
    ./bin/kubelog {{ ARGS }}