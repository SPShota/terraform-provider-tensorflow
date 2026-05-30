package python

import (
	"math"
	"testing"
)

func TestLiteral(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "none", value: nil, want: "None"},
		{name: "string", value: "hello\nworld", want: "\"hello\\nworld\""},
		{name: "true", value: true, want: "True"},
		{name: "false", value: false, want: "False"},
		{name: "int", value: 42, want: "42"},
		{name: "float", value: 1.25, want: "1.25"},
		{name: "slice", value: []int{1, 2, 3}, want: "[1, 2, 3]"},
		{name: "nested", value: []any{"x", []any{1, false}}, want: "[\"x\", [1, False]]"},
		{name: "map sorted", value: map[string]any{"b": 2, "a": 1}, want: "{\"a\": 1, \"b\": 2}"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Literal(tt.value)
			if err != nil {
				t.Fatalf("Literal() returned error: %v", err)
			}

			if got.Code() != tt.want {
				t.Fatalf("Literal().Code() = %q, want %q", got.Code(), tt.want)
			}
		})
	}
}

func TestLiteralRejectsUnsupportedValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{name: "nan", value: math.NaN()},
		{name: "inf", value: math.Inf(1)},
		{name: "non string map key", value: map[int]any{1: "x"}},
		{name: "struct", value: struct{ Name string }{Name: "x"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := Literal(tt.value); err == nil {
				t.Fatalf("Literal() returned nil error")
			}
		})
	}
}
