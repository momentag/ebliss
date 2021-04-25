package logical

import (
	"context"
	"fmt"

	"github.com/momentag/ebliss/sdk/helpers/jsonutils"
)

// Storage is a key-value storage backend interface
type Storage interface {
	List(context.Context, string) ([]string, error)
	Get(context.Context, string) (*StorageEntry, error)
	Put(context.Context, *StorageEntry) error
	Delete(context.Context, string) error
}

// StorageEntry is an entry that maps a key to a value
type StorageEntry struct {
	Key   string
	Value []byte
}

func (e *StorageEntry) DecodeJSON(out interface{}) error {
	return jsonutils.DecodeJSON(e.Value, out)
}

func StorageEntryJSON(k string, v interface{}) (*StorageEntry, error) {
	encodedBytes, err := jsonutils.EncodeJSON(v)
	if err != nil {
		return nil, fmt.Errorf("failed to encode entry - %v", err.Error())
	}
	return &StorageEntry{
		Key:   k,
		Value: encodedBytes,
	}, nil
}
