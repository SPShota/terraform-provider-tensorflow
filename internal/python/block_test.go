package python

import "testing"

func TestReturn(t *testing.T) {
	t.Parallel()

	expr, err := Reference("loss")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	statement, err := Return(expr)
	if err != nil {
		t.Fatalf("Return() returned error: %v", err)
	}

	if statement.Code() != "return loss" {
		t.Fatalf("Return().Code() = %q, want %q", statement.Code(), "return loss")
	}
}

func TestWith(t *testing.T) {
	t.Parallel()

	contextExpr, err := CallNamed("tf.GradientTape", nil, nil)
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	body, err := RawStatement("loss = model(x)")
	if err != nil {
		t.Fatalf("RawStatement() returned error: %v", err)
	}

	statement, err := With(contextExpr, "tape", []Statement{body})
	if err != nil {
		t.Fatalf("With() returned error: %v", err)
	}

	want := "with tf.GradientTape() as tape:\n\tloss = model(x)"
	if statement.Code() != want {
		t.Fatalf("With().Code() = %q, want %q", statement.Code(), want)
	}
}

func TestFunctionDef(t *testing.T) {
	t.Parallel()

	decorator, err := Reference("tf.function")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	returnValue, err := Reference("loss")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	returnStatement, err := Return(returnValue)
	if err != nil {
		t.Fatalf("Return() returned error: %v", err)
	}

	statement, err := FunctionDef("train_step", []string{"x", "y"}, []Expression{decorator}, []Statement{returnStatement})
	if err != nil {
		t.Fatalf("FunctionDef() returned error: %v", err)
	}

	want := "@tf.function\ndef train_step(x, y):\n\treturn loss"
	if statement.Code() != want {
		t.Fatalf("FunctionDef().Code() = %q, want %q", statement.Code(), want)
	}
}

func TestBlockWithEmptyBodyUsesPass(t *testing.T) {
	t.Parallel()

	statement, err := Block("if True", nil)
	if err != nil {
		t.Fatalf("Block() returned error: %v", err)
	}

	want := "if True:\n\tpass"
	if statement.Code() != want {
		t.Fatalf("Block().Code() = %q, want %q", statement.Code(), want)
	}
}
