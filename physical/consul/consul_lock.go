package consul

import "github.com/hashicorp/consul/api"

type ConsulLock struct {
	client          *api.Client
	key             string
	lock            *api.Lock
	consistencyMode string
}

func (c *ConsulLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	return c.lock.Lock(stopCh)
}

func (c *ConsulLock) Unlock() error {
	return c.lock.Unlock()
}

func (c *ConsulLock) Value() (bool, string, error) {
	kv := c.client.KV()

	var queryOptions *api.QueryOptions
	if c.consistencyMode == ConsistencyModeStrong {
		queryOptions = &api.QueryOptions{RequireConsistent: true}
	}

	pair, _, err := kv.Get(c.key, queryOptions)
	if err != nil {
		return false, "", err
	}
	if pair == nil {
		return false, "", nil
	}

	held := pair.Session != ""
	value := string(pair.Value)
	return held, value, nil
}
