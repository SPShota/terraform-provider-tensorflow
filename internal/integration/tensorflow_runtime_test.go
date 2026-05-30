package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	tftfpython "github.com/SPShota/terraform-provider-tensorflow/internal/python"
)

func TestGeneratedProgramRunsWithTensorFlow(t *testing.T) {
	if os.Getenv("TF_TF_INTEGRATION") != "1" {
		t.Skip("set TF_TF_INTEGRATION=1 to run TensorFlow integration tests")
	}

	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not available")
	}

	if output, err := exec.Command(python, "-c", "import tensorflow").CombinedOutput(); err != nil {
		t.Skipf("tensorflow is not importable: %v\n%s", err, output)
	}

	content := generatedRuntimeProgram(t)
	path := filepath.Join(t.TempDir(), "generated.py")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write generated Python: %v", err)
	}

	cmd := exec.Command(python, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated TensorFlow program failed: %v\n%s\nGenerated Python:\n%s", err, output, content)
	}

	if got := strings.TrimSpace(string(output)); got != "6.0" {
		t.Fatalf("generated TensorFlow program output = %q, want %q\nGenerated Python:\n%s", got, "6.0", content)
	}
}

func generatedRuntimeProgram(t *testing.T) string {
	t.Helper()

	values, err := tftfpython.Literal([]float64{1, 2, 3})
	if err != nil {
		t.Fatalf("Literal() returned error: %v", err)
	}

	dtype, err := tftfpython.Reference("tf.float32")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	constant, err := tftfpython.CallNamed("tf.constant", []tftfpython.Expression{values}, []tftfpython.KeywordArgument{
		{Name: "dtype", Value: dtype},
	})
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	assign, err := tftfpython.Assign("x", constant)
	if err != nil {
		t.Fatalf("Assign() returned error: %v", err)
	}

	ref, err := tftfpython.Reference("x")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	reduceSum, err := tftfpython.CallNamed("tf.reduce_sum", []tftfpython.Expression{ref}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	numpyValue, err := tftfpython.Attribute(reduceSum, "numpy")
	if err != nil {
		t.Fatalf("Attribute() returned error: %v", err)
	}

	numpyCall, err := tftfpython.Call(numpyValue, nil, nil)
	if err != nil {
		t.Fatalf("Call() returned error: %v", err)
	}

	printCall, err := tftfpython.CallNamed("print", []tftfpython.Expression{numpyCall}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printStatement, err := tftfpython.ExpressionStatement(printCall)
	if err != nil {
		t.Fatalf("ExpressionStatement() returned error: %v", err)
	}

	content, err := (tftfpython.Program{Statements: []tftfpython.Statement{assign, printStatement}}).Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	return content
}
