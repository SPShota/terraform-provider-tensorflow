package python

import "testing"

func TestCallNamed(t *testing.T) {
	t.Parallel()

	arg, err := Literal([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("Literal() returned error: %v", err)
	}

	dtype, err := Reference("tf.float32")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	expr, err := CallNamed("tf.constant", []Expression{arg}, []KeywordArgument{
		{Name: "dtype", Value: dtype},
		{Name: "name", Value: mustLiteral(t, "x")},
	})
	if err != nil {
		t.Fatalf("CallNamed() returned error: %v", err)
	}

	want := "tf.constant([1, 2, 3], dtype=tf.float32, name=\"x\")"
	if expr.Code() != want {
		t.Fatalf("CallNamed().Code() = %q, want %q", expr.Code(), want)
	}
}

func TestCallRejectsDuplicateKeywordArguments(t *testing.T) {
	t.Parallel()

	expr, err := Reference("tf.constant")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	one := mustLiteral(t, 1)
	if _, err := Call(expr, nil, []KeywordArgument{
		{Name: "name", Value: one},
		{Name: "name", Value: one},
	}); err == nil {
		t.Fatalf("Call() returned nil error")
	}
}

func mustLiteral(t *testing.T, value any) Expression {
	t.Helper()

	expr, err := Literal(value)
	if err != nil {
		t.Fatalf("Literal() returned error: %v", err)
	}

	return expr
}
