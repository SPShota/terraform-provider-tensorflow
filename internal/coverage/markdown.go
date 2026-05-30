package coverage

import (
	"fmt"
	"io"
)

func WriteMarkdown(w io.Writer, report Report) error {
	if _, err := fmt.Fprintf(w, "# TF.tf API Coverage\n\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "- Source module: `%s`\n", report.SourceModule); err != nil {
		return err
	}
	if report.SourceVersion != "" {
		if _, err := fmt.Fprintf(w, "- Source version: `%s`\n", report.SourceVersion); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "- Root: `%s`\n", report.Root); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "- Wrappable symbols: `%d`\n", report.WrappableTotal); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "- Covered symbols: `%d`\n", report.CoveredTotal); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "- Coverage: `%.2f%%`\n", report.CoveragePercent); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "- Generated wrappers: `%d`\n", report.GeneratedWrappers); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "- Raw ops covered by `tensorflow_raw_op`: `%d`\n\n", report.RawOpsCovered); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "## Namespaces"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Namespace | Covered | Total | Coverage |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | ---: | ---: | ---: |"); err != nil {
		return err
	}
	for _, namespace := range report.Namespaces {
		if _, err := fmt.Fprintf(w, "| `%s` | %d | %d | %.2f%% |\n", namespace.Namespace, namespace.CoveredTotal, namespace.WrappableTotal, namespace.CoveragePercent); err != nil {
			return err
		}
	}

	if len(report.Missing) > 0 {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "## Missing Symbols"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
		for _, symbol := range report.Missing {
			if _, err := fmt.Fprintf(w, "- `%s` (%s)\n", symbol.Path, symbol.Kind); err != nil {
				return err
			}
		}
	}

	return nil
}
