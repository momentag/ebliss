package physical

type Entry struct {
	Key   *Variable
	Value []byte

	KeyHash   []byte
	ValueHash []byte
}
