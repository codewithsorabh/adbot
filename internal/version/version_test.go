package version

import (
	"testing"
)

func TestVersionFunctions(t *testing.T) {
	orig := Version
	defer func() { Version = orig }()

	Version = "1.2.3"
	if Get() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", Get())
	}
	if String() != "v1.2.3" {
		t.Errorf("expected v1.2.3, got %s", String())
	}
	if FullInfo() != "adbot v1.2.3" {
		t.Errorf("expected 'adbot v1.2.3', got %s", FullInfo())
	}

	Version = "dev"
	if Get() != "dev" {
		t.Errorf("expected dev, got %s", Get())
	}
	if String() != "dev" {
		t.Errorf("expected dev, got %s", String())
	}

	Version = ""
	if Get() != "dev" {
		t.Errorf("expected dev for empty version, got %s", Get())
	}
}
