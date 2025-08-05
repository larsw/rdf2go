package rdf2go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	rdf "github.com/deiu/gon3"
	jsonld "github.com/linkeddata/gojsonld"
)

// Dataset structure holds multiple named graphs
type Dataset struct {
	quads      map[*Quad]bool
	httpClient *http.Client
	uri        string
	term       Term
}

// NewDataset creates a Dataset object
func NewDataset(uri string, skipVerify ...bool) *Dataset {
	skip := false
	if len(skipVerify) > 0 {
		skip = skipVerify[0]
	}
	d := &Dataset{
		quads:      make(map[*Quad]bool),
		httpClient: NewHttpClient(skip),
		uri:        uri,
		term:       NewResource(uri),
	}
	return d
}

// Len returns the length of the dataset as number of quads
func (d *Dataset) Len() int {
	return len(d.quads)
}

// Term returns a Dataset Term object
func (d *Dataset) Term() Term {
	return d.term
}

// URI returns a Dataset URI object
func (d *Dataset) URI() string {
	return d.uri
}

// Add is used to add a Quad object to the dataset
func (d *Dataset) Add(q *Quad) {
	d.quads[q] = true
}

// AddQuad is used to add a quad made of individual S, P, O, G objects
func (d *Dataset) AddQuad(s Term, p Term, o Term, g Term) {
	d.quads[NewQuad(s, p, o, g)] = true
}

// AddTriple is used to add a triple to the default graph (G = nil)
func (d *Dataset) AddTriple(s Term, p Term, o Term) {
	d.quads[NewQuad(s, p, o, nil)] = true
}

// Remove is used to remove a Quad object
func (d *Dataset) Remove(q *Quad) {
	delete(d.quads, q)
}

// IterQuads provides a channel containing all the quads in the dataset.
func (d *Dataset) IterQuads() (ch chan *Quad) {
	ch = make(chan *Quad, len(d.quads))
	for quad := range d.quads {
		ch <- quad
	}
	close(ch)
	return ch
}

// GetGraph returns a Graph containing all triples for a specific named graph
func (d *Dataset) GetGraph(graphName Term) *Graph {
	g := NewGraph(d.uri)
	for quad := range d.IterQuads() {
		// Handle default graph (nil) vs named graphs
		if graphName == nil && quad.Graph == nil {
			g.Add(quad.ToTriple())
		} else if graphName != nil && quad.Graph != nil && graphName.Equal(quad.Graph) {
			g.Add(quad.ToTriple())
		}
	}
	return g
}

// GetDefaultGraph returns a Graph containing all triples in the default graph (Graph = nil)
func (d *Dataset) GetDefaultGraph() *Graph {
	return d.GetGraph(nil)
}

// GetNamedGraphs returns a list of all named graph identifiers in the dataset
func (d *Dataset) GetNamedGraphs() []Term {
	graphNames := make(map[string]Term)
	for quad := range d.IterQuads() {
		if quad.Graph != nil {
			graphNames[quad.Graph.String()] = quad.Graph
		}
	}
	
	var result []Term
	for _, graph := range graphNames {
		result = append(result, graph)
	}
	return result
}

// One returns one quad based on a quad pattern of S, P, O, G objects
func (d *Dataset) One(s Term, p Term, o Term, g Term) *Quad {
	for quad := range d.IterQuads() {
		if s != nil && !quad.Subject.Equal(s) {
			continue
		}
		if p != nil && !quad.Predicate.Equal(p) {
			continue
		}
		if o != nil && !quad.Object.Equal(o) {
			continue
		}
		if g != nil && (quad.Graph == nil || !quad.Graph.Equal(g)) {
			continue
		}
		if g == nil && quad.Graph != nil {
			continue
		}
		return quad
	}
	return nil
}

// All returns all quads that match a given pattern of S, P, O, G objects
func (d *Dataset) All(s Term, p Term, o Term, g Term) []*Quad {
	var quads []*Quad
	for quad := range d.IterQuads() {
		if s != nil && !quad.Subject.Equal(s) {
			continue
		}
		if p != nil && !quad.Predicate.Equal(p) {
			continue
		}
		if o != nil && !quad.Object.Equal(o) {
			continue
		}
		if g != nil && (quad.Graph == nil || !quad.Graph.Equal(g)) {
			continue
		}
		if g == nil && quad.Graph != nil {
			continue
		}
		quads = append(quads, quad)
	}
	return quads
}

// String returns the NQuads representation of the dataset
func (d *Dataset) String() string {
	var toString string
	for quad := range d.IterQuads() {
		toString += quad.String() + "\n"
	}
	return toString
}

// Parse is used to parse RDF data from a reader, using the provided mime type
func (d *Dataset) Parse(reader io.Reader, mime string) error {
	parserName := mimeParser[mime]
	if len(parserName) == 0 {
		parserName = "guess"
	}
	
	if parserName == "trig" {
		return d.parseTrig(reader)
	} else if parserName == "jsonld" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		jsonData, err := jsonld.ReadJSON(buf.Bytes())
		if err != nil {
			return err
		}
		options := &jsonld.Options{}
		options.Base = ""
		options.ProduceGeneralizedRdf = false
		dataSet, err := jsonld.ToRDF(jsonData, options)
		if err != nil {
			return err
		}
		for t := range dataSet.IterTriples() {
			d.AddTriple(jterm2term(t.Subject), jterm2term(t.Predicate), jterm2term(t.Object))
		}
	} else if parserName == "turtle" {
		parser, err := rdf.NewParser(d.uri).Parse(reader)
		if err != nil {
			return err
		}
		for s := range parser.IterTriples() {
			d.AddTriple(rdf2term(s.Subject), rdf2term(s.Predicate), rdf2term(s.Object))
		}
	} else {
		return errors.New(parserName + " is not supported by the parser")
	}
	return nil
}

// parseTrig parses TriG format - simplified implementation
func (d *Dataset) parseTrig(reader io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	content := buf.String()
	
	// This is a simplified TriG parser. A full implementation would require
	// a proper grammar parser, but this handles basic TriG syntax
	lines := strings.Split(content, "\n")
	var currentGraph Term = nil // Default graph
	var currentTripleLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Handle prefix declarations
		if strings.HasPrefix(line, "@prefix") {
			// TODO: Handle prefixes properly - for now skip
			continue
		}
		
		// Handle graph declarations like { or <graphname> {
		if strings.Contains(line, "{") {
			parts := strings.Split(line, "{")
			if len(parts) > 1 {
				graphPart := strings.TrimSpace(parts[0])
				if graphPart == "" {
					currentGraph = nil // Default graph
				} else {
					// Parse graph name
					currentGraph = parseGraphName(graphPart)
				}
			}
			continue
		}
		
		// Handle end of graph
		if strings.Contains(line, "}") {
			// Process any remaining triple lines
			if len(currentTripleLines) > 0 {
				d.processTripleLines(currentTripleLines, currentGraph)
				currentTripleLines = nil
			}
			currentGraph = nil // Reset to default graph
			continue
		}
		
		// Collect lines for turtle-style parsing within graph blocks
		if line != "" {
			currentTripleLines = append(currentTripleLines, line)
			// If line ends with '.', process the collected lines
			if strings.HasSuffix(line, ".") {
				d.processTripleLines(currentTripleLines, currentGraph)
				currentTripleLines = nil
			}
		}
	}
	
	// Process any remaining lines
	if len(currentTripleLines) > 0 {
		d.processTripleLines(currentTripleLines, currentGraph)
	}
	
	return nil
}

// processTripleLines processes a set of lines that form turtle-style statements
func (d *Dataset) processTripleLines(lines []string, currentGraph Term) {
	// Join all lines and parse as turtle-style content
	content := strings.Join(lines, "\n")
	
	// Use the gon3 parser to parse this as turtle content
	reader := strings.NewReader(content)
	parser, err := rdf.NewParser(d.uri).Parse(reader)
	if err != nil {
		return // Skip invalid content
	}
	
	for s := range parser.IterTriples() {
		d.AddQuad(rdf2term(s.Subject), rdf2term(s.Predicate), rdf2term(s.Object), currentGraph)
	}
}

// parseGraphName parses a graph name from TriG syntax
func parseGraphName(graphStr string) Term {
	graphStr = strings.TrimSpace(graphStr)
	if strings.HasPrefix(graphStr, "<") && strings.HasSuffix(graphStr, ">") {
		return NewResource(graphStr[1 : len(graphStr)-1])
	}
	// TODO: Handle prefixed names, blank nodes, etc.
	return NewResource(graphStr)
}

// Serialize serializes the dataset to a writer in the specified format
func (d *Dataset) Serialize(w io.Writer, mime string) error {
	serializerName := mimeSerializer[mime]
	if serializerName == "trig" {
		return d.serializeTrig(w)
	} else if serializerName == "jsonld" {
		return d.serializeJSONLD(w)
	}
	// Default to NQuads
	return d.serializeNQuads(w)
}

// serializeTrig serializes to TriG format
func (d *Dataset) serializeTrig(w io.Writer) error {
	// Group quads by graph
	graphQuads := make(map[string][]*Quad)
	var defaultGraphQuads []*Quad
	
	for quad := range d.IterQuads() {
		if quad.Graph == nil {
			defaultGraphQuads = append(defaultGraphQuads, quad)
		} else {
			graphName := quad.Graph.String()
			graphQuads[graphName] = append(graphQuads[graphName], quad)
		}
	}
	
	// Write default graph first
	if len(defaultGraphQuads) > 0 {
		fmt.Fprintln(w, "{")
		for _, quad := range defaultGraphQuads {
			fmt.Fprintf(w, "  %s %s %s .\n", 
				encodeTerm(quad.Subject), 
				encodeTerm(quad.Predicate), 
				encodeTerm(quad.Object))
		}
		fmt.Fprintln(w, "}")
	}
	
	// Write named graphs
	for graphName, quads := range graphQuads {
		fmt.Fprintf(w, "\n%s {\n", graphName)
		for _, quad := range quads {
			fmt.Fprintf(w, "  %s %s %s .\n", 
				encodeTerm(quad.Subject), 
				encodeTerm(quad.Predicate), 
				encodeTerm(quad.Object))
		}
		fmt.Fprintln(w, "}")
	}
	
	return nil
}

// serializeNQuads serializes to NQuads format (default)
func (d *Dataset) serializeNQuads(w io.Writer) error {
	for quad := range d.IterQuads() {
		fmt.Fprintln(w, quad.String())
	}
	return nil
}

// serializeJSONLD serializes to JSON-LD format with named graphs
func (d *Dataset) serializeJSONLD(w io.Writer) error {
	// Create a JSON-LD compatible structure
	result := make(map[string]interface{})
	
	// Handle default graph
	defaultGraph := d.GetDefaultGraph()
	if defaultGraph.Len() > 0 {
		var defaultTriples []map[string]interface{}
		subjectMap := make(map[string]map[string]interface{})
		
		for triple := range defaultGraph.IterTriples() {
			subjectID := termToJSONLDID(triple.Subject)
			predicateID := termToJSONLDID(triple.Predicate)
			objectValue := termToJSONLDValue(triple.Object)
			
			if _, exists := subjectMap[subjectID]; !exists {
				subjectMap[subjectID] = map[string]interface{}{
					"@id": subjectID,
				}
			}
			
			// Handle multiple values for the same predicate
			if existing, exists := subjectMap[subjectID][predicateID]; exists {
				// Convert to array if not already
				if arr, isArray := existing.([]interface{}); isArray {
					subjectMap[subjectID][predicateID] = append(arr, objectValue)
				} else {
					subjectMap[subjectID][predicateID] = []interface{}{existing, objectValue}
				}
			} else {
				subjectMap[subjectID][predicateID] = objectValue
			}
		}
		
		for _, subjectData := range subjectMap {
			defaultTriples = append(defaultTriples, subjectData)
		}
		result["@graph"] = defaultTriples
	}
	
	// Handle named graphs
	namedGraphs := d.GetNamedGraphs()
	for _, graphName := range namedGraphs {
		graph := d.GetGraph(graphName)
		if graph.Len() > 0 {
			var graphTriples []map[string]interface{}
			subjectMap := make(map[string]map[string]interface{})
			
			for triple := range graph.IterTriples() {
				subjectID := termToJSONLDID(triple.Subject)
				predicateID := termToJSONLDID(triple.Predicate)
				objectValue := termToJSONLDValue(triple.Object)
				
				if _, exists := subjectMap[subjectID]; !exists {
					subjectMap[subjectID] = map[string]interface{}{
						"@id": subjectID,
					}
				}
				
				// Handle multiple values for the same predicate
				if existing, exists := subjectMap[subjectID][predicateID]; exists {
					// Convert to array if not already
					if arr, isArray := existing.([]interface{}); isArray {
						subjectMap[subjectID][predicateID] = append(arr, objectValue)
					} else {
						subjectMap[subjectID][predicateID] = []interface{}{existing, objectValue}
					}
				} else {
					subjectMap[subjectID][predicateID] = objectValue
				}
			}
			
			for _, subjectData := range subjectMap {
				graphTriples = append(graphTriples, subjectData)
			}
			
			graphNameID := termToJSONLDID(graphName)
			result[graphNameID] = map[string]interface{}{
				"@graph": graphTriples,
			}
		}
	}
	
	// Use json.NewEncoder to avoid HTML escaping
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// termToJSONLDID converts a term to a JSON-LD @id value
func termToJSONLDID(term Term) string {
	switch t := term.(type) {
	case *Resource:
		return t.URI
	case *BlankNode:
		return "_:" + t.ID
	default:
		return term.String()
	}
}

// termToJSONLDValue converts a term to a JSON-LD value
func termToJSONLDValue(term Term) interface{} {
	switch t := term.(type) {
	case *Resource:
		return map[string]string{"@id": t.URI}
	case *BlankNode:
		return map[string]string{"@id": "_:" + t.ID}
	case *Literal:
		result := map[string]string{"@value": t.Value}
		if len(t.Language) > 0 {
			result["@language"] = t.Language
		}
		if t.Datatype != nil {
			result["@type"] = termToJSONLDID(t.Datatype)
		}
		return result
	default:
		return term.String()
	}
}

// LoadURI loads RDF data from a specific URI into the dataset
func (d *Dataset) LoadURI(uri string) error {
	doc := defrag(uri)
	q, err := http.NewRequest("GET", doc, nil)
	if err != nil {
		return err
	}
	if len(d.uri) == 0 {
		d.uri = doc
	}
	q.Header.Set("Accept", "application/trig;q=1,text/turtle;q=0.8,application/ld+json;q=0.5")
	r, err := d.httpClient.Do(q)
	if err != nil {
		return err
	}
	if r != nil {
		defer r.Body.Close()
		if r.StatusCode == 200 {
			d.Parse(r.Body, r.Header.Get("Content-Type"))
		} else {
			return fmt.Errorf("Could not fetch dataset from %s - HTTP %d", uri, r.StatusCode)
		}
	}
	return nil
}

// Merge merges another dataset into this one
func (d *Dataset) Merge(toMerge *Dataset) {
	for quad := range toMerge.IterQuads() {
		d.Add(quad)
	}
}
