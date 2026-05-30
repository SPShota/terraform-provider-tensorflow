package manifest

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGenerateWithFakeTensorFlowModule(t *testing.T) {
	t.Parallel()

	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not available")
	}

	moduleDir := filepath.Join(t.TempDir(), "tensorflow")
	if err := os.MkdirAll(moduleDir, 0o755); err != nil {
		t.Fatalf("create fake module dir: %v", err)
	}

	module := `__version__ = "fake"

def constant(value, dtype=None):
    return value

class math:
    @staticmethod
    def reduce_sum(value):
        return sum(value)
`
	if err := os.WriteFile(filepath.Join(moduleDir, "__init__.py"), []byte(module), 0o600); err != nil {
		t.Fatalf("write fake module: %v", err)
	}

	env := append(os.Environ(), "PYTHONPATH="+filepath.Dir(moduleDir))
	m, err := Generate(context.Background(), GenerateOptions{
		PythonBin: python,
		Module:    "tensorflow",
		Root:      "tf",
		MaxDepth:  2,
		Env:       env,
	})
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	if m.SourceVersion != "fake" {
		t.Fatalf("SourceVersion = %q, want fake", m.SourceVersion)
	}

	paths := make(map[string]Symbol, len(m.Symbols))
	for _, symbol := range m.Symbols {
		paths[symbol.Path] = symbol
	}

	for _, path := range []string{"tf", "tf.constant", "tf.math", "tf.math.reduce_sum"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("expected symbol %q in manifest; got %#v", path, paths)
		}
	}

	if got := paths["tf.constant"].Kind; got != "function" {
		t.Fatalf("tf.constant kind = %q, want function", got)
	}
}

func TestGenerateRejectsNegativeDepth(t *testing.T) {
	t.Parallel()

	_, err := Generate(context.Background(), GenerateOptions{MaxDepth: -1})
	if err == nil {
		t.Fatalf("Generate() returned nil error")
	}
}

func TestGenerateReportsPythonErrors(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("shell path assumptions are Unix-specific")
	}

	_, err := Generate(context.Background(), GenerateOptions{
		PythonBin: "/bin/sh",
		Module:    "missing",
		Root:      "tf",
		MaxDepth:  1,
	})
	if err == nil {
		t.Fatalf("Generate() returned nil error")
	}
}
