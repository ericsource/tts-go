package azure_tts

import (
	"strconv"
	"strings"
)

// Version is the Edge TTS version string.
const Version = "1.0.0"

// ParseVersion parses the version string and returns a tuple of integers.
func ParseVersion() (int, int, int) {
	parts := strings.Split(Version, ".")
	if len(parts) != 3 {
		return 0, 0, 0
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0
	}

	return major, minor, patch
}
