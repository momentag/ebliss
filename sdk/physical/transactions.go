package physical

import "context"

type TxnEntry struct {
	Operation Operation
	Entry     *Entry
}

type Transactional interface {
	Transaction(context.Context, []*TxnEntry) error
}

type TransactionalBackend interface {
	Backend
	Transactional
}
