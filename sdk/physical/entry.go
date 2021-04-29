package physical

import (
	"hash"
)

type Entry struct {
	Key   *Variable
	Value []byte

	KeyHash   hash.Hash
	ValueHash hash.Hash
}
