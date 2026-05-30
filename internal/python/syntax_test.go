package python

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGeneratedProgramCompilesWithPython(t *testing.T) {
	t.Parallel()

	values, err := Literal([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("Literal() returned error: %v", err)
	}

	dtype, err := Reference("tf.float32")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	constant, err := CallNamed("tf.constant", []Expression{values}, []KeywordArgument{
		{Name: "dtype", Value: dtype},
	})
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	assign, err := Assign("x", constant)
	if err != nil {
		t.Fatalf("Assign() returned error: %v", err)
	}

	ref, err := Reference("x")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	reduceSum, err := CallNamed("tf.reduce_sum", []Expression{ref}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printCall, err := CallNamed("print", []Expression{reduceSum}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printStatement, err := ExpressionStatement(printCall)
	if err != nil {
		t.Fatalf("ExpressionStatement() returned error: %v", err)
	}

	content, err := Program{Statements: []Statement{assign, printStatement}}.Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	compilePython(t, content)
}

func TestCustomProgramCompilesWithPython(t *testing.T) {
	t.Parallel()

	statement, err := RawStatement("print(Path.cwd())")
	if err != nil {
		t.Fatalf("RawStatement() returned error: %v", err)
	}

	content, err := Program{
		Header:     "# custom",
		Imports:    []string{"from pathlib import Path"},
		Statements: []Statement{statement},
	}.Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	compilePython(t, content)
}

func TestBlockProgramCompilesWithPython(t *testing.T) {
	t.Parallel()

	value, err := Reference("x")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	returnStatement, err := Return(value)
	if err != nil {
		t.Fatalf("Return() returned error: %v", err)
	}

	functionStatement, err := FunctionDef("identity", []string{"x"}, nil, []Statement{returnStatement})
	if err != nil {
		t.Fatalf("FunctionDef() returned error: %v", err)
	}

	call, err := CallNamed("identity", []Expression{mustSyntaxLiteral(t, 1)}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printCall, err := CallNamed("print", []Expression{call}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printStatement, err := ExpressionStatement(printCall)
	if err != nil {
		t.Fatalf("ExpressionStatement() returned error: %v", err)
	}

	content, err := Program{Imports: []string{}, Statements: []Statement{functionStatement, printStatement}}.Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	compilePython(t, content)
}

func TestWithProgramCompilesWithPython(t *testing.T) {
	t.Parallel()

	contextExpr, err := CallNamed("nullcontext", []Expression{mustSyntaxLiteral(t, 1)}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printCall, err := CallNamed("print", []Expression{mustSyntaxExpression(t, "value")}, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	printStatement, err := ExpressionStatement(printCall)
	if err != nil {
		t.Fatalf("ExpressionStatement() returned error: %v", err)
	}

	withStatement, err := With(contextExpr, "value", []Statement{printStatement})
	if err != nil {
		t.Fatalf("With() returned error: %v", err)
	}

	content, err := Program{
		Imports:    []string{"from contextlib import nullcontext"},
		Statements: []Statement{withStatement},
	}.Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	compilePython(t, content)
}

func compilePython(t *testing.T, content string) {
	t.Helper()

	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not available")
	}

	path := filepath.Join(t.TempDir(), "generated.py")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write generated Python: %v", err)
	}

	cmd := exec.Command(python, "-m", "py_compile", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("python syntax check failed: %v\n%s\nGenerated Python:\n%s", err, output, content)
	}
}

func mustSyntaxLiteral(t *testing.T, value any) Expression {
	t.Helper()

	expr, err := Literal(value)
	if err != nil {
		t.Fatalf("Literal() returned error: %v", err)
	}

	return expr
}

func mustSyntaxExpression(t *testing.T, code string) Expression {
	t.Helper()

	expr, err := RawExpression(code)
	if err != nil {
		t.Fatalf("RawExpression() returned error: %v", err)
	}

	return expr
}
