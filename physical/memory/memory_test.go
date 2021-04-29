package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	physical "github.com/momentag/ebliss/sdk/physical"
)

func TestInMemoryBackend(t *testing.T) {

	inmem, err := NewInMemoryBackend()
	assert.Nil(t, err)

	variable := &physical.Variable{
		Name:       "key",
		Implements: physical.String,
	}

	entry := variable.NewEntry("value")

	err = inmem.Put(nil, entry)
	assert.Nil(t, err)

	actual, err := inmem.Get(nil, variable)
	assert.Nil(t, err)
	assert.Equal(t, entry, actual)

}
