package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	EnvFile     = ".env"
	DockerFile  = "Dockerfile"
	ComposeFile = "docker-compose.yml"

	DistDir    = "dist"
	OutputFile = "dist/adbot"
)

// Parameters holds configuration settings collected from the setup wizard.
type Parameters struct {
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	BotToken         string
	AdminIDs         string
	EncryptionKey    string
}

// GenerateEnv creates the .env file contents.
func GenerateEnv(p Parameters) string {
	return fmt.Sprintf(`#
# ADBOT ENVIRONMENT CONFIGURATION
# Generated automatically by installer
#

# Debug
DEBUG=True
MAX_CONCURRENT_WORKERS=2000
SCHEDULER_INTERVAL=60

# PostgreSQL
POSTGRES_DB=%s
POSTGRES_USER=%s
POSTGRES_PASSWORD=%s
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Telegram Bot
TELEGRAM_BOT_TOKEN=%s

# Telegram Admin IDs
ADMIN_ID=%s

# Auto Delete
AUTO_DELETE=false
AUTO_DELETE_TIME=120

# Encryption
ENCRYPTION_KEY=%s
`,
		p.PostgresDB,
		p.PostgresUser,
		p.PostgresPassword,
		p.BotToken,
		p.AdminIDs,
		p.EncryptionKey,
	)
}

// GenerateDockerfile creates the Dockerfile contents for running ADBOT container.
func GenerateDockerfile() string {
	return `FROM debian:trixie-slim

WORKDIR /app

COPY ./dist/adbot ./dist/adbot
COPY .env ./.env

RUN chmod +x ./dist/adbot

ENTRYPOINT ["./dist/adbot"]
`
}

// GenerateDockerCompose creates the docker-compose.yml file contents.
func GenerateDockerCompose() string {
	return `services:

  postgres:
    image: ghcr.io/microsoft/pg_durable:pg17
    container_name: postgres
    restart: unless-stopped

    command:
      [
        "postgres",
        "-c",
        "shared_preload_libraries=pg_durable",
        "-c",
        "pg_durable.database=${POSTGRES_DB}",
        "-c",
        "pg_durable.worker_role=${POSTGRES_USER}"
      ]

    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      PGUSER: ${POSTGRES_USER}
      PGDATABASE: ${POSTGRES_DB}
      TZ: Asia/Kolkata

    ports:
      - "5432:5432"

    volumes:
      - postgres_data:/var/lib/postgresql/data

    healthcheck:
      test:
        [
          "CMD-SHELL",
          "pg_isready -U $${POSTGRES_USER} -d $${POSTGRES_DB}"
        ]
      interval: 5s
      timeout: 5s
      retries: 12
      start_period: 10s

    networks:
      - backend


  redis:
    image: redis:7-alpine
    container_name: redis
    restart: unless-stopped

    command: redis-server --appendonly yes

    ports:
      - "6379:6379"

    volumes:
      - redis_data:/data

    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 10

    networks:
      - backend


  bot:
    build:
      context: .
      dockerfile: Dockerfile

    image: adbot:latest
    container_name: aibot
    restart: unless-stopped

    environment:
      DATABASE_URL: postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}
      REDIS_URL: redis://redis:6379

    depends_on:

      postgres:
        condition: service_healthy

      redis:
        condition: service_healthy

    networks:
      - backend


volumes:

  postgres_data:
  redis_data:


networks:

  backend:
    name: backend
`
}

// WriteFile writes a file with 0644 permissions, creating parent directories if needed.
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0644)
}

// WriteFileSecure writes a file with specific permissions, creating parent directories if needed.
func WriteFileSecure(path string, data []byte, permission os.FileMode) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(path, data, permission); err != nil {
		return err
	}
	return os.Chmod(path, permission)
}
