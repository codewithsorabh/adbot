package embedded

import (
	_ "embed"
)

// AdbotDocker contains the embedded binary payload.
//
//go:embed adbot-docker
var AdbotDocker []byte
