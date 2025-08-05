package rdf2go

import (
	"bytes"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

var (
	testDatasetUri = "https://example.org/dataset"
	simpleTrig = `{
  <#me> <http://xmlns.com/foaf/0.1/name> "Test" .
  <#me> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .
}

<http://example.org/graph1> {
  <#alice> <http://xmlns.com/foaf/0.1/name> "Alice" .
  <#alice> <http://xmlns.com/foaf/0.1/knows> <#bob> .
}

<http://example.org/graph2> {
  <#bob> <http://xmlns.com/foaf/0.1/name> "Bob" .
}`
)

func TestNewDataset(t *testing.T) {
	d := NewDataset(testDatasetUri)
	assert.Equal(t, testDatasetUri, d.URI())
	assert.Equal(t, 0, d.Len())
	assert.Equal(t, NewResource(testDatasetUri), d.Term())
}

func TestDatasetAdd(t *testing.T) {
	d := NewDataset(testDatasetUri)
	quad := NewQuad(NewResource("a"), NewResource("b"), NewResource("c"), NewResource("g"))
	d.Add(quad)
	assert.Equal(t, 1, d.Len())
	d.Remove(quad)
	assert.Equal(t, 0, d.Len())
}

func TestDatasetAddQuad(t *testing.T) {
	d := NewDataset(testDatasetUri)
	d.AddQuad(NewResource("a"), NewResource("b"), NewResource("c"), NewResource("g"))
	assert.Equal(t, 1, d.Len())
}

func TestDatasetAddTriple(t *testing.T) {
	d := NewDataset(testDatasetUri)
	d.AddTriple(NewResource("a"), NewResource("b"), NewResource("c"))
	assert.Equal(t, 1, d.Len())
	
	// Verify it was added to default graph
	quad := d.One(NewResource("a"), NewResource("b"), NewResource("c"), nil)
	assert.NotNil(t, quad)
	assert.Nil(t, quad.Graph)
}

func TestDatasetGetGraph(t *testing.T) {
	d := NewDataset(testDatasetUri)
	graphName := NewResource("http://example.org/graph1")
	
	// Add some quads to different graphs
	d.AddQuad(NewResource("a"), NewResource("b"), NewResource("c"), graphName)
	d.AddQuad(NewResource("d"), NewResource("e"), NewResource("f"), graphName)
	d.AddTriple(NewResource("x"), NewResource("y"), NewResource("z")) // default graph
	
	// Get the named graph
	namedGraph := d.GetGraph(graphName)
	assert.Equal(t, 2, namedGraph.Len())
	
	// Get the default graph
	defaultGraph := d.GetDefaultGraph()
	assert.Equal(t, 1, defaultGraph.Len())
}

func TestDatasetGetNamedGraphs(t *testing.T) {
	d := NewDataset(testDatasetUri)
	graph1 := NewResource("http://example.org/graph1")
	graph2 := NewResource("http://example.org/graph2")
	
	d.AddQuad(NewResource("a"), NewResource("b"), NewResource("c"), graph1)
	d.AddQuad(NewResource("d"), NewResource("e"), NewResource("f"), graph2)
	d.AddTriple(NewResource("x"), NewResource("y"), NewResource("z")) // default graph
	
	namedGraphs := d.GetNamedGraphs()
	assert.Equal(t, 2, len(namedGraphs))
	
	// Check that both graphs are present (order may vary)
	foundGraph1 := false
	foundGraph2 := false
	for _, g := range namedGraphs {
		if g.Equal(graph1) {
			foundGraph1 = true
		}
		if g.Equal(graph2) {
			foundGraph2 = true
		}
	}
	assert.True(t, foundGraph1)
	assert.True(t, foundGraph2)
}

func TestDatasetOne(t *testing.T) {
	d := NewDataset(testDatasetUri)
	graph1 := NewResource("http://example.org/graph1")
	
	d.AddQuad(NewResource("a"), NewResource("b"), NewResource("c"), graph1)
	d.AddTriple(NewResource("a"), NewResource("b"), NewResource("d")) // default graph
	
	// Find quad in named graph
	quad1 := d.One(NewResource("a"), NewResource("b"), NewResource("c"), graph1)
	assert.NotNil(t, quad1)
	assert.Equal(t, graph1, quad1.Graph)
	
	// Find quad in default graph
	quad2 := d.One(NewResource("a"), NewResource("b"), NewResource("d"), nil)
	assert.NotNil(t, quad2)
	assert.Nil(t, quad2.Graph)
	
	// Should not find quad in wrong graph
	quad3 := d.One(NewResource("a"), NewResource("b"), NewResource("c"), nil)
	assert.Nil(t, quad3)
}

func TestDatasetAll(t *testing.T) {
	d := NewDataset(testDatasetUri)
	graph1 := NewResource("http://example.org/graph1")
	
	d.AddQuad(NewResource("a"), NewResource("b"), NewResource("c"), graph1)
	d.AddQuad(NewResource("a"), NewResource("b"), NewResource("d"), graph1)
	d.AddTriple(NewResource("a"), NewResource("b"), NewResource("e")) // default graph
	
	// Find all quads with specific subject and predicate in named graph
	quads1 := d.All(NewResource("a"), NewResource("b"), nil, graph1)
	assert.Equal(t, 2, len(quads1))
	
	// Find all quads with specific subject and predicate in default graph
	quads2 := d.All(NewResource("a"), NewResource("b"), nil, nil)
	assert.Equal(t, 1, len(quads2))
}

func TestDatasetString(t *testing.T) {
	d := NewDataset(testDatasetUri)
	quad := NewQuad(NewResource("a"), NewResource("b"), NewResource("c"), NewResource("g"))
	d.Add(quad)
	expected := "<a> <b> <c> <g> .\n"
	assert.Equal(t, expected, d.String())
}

func TestDatasetParseTrig(t *testing.T) {
	d := NewDataset(testDatasetUri)
	err := d.Parse(strings.NewReader(simpleTrig), "application/trig")
	assert.NoError(t, err)
	assert.True(t, d.Len() > 0)
	
	// Check default graph
	defaultGraph := d.GetDefaultGraph()
	assert.True(t, defaultGraph.Len() > 0)
	
	// Check named graphs
	namedGraphs := d.GetNamedGraphs()
	assert.True(t, len(namedGraphs) > 0)
}

func TestDatasetSerializeTrig(t *testing.T) {
	d := NewDataset(testDatasetUri)
	graph1 := NewResource("http://example.org/graph1")
	
	// Add some data
	d.AddTriple(NewResource("http://example.org/alice"), NewResource("http://xmlns.com/foaf/0.1/name"), NewLiteral("Alice"))
	d.AddQuad(NewResource("http://example.org/bob"), NewResource("http://xmlns.com/foaf/0.1/name"), NewLiteral("Bob"), graph1)
	
	var buf bytes.Buffer
	err := d.Serialize(&buf, "application/trig")
	assert.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "Alice")
	assert.Contains(t, output, "Bob")
	assert.Contains(t, output, "{")
	assert.Contains(t, output, "}")
}

func TestDatasetSerializeNQuads(t *testing.T) {
	d := NewDataset(testDatasetUri)
	quad := NewQuad(NewResource("a"), NewResource("b"), NewResource("c"), NewResource("g"))
	d.Add(quad)
	
	var buf bytes.Buffer
	err := d.Serialize(&buf, "application/n-quads")
	assert.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "<a> <b> <c> <g> .")
}

func TestDatasetMerge(t *testing.T) {
	d1 := NewDataset(testDatasetUri)
	d2 := NewDataset(testDatasetUri)
	
	d1.AddTriple(NewResource("a"), NewResource("b"), NewResource("c"))
	d2.AddQuad(NewResource("d"), NewResource("e"), NewResource("f"), NewResource("g"))
	
	d1.Merge(d2)
	assert.Equal(t, 2, d1.Len())
}
