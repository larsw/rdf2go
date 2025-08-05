package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/deiu/rdf2go"
)

func main() {
	// Example TriG data with multiple named graphs
	trigData := `@prefix foaf: <http://xmlns.com/foaf/0.1/> .

# Default graph
{
	<#me> foaf:name "John Doe" ;
		foaf:age 30 .
}

# Named graph 1
<http://example.org/friends> {
	<#alice> foaf:name "Alice Smith" ;
		foaf:knows <#bob> .
	<#bob> foaf:name "Bob Jones" .
}

# Named graph 2  
<http://example.org/work> {
	<#me> foaf:workplaceHomepage <http://example.org/company> ;
		foaf:title "Software Engineer" .
}`

	fmt.Println("=== TriG Support Example ===")
	fmt.Println("Original TriG data:")
	fmt.Println(trigData)
	fmt.Println()

	// Parse TriG into a dataset
	dataset := rdf2go.NewDataset("https://example.org/")
	err := dataset.Parse(strings.NewReader(trigData), "application/trig")
	if err != nil {
		fmt.Printf("Error parsing TriG: %v\n", err)
		return
	}

	fmt.Printf("Dataset contains %d quads\n", dataset.Len())

	// Get the default graph
	defaultGraph := dataset.GetDefaultGraph()
	fmt.Printf("Default graph contains %d triples\n", defaultGraph.Len())

	// Get named graphs
	namedGraphs := dataset.GetNamedGraphs()
	fmt.Printf("Found %d named graphs:\n", len(namedGraphs))
	for _, graph := range namedGraphs {
		fmt.Printf("  - %s\n", graph.String())
	}

	fmt.Println("\n=== Serializing back to TriG ===")
	err = dataset.Serialize(os.Stdout, "application/trig")
	if err != nil {
		fmt.Printf("Error serializing TriG: %v\n", err)
		return
	}

	fmt.Println("\n=== Working with individual graphs ===")
	
	// Get a specific named graph
	friendsGraph := dataset.GetGraph(rdf2go.NewResource("http://example.org/friends"))
	fmt.Printf("Friends graph contains %d triples:\n", friendsGraph.Len())
	
	// Show all triples in the friends graph
	for triple := range friendsGraph.IterTriples() {
		fmt.Printf("  %s\n", triple.String())
	}

	fmt.Println("\n=== Converting Graph to TriG ===")
	
	// Create a simple graph and serialize it as TriG
	simpleGraph := rdf2go.NewGraph("https://example.org/")
	simpleGraph.AddTriple(
		rdf2go.NewResource("https://example.org/person1"),
		rdf2go.NewResource("http://xmlns.com/foaf/0.1/name"),
		rdf2go.NewLiteral("Test Person"),
	)
	
	fmt.Println("Simple graph as TriG:")
	err = simpleGraph.Serialize(os.Stdout, "application/trig")
	if err != nil {
		fmt.Printf("Error serializing graph as TriG: %v\n", err)
		return
	}
}
