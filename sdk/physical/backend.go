package physical

import (
	"context"
	"strings"

	"github.com/momentag/ebliss/sdk/resources"
)

type Backend interface {
	Put(ctx context.Context, entry *Entry) error
	Get(ctx context.Context, key *resources.Variable) (*Entry, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
}

type Lock interface {
	Lock(stopCh <-chan struct{}) (<-chan struct{}, error)
	Unlock() error
	Value() (bool, string, error)
}

type HABackend interface {
	LockWith(key, value string) (Lock, error)
	HAEnabled() bool
}

type RedirectDetect interface {
	DetectHostAddr() (string, error)
}

func Prefixes(s string) []string {
	components := strings.Split(s, "/")
	var result []string
	for i := 1; i < len(components); i++ {
		result = append(result, strings.Join(components[:i], "/"))
	}
	return result
}
