package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/deiu/rdf2go"
)

func main() {
	// Create a simple dataset
	dataset := rdf2go.NewDataset("http://example.org/")
	
	// Add data to default graph
	dataset.AddTriple(
		rdf2go.NewResource("http://example.org/org"),
		rdf2go.NewResource("http://xmlns.com/foaf/0.1/name"),
		rdf2go.NewLiteral("Example Organization"),
	)
	
	// Add data to named graph
	metadataGraph := rdf2go.NewResource("http://example.org/metadata")
	dataset.AddQuad(
		rdf2go.NewResource("http://example.org/dataset1"),
		rdf2go.NewResource("http://purl.org/dc/terms/created"),
		rdf2go.NewLiteral("2025-08-04"),
		metadataGraph,
	)
	dataset.AddQuad(
		rdf2go.NewResource("http://example.org/dataset1"),
		rdf2go.NewResource("http://purl.org/dc/terms/creator"),
		rdf2go.NewLiteral("Data Team"),
		metadataGraph,
	)
	
	fmt.Println("=== Testing Improved JSON-LD Output ===")
	
	// Test TriG output
	fmt.Println("As TriG:")
	var trigBuffer bytes.Buffer
	dataset.Serialize(&trigBuffer, "application/trig")
	fmt.Println(trigBuffer.String())
	
	// Test JSON-LD output
	fmt.Println("As JSON-LD:")
	var jsonldBuffer bytes.Buffer
	dataset.Serialize(&jsonldBuffer, "application/ld+json")
	fmt.Println(jsonldBuffer.String())
	
	// Test that we can parse it back
	fmt.Println("=== Round-trip Test ===")
	dataset2 := rdf2go.NewDataset("http://example.org/")
	err := dataset2.Parse(strings.NewReader(jsonldBuffer.String()), "application/ld+json")
	if err != nil {
		fmt.Printf("Error parsing JSON-LD: %v\n", err)
	} else {
		fmt.Printf("✓ Successfully parsed JSON-LD back into dataset with %d quads\n", dataset2.Len())
	}
}
