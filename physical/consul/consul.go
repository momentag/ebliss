package consul

import (
	"github.com/momentag/ebliss/sdk/physical"
)

const (
	ConsistencyModeDefault = "default"
	ConsistencyModeStrong  = "strong"
)

var (
	_ physical.Backend       = (*ConsulBackend)(nil)
	_ physical.HABackend     = (*ConsulBackend)(nil)
	_ physical.Lock          = (*ConsulLock)(nil)
	_ physical.Transactional = (*ConsulBackend)(nil)
)

type ConsulBackend struct {
	client api.client
}

type ConsulLock struct {
}
