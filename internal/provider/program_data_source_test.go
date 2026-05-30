package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestProgramDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := NewProgramDataSource()
	var resp datasource.MetadataResponse

	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tensorflow"}, &resp)

	if resp.TypeName != "tensorflow_program" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "tensorflow_program")
	}
}

func TestProgramDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := NewProgramDataSource()
	var resp datasource.SchemaResponse

	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	for _, name := range []string{"header", "imports", "statements", "content"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("schema is missing %q attribute", name)
		}
	}
}
