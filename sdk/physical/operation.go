package physical

type Operation string

const (
	DeleteOp Operation = "delete"
	GetOp              = "get"
	ListOp             = "list"
	PutOp              = "put"
)
