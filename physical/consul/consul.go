package consul

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/momentag/ebliss/sdk/helpers/parseutils"
	"github.com/momentag/ebliss/sdk/physical"
	"github.com/rs/zerolog"
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
	client          *api.Client
	path            string
	kv              *api.KV
	permitPool      *physical.PermitPool
	consistencyMode string
	sessionTTL      string
	lockWaitTime    time.Duration
}

func NewConsulBackend(conf map[string]string, logger zerolog.Logger) (physical.Backend, error) {

	path, ok := conf["path"]
	if !ok {
		path = "ebliss/"
	}
	logger.Debug().Str("path", path).Bool("ok", ok).Msg("consul backend config path set")

	if !strings.HasSuffix(path, "/") {
		logger.Warn().Str("path", path).Msg("path does not have forward slash, appending it")
		path += "/"
	}

	if strings.HasPrefix(path, "/") {
		logger.Warn().Str("path", path).Msg("path has prefix, trimming it")
		path = strings.TrimPrefix(path, "/")
	}

	sessionTTL := api.DefaultLockSessionTTL
	sessionTTLStr, ok := conf["session_ttl"]
	if ok {
		_, err := parseutils.ParseDurationSecond(sessionTTLStr)
		if err != nil {
			return nil, fmt.Errorf("invalid session_ttl: %v", err.Error())
		}
		sessionTTL = sessionTTLStr
		logger.Debug().Str("session_ttl", sessionTTL).Msg("config session_ttl set")
	}

	lockWaitTime := api.DefaultLockWaitTime
	lockWaitTimeRaw, ok := conf["lock_wait_time"]
	if ok {
		d, err := parseutils.ParseDurationSecond(lockWaitTimeRaw)
		if err != nil {
			return nil, fmt.Errorf("cannot parse lock_wait_time: %v", err.Error())
		}
		lockWaitTime = d
		logger.Debug().Dur("lock_wait_time", lockWaitTime).Msg("config lock_wait_time set")
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		mpi, err := strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse max_parallel: %v", err.Error())
		}
		maxParInt = mpi
		logger.Debug().Int("max_parallel", maxParInt)
	}

	c := &ConsulBackend{
		client:          nil,
		path:            "",
		kv:              nil,
		permitPool:      nil,
		consistencyMode: "",
		sessionTTL:      "",
		lockWaitTime:    0,
	}

	return c, nil
}

func (c ConsulBackend) Put(ctx context.Context, entry *physical.Entry) error {
	panic("implement me")
}

func (c ConsulBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	panic("implement me")
}

func (c ConsulBackend) Delete(ctx context.Context, key string) error {
	panic("implement me")
}

func (c ConsulBackend) List(ctx context.Context, prefix string) ([]string, error) {
	panic("implement me")
}

func (c ConsulBackend) Transaction(ctx context.Context, i *[]physical.TxnEntry) error {
	panic("implement me")
}

func (c ConsulBackend) LockWith(key, value string) (physical.Lock, error) {
	panic("implement me")
}

func (c ConsulBackend) HAEnabled() bool {
	return true
}
