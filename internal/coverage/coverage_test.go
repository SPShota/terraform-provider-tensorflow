package coverage

import (
	"bytes"
	"strings"
	"testing"

	"github.com/SPShota/terraform-provider-tensorflow/internal/manifest"
)

func TestBuild(t *testing.T) {
	t.Parallel()

	report, err := Build(testManifest(), Options{IncludeRawOps: true})
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	if report.WrappableTotal != 4 {
		t.Fatalf("WrappableTotal = %d, want 4", report.WrappableTotal)
	}
	if report.CoveredTotal != 3 {
		t.Fatalf("CoveredTotal = %d, want 3", report.CoveredTotal)
	}
	if report.RawOpsCovered != 1 {
		t.Fatalf("RawOpsCovered = %d, want 1", report.RawOpsCovered)
	}
	if len(report.Missing) != 1 {
		t.Fatalf("len(Missing) = %d, want 1: %#v", len(report.Missing), report.Missing)
	}
}

func TestBuildWithoutRawOps(t *testing.T) {
	t.Parallel()

	report, err := Build(testManifest(), Options{IncludeRawOps: false})
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	for _, symbol := range report.Missing {
		if symbol.Path == "tf.raw_ops.AddV2" {
			return
		}
	}

	t.Fatalf("expected tf.raw_ops.AddV2 to be missing: %#v", report.Missing)
}

func TestWriteMarkdown(t *testing.T) {
	t.Parallel()

	report, err := Build(testManifest(), Options{IncludeRawOps: true})
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	var buf bytes.Buffer
	if err := WriteMarkdown(&buf, report); err != nil {
		t.Fatalf("WriteMarkdown() returned error: %v", err)
	}

	text := buf.String()
	for _, want := range []string{"# TF.tf API Coverage", "| `tf` |", "`tf.missing`"} {
		if !strings.Contains(text, want) {
			t.Fatalf("markdown does not contain %q:\n%s", want, text)
		}
	}
}

func TestRegisteredWrappers(t *testing.T) {
	t.Parallel()

	wrappers, err := RegisteredWrappers()
	if err != nil {
		t.Fatalf("RegisteredWrappers() returned error: %v", err)
	}

	if len(wrappers) == 0 {
		t.Fatalf("RegisteredWrappers() returned no wrappers")
	}
}

func testManifest() manifest.Manifest {
	return manifest.Manifest{
		SchemaVersion: manifest.CurrentSchemaVersion,
		SourceModule:  "tensorflow",
		SourceVersion: "fake",
		Root:          "tf",
		Symbols: []manifest.Symbol{
			{Path: "tf", Kind: "module"},
			{Path: "tf.constant", Kind: "function"},
			{Path: "tf.keras", Kind: "module"},
			{Path: "tf.keras.Sequential", Kind: "class"},
			{Path: "tf.missing", Kind: "function"},
			{Path: "tf.raw_ops", Kind: "module"},
			{Path: "tf.raw_ops.AddV2", Kind: "function"},
			{Path: "tf.value", Kind: "value"},
		},
	}
}
