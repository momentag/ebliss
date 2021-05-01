package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/momentag/ebliss/sdk/physical"
)

func TestInMemoryBackend(t *testing.T) {
	t.Parallel()

	inmem, err := NewInMemoryBackend()
	assert.Nil(t, err)

	variable := &physical.Variable{
		Name:       "key",
		Implements: physical.String,
	}

	t.Run("un-hashed entries", func(t *testing.T) {
		t.Run("PUT", func(t *testing.T) {
			if entry, ok := variable.NewEntry("value"); ok {
				err = inmem.Put(nil, &entry)
				assert.True(t, ok)
			}
			assert.Nil(t, err)
		})
		t.Run("GET", func(t *testing.T) {
			actual, err := inmem.Get(nil, variable)
			assert.Nil(t, err)
			assert.Equal(t, []byte("value"), actual.Value)
		})
		t.Run("DELETE", func(t *testing.T) {
			err := inmem.Delete(nil, variable.Name)
			assert.Nil(t, err)
		})
	})

	t.Run("hashed entries", func(t *testing.T) {
		t.Run("PUT", func(t *testing.T) {
			if entry, ok := variable.NewHashedEntry("value"); ok {
				err = inmem.Put(nil, &entry)
				assert.Nil(t, err)
			}
		})
	})

}
