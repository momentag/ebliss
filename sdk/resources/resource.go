package resources

import (
	dag2 "github.com/heimdalr/dag"

	"github.com/momentag/ebliss/sdk/physical"
)

const (
	OneToOne   = byte(1)
	OneToMany  = byte(2)
	ManyToMany = byte(3)
)

const (
	String          = byte(1)
	EncryptedString = byte(2)
	Integer         = byte(3)
	Decimal         = byte(4)
	Timestamp       = byte(5)
	Boolean         = byte(6)
	Blob            = byte(7)
)

type Relationship struct {
	Kind     byte
	Resource *Resource
}

type Schema struct {
	Variables     []*Variable
	Relationships []*Relationship
}

type Resource struct {
	Name     string
	Schema   *Schema
	Backends []*physical.Backend
}

type ResourceList []*Resource

func CreateSchemaDAG(list ResourceList) (*dag2.DAG, map[string]string) {
	dag := dag2.NewDAG()
	var gmp = make(map[string]string)

	// Add physical
	for _, resource := range list {
		if id, err := dag.AddVertex(resource); err != nil {
			panic(err)
		} else {
			gmp[resource.Name] = id
		}
	}

	// Connect nodes
	for _, resource := range list {
		srcId := gmp[resource.Name]
		for _, relationship := range resource.Schema.Relationships {
			targetId := gmp[relationship.Resource.Name]
			if err := dag.AddEdge(srcId, targetId); err != nil {
				panic(err)
			}
		}
	}

	return dag, gmp
}
