# Canvus CLI

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/jaypaulb/canvus-cli)](https://github.com/jaypaulb/canvus-cli/releases)

Command-line interface for the Canvus collaborative workspace API.

## Overview

The Canvus CLI provides a powerful command-line interface for managing Canvus workspaces, canvases, widgets, users, and more. Built on top of the [Canvus Go SDK](https://github.com/jaypaulb/Canvus-Go-API), it offers complete access to the Canvus API from your terminal.

## Installation

### Using Go

```bash
go install github.com/jaypaulb/canvus-cli/cmd/canvus@latest
```

### Binary Download

Download the pre-built binary for your platform from the [Releases](https://github.com/jaypaulb/canvus-cli/releases) page.

#### macOS (Apple Silicon)
```bash
curl -L https://github.com/jaypaulb/canvus-cli/releases/latest/download/canvus_darwin_arm64.tar.gz | tar xz
sudo mv canvus /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -L https://github.com/jaypaulb/canvus-cli/releases/latest/download/canvus_darwin_amd64.tar.gz | tar xz
sudo mv canvus /usr/local/bin/
```

#### Linux (x86_64)
```bash
curl -L https://github.com/jaypaulb/canvus-cli/releases/latest/download/canvus_linux_amd64.tar.gz | tar xz
sudo mv canvus /usr/local/bin/
```

#### Linux (ARM64)
```bash
curl -L https://github.com/jaypaulb/canvus-cli/releases/latest/download/canvus_linux_arm64.tar.gz | tar xz
sudo mv canvus /usr/local/bin/
```

#### Windows
Download `canvus_windows_amd64.zip` from the releases page and add the binary to your PATH.

## Quick Start

### 1. Configure Your Environment

Set your Canvus server URL and authentication:

```bash
export CANVUS_URL="https://your-canvus-server.com"
export CANVUS_API_KEY="your-api-key"
```

Or use command-line flags:
```bash
canvus --url https://your-canvus-server.com --api-key your-api-key <command>
```

### 2. Verify Connection

```bash
canvus system info
```

### 3. List Canvases

```bash
canvus canvas list
```

### 4. Get Canvas Details

```bash
canvus canvas get <canvas-id>
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CANVUS_URL` | Base URL of your Canvus server | Yes |
| `CANVUS_API_KEY` | API key for authentication | Yes* |
| `CANVUS_USERNAME` | Username for login authentication | No |
| `CANVUS_PASSWORD` | Password for login authentication | No |
| `CANVUS_INSECURE` | Skip TLS verification (not recommended) | No |

*Required if not using username/password authentication

### Configuration File

Create `~/.canvus/config.yaml`:

```yaml
url: https://your-canvus-server.com
api_key: your-api-key
# Or use login credentials:
# username: your-username
# password: your-password
insecure: false
```

### Command-Line Flags

Global flags available for all commands:

```
--url string        Canvus server URL
--api-key string    API key for authentication
--username string   Username for authentication
--password string   Password for authentication
--insecure          Skip TLS certificate verification
--output string     Output format: json, yaml, table (default "table")
--verbose           Enable verbose output
--help              Show help for command
```

## Commands

### System Commands

```bash
# Get system information
canvus system info

# Get license information
canvus system license
```

### Canvas Commands

```bash
# List all canvases
canvus canvas list

# Get canvas details
canvus canvas get <canvas-id>

# Create a new canvas
canvus canvas create --name "My Canvas"

# Update a canvas
canvus canvas update <canvas-id> --name "New Name"

# Delete a canvas
canvus canvas delete <canvas-id>

# Copy a canvas
canvus canvas copy <canvas-id> --name "Canvas Copy"

# Export a canvas
canvus canvas export <canvas-id> --output ./export/

# Import a canvas
canvus canvas import ./export/canvas.json
```

### Widget Commands

```bash
# List widgets on a canvas
canvus widget list <canvas-id>

# Get widget details
canvus widget get <canvas-id> <widget-id>

# Create a note widget
canvus widget create <canvas-id> --type note --text "Hello World"

# Create a browser widget
canvus widget create <canvas-id> --type browser --url "https://example.com"

# Update a widget
canvus widget update <canvas-id> <widget-id> --text "Updated text"

# Delete a widget
canvus widget delete <canvas-id> <widget-id>

# Search widgets across all canvases
canvus widget search --text "keyword"
```

### User Commands

```bash
# List all users
canvus user list

# Get user details
canvus user get <user-id>

# Create a user
canvus user create --name "John Doe" --username "jdoe" --email "jdoe@example.com"

# Update a user
canvus user update <user-id> --name "Jane Doe"

# Delete a user
canvus user delete <user-id>

# Create an API token for a user
canvus user token create <user-id> --name "CI Token"

# List user tokens
canvus user token list <user-id>
```

### Group Commands

```bash
# List all groups
canvus group list

# Get group details
canvus group get <group-id>

# Create a group
canvus group create --name "Developers"

# Add user to group
canvus group add-user <group-id> <user-id>

# Remove user from group
canvus group remove-user <group-id> <user-id>
```

## Output Formats

### Table (Default)

```bash
canvus canvas list
```

```
ID                                   NAME            CREATED
550e8400-e29b-41d4-a716-446655440000 Project Alpha   2024-01-15
550e8400-e29b-41d4-a716-446655440001 Project Beta    2024-01-16
```

### JSON

```bash
canvus canvas list --output json
```

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Project Alpha",
    "created": "2024-01-15T10:30:00Z"
  }
]
```

### YAML

```bash
canvus canvas list --output yaml
```

```yaml
- id: 550e8400-e29b-41d4-a716-446655440000
  name: Project Alpha
  created: 2024-01-15T10:30:00Z
```

## Examples

### Batch Create Widgets

```bash
# Create multiple note widgets from a file
for note in $(cat notes.txt); do
  canvus widget create <canvas-id> --type note --text "$note"
done
```

### Export and Import Canvas

```bash
# Export canvas with all assets
canvus canvas export abc123 --output ./backup/ --include-assets

# Import to a new server
CANVUS_URL=https://new-server.com canvus canvas import ./backup/canvas.json
```

### Automated User Provisioning

```bash
# Create users from CSV
while IFS=, read -r name username email; do
  canvus user create --name "$name" --username "$username" --email "$email"
done < users.csv
```

### Search Widgets Across Canvases

```bash
# Find all widgets containing "TODO"
canvus widget search --text "TODO" --output json | jq '.[] | {canvas: .canvas_id, id: .id, text: .text}'
```

## Scripting

The CLI is designed for scripting and automation:

```bash
#!/bin/bash
set -e

# Get canvas ID by name
CANVAS_ID=$(canvus canvas list --output json | jq -r '.[] | select(.name=="Project") | .id')

# Create widgets on the canvas
canvus widget create "$CANVAS_ID" --type note --text "Task 1" --x 100 --y 100
canvus widget create "$CANVAS_ID" --type note --text "Task 2" --x 300 --y 100

# List created widgets
canvus widget list "$CANVAS_ID"
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication error |
| 3 | Resource not found |
| 4 | Permission denied |
| 5 | Invalid input |

## Troubleshooting

### Connection Refused

```
Error: connection refused
```

Check that:
- `CANVUS_URL` is correct
- The Canvus server is running
- No firewall is blocking the connection

### Authentication Failed

```
Error: authentication failed: invalid API key
```

Verify your API key:
- Check `CANVUS_API_KEY` is set correctly
- Ensure the API key has not expired
- Verify the key has appropriate permissions

### Certificate Errors

```
Error: x509: certificate signed by unknown authority
```

Options:
1. Add the CA certificate to your system trust store
2. Use `--insecure` flag (not recommended for production)

### Verbose Mode

Enable verbose output to see request/response details:

```bash
canvus --verbose canvas list
```

## SDK Documentation

For programmatic access, see the [Canvus Go SDK](https://github.com/jaypaulb/Canvus-Go-API):

- [Getting Started](https://github.com/jaypaulb/Canvus-Go-API/blob/main/docs/GETTING_STARTED.md)
- [API Reference](https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus)
- [Best Practices](https://github.com/jaypaulb/Canvus-Go-API/blob/main/docs/BEST_PRACTICES.md)
- [Examples](https://github.com/jaypaulb/Canvus-Go-API/tree/main/examples)

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- [GitHub Issues](https://github.com/jaypaulb/canvus-cli/issues) - Bug reports and feature requests
- [SDK Repository](https://github.com/jaypaulb/Canvus-Go-API) - For SDK-related issues
