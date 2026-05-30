package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/SPShota/terraform-provider-tensorflow/internal/coverage"
	"github.com/SPShota/terraform-provider-tensorflow/internal/manifest"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "report":
		if err := report(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func report(args []string) error {
	flags := flag.NewFlagSet("report", flag.ContinueOnError)
	input := flags.String("input", "", "Input manifest JSON path")
	output := flags.String("output", "", "Output path. Defaults to stdout")
	format := flags.String("format", "markdown", "Output format: markdown or json")
	includeRawOps := flags.Bool("include-raw-ops", true, "Treat tf.raw_ops.* as covered by tensorflow_raw_op")
	limitMissing := flags.Int("limit-missing", 100, "Maximum missing symbols to include. Use 0 for no limit")
	limitCovered := flags.Int("limit-covered", 0, "Maximum covered symbols to include. Use 0 for no limit")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if *input == "" {
		return fmt.Errorf("-input is required")
	}

	inputFile, err := os.Open(*input)
	if err != nil {
		return fmt.Errorf("open input manifest: %w", err)
	}
	defer inputFile.Close()

	m, err := manifest.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("decode input manifest: %w", err)
	}

	result, err := coverage.Build(m, coverage.Options{
		IncludeRawOps: *includeRawOps,
		LimitMissing:  *limitMissing,
		LimitCovered:  *limitCovered,
	})
	if err != nil {
		return err
	}

	var writer io.Writer = os.Stdout
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer file.Close()
		writer = file
	}

	switch *format {
	case "markdown", "md":
		return coverage.WriteMarkdown(writer, result)
	case "json":
		encoder := json.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	default:
		return fmt.Errorf("unsupported format %q", *format)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  tftf-coverage report -input tf-manifest.json [flags]

Flags:
  -input string          Input manifest JSON path
  -output string         Output path. Defaults to stdout
  -format string         Output format: markdown or json (default "markdown")
  -include-raw-ops       Treat tf.raw_ops.* as covered by tensorflow_raw_op (default true)
  -limit-missing int     Maximum missing symbols to include. Use 0 for no limit (default 100)
  -limit-covered int     Maximum covered symbols to include. Use 0 for no limit
`)
}
