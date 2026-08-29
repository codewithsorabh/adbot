package validate

import (
	"testing"
)

func TestPostgresPassword(t *testing.T) {
	tests := []struct {
		name    string
		pwd     string
		wantErr bool
	}{
		{"valid", "mysecurepassword123", false},
		{"empty", "", true},
		{"too short", "short", true},
		{"contains newline", "pass\nword123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PostgresPassword(tt.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresPassword(%q) error = %v, wantErr %v", tt.pwd, err, tt.wantErr)
			}
		})
	}
}

func TestAdminIDs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"single valid", "123456789", false},
		{"multiple valid", "123456789, 987654321", false},
		{"empty", "", true},
		{"non-numeric", "123, abc", true},
		{"duplicate", "123456, 123456", true},
		{"negative", "-12345", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AdminIDs(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminIDs(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBotToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{"valid token", "123456789:ABCdefGHIjklMNOpqrsTUVwxyz12345", true},
		{"invalid empty", "", false},
		{"invalid format", "invalid_token", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BotToken(tt.token)
			if got != tt.want {
				t.Errorf("BotToken(%q) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}
