# CLI Tool

A simple powerful command-line interface (CLI) tool built in Go that provides a robust framework for executing commands with arguments, options, and dynamic help.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

This CLI tool is designed to provide a flexible and extensible command-line interface for interacting with various functionalities. It supports:

- Subcommands with arguments and options.
- Dynamic help generation for each command.
- Persistent configuration and history support.
- Cross-platform compatibility.

The tool is implemented in Go, ensuring high performance and ease of distribution across different operating systems.

---

## Features

- **Subcommands**: Execute specific actions using subcommands (e.g., `cli version`, `cli help`).
- **Arguments and Options**: Pass arguments and options to customize behavior.
- **Help System**: Comprehensive help system with detailed descriptions for each command.
- **Error Handling**: Graceful error handling with descriptive messages.
- **Extensibility**: Easily add new commands or modify existing ones.
- **Cross-Platform**: Works seamlessly on Windows, macOS, and Linux.

---

## Installation

### Prerequisites

- Go 1.24+ installed on your system.
- Basic knowledge of Go and command-line tools.

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/yassirdeveloper/cli.git
   cd cli
   ```

2. Build the binary:
   ```bash
   go build -o cli .
   ```

3. (Optional) Install the binary globally:
   ```bash
   mv cli /usr/local/bin/cli
   ```

4. Verify the installation:
   ```bash
   cli version
   ```

---

## Usage

Run the CLI tool using the following format:

```bash
cli [command] [arguments] [options]
```

### Examples

1. Display the version:
   ```bash
   cli version
   ```

2. Get help for all commands:
   ```bash
   cli help
   ```

3. Get help for a specific command:
   ```bash
   cli help exit
   ```

4. Exit the application:
   ```bash
   cli exit
   ```

---

## Commands

### `version`
Displays the current version of the CLI tool.

**Usage:**
```bash
cli version
```

**Example Output:**
```
v1.0.0
```

---

### `help`
Provides help information for commands.

**Usage:**
```bash
cli help [command]
```

- Without arguments: Lists help for all commands.
- With a command name: Provides detailed help for the specified command.

**Example:**
```bash
cli help exit
```

**Output:**
```
Exit the application.
Usage: exit
```

---

### `exit`
Exits the application.

**Usage:**
```bash
cli exit
```

---

## Development

To contribute to this project or extend its functionality, follow these steps:

1. Fork the repository and clone it locally.
2. Set up your development environment:
   ```bash
   go mod tidy
   ```
3. Run tests to ensure everything works:
   ```bash
   go test ./...
   ```
4. Make your changes and ensure they pass all tests.
5. Submit a pull request with a clear description of your changes.

---

## Contributing

Contributions are welcome! If you'd like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Write tests for your changes.
4. Ensure all tests pass (`go test ./...`).
5. Submit a pull request with a clear description of your changes.

---

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute it as needed.
