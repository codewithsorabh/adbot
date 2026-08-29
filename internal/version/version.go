package version

import "fmt"

// Version is the current version of ADBOT.
// It can be injected at build time with:
//
//	-ldflags="-X 'github.com/codewithsorabh/adbot/internal/version.Version=1.0.0'"
//
// Falls back to "dev" if built without that flag.
var Version = "1.0.0"

// Get returns the raw version string.
func Get() string {
	if Version == "" {
		return "dev"
	}
	return Version
}

// String returns formatted version string e.g. "v1.0.0".
func String() string {
	v := Get()
	if v == "dev" {
		return "dev"
	}
	return fmt.Sprintf("v%s", v)
}

// FullInfo returns a user-friendly version description.
func FullInfo() string {
	return fmt.Sprintf("adbot %s", String())
}
