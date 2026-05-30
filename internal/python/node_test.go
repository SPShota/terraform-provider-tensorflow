package python

import "testing"

func TestReferenceAttributeAndAssign(t *testing.T) {
	t.Parallel()

	tf, err := Reference("tf")
	if err != nil {
		t.Fatalf("Reference() returned error: %v", err)
	}

	keras, err := Attribute(tf, "keras")
	if err != nil {
		t.Fatalf("Attribute() returned error: %v", err)
	}

	if keras.Code() != "tf.keras" {
		t.Fatalf("Attribute().Code() = %q, want %q", keras.Code(), "tf.keras")
	}

	stmt, err := Assign("model", keras)
	if err != nil {
		t.Fatalf("Assign() returned error: %v", err)
	}

	if stmt.Code() != "model = tf.keras" {
		t.Fatalf("Assign().Code() = %q, want %q", stmt.Code(), "model = tf.keras")
	}
}

func TestExpressionStatement(t *testing.T) {
	t.Parallel()

	expr, err := RawExpression("print(x)")
	if err != nil {
		t.Fatalf("RawExpression() returned error: %v", err)
	}

	stmt, err := ExpressionStatement(expr)
	if err != nil {
		t.Fatalf("ExpressionStatement() returned error: %v", err)
	}

	if stmt.Code() != "print(x)" {
		t.Fatalf("ExpressionStatement().Code() = %q, want %q", stmt.Code(), "print(x)")
	}
}
