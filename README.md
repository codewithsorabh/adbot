# ADBOT

Telegram Multi-Account Auto-Forwarder & Management System.

## 1. Install

Run this command on your Linux server:

```bash
curl -fsSL https://raw.githubusercontent.com/codewithsorabh/adbot/main/install.sh | bash

```

## 2. Setup

Run the installer:

```bash
installer
```

The installer will prompt you to enter your configuration details:

- **PostgreSQL database name**
- **PostgreSQL username**
- **PostgreSQL password**
- **Telegram Bot Token**
- **Admin ID(s)**
- **Encryption Key** _(Choose between Auto generate or Manual entry)_

Example interactive session:

```text
PostgreSQL Configuration
------------------------
PostgreSQL database name [mybot_db]: mybot_db
PostgreSQL username [mybot_user]: mybot_user
PostgreSQL password:

Telegram Configuration
----------------------
Telegram Bot Token: 123456789:AAxxxxxxxxxxxxxxxx
Admin ID(s), comma-separated: 987654321,123456789

Encryption Key
--------------
1. Auto generate
2. Enter manually

Select option (1/2) [1]: 1
[OK] Encryption key generated

```

After setup completes, the installer creates:

```text
.
├── .env
├── Dockerfile
├── docker-compose.yml
└── dist/
    └── adbot

```

## 3. Run Docker

Build and start all ADBOT services in background mode:

```bash
docker compose up -d --build

```

Manage services:

```bash
# Check container statuses
docker compose ps

# Restart the ADBOT service
docker compose restart bot

# Stop all services
docker compose down

```

## 4. Logs

View ADBOT core application logs:

```bash
docker compose logs -f bot

```

Other useful log commands:

```bash
# View all service logs (Bot, Postgres, Redis)
docker compose logs -f

# View PostgreSQL logs
docker compose logs -f postgres

# View Redis logs
docker compose logs -f redis

```

## Requirements

- **OS:** Linux x86_64 / AMD64
- **Dependencies:** Docker, Docker Compose, `curl`
- **Network:** Active Internet Connection

## Security

Never share or commit your `.env` file. It contains sensitive credentials:

- PostgreSQL password
- Telegram Bot Token
- Admin IDs
- Encryption key

## License
