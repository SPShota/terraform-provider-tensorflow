package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

const CurrentSchemaVersion = 1

type Manifest struct {
	SchemaVersion     int      `json:"schema_version"`
	SourceModule      string   `json:"source_module"`
	SourceVersion     string   `json:"source_version,omitempty"`
	GeneratedBy       string   `json:"generated_by"`
	Root              string   `json:"root"`
	DocumentationBase string   `json:"documentation_base"`
	Symbols           []Symbol `json:"symbols"`
}

type Symbol struct {
	Path      string   `json:"path"`
	Kind      string   `json:"kind"`
	Module    string   `json:"module,omitempty"`
	Signature string   `json:"signature,omitempty"`
	DocURL    string   `json:"doc_url,omitempty"`
	Children  []string `json:"children,omitempty"`
}

func Decode(r io.Reader) (Manifest, error) {
	var m Manifest
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&m); err != nil {
		return Manifest{}, err
	}

	if err := m.Validate(); err != nil {
		return Manifest{}, err
	}

	return m, nil
}

func (m Manifest) Encode(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m)
}

func (m Manifest) Validate() error {
	if m.SchemaVersion != CurrentSchemaVersion {
		return fmt.Errorf("unsupported schema version %d", m.SchemaVersion)
	}
	if m.SourceModule == "" {
		return fmt.Errorf("source module must not be empty")
	}
	if m.Root == "" {
		return fmt.Errorf("root must not be empty")
	}

	seen := make(map[string]struct{}, len(m.Symbols))
	for _, symbol := range m.Symbols {
		if symbol.Path == "" {
			return fmt.Errorf("symbol path must not be empty")
		}
		if symbol.Kind == "" {
			return fmt.Errorf("symbol %q kind must not be empty", symbol.Path)
		}
		if _, ok := seen[symbol.Path]; ok {
			return fmt.Errorf("duplicate symbol path %q", symbol.Path)
		}
		seen[symbol.Path] = struct{}{}
	}

	if !sort.SliceIsSorted(m.Symbols, func(i, j int) bool {
		return m.Symbols[i].Path < m.Symbols[j].Path
	}) {
		return fmt.Errorf("symbols must be sorted by path")
	}

	return nil
}
