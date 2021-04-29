package memory

import (
	"context"
	"time"

	"github.com/armon/go-metrics"
	"github.com/dgraph-io/badger/v3"

	"github.com/momentag/ebliss/sdk/physical"
)

var (
	_ physical.Backend       = (*InMemoryBackend)(nil)
	_ physical.HABackend     = (*InMemoryBackend)(nil)
	_ physical.Transactional = (*InMemoryBackend)(nil)
)

type InMemoryBackend struct {
	db         *badger.DB
	permitPool *physical.PermitPool
}

func (i InMemoryBackend) Transaction(ctx context.Context, entries []*physical.TxnEntry) error {
	panic("implement me")
}

func (i InMemoryBackend) LockWith(key, value string) (physical.Lock, error) {
	panic("implement me")
}

func (i InMemoryBackend) HAEnabled() bool {
	return false
}

func (i InMemoryBackend) Put(_ context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"inmem", "put"}, time.Now())
	i.permitPool.Acquire()
	defer i.permitPool.Release()
	return i.db.Update(func(txn *badger.Txn) error {
		entry := &badger.Entry{
			Key:       []byte(entry.Key.Name),
			Value:     entry.Value,
			ExpiresAt: 0,
			UserMeta:  0,
		}
		return txn.SetEntry(entry)
	})
}

func (i InMemoryBackend) Get(_ context.Context, key *physical.Variable) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"inmem", "get"}, time.Now())
	i.permitPool.Acquire()
	defer i.permitPool.Release()
	var val []byte
	err := i.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key.Name))
		if err != nil {
			return err
		}
		return item.Value(func(local []byte) error {
			val = local
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return &physical.Entry{
		Key:       key,
		Value:     val,
		KeyHash:   nil,
		ValueHash: nil,
	}, nil
}

func (i InMemoryBackend) Delete(_ context.Context, key string) error {
	defer metrics.MeasureSince([]string{"inmem", "delete"}, time.Now())
	i.permitPool.Acquire()
	defer i.permitPool.Release()
	return i.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (i InMemoryBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"inmem", "list"}, time.Now())
	i.permitPool.Acquire()
	defer i.permitPool.Release()
	var keys []string
	err := i.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func NewInMemoryBackend() (physical.Backend, error) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	pool := physical.NewPermitPool(1)
	c := &InMemoryBackend{
		db:         db,
		permitPool: pool,
	}
	return c, nil
}
