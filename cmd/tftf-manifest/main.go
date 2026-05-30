package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/SPShota/terraform-provider-tensorflow/internal/manifest"
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
	pythonBin := flags.String("python", "python3", "Python executable used for TensorFlow introspection")
	module := flags.String("module", "tensorflow", "Python module to introspect")
	root := flags.String("root", "tf", "Root symbol name written to the manifest")
	maxDepth := flags.Int("max-depth", 3, "Maximum namespace traversal depth")
	output := flags.String("output", "", "Output path. Defaults to stdout")
	if err := flags.Parse(args); err != nil {
		return err
	}

	m, err := manifest.Generate(context.Background(), manifest.GenerateOptions{
		PythonBin: *pythonBin,
		Module:    *module,
		Root:      *root,
		MaxDepth:  *maxDepth,
		Env:       os.Environ(),
	})
	if err != nil {
		return err
	}

	if *output == "" {
		return m.Encode(os.Stdout)
	}

	file, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()

	return m.Encode(file)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  tftf-manifest generate [flags]

Flags:
  -python string     Python executable used for TensorFlow introspection (default "python3")
  -module string     Python module to introspect (default "tensorflow")
  -root string       Root symbol name written to the manifest (default "tf")
  -max-depth int     Maximum namespace traversal depth (default 3)
  -output string     Output path. Defaults to stdout
`)
}
