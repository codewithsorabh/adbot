# ADBOT

[![Release](https://img.shields.io/github/v/release/codewithsorabh/adbot?include_prereleases)](https://github.com/codewithsorabh/adbot/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/codewithsorabh/adbot)](https://goreportcard.com/report/github.com/codewithsorabh/adbot)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**ADBOT** is an automated deployment wizard and management utility for the ADBOT Telegram service. It sets up database credentials, Telegram bot tokens, secure encryption keys, and provisions Docker Compose services (PostgreSQL, Redis, and ADBOT container) with a single command.

---

## Table of Contents

- [Supported Platforms](#supported-platforms)
- [Quick Installation](#quick-installation)
  - [Install Latest Release](#install-latest-release)
  - [Install Specific Version](#install-specific-version)
- [Usage](#usage)
  - [Interactive Setup Wizard](#interactive-setup-wizard)
  - [Running the Services](#running-the-services)
  - [Managing Containers & Logs](#managing-containers--logs)
  - [Checking for Updates](#checking-for-updates)
- [Local Development](#local-development)
  - [Prerequisites](#prerequisites)
  - [Build and Run](#build-and-run)
  - [Running Tests](#running-tests)
- [Docker Usage](#docker-usage)
  - [Building Docker Image from Source](#building-docker-image-from-source)
- [Release Process](#release-process)
  - [Versioning](#versioning)
  - [Creating a Release](#creating-a-release)
- [Security](#security)
- [License](#license)

---

## Supported Platforms

ADBOT is distributed as statically-compiled native binaries (`CGO_ENABLED=0`) for the following operating systems and architectures:

| Operating System | Architecture | Binary Asset Name |
| :--- | :--- | :--- |
| **Linux** | `amd64` (x86_64) | `adbot-v<VERSION>-linux-amd64` |
| **Linux** | `arm64` (aarch64) | `adbot-v<VERSION>-linux-arm64` |
| **macOS (Darwin)** | `amd64` (Intel) | `adbot-v<VERSION>-darwin-amd64` |
| **macOS (Darwin)** | `arm64` (Apple Silicon) | `adbot-v<VERSION>-darwin-arm64` |
| **Windows** | `amd64` (x86_64) | `adbot-v<VERSION>-windows-amd64.exe` |

---

## Quick Installation

### Install Latest Release

Run the official installer script to automatically detect your OS/architecture, download the latest release binary from GitHub Releases, verify its SHA256 checksum, and configure your `$PATH`:

```bash
curl -fsSL https://raw.githubusercontent.com/codewithsorabh/adbot/main/install.sh | bash
```

### Install Specific Version

You can specify a target version using the `VERSION` environment variable:

```bash
curl -fsSL https://raw.githubusercontent.com/codewithsorabh/adbot/main/install.sh | VERSION=1.0.0 bash
```

---

## Usage

### Interactive Setup Wizard

Once installed, start the interactive setup wizard:

```bash
adbot
```

The wizard prompts for:
1. **PostgreSQL Configuration**: Database name, username, and password.
2. **Telegram Bot Token**: The bot API token provided by [@BotFather](https://t.me/BotFather).
3. **Telegram Admin IDs**: Comma-separated numeric Telegram user IDs.
4. **Encryption Key**: Automatically generated 32-byte Fernet key or user-provided key.

The installer generates the project layout:
```
adbot/
├── .env
├── Dockerfile
├── docker-compose.yml
└── dist/
    └── adbot
```

### Running the Services

Navigate into the generated directory and launch the background Docker services:

```bash
cd adbot
docker compose up -d --build
```

### Managing Containers & Logs

```bash
# Check status of running containers
docker compose ps

# View live application logs
docker compose logs -f bot

# View PostgreSQL or Redis logs
docker compose logs -f postgres
docker compose logs -f redis

# Restart the bot container
docker compose restart bot

# Stop all services
docker compose down
```

### Checking for Updates

Check your current version:
```bash
adbot --version
```

Check and download the latest release from GitHub:
```bash
adbot --check-update
```

---

## Local Development

### Prerequisites
- [Go](https://go.dev/dl/) 1.22 or higher
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/)
- Git

### Build and Run

Clone the repository:
```bash
git clone https://github.com/codewithsorabh/adbot.git
cd adbot
```

Download Go dependencies:
```bash
go mod tidy
```

Build the local binary:
```bash
go build -o bin/adbot ./cmd/adbot
```

Run directly:
```bash
go run ./cmd/adbot
```

### Running Tests

Run all unit tests:
```bash
go test -v ./...
```

Run code formatting and static analysis:
```bash
gofmt -w .
go vet ./...
```

---

## Docker Usage

### Building Docker Image from Source

You can build a minimal multi-stage Docker image directly from the repository source:

```bash
docker build -t adbot:latest .
```

Run the container:
```bash
docker run --rm -it adbot:latest --version
```

---

## Release Process

### Versioning

ADBOT follows [Semantic Versioning (SemVer)](https://semver.org/). The root `VERSION` file defines the current version in `MAJOR.MINOR.PATCH` format (e.g., `1.0.0`).

### Creating a Release

GitHub Actions automatically builds multi-platform binaries, generates SHA256 checksums, and publishes the GitHub Release whenever a Git tag matching `v*` is pushed.

To create a new release:

1. Update the version number in `VERSION`:
   ```bash
   echo "1.0.1" > VERSION
   ```

2. Commit the change:
   ```bash
   git add VERSION
   git commit -m "release: v1.0.1"
   ```

3. Create the release tag:
   ```bash
   git tag v1.0.1
   ```

4. Push the branch and the tag to GitHub:
   ```bash
   git push origin main
   git push origin v1.0.1
   ```

5. GitHub Actions builds all binaries and attaches them to the new [GitHub Release](https://github.com/codewithsorabh/adbot/releases).

---

## Security

- Never commit your `.env` file or sensitive credentials to version control.
- Generated `.env` files and encryption keys are created with restricted file permissions (`0600`).
- Checksums are provided in `checksums.txt` for every release binary to verify download integrity.

---

## License

This project is licensed under the [MIT License](LICENSE).
