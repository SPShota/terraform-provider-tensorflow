package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestExpressionDataSourceMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dataSource datasource.DataSource
		want       string
	}{
		{name: "literal", dataSource: NewLiteralDataSource(), want: "tensorflow_literal"},
		{name: "ref", dataSource: NewRefDataSource(), want: "tensorflow_ref"},
		{name: "attr", dataSource: NewAttrDataSource(), want: "tensorflow_attr"},
		{name: "call", dataSource: NewCallDataSource(), want: "tensorflow_call"},
		{name: "assign", dataSource: NewAssignDataSource(), want: "tensorflow_assign"},
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

func TestExpressionDataSourceSchemas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dataSource datasource.DataSource
		attributes []string
	}{
		{name: "literal", dataSource: NewLiteralDataSource(), attributes: []string{"value_json", "expression"}},
		{name: "ref", dataSource: NewRefDataSource(), attributes: []string{"name", "expression"}},
		{name: "attr", dataSource: NewAttrDataSource(), attributes: []string{"receiver", "name", "expression"}},
		{name: "call", dataSource: NewCallDataSource(), attributes: []string{"function", "args", "kwargs", "expression", "statement"}},
		{name: "assign", dataSource: NewAssignDataSource(), attributes: []string{"name", "value", "expression", "statement"}},
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

func TestLiteralExpressionFromJSON(t *testing.T) {
	t.Parallel()

	expr, err := literalExpressionFromJSON(`{"shape":[2,3],"name":"x","trainable":true}`)
	if err != nil {
		t.Fatalf("literalExpressionFromJSON() returned error: %v", err)
	}

	want := "{\"name\": \"x\", \"shape\": [2, 3], \"trainable\": True}"
	if expr.Code() != want {
		t.Fatalf("Code() = %q, want %q", expr.Code(), want)
	}
}

func TestLiteralExpressionFromJSONRejectsTrailingValues(t *testing.T) {
	t.Parallel()

	if _, err := literalExpressionFromJSON(`1 2`); err == nil {
		t.Fatalf("literalExpressionFromJSON() returned nil error")
	}
}
