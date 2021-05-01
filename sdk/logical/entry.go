package logical

type Entry struct {
	Key   Variable
	Value Value
}

type Document struct {
	Key     string
	Entries []*Entry
}

type Collection []*Document

type Record struct {
	Key         string
	Collections []*Collection
}

type Table []*Record
