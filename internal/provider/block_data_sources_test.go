package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestBlockDataSourceMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dataSource datasource.DataSource
		want       string
	}{
		{name: "return", dataSource: NewReturnDataSource(), want: "tensorflow_return"},
		{name: "with", dataSource: NewWithDataSource(), want: "tensorflow_with"},
		{name: "function", dataSource: NewFunctionDataSource(), want: "tensorflow_function"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resp datasource.MetadataResponse
			tt.dataSource.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tensorflow"}, &resp)

			if resp.TypeName != tt.want {
				t.Fatalf("TypeName = %q, want %q", resp.TypeName, tt.want)
			}
		})
	}
}

func TestBlockDataSourceSchemas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dataSource datasource.DataSource
		attributes []string
	}{
		{name: "return", dataSource: NewReturnDataSource(), attributes: []string{"value", "statement"}},
		{name: "with", dataSource: NewWithDataSource(), attributes: []string{"context", "alias", "statements", "statement"}},
		{name: "function", dataSource: NewFunctionDataSource(), attributes: []string{"name", "args", "decorators", "statements", "expression", "statement"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resp datasource.SchemaResponse
			tt.dataSource.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

			for _, attribute := range tt.attributes {
				if _, ok := resp.Schema.Attributes[attribute]; !ok {
					t.Fatalf("schema is missing %q attribute", attribute)
				}
			}
		})
	}
}
