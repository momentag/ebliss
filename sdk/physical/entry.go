package physical

import "github.com/momentag/ebliss/sdk/resources"

type Entry struct {
	Key   *resources.Variable
	Value []byte

	KeyHash   []byte
	ValueHash []byte
}
