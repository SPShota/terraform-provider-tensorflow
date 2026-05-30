package manifest

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

//go:embed introspect.py
var introspectScript string

type GenerateOptions struct {
	PythonBin string
	Module    string
	Root      string
	MaxDepth  int
	Env       []string
}

func Generate(ctx context.Context, opts GenerateOptions) (Manifest, error) {
	pythonBin := opts.PythonBin
	if pythonBin == "" {
		pythonBin = "python3"
	}

	module := opts.Module
	if module == "" {
		module = "tensorflow"
	}

	root := opts.Root
	if root == "" {
		root = "tf"
	}

	maxDepth := opts.MaxDepth
	if maxDepth == 0 {
		maxDepth = 3
	}
	if maxDepth < 0 {
		return Manifest{}, fmt.Errorf("max depth must be >= 0")
	}

	cmd := exec.CommandContext(ctx, pythonBin, "-c", introspectScript, "--module", module, "--root", root, "--max-depth", strconv.Itoa(maxDepth))
	cmd.Env = opts.Env
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Manifest{}, fmt.Errorf("generate manifest: %w: %s", err, bytes.TrimSpace(exitErr.Stderr))
		}
		return Manifest{}, fmt.Errorf("generate manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(output, &m); err != nil {
		return Manifest{}, fmt.Errorf("decode manifest JSON: %w", err)
	}
	if err := m.Validate(); err != nil {
		return Manifest{}, err
	}

	return m, nil
}
