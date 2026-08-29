package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var botTokenPattern = regexp.MustCompile(`^[0-9]{6,12}:[A-Za-z0-9_-]{20,}$`)

// PostgresPassword validates PostgreSQL password requirements.
func PostgresPassword(password string) error {
	if password == "" {
		return fmt.Errorf("PostgreSQL password cannot be empty")
	}

	if len(password) < 8 {
		return fmt.Errorf("PostgreSQL password must be at least 8 characters")
	}

	if len(password) > 256 {
		return fmt.Errorf("PostgreSQL password is too long")
	}

	if strings.ContainsAny(password, "\r\n") {
		return fmt.Errorf("PostgreSQL password cannot contain newline characters")
	}

	return nil
}

// AdminIDs validates comma-separated Telegram numeric user IDs.
func AdminIDs(value string) error {
	if value == "" {
		return fmt.Errorf("Admin ID cannot be empty")
	}

	parts := strings.Split(value, ",")
	seen := make(map[int64]bool)

	for _, part := range parts {
		idString := strings.TrimSpace(part)

		if idString == "" {
			return fmt.Errorf("empty Admin ID found")
		}

		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			return fmt.Errorf("Admin ID '%s' is not numeric", idString)
		}

		if id <= 0 {
			return fmt.Errorf("Admin ID must be greater than 0")
		}

		if seen[id] {
			return fmt.Errorf("duplicate Admin ID: %d", id)
		}

		seen[id] = true
	}

	return nil
}

// BotToken validates Telegram bot token format.
func BotToken(token string) bool {
	if token == "" || len(token) > 256 {
		return false
	}
	return botTokenPattern.MatchString(token)
}
