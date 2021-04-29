package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryBackend(t *testing.T) {

	inmem, err := NewInMemoryBackend()
	assert.Nil(t, err)

	variable := &physical.Variable{
		Name:       "key",
		Implements: resources.String,
	}

	entry, err := variable.NewEntry("value")
	assert.NotNil(t, err)

	err = inmem.Put(nil, entry)
	assert.Nil(t, err)

	actual, err := inmem.Get(nil, variable)
	assert.Nil(t, err)
	assert.Equal(t, entry, actual)

}
