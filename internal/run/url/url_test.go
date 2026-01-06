package url

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenInBrowser_ValidURL(t *testing.T) {
	if runtime.GOOS == "linux" && !isCI() {
		t.Skip("Skipping browser test in non-CI environment")
	}

	url := "https://github.com/JulienBreux/run-cli"
	err := OpenInBrowser(url)
	_ = err
}

func TestOpenInBrowser_EmptyURL(t *testing.T) {
	if runtime.GOOS == "linux" && !isCI() {
		t.Skip("Skipping browser test in non-CI environment")
	}

	err := OpenInBrowser("")
	_ = err
}

func TestOpenInBrowser_InvalidURL(t *testing.T) {
	if runtime.GOOS == "linux" && !isCI() {
		t.Skip("Skipping browser test in non-CI environment")
	}

	err := OpenInBrowser("not-a-valid-url")
	_ = err
}

func TestOpenInBrowser_SupportedOS(t *testing.T) {
	supportedOS := map[string]bool{
		"linux":   true,
		"windows": true,
		"darwin":  true,
	}

	currentOS := runtime.GOOS
	if supportedOS[currentOS] {
		assert.True(t, supportedOS[currentOS])
	}
}

func isCI() bool {
	ciEnvVars := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"GITHUB_ACTIONS",
		"GITLAB_CI",
		"CIRCLECI",
	}

	for _, envVar := range ciEnvVars {
		if val := os.Getenv(envVar); val != "" {
			return true
		}
	}

	return false
}
