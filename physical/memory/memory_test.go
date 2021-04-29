package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryBackend(t *testing.T) {
	_, err := NewInMemoryBackend()
	assert.Nil(t, err)
}
