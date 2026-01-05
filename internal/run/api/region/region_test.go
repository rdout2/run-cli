package region

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	regions := List()
	assert.NotEmpty(t, regions)
	assert.Contains(t, regions, "us-central1")
	assert.Contains(t, regions, "europe-west1")
}
