# Kubelog

Kubelog is a CLI tool to fetch and enhance Kubernetes pod logs. It simplifies the retrieval and parsing of Kubernetes pod logs, providing enhanced formatting and filtering options for efficient troubleshooting.

## Installation

### Prerequisites

-   Go 1.16 or later
-   Access to a Kubernetes cluster
-   kubectl configured with the appropriate context

### Steps

1. Clone the repository:

    ```
    git clone https://github.com/dantech2000/kubelog.git
    ```

2. Navigate to the project directory:

    ```
    cd kubelog
    ```

3. Build the binary:

    ```
    go build -o kubelog
    ```

4. (Optional) Move the binary to a directory in your PATH:
    ```
    sudo mv kubelog /usr/local/bin/
    ```

## Usage

### Fetching Logs

To fetch logs from a pod:

```
kubelog logs [pod-name] -n [namespace]
```

Options:

-   `-n, --namespace`: Specify the Kubernetes namespace (default is "default")
-   `-c, --container`: Specify the container name (if pod has multiple containers)
-   `-f, --follow`: Follow the log output (similar to `tail -f`)

Example:

```
kubelog logs my-pod -n my-namespace -c my-container -f
```

### Listing Containers

To list containers in a pod:

```
kubelog containers [pod-name] -n [namespace]
```

Options:

-   `-n, --namespace`: Specify the Kubernetes namespace (default is "default")
-   `-c, --container`: Specify the container name (if pod has multiple containers)

Example:

```
kubelog containers my-pod -n my-namespace -c my-container
```
