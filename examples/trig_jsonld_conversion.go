package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/deiu/rdf2go"
)

func main() {
	fmt.Println("=== TriG ⇄ JSON-LD Conversion Example ===")
	
	// Example TriG data with multiple named graphs
	trigData := `# Default graph
{
	<http://example.org/person1> <http://xmlns.com/foaf/0.1/name> "Alice Johnson" .
	<http://example.org/person1> <http://xmlns.com/foaf/0.1/age> "28" .
}

# Friends graph
<http://example.org/graphs/friends> {
	<http://example.org/person1> <http://xmlns.com/foaf/0.1/knows> <http://example.org/person2> .
	<http://example.org/person2> <http://xmlns.com/foaf/0.1/name> "Bob Smith" .
	<http://example.org/person2> <http://xmlns.com/foaf/0.1/age> "30" .
}

# Work graph
<http://example.org/graphs/work> {
	<http://example.org/person1> <http://xmlns.com/foaf/0.1/workplaceHomepage> <http://example.org/company> .
	<http://example.org/person1> <http://xmlns.com/foaf/0.1/title> "Software Engineer" .
	<http://example.org/company> <http://xmlns.com/foaf/0.1/name> "Tech Corp" .
}`

	fmt.Println("1. Original TriG Data:")
	fmt.Println("=====================")
	fmt.Println(trigData)
	fmt.Println()

	// Parse TriG into a dataset
	dataset := rdf2go.NewDataset("http://example.org/base")
	err := dataset.Parse(strings.NewReader(trigData), "application/trig")
	if err != nil {
		fmt.Printf("Error parsing TriG: %v\n", err)
		return
	}

	fmt.Printf("✓ Parsed dataset contains %d quads\n", dataset.Len())
	fmt.Printf("✓ Default graph: %d triples\n", dataset.GetDefaultGraph().Len())
	fmt.Printf("✓ Named graphs: %d\n", len(dataset.GetNamedGraphs()))
	fmt.Println()

	// Convert TriG to JSON-LD
	fmt.Println("2. Converting TriG → JSON-LD:")
	fmt.Println("=============================")
	
	var jsonldBuffer bytes.Buffer
	err = dataset.Serialize(&jsonldBuffer, "application/ld+json")
	if err != nil {
		fmt.Printf("Error serializing to JSON-LD: %v\n", err)
		return
	}
	
	jsonldData := jsonldBuffer.String()
	fmt.Println(jsonldData)
	fmt.Println()

	// Parse the JSON-LD back into a new dataset
	fmt.Println("3. Converting JSON-LD → TriG:")
	fmt.Println("=============================")
	
	dataset2 := rdf2go.NewDataset("http://example.org/base")
	err = dataset2.Parse(strings.NewReader(jsonldData), "application/ld+json")
	if err != nil {
		fmt.Printf("Error parsing JSON-LD: %v\n", err)
		return
	}

	fmt.Printf("✓ Round-trip dataset contains %d quads\n", dataset2.Len())
	
	// Serialize back to TriG
	var trigBuffer bytes.Buffer
	err = dataset2.Serialize(&trigBuffer, "application/trig")
	if err != nil {
		fmt.Printf("Error serializing to TriG: %v\n", err)
		return
	}
	
	fmt.Println("Round-trip TriG output:")
	fmt.Println(trigBuffer.String())

	// Demonstrate working with individual formats
	fmt.Println("4. Format-Specific Operations:")
	fmt.Println("=============================")

	// Show how to work with specific named graphs
	friendsGraphName := rdf2go.NewResource("http://example.org/graphs/friends")
	friendsGraph := dataset.GetGraph(friendsGraphName)
	
	fmt.Printf("Friends graph contains %d triples:\n", friendsGraph.Len())
	for triple := range friendsGraph.IterTriples() {
		fmt.Printf("  %s\n", triple.String())
	}
	fmt.Println()

	// Convert individual graph to different formats
	fmt.Println("Friends graph as Turtle:")
	var turtleBuffer bytes.Buffer
	err = friendsGraph.Serialize(&turtleBuffer, "text/turtle")
	if err != nil {
		fmt.Printf("Error serializing to Turtle: %v\n", err)
		return
	}
	fmt.Println(turtleBuffer.String())

	fmt.Println("Friends graph as TriG (default graph):")
	var trigGraphBuffer bytes.Buffer
	err = friendsGraph.Serialize(&trigGraphBuffer, "application/trig")
	if err != nil {
		fmt.Printf("Error serializing to TriG: %v\n", err)
		return
	}
	fmt.Println(trigGraphBuffer.String())

	// Demonstrate creating a dataset from multiple graphs
	fmt.Println("5. Building Dataset from Graphs:")
	fmt.Println("===============================")
	
	// Create a new dataset and add graphs programmatically
	newDataset := rdf2go.NewDataset("http://example.org/new")
	
	// Add some triples to default graph
	newDataset.AddTriple(
		rdf2go.NewResource("http://example.org/org"),
		rdf2go.NewResource("http://xmlns.com/foaf/0.1/name"),
		rdf2go.NewLiteral("Example Organization"),
	)
	
	// Add some quads to named graphs
	metadataGraph := rdf2go.NewResource("http://example.org/metadata")
	newDataset.AddQuad(
		rdf2go.NewResource("http://example.org/dataset1"),
		rdf2go.NewResource("http://purl.org/dc/terms/created"),
		rdf2go.NewLiteral("2025-08-04"),
		metadataGraph,
	)
	newDataset.AddQuad(
		rdf2go.NewResource("http://example.org/dataset1"),
		rdf2go.NewResource("http://purl.org/dc/terms/creator"),
		rdf2go.NewLiteral("Data Team"),
		metadataGraph,
	)
	
	fmt.Printf("New dataset contains %d quads\n", newDataset.Len())
	
	// Show both formats
	fmt.Println("\nAs TriG:")
	newDataset.Serialize(&bytes.Buffer{}, "application/trig")
	var newTrigBuffer bytes.Buffer
	newDataset.Serialize(&newTrigBuffer, "application/trig")
	fmt.Println(newTrigBuffer.String())
	
	fmt.Println("As JSON-LD:")
	var newJsonldBuffer bytes.Buffer
	newDataset.Serialize(&newJsonldBuffer, "application/ld+json")
	fmt.Println(newJsonldBuffer.String())

	fmt.Println("\n✓ Conversion complete! The library seamlessly handles:")
	fmt.Println("  • TriG ↔ JSON-LD conversion")
	fmt.Println("  • Named graph preservation")
	fmt.Println("  • Individual graph extraction")
	fmt.Println("  • Multiple serialization formats")
}
