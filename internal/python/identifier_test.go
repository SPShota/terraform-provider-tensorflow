package python

import "testing"

func TestValidateIdentifier(t *testing.T) {
	t.Parallel()

	tests := map[string]bool{
		"_value": true,
		"layer1": true,
		"tf":     true,
		"":       false,
		"1x":     false,
		"x-y":    false,
		"class":  false,
	}

	for name, wantValid := range tests {
		name := name
		wantValid := wantValid
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := ValidateIdentifier(name)
			if (err == nil) != wantValid {
				t.Fatalf("ValidateIdentifier(%q) valid = %v, want %v; err = %v", name, err == nil, wantValid, err)
			}
		})
	}
}

func TestValidateDottedIdentifier(t *testing.T) {
	t.Parallel()

	tests := map[string]bool{
		"tf":                    true,
		"tf.keras.Sequential":   true,
		"tf.keras.layers.Dense": true,
		".tf":                   false,
		"tf.":                   false,
		"tf.class":              false,
		"tf.reduce-sum":         false,
	}

	for name, wantValid := range tests {
		name := name
		wantValid := wantValid
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := ValidateDottedIdentifier(name)
			if (err == nil) != wantValid {
				t.Fatalf("ValidateDottedIdentifier(%q) valid = %v, want %v; err = %v", name, err == nil, wantValid, err)
			}
		})
	}
}
