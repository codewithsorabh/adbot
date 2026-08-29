package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codewithsorabh/adbot/internal/validate"
	"golang.org/x/term"
)

// Ask prompts for user input with an optional default value.
func Ask(reader *bufio.Reader, label, defaultValue string) string {
	for {
		if defaultValue != "" {
			fmt.Printf("%s [%s]: ", label, defaultValue)
		} else {
			fmt.Printf("%s: ", label)
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			Fatal("Failed to read input", err)
		}

		input = strings.TrimSpace(input)
		if input == "" && defaultValue != "" {
			return defaultValue
		}
		if input != "" {
			return input
		}
	}
}

// AskRequired prompts for non-empty user input.
func AskRequired(reader *bufio.Reader, label, defaultValue string) string {
	for {
		value := Ask(reader, label, defaultValue)
		if value != "" {
			return value
		}
	}
}

// AskPassword prompts for a password securely.
func AskPassword(label string) string {
	for {
		fmt.Printf("%s: ", label)

		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err != nil {
			reader := bufio.NewReader(os.Stdin)
			password, readErr := reader.ReadString('\n')
			if readErr != nil {
				Fatal("Failed to read password", readErr)
			}
			password = strings.TrimSpace(password)
			if err := validate.PostgresPassword(password); err != nil {
				fmt.Printf("[ERROR] %v\n", err)
				continue
			}
			return password
		}

		password := strings.TrimSpace(string(passwordBytes))
		if err := validate.PostgresPassword(password); err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			continue
		}
		return password
	}
}

// AskSecret prompts for sensitive input with hidden terminal echo.
func AskSecret(label string) string {
	fmt.Printf("%s: ", label)

	value, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		fallback, readErr := reader.ReadString('\n')
		if readErr != nil {
			Fatal("Failed to read secret", readErr)
		}
		return strings.TrimSpace(fallback)
	}

	return strings.TrimSpace(string(value))
}

// AskBotToken prompts for a Telegram Bot Token and validates it.
func AskBotToken() string {
	for {
		token := AskSecret("Telegram Bot Token")
		if token == "" {
			continue
		}
		if !validate.BotToken(token) {
			fmt.Println("[ERROR] Invalid Telegram Bot Token format. Example: 123456789:ABCdefGHIjklMNOpqrsTUVwxyz12345")
			continue
		}
		return token
	}
}

// AskAdminIDs prompts for Telegram Admin IDs and validates them.
func AskAdminIDs() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		value := Ask(reader, "Admin ID(s), comma-separated", "")
		if err := validate.AdminIDs(value); err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			continue
		}
		return value
	}
}

// Fatal prints an error and terminates execution with exit code 1.
func Fatal(message string, err error) {
	fmt.Printf("[ERROR] %s\n", message)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
	os.Exit(1)
}

// FatalMessage prints an error message and terminates execution with exit code 1.
func FatalMessage(message string) {
	fmt.Printf("[ERROR] %s\n", message)
	os.Exit(1)
}
