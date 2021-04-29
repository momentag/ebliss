package physical

import (
	"hash"

	"github.com/momentag/ebliss/sdk/resources"
)

type Entry struct {
	Key   *resources.Variable
	Value []byte

	KeyHash   hash.Hash
	ValueHash hash.Hash
}
