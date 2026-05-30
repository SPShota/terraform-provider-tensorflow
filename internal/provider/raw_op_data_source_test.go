package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestRawOpDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := NewRawOpDataSource()
	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tensorflow"}, &resp)

	if resp.TypeName != "tensorflow_raw_op" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "tensorflow_raw_op")
	}
}

func TestRawOpDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := NewRawOpDataSource()
	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	for _, name := range []string{"op", "args", "kwargs", "expression", "statement"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("schema is missing %q attribute", name)
		}
	}
}
