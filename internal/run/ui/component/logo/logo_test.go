package logo_test

import (
	"strings"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/ui/component/logo"
)

func TestString(t *testing.T) {
	s := logo.String()

	if !strings.Contains(s, "____  _   _ _   _") {
		t.Errorf("logo.String() should contain the logo art")
	}
}
