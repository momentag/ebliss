package logical

type Entry struct {
	Key   *Variable
	Value *Value
}

type Document []*Entry

type Collection []*Document

type Record map[Variable]*Collection

type Table []*Record
