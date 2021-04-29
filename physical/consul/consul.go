package consul

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
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

	consistencyMode, ok := conf["consistency_mode"]
	if ok {
		switch consistencyMode {
		case ConsistencyModeStrong, ConsistencyModeDefault:
		default:
			return nil, fmt.Errorf("invalid consistency mode %v", consistencyMode)
		}
	} else {
		consistencyMode = ConsistencyModeDefault
	}
	logger.Debug().Str("consistency_mode", consistencyMode).Msg("consistency_mode set")

	consulConfig := api.DefaultConfig()
	consulConfig.Transport.MaxIdleConnsPerHost = 64

	if addr, ok := conf["address"]; ok {
		consulConfig.Address = addr
		logger.Debug().Str("consul_addr", addr).Msg("consul address set")

		parts := strings.SplitN(addr, "://", 2)
		if len(parts) == 2 {
			if parts[0] == "http" || parts[0] == "https" {
				consulConfig.Scheme = parts[0]
				consulConfig.Address = parts[1]
				logger.Debug().Str("consul_scheme", consulConfig.Scheme).Str("consul_addr", consulConfig.Address).Msg("consul scheme and address set")
			}
		}
	}

	if scheme, ok := conf["scheme"]; ok {
		consulConfig.Scheme = scheme
		logger.Debug().Str("consul_scheme", consulConfig.Scheme).Msg("consul scheme set")
	}

	if token, ok := conf["token"]; ok {
		consulConfig.Token = token
		logger.Debug().Str("consul_token", token).Msg("consul token set")
	}

	if consulConfig.Scheme == "https" {
		if tlsClientConfig, err := tlsutils.SetupTLSConfig(conf, consulConfig.Address); err != nil {
			return nil, err
		} else {
			consulConfig.Transport.TLSClientConfig = tlsClientConfig
			logger.Debug().Msg("configured TLS")
		}
	}

	consulConfig.HttpClient = &http.Client{Transport: consulConfig.Transport}
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create consul client - %v", err.Error())
	}

	c := &ConsulBackend{
		client:          client,
		path:            path,
		kv:              client.KV(),
		permitPool:      physical.NewPermitPool(maxParInt),
		consistencyMode: consistencyMode,
		sessionTTL:      sessionTTL,
		lockWaitTime:    lockWaitTime,
	}

	return c, nil
}

func (c ConsulBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"consul", "put"}, time.Now())
	c.permitPool.Acquire()
	defer c.permitPool.Release()

	pair := &api.KVPair{
		Key:   c.path + entry.Key.Name,
		Value: entry.Value,
	}

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)

	_, err := c.kv.Put(pair, writeOpts)
	if err != nil {
		return fmt.Errorf("could not PUT to consul - %v", err.Error())
	}

	return nil
}

func (c ConsulBackend) Get(ctx context.Context, key *physical.Variable) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"consul", "get"}, time.Now())
	c.permitPool.Acquire()
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	if c.consistencyMode == ConsistencyModeStrong {
		queryOpts.RequireConsistent = true
	}

	pair, _, err := c.kv.Get(c.path+key.Name, queryOpts)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, nil
	}

	ent := &physical.Entry{
		Key:   key,
		Value: pair.Value,
	}

	return ent, nil
}

func (c ConsulBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"consul", "delete"}, time.Now())
	c.permitPool.Acquire()
	defer c.permitPool.Release()

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)

	_, err := c.kv.Delete(c.path+key, writeOpts)
	return err
}

func (c ConsulBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"consul", "delete"}, time.Now())

	scan := c.path + prefix
	if strings.HasSuffix(scan, "//") {
		scan = scan[:len(scan)-1]
	}

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	out, _, err := c.kv.Keys(scan, "/", queryOpts)
	for idx, val := range out {
		out[idx] = strings.TrimPrefix(val, scan)
	}

	return out, err
}

func (c ConsulBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if len(txns) == 0 {
		return nil
	}
	defer metrics.MeasureSince([]string{"consul", "transaction"}, time.Now())

	ops := make(api.KVTxnOps, 0, len(txns))

	for _, op := range txns {
		consulOperation := &api.KVTxnOp{
			Key: c.path + op.Entry.Key.Name,
		}
		switch op.Operation {
		case physical.DeleteOp:
			consulOperation.Verb = api.KVDelete
		case physical.PutOp:
			consulOperation.Verb = api.KVSet
			consulOperation.Value = op.Entry.Value
		default:
			return fmt.Errorf("%q is not supported transaction operation", op.Operation)
		}
		ops = append(ops, consulOperation)
	}

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	ok, resp, _, err := c.kv.Txn(ops, queryOpts)
	if err != nil {
		return fmt.Errorf("could not perform operation %v", err.Error())
	}
	if ok && len(resp.Errors) == 0 {
		return nil
	}

	var errors []error
	for _, err := range resp.Errors {
		errors = append(errors, fmt.Errorf("error for op %v - happened %v", err.OpIndex, err.What))
	}

	return fmt.Errorf("errors: %v", errors)
}

func (c ConsulBackend) LockWith(key, value string) (physical.Lock, error) {
	opts := &api.LockOptions{
		Key:            c.path + key,
		Value:          []byte(value),
		SessionName:    "Ebliss Lock",
		SessionTTL:     c.sessionTTL,
		MonitorRetries: 5,
		LockWaitTime:   c.lockWaitTime,
	}
	lock, err := c.client.LockOpts(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul lock %v", err.Error())
	}
	cl := &ConsulLock{
		client:          c.client,
		key:             c.path + key,
		lock:            lock,
		consistencyMode: c.consistencyMode,
	}
	return cl, nil
}

func (c ConsulBackend) HAEnabled() bool {
	return true
}
