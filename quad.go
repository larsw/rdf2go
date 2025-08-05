package rdf2go

import (
	"fmt"
)

// Quad contains a subject, predicate, object and graph term (named graph).
type Quad struct {
	Subject   Term
	Predicate Term
	Object    Term
	Graph     Term // Named graph - can be nil for default graph
}

// NewQuad returns a new quad with the given subject, predicate, object and graph.
func NewQuad(subject Term, predicate Term, object Term, graph Term) (quad *Quad) {
	return &Quad{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Graph:     graph,
	}
}

// NewTripleQuad returns a new quad from a triple (with nil graph - default graph).
func NewTripleQuad(triple *Triple) (quad *Quad) {
	return &Quad{
		Subject:   triple.Subject,
		Predicate: triple.Predicate,
		Object:    triple.Object,
		Graph:     nil,
	}
}

// ToTriple returns a Triple from this Quad (ignoring the graph).
func (quad Quad) ToTriple() *Triple {
	return NewTriple(quad.Subject, quad.Predicate, quad.Object)
}

// String returns the NQuads representation of this quad.
func (quad Quad) String() (str string) {
	subjStr := "nil"
	if quad.Subject != nil {
		subjStr = quad.Subject.String()
	}

	predStr := "nil"
	if quad.Predicate != nil {
		predStr = quad.Predicate.String()
	}

	objStr := "nil"
	if quad.Object != nil {
		objStr = quad.Object.String()
	}

	if quad.Graph != nil {
		return fmt.Sprintf("%s %s %s %s .", subjStr, predStr, objStr, quad.Graph.String())
	}
	return fmt.Sprintf("%s %s %s .", subjStr, predStr, objStr)
}

// Equal returns this quad is equivalent to the argument.
func (quad Quad) Equal(other *Quad) bool {
	sameTriple := quad.Subject.Equal(other.Subject) &&
		quad.Predicate.Equal(other.Predicate) &&
		quad.Object.Equal(other.Object)
	
	// Handle nil graphs
	if quad.Graph == nil && other.Graph == nil {
		return sameTriple
	}
	if quad.Graph == nil || other.Graph == nil {
		return false
	}
	
	return sameTriple && quad.Graph.Equal(other.Graph)
}
