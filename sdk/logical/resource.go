package logical

type Entity struct {
	Name   string
	Schema []*Variable
}

type Relationship struct {
	Owning     *Entity
	Target     *Entity
	Descriptor string
}

type Resource struct {
	Entity        *Entity
	Relationships []*Relationship
}
