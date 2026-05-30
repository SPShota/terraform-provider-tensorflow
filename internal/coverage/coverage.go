package coverage

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SPShota/terraform-provider-tensorflow/internal/manifest"
	"github.com/SPShota/terraform-provider-tensorflow/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type Report struct {
	SourceModule      string             `json:"source_module"`
	SourceVersion     string             `json:"source_version,omitempty"`
	Root              string             `json:"root"`
	WrappableTotal    int                `json:"wrappable_total"`
	CoveredTotal      int                `json:"covered_total"`
	CoveragePercent   float64            `json:"coverage_percent"`
	GeneratedWrappers int                `json:"generated_wrappers"`
	RawOpsCovered     int                `json:"raw_ops_covered"`
	Namespaces        []NamespaceSummary `json:"namespaces"`
	Missing           []SymbolCoverage   `json:"missing"`
	Covered           []SymbolCoverage   `json:"covered"`
}

type NamespaceSummary struct {
	Namespace       string  `json:"namespace"`
	WrappableTotal  int     `json:"wrappable_total"`
	CoveredTotal    int     `json:"covered_total"`
	CoveragePercent float64 `json:"coverage_percent"`
}

type SymbolCoverage struct {
	Path      string `json:"path"`
	Kind      string `json:"kind"`
	CoveredBy string `json:"covered_by,omitempty"`
}

type Options struct {
	IncludeRawOps bool
	LimitMissing  int
	LimitCovered  int
}

type wrapperSpecProvider interface {
	WrapperSpec() provider.WrapperDataSourceSpec
}

func Build(m manifest.Manifest, opts Options) (Report, error) {
	if err := m.Validate(); err != nil {
		return Report{}, err
	}

	wrappers, err := RegisteredWrappers()
	if err != nil {
		return Report{}, err
	}

	wrapperByPath := make(map[string]provider.WrapperDataSourceSpec, len(wrappers))
	for _, wrapper := range wrappers {
		wrapperByPath[wrapper.Function] = wrapper
	}

	report := Report{
		SourceModule:      m.SourceModule,
		SourceVersion:     m.SourceVersion,
		Root:              m.Root,
		GeneratedWrappers: len(wrappers),
	}

	namespaceTotals := map[string]*NamespaceSummary{}
	for _, symbol := range m.Symbols {
		if !isWrappable(symbol) {
			continue
		}

		report.WrappableTotal++
		namespace := namespaceFor(m.Root, symbol.Path)
		summary := namespaceTotals[namespace]
		if summary == nil {
			summary = &NamespaceSummary{Namespace: namespace}
			namespaceTotals[namespace] = summary
		}
		summary.WrappableTotal++

		coveredBy := ""
		if wrapper, ok := wrapperByPath[symbol.Path]; ok {
			coveredBy = "tensorflow_" + wrapper.TypeNameSuffix
		} else if opts.IncludeRawOps && isRawOp(m.Root, symbol.Path) {
			coveredBy = "tensorflow_raw_op"
			report.RawOpsCovered++
		}

		coverage := SymbolCoverage{
			Path:      symbol.Path,
			Kind:      symbol.Kind,
			CoveredBy: coveredBy,
		}
		if coveredBy == "" {
			report.Missing = append(report.Missing, coverage)
			continue
		}

		report.CoveredTotal++
		summary.CoveredTotal++
		report.Covered = append(report.Covered, coverage)
	}

	report.CoveragePercent = percent(report.CoveredTotal, report.WrappableTotal)
	report.Namespaces = namespaceSummaries(namespaceTotals)
	limitSymbols(&report.Missing, opts.LimitMissing)
	limitSymbols(&report.Covered, opts.LimitCovered)

	return report, nil
}

func RegisteredWrappers() ([]provider.WrapperDataSourceSpec, error) {
	return wrappersFromConstructors(provider.GeneratedDataSources())
}

func wrappersFromConstructors(constructors []func() datasource.DataSource) ([]provider.WrapperDataSourceSpec, error) {
	wrappers := make([]provider.WrapperDataSourceSpec, 0, len(constructors))
	seen := map[string]struct{}{}

	for _, constructor := range constructors {
		dataSource := constructor()
		specProvider, ok := dataSource.(wrapperSpecProvider)
		if !ok {
			return nil, fmt.Errorf("generated data source does not expose wrapper spec: %T", dataSource)
		}

		spec := specProvider.WrapperSpec()
		if spec.Function == "" {
			return nil, fmt.Errorf("wrapper function must not be empty")
		}
		if _, ok := seen[spec.Function]; ok {
			return nil, fmt.Errorf("duplicate wrapper function %q", spec.Function)
		}
		seen[spec.Function] = struct{}{}
		wrappers = append(wrappers, spec)
	}

	sort.Slice(wrappers, func(i, j int) bool {
		return wrappers[i].Function < wrappers[j].Function
	})

	return wrappers, nil
}

func namespaceSummaries(values map[string]*NamespaceSummary) []NamespaceSummary {
	namespaces := make([]NamespaceSummary, 0, len(values))
	for _, summary := range values {
		summary.CoveragePercent = percent(summary.CoveredTotal, summary.WrappableTotal)
		namespaces = append(namespaces, *summary)
	}

	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].Namespace < namespaces[j].Namespace
	})

	return namespaces
}

func isWrappable(symbol manifest.Symbol) bool {
	switch symbol.Kind {
	case "function", "class", "callable":
		return true
	default:
		return false
	}
}

func isRawOp(root, path string) bool {
	return strings.HasPrefix(path, root+".raw_ops.")
}

func namespaceFor(root, path string) string {
	prefix := root + "."
	if path == root || !strings.HasPrefix(path, prefix) {
		return root
	}

	remainder := strings.TrimPrefix(path, prefix)
	parts := strings.Split(remainder, ".")
	if len(parts) == 0 || parts[0] == "" {
		return root
	}
	if len(parts) == 1 {
		return root
	}

	if len(parts) > 1 && parts[0] == "keras" {
		return root + ".keras." + parts[1]
	}
	if len(parts) > 2 && parts[0] == "data" && parts[1] == "Dataset" {
		return root + ".data.Dataset"
	}

	return root + "." + parts[0]
}

func percent(covered, total int) float64 {
	if total == 0 {
		return 0
	}

	return float64(covered) / float64(total) * 100
}

func limitSymbols(symbols *[]SymbolCoverage, limit int) {
	if limit <= 0 || len(*symbols) <= limit {
		return
	}

	*symbols = (*symbols)[:limit]
}
