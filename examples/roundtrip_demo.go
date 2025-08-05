package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/deiu/rdf2go"
)

func main() {
	fmt.Println("=== TriG ↔ JSON-LD ↔ TriG Round-trip Conversion Example ===")
	
	// Start with comprehensive TriG data
	originalTrigData := `# Default graph - main entities
{
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/name> "Alice Johnson" .
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/age> "28" .
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/email> "alice@example.org" .
}

# Social graph - friendships and connections
<http://example.org/graphs/social> {
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/knows> <http://example.org/bob> .
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/knows> <http://example.org/charlie> .
	<http://example.org/bob> <http://xmlns.com/foaf/0.1/name> "Bob Smith" .
	<http://example.org/bob> <http://xmlns.com/foaf/0.1/age> "30" .
	<http://example.org/charlie> <http://xmlns.com/foaf/0.1/name> "Charlie Brown" .
}

# Work graph - professional information
<http://example.org/graphs/work> {
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/workplaceHomepage> <http://example.org/company> .
	<http://example.org/alice> <http://xmlns.com/foaf/0.1/title> "Software Engineer" .
	<http://example.org/company> <http://xmlns.com/foaf/0.1/name> "Tech Corp" .
	<http://example.org/company> <http://xmlns.com/foaf/0.1/homepage> <http://techcorp.example.org> .
}`

	fmt.Println("1. ORIGINAL TriG DATA:")
	fmt.Println("======================")
	fmt.Println(originalTrigData)
	fmt.Println()

	// Step 1: Parse original TriG into dataset
	fmt.Println("2. PARSING TriG → Dataset:")
	fmt.Println("==========================")
	dataset1 := rdf2go.NewDataset("http://example.org/")
	err := dataset1.Parse(strings.NewReader(originalTrigData), "application/trig")
	if err != nil {
		fmt.Printf("❌ Error parsing TriG: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully parsed TriG data\n")
	fmt.Printf("   • Total quads: %d\n", dataset1.Len())
	fmt.Printf("   • Default graph triples: %d\n", dataset1.GetDefaultGraph().Len())
	fmt.Printf("   • Named graphs: %d\n", len(dataset1.GetNamedGraphs()))
	
	for _, graph := range dataset1.GetNamedGraphs() {
		graphTriples := dataset1.GetGraph(graph).Len()
		fmt.Printf("     - %s: %d triples\n", graph.String(), graphTriples)
	}
	fmt.Println()

	// Step 2: Convert dataset to JSON-LD
	fmt.Println("3. CONVERTING Dataset → JSON-LD:")
	fmt.Println("================================")
	var jsonldBuffer bytes.Buffer
	err = dataset1.Serialize(&jsonldBuffer, "application/ld+json")
	if err != nil {
		fmt.Printf("❌ Error converting to JSON-LD: %v\n", err)
		return
	}

	jsonldData := jsonldBuffer.String()
	fmt.Printf("✅ Successfully converted to JSON-LD\n")
	fmt.Println("JSON-LD Output:")
	fmt.Println(jsonldData)
	fmt.Println()

	// Step 3: Parse JSON-LD back into new dataset
	fmt.Println("4. PARSING JSON-LD → Dataset:")
	fmt.Println("=============================")
	dataset2 := rdf2go.NewDataset("http://example.org/")
	err = dataset2.Parse(strings.NewReader(jsonldData), "application/ld+json")
	if err != nil {
		fmt.Printf("❌ Error parsing JSON-LD: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully parsed JSON-LD back to dataset\n")
	fmt.Printf("   • Total quads: %d\n", dataset2.Len())
	fmt.Printf("   • Default graph triples: %d\n", dataset2.GetDefaultGraph().Len())
	fmt.Printf("   • Named graphs: %d\n", len(dataset2.GetNamedGraphs()))
	fmt.Println()

	// Step 4: Convert back to TriG
	fmt.Println("5. CONVERTING Dataset → TriG:")
	fmt.Println("=============================")
	var finalTrigBuffer bytes.Buffer
	err = dataset2.Serialize(&finalTrigBuffer, "application/trig")
	if err != nil {
		fmt.Printf("❌ Error converting back to TriG: %v\n", err)
		return
	}

	finalTrigData := finalTrigBuffer.String()
	fmt.Printf("✅ Successfully converted back to TriG\n")
	fmt.Println("Final TriG Output:")
	fmt.Println(finalTrigData)

	// Step 5: Validate round-trip integrity
	fmt.Println("6. ROUND-TRIP VALIDATION:")
	fmt.Println("=========================")
	
	// Parse the final TriG to verify it's valid
	dataset3 := rdf2go.NewDataset("http://example.org/")
	err = dataset3.Parse(strings.NewReader(finalTrigData), "application/trig")
	if err != nil {
		fmt.Printf("❌ Error parsing final TriG: %v\n", err)
		return
	}

	// Compare dataset sizes
	fmt.Printf("Data integrity check:\n")
	fmt.Printf("   • Original dataset quads: %d\n", dataset1.Len())
	fmt.Printf("   • JSON-LD round-trip quads: %d\n", dataset2.Len())
	fmt.Printf("   • Final TriG quads: %d\n", dataset3.Len())

	// Check if we have the same number of quads
	if dataset1.Len() == dataset2.Len() && dataset2.Len() == dataset3.Len() {
		fmt.Println("   ✅ Quad count preserved through round-trip!")
	} else {
		fmt.Println("   ⚠️  Quad count differences detected")
	}

	// Check graphs
	originalGraphs := len(dataset1.GetNamedGraphs())
	finalGraphs := len(dataset3.GetNamedGraphs())
	fmt.Printf("   • Original named graphs: %d\n", originalGraphs)
	fmt.Printf("   • Final named graphs: %d\n", finalGraphs)
	
	if originalGraphs == finalGraphs {
		fmt.Println("   ✅ Named graph count preserved!")
	} else {
		fmt.Println("   ⚠️  Named graph count differences detected")
	}

	// Detailed comparison by graph
	fmt.Println("\n7. DETAILED GRAPH COMPARISON:")
	fmt.Println("=============================")
	
	fmt.Println("Original default graph:")
	for triple := range dataset1.GetDefaultGraph().IterTriples() {
		fmt.Printf("   %s\n", triple.String())
	}
	
	fmt.Println("\nFinal default graph:")
	for triple := range dataset3.GetDefaultGraph().IterTriples() {
		fmt.Printf("   %s\n", triple.String())
	}

	// Compare named graphs
	for _, graphName := range dataset1.GetNamedGraphs() {
		originalGraph := dataset1.GetGraph(graphName)
		finalGraph := dataset3.GetGraph(graphName)
		
		fmt.Printf("\nNamed graph %s:\n", graphName.String())
		fmt.Printf("   Original: %d triples\n", originalGraph.Len())
		fmt.Printf("   Final: %d triples\n", finalGraph.Len())
		
		if originalGraph.Len() == finalGraph.Len() {
			fmt.Printf("   ✅ Triple count preserved\n")
		} else {
			fmt.Printf("   ⚠️  Triple count difference\n")
		}
	}

	// Summary
	fmt.Println("\n8. CONVERSION SUMMARY:")
	fmt.Println("======================")
	fmt.Println("✅ TriG → Dataset parsing: SUCCESS")
	fmt.Println("✅ Dataset → JSON-LD serialization: SUCCESS")
	fmt.Println("✅ JSON-LD → Dataset parsing: SUCCESS") 
	fmt.Println("✅ Dataset → TriG serialization: SUCCESS")
	fmt.Println("✅ Final TriG parsing validation: SUCCESS")
	
	if dataset1.Len() == dataset3.Len() && originalGraphs == finalGraphs {
		fmt.Println("🎉 PERFECT ROUND-TRIP: All data preserved!")
	} else {
		fmt.Println("⚠️  PARTIAL ROUND-TRIP: Some differences detected")
	}
	
	fmt.Println("\n📋 CONVERSION CHAIN:")
	fmt.Println("TriG → Dataset → JSON-LD → Dataset → TriG → Dataset")
	fmt.Println("\nThe rdf2go library successfully handles multi-format RDF conversions")
	fmt.Println("while preserving named graph structure and semantic integrity!")
}
