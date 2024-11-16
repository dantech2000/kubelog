# Kubelog

Kubelog is a CLI tool to fetch and enhance Kubernetes pod logs. It simplifies the retrieval and parsing of Kubernetes pod logs, providing enhanced formatting and filtering options for efficient troubleshooting.

## Features

- 🎯 **Smart Log Parsing**
  - Automatic detection of JSON and plain text log formats
  - Intelligent timestamp parsing across multiple formats
  - Log level detection (DEBUG, INFO, WARN, ERROR, FATAL)
  - Structured field parsing for JSON logs

- 🎨 **Beautiful Output Formatting**
  - Color-coded log levels and timestamps
  - Consistent timestamp formatting
  - Highlighted error and warning messages
  - Clean key-value formatting for JSON fields
  - Logger type identification (e.g., logrus, zap)

- 🚀 **Kubernetes Integration**
  - Easy container selection with interactive prompts
  - Support for multi-container pods
  - Previous container logs with `-p` flag
  - Real-time log following with `-f` flag
  - Container status indicators

- ⚡ **Performance**
  - Efficient log streaming
  - Smart container name completion
  - Optimized log parsing

## Installation

### Using Homebrew

```bash
brew tap dantech2000/tap
brew install kubelog
```

### Manual Installation

### Prerequisites

- Go 1.16 or later
- Access to a Kubernetes cluster
- kubectl configured with the appropriate context
- [just](https://github.com/casey/just) command runner (optional, but recommended)

### Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/dantech2000/kubelog.git
    ```

2. Navigate to the project directory:

    ```bash
    cd kubelog
    ```

3. Build the binary:

    ```bash
    # Using just (recommended, includes version information)
    just build-version

    # Or using Go directly
    go build -o bin/kubelog main.go
    ```

4. (Optional) Move the binary to a directory in your PATH:
    ```bash
    sudo mv bin/kubelog /usr/local/bin/
    ```

## Usage

### Fetching Logs

To fetch logs from a pod:

```bash
kubelog logs [pod-name] -n [namespace]
```

Options:

- `-n, --namespace`: Specify the Kubernetes namespace (default is "default")
- `-c, --container`: Specify the container name (if pod has multiple containers)
- `-f, --follow`: Follow the log output (similar to `tail -f`)
- `-l, --level`: Filter logs by level (DEBUG, INFO, WARN, ERROR)

Example:

```bash
kubelog logs my-pod -n my-namespace -c my-container -f -l INFO
```

### Listing Containers

To list containers in a pod:

```bash
kubelog containers [pod-name] -n [namespace]
```

Options:

- `-n, --namespace`: Specify the Kubernetes namespace (default is "default")

Example:

```bash
kubelog containers my-pod -n my-namespace
```

### Version Information

To display version information:

```bash
kubelog version
```

Options:

- `-s, --short`: Display only the version number
- `-o, --output`: Output format (json or yaml)

Examples:

```bash
# Display full version information
kubelog version

# Display only version number
kubelog version --short

# Get version info in JSON format
kubelog version --output json

# Get version info in YAML format
kubelog version --output yaml
```

## Development

### Available Make Commands

```bash
just --list
```

Common commands:

- `just build-version`: Build with version information
- `just test`: Run tests
- `just lint`: Run linter
- `just fmt`: Format code
- `just clean`: Clean build artifacts
- `just deps`: Install dependencies
- `just cross-compile`: Build for multiple platforms

### Creating a Release

1. Update the version in `lib/version.go`

2. Commit your changes:
    ```bash
    git add .
    git commit -m "Bump version to X.Y.Z"
    ```

3. Create and push a new tag:
    ```bash
    git tag -a vX.Y.Z -m "Release vX.Y.Z"
    git push origin vX.Y.Z
    ```

This will trigger the GitHub Actions workflow which will:
- Build the project
- Create a GitHub release
- Upload the binaries
- Update the Homebrew tap

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
