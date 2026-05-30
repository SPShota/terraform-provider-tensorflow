package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SPShota/terraform-provider-tensorflow/internal/manifest"
	"github.com/SPShota/terraform-provider-tensorflow/internal/wrappergen"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "generate":
		if err := generate(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func generate(args []string) error {
	flags := flag.NewFlagSet("generate", flag.ContinueOnError)
	input := flags.String("input", "", "Input manifest JSON path")
	output := flags.String("output", "", "Output Go file path. Defaults to stdout")
	packageName := flags.String("package", "provider", "Generated Go package name")
	functionName := flags.String("function", "GeneratedDataSources", "Generated Go function name")
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

	source, err := wrappergen.Generate(m, wrappergen.Options{
		PackageName:  *packageName,
		FunctionName: *functionName,
	})
	if err != nil {
		return err
	}

	if *output == "" {
		_, err := os.Stdout.Write(source)
		return err
	}

	return os.WriteFile(*output, source, 0o600)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  tftf-wrappergen generate -input tf-manifest.json [flags]

Flags:
  -input string      Input manifest JSON path
  -output string     Output Go file path. Defaults to stdout
  -package string    Generated Go package name (default "provider")
  -function string   Generated Go function name (default "GeneratedDataSources")
`)
}
