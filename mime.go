package rdf2go

import (
	"regexp"
)

var mimeParser = map[string]string{
	"text/turtle":               "turtle",
	"application/trig":          "trig",
	"application/ld+json":       "jsonld",
	"application/sparql-update": "internal",
}

var mimeSerializer = map[string]string{
	"application/ld+json": "jsonld",
	"application/trig":    "trig",
	"text/html":           "internal",
}

var mimeRdfExt = map[string]string{
	".ttl":    "text/turtle",
	".trig":   "application/trig",
	".n3":     "text/n3",
	".rdf":    "application/rdf+xml",
	".jsonld": "application/ld+json",
}

var rdfExtensions = []string{
	".ttl",
	".trig",
	".n3",
	".rdf",
	".jsonld",
}

var (
	serializerMimes = []string{}
	validMimeType   = regexp.MustCompile(`^\w+/\w+$`)
)
