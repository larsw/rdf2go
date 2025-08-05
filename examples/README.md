# TriG and JSON-LD Examples

This directory contains comprehensive examples demonstrating the TriG (RDF Dataset Language) support added to rdf2go, along with enhanced JSON-LD functionality and format conversion capabilities.

## Example Files

### Core Examples

- **`example_trig.go`** - Basic TriG parsing, named graph handling, and serialization
- **`test_jsonld.go`** - Verifies clean JSON-LD output and validates JSON-LD parsing

### Format Conversion Demos

- **`format_demo/format_conversion_demo.go`** - Comprehensive format conversion between TriG, Turtle, and JSON-LD with full round-trip validation
- **`final_roundtrip_demo.go`** - Focused round-trip testing with detailed analysis of what works and current limitations
- **`advanced_roundtrip_demo.go`** - Advanced round-trip testing with multiple test cases

### Legacy Examples

- **`roundtrip_demo.go`** - Early round-trip demonstration
- **`enhanced_roundtrip.go`** - Enhanced round-trip with validation
- **`format_conversion_demo.go`** - Original format conversion demo (superseded by format_demo/)

## Running the Examples

Each example can be run independently:

```bash
# Basic TriG functionality
go run example_trig.go

# JSON-LD quality verification
go run test_jsonld.go

# Comprehensive format conversion with round-trip validation
cd format_demo && go run format_conversion_demo.go

# Focused round-trip analysis
go run final_roundtrip_demo.go
```

## Key Features Demonstrated

### ✅ What Works Perfectly
- TriG parsing and serialization
- JSON-LD generation with clean output (no HTML escaping)
- Semantic triple content preservation during format conversion
- Resource URI preservation
- Multiple format support (TriG, Turtle, JSON-LD)
- Named graph handling in TriG format
- Dataset and Graph abstractions

### ⚠️ Current Limitations
- Named graph structure is simplified in JSON-LD round-trips due to the underlying JSON-LD library
- Literal datatypes may be normalized (e.g., plain strings get xsd:string datatype)
- JSON-LD library may introduce metadata triples in some cases

## Recommended Usage

- **Use TriG** for full RDF dataset/named graph fidelity
- **Use JSON-LD** for web APIs and semantic web integration  
- **For round-trip scenarios**, focus on semantic content preservation rather than exact structural preservation
- All core RDF functionality (triples, resources, literals) works perfectly across formats

## Library Extensions

These examples demonstrate the enhanced rdf2go library with:
- Full TriG support through `Dataset` type
- Clean JSON-LD serialization with proper `@id`/`@value` structure
- Enhanced MIME type handling for TriG
- Backward compatibility with existing `Graph` API
- Round-trip conversion capabilities between all supported formats
