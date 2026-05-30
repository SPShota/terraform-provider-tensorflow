package manifest

import (
	"bytes"
	"strings"
	"testing"
)

func TestManifestEncodeDecode(t *testing.T) {
	t.Parallel()

	input := Manifest{
		SchemaVersion:     CurrentSchemaVersion,
		SourceModule:      "tensorflow",
		SourceVersion:     "2.16.1",
		GeneratedBy:       "test",
		Root:              "tf",
		DocumentationBase: "https://www.tensorflow.org/api_docs/python",
		Symbols: []Symbol{
			{Path: "tf", Kind: "module", DocURL: "https://www.tensorflow.org/api_docs/python/tf"},
			{Path: "tf.constant", Kind: "function", Signature: "(value, dtype=None, shape=None, name='Const')"},
		},
	}

	var buf bytes.Buffer
	if err := input.Encode(&buf); err != nil {
		t.Fatalf("Encode() returned error: %v", err)
	}

	got, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode() returned error: %v", err)
	}

	if got.SourceVersion != input.SourceVersion {
		t.Fatalf("SourceVersion = %q, want %q", got.SourceVersion, input.SourceVersion)
	}
}

func TestManifestValidateRejectsDuplicateSymbols(t *testing.T) {
	t.Parallel()

	m := Manifest{
		SchemaVersion: CurrentSchemaVersion,
		SourceModule:  "tensorflow",
		Root:          "tf",
		Symbols: []Symbol{
			{Path: "tf", Kind: "module"},
			{Path: "tf", Kind: "module"},
		},
	}

	if err := m.Validate(); err == nil {
		t.Fatalf("Validate() returned nil error")
	}
}

func TestManifestValidateRejectsUnsortedSymbols(t *testing.T) {
	t.Parallel()

	m := Manifest{
		SchemaVersion: CurrentSchemaVersion,
		SourceModule:  "tensorflow",
		Root:          "tf",
		Symbols: []Symbol{
			{Path: "tf.zeros", Kind: "function"},
			{Path: "tf", Kind: "module"},
		},
	}

	err := m.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}
	if !strings.Contains(err.Error(), "sorted") {
		t.Fatalf("Validate() error = %q, want sorted error", err.Error())
	}
}
