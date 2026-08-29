package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateEncryptionKey generates a secure 32-byte Fernet-compatible URL-safe base64 encryption key.
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

// ValidateEncryptionKey validates that the key is a valid URL-safe base64 string decoding to 32 bytes.
func ValidateEncryptionKey(key string) error {
	if key == "" {
		return fmt.Errorf("encryption key cannot be empty")
	}

	decoded, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		// Also allow raw URL encoding (without padding).
		decoded, err = base64.RawURLEncoding.DecodeString(key)
		if err != nil {
			return fmt.Errorf("invalid Base64 encoding")
		}
	}

	if len(decoded) != 32 {
		return fmt.Errorf("decoded encryption key must be exactly 32 bytes, got %d bytes", len(decoded))
	}

	return nil
}
