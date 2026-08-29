package crypto

import (
	"testing"
)

func TestGenerateAndValidateEncryptionKey(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey failed: %v", err)
	}

	if err := ValidateEncryptionKey(key); err != nil {
		t.Fatalf("ValidateEncryptionKey failed for generated key: %v", err)
	}
}

func TestValidateEncryptionKey_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"empty", "", true},
		{"too short", "c2hvcnQ=", true},
		{"invalid characters", "???notbase64???", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEncryptionKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEncryptionKey(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
			}
		})
	}
}
