package wrappergen

import (
	"strings"
	"testing"

	"github.com/SPShota/terraform-provider-tensorflow/internal/manifest"
)

func TestWrappers(t *testing.T) {
	t.Parallel()

	wrappers, err := Wrappers(testManifest())
	if err != nil {
		t.Fatalf("Wrappers() returned error: %v", err)
	}

	want := []Wrapper{
		{TypeNameSuffix: "constant", Function: "tf.constant", DocURL: "https://www.tensorflow.org/api_docs/python/tf/constant"},
		{TypeNameSuffix: "keras_sequential", Function: "tf.keras.Sequential", DocURL: "https://www.tensorflow.org/api_docs/python/tf/keras/Sequential"},
		{TypeNameSuffix: "math_reduce_sum", Function: "tf.math.reduce_sum", DocURL: "https://www.tensorflow.org/api_docs/python/tf/math/reduce_sum"},
	}

	if len(wrappers) != len(want) {
		t.Fatalf("len(Wrappers()) = %d, want %d: %#v", len(wrappers), len(want), wrappers)
	}

	for i := range want {
		if wrappers[i] != want[i] {
			t.Fatalf("wrapper[%d] = %#v, want %#v", i, wrappers[i], want[i])
		}
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	source, err := Generate(testManifest(), Options{})
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	text := string(source)
	for _, want := range []string{
		"func GeneratedDataSources() []func() datasource.DataSource",
		`TypeNameSuffix: "constant"`,
		`Function:       "tf.constant"`,
		`DocURL:         "https://www.tensorflow.org/api_docs/python/tf/constant"`,
		`TypeNameSuffix: "math_reduce_sum"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("generated source does not contain %q:\n%s", want, text)
		}
	}
}

func TestTypeNameSuffix(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"tf.constant":                 "constant",
		"tf.math.reduce_sum":          "math_reduce_sum",
		"tf.keras.layers.Conv2D":      "keras_layers_conv2_d",
		"tf.experimental.numpy.zeros": "experimental_numpy_zeros",
	}

	for path, want := range tests {
		path := path
		want := want
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			got, err := TypeNameSuffix("tf", path)
			if err != nil {
				t.Fatalf("TypeNameSuffix() returned error: %v", err)
			}
			if got != want {
				t.Fatalf("TypeNameSuffix() = %q, want %q", got, want)
			}
		})
	}
}

func TestTypeNameSuffixRejectsRoot(t *testing.T) {
	t.Parallel()

	if _, err := TypeNameSuffix("tf", "tf"); err == nil {
		t.Fatalf("TypeNameSuffix() returned nil error")
	}
}

func testManifest() manifest.Manifest {
	return manifest.Manifest{
		SchemaVersion:     manifest.CurrentSchemaVersion,
		SourceModule:      "tensorflow",
		SourceVersion:     "fake",
		GeneratedBy:       "test",
		Root:              "tf",
		DocumentationBase: "https://www.tensorflow.org/api_docs/python",
		Symbols: []manifest.Symbol{
			{Path: "tf", Kind: "module"},
			{Path: "tf.constant", Kind: "function", DocURL: "https://www.tensorflow.org/api_docs/python/tf/constant"},
			{Path: "tf.keras", Kind: "module", DocURL: "https://www.tensorflow.org/api_docs/python/tf/keras"},
			{Path: "tf.keras.Sequential", Kind: "class", DocURL: "https://www.tensorflow.org/api_docs/python/tf/keras/Sequential"},
			{Path: "tf.math", Kind: "module", DocURL: "https://www.tensorflow.org/api_docs/python/tf/math"},
			{Path: "tf.math.reduce_sum", Kind: "function", DocURL: "https://www.tensorflow.org/api_docs/python/tf/math/reduce_sum"},
			{Path: "tf.version", Kind: "value", DocURL: "https://www.tensorflow.org/api_docs/python/tf/version"},
		},
	}
}
