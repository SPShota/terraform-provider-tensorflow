package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestWrapperDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := NewWrapperDataSource(WrapperDataSourceSpec{
		TypeNameSuffix: "constant",
		Function:       "tf.constant",
	})

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tensorflow"}, &resp)

	if resp.TypeName != "tensorflow_constant" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "tensorflow_constant")
	}
}

func TestWrapperDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := NewWrapperDataSource(WrapperDataSourceSpec{
		TypeNameSuffix: "constant",
		Function:       "tf.constant",
		DocURL:         "https://www.tensorflow.org/api_docs/python/tf/constant",
	})

	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	for _, name := range []string{"args", "kwargs", "expression", "statement"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("schema is missing %q attribute", name)
		}
	}
}
