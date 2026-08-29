package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codewithsorabh/adbot/internal/config"
	"github.com/codewithsorabh/adbot/internal/crypto"
	"github.com/codewithsorabh/adbot/internal/embedded"
	"github.com/codewithsorabh/adbot/internal/ui"
	"github.com/codewithsorabh/adbot/internal/updater"
	"github.com/codewithsorabh/adbot/internal/version"
)

func main() {
	// ========================================================================
	// CLI flags
	// ========================================================================
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Println(version.FullInfo())
			os.Exit(0)

		case "--check-update":
			updater.CheckForUpdate(true)
			os.Exit(0)

		case "--help", "-h":
			printUsage()
			os.Exit(0)
		}
	}

	// ========================================================================
	// Update check (best-effort, non-blocking on failure)
	// ========================================================================
	updater.CheckForUpdate(false)

	// ========================================================================
	// Create and enter ADBOT project directory
	// ========================================================================
	currentDir, err := os.Getwd()
	if err != nil {
		ui.Fatal("Failed to get current directory", err)
	}

	projectDir := filepath.Join(currentDir, "adbot")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		ui.Fatal("Failed to create adbot directory", err)
	}

	if err := os.Chdir(projectDir); err != nil {
		ui.Fatal("Failed to enter adbot directory", err)
	}

	reader := bufio.NewReader(os.Stdin)

	// ========================================================================
	// PostgreSQL Configuration
	// ========================================================================
	postgresDB := ui.AskRequired(
		reader,
		"PostgreSQL database name",
		"mybot_db",
	)

	postgresUser := ui.AskRequired(
		reader,
		"PostgreSQL username",
		"mybot_user",
	)

	postgresPassword := ui.AskPassword(
		"PostgreSQL password",
	)

	// ========================================================================
	// Telegram Configuration
	// ========================================================================
	botToken := ui.AskBotToken()
	adminIDs := ui.AskAdminIDs()

	// ========================================================================
	// Encryption Key
	// ========================================================================
	var encryptionKey string
	for {
		option := ui.Ask(
			reader,
			"Encryption key: 1) Auto generate  2) Enter manually",
			"1",
		)

		switch option {
		case "1":
			key, err := crypto.GenerateEncryptionKey()
			if err != nil {
				ui.Fatal("Failed to generate encryption key", err)
			}
			encryptionKey = key

		case "2":
			key := ui.AskSecret("Enter encryption key")
			if err := crypto.ValidateEncryptionKey(key); err != nil {
				fmt.Printf("[ERROR] %v\n", err)
				continue
			}
			encryptionKey = key

		default:
			continue
		}
		break
	}

	// ========================================================================
	// Installing
	// ========================================================================
	fmt.Println("Installing...")

	// ========================================================================
	// Create .env
	// ========================================================================
	envParams := config.Parameters{
		PostgresDB:       postgresDB,
		PostgresUser:     postgresUser,
		PostgresPassword: postgresPassword,
		BotToken:         botToken,
		AdminIDs:         adminIDs,
		EncryptionKey:    encryptionKey,
	}

	envContent := config.GenerateEnv(envParams)
	if err := config.WriteFileSecure(
		config.EnvFile,
		[]byte(envContent),
		0600,
	); err != nil {
		ui.Fatal("Failed to create .env", err)
	}

	// ========================================================================
	// Create Dockerfile
	// ========================================================================
	dockerContent := config.GenerateDockerfile()
	if err := config.WriteFile(
		config.DockerFile,
		[]byte(dockerContent),
	); err != nil {
		ui.Fatal("Failed to create Dockerfile", err)
	}

	// ========================================================================
	// Create docker-compose.yml
	// ========================================================================
	composeContent := config.GenerateDockerCompose()
	if err := config.WriteFile(
		config.ComposeFile,
		[]byte(composeContent),
	); err != nil {
		ui.Fatal("Failed to create docker-compose.yml", err)
	}

	// ========================================================================
	// Extract embedded payload to ./dist/adbot
	// ========================================================================
	if len(embedded.AdbotDocker) == 0 {
		ui.FatalMessage("Embedded adbot-docker payload is empty.")
	}

	if err := os.MkdirAll(config.DistDir, 0755); err != nil {
		ui.Fatal("Failed to create dist directory", err)
	}

	if err := config.WriteFileSecure(
		config.OutputFile,
		embedded.AdbotDocker,
		0755,
	); err != nil {
		ui.Fatal("Failed to create ./dist/adbot", err)
	}

	if err := os.Chmod(config.OutputFile, 0755); err != nil {
		ui.Fatal("Failed to make ./dist/adbot executable", err)
	}

	// ========================================================================
	// Complete
	// ========================================================================
	fmt.Println("Installation complete.")
	fmt.Printf("Project: %s\n", projectDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  cd adbot")
	fmt.Println("  docker compose up -d --build")
}

func printUsage() {
	fmt.Println("ADBOT - Interactive deployment wizard & installer")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  adbot [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -v, --version       Show version information")
	fmt.Println("      --check-update  Check and install updates from GitHub Releases")
	fmt.Println("  -h, --help          Show help")
}
