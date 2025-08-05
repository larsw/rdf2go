package rdf2go

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewQuad(t *testing.T) {
	s := NewResource("http://example.org/subject")
	p := NewResource("http://example.org/predicate")
	o := NewLiteral("object")
	g := NewResource("http://example.org/graph")
	
	quad := NewQuad(s, p, o, g)
	
	assert.Equal(t, s, quad.Subject)
	assert.Equal(t, p, quad.Predicate)
	assert.Equal(t, o, quad.Object)
	assert.Equal(t, g, quad.Graph)
}

func TestNewTripleQuad(t *testing.T) {
	triple := NewTriple(NewResource("a"), NewResource("b"), NewResource("c"))
	quad := NewTripleQuad(triple)
	
	assert.Equal(t, triple.Subject, quad.Subject)
	assert.Equal(t, triple.Predicate, quad.Predicate)
	assert.Equal(t, triple.Object, quad.Object)
	assert.Nil(t, quad.Graph)
}

func TestQuadToTriple(t *testing.T) {
	s := NewResource("http://example.org/subject")
	p := NewResource("http://example.org/predicate")
	o := NewLiteral("object")
	g := NewResource("http://example.org/graph")
	
	quad := NewQuad(s, p, o, g)
	triple := quad.ToTriple()
	
	assert.Equal(t, s, triple.Subject)
	assert.Equal(t, p, triple.Predicate)
	assert.Equal(t, o, triple.Object)
}

func TestQuadString(t *testing.T) {
	// Test quad with graph
	quad1 := NewQuad(NewResource("a"), NewResource("b"), NewResource("c"), NewResource("g"))
	assert.Equal(t, "<a> <b> <c> <g> .", quad1.String())
	
	// Test quad without graph (default graph)
	quad2 := NewQuad(NewResource("a"), NewResource("b"), NewResource("c"), nil)
	assert.Equal(t, "<a> <b> <c> .", quad2.String())
}

func TestQuadEqual(t *testing.T) {
	s := NewResource("a")
	p := NewResource("b")
	o := NewResource("c")
	g := NewResource("g")
	
	quad1 := NewQuad(s, p, o, g)
	quad2 := NewQuad(s, p, o, g)
	quad3 := NewQuad(s, p, o, nil)
	quad4 := NewQuad(s, p, o, nil)
	
	assert.True(t, quad1.Equal(quad2))
	assert.True(t, quad3.Equal(quad4))
	assert.False(t, quad1.Equal(quad3))
	assert.False(t, quad3.Equal(quad1))
}
