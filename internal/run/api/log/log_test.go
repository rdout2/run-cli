package log

import (
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"github.com/stretchr/testify/assert"
)

func TestFormatEntry(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:00Z")
	entry := &logging.Entry{
		Timestamp: ts,
		Payload:   "Log message",
	}

	result := formatEntry(entry)

	// Time format is 15:04:05. In UTC (implied by Z) it is 10:00:00.
	// However, Format uses local time? No, it uses the time object's location.
	// Parse returns time in UTC if Z is present.
	
	assert.Contains(t, result, "10:00:00")
	assert.Contains(t, result, "Log message")
	assert.Equal(t, "[10:00:00] Log message", result)
}
