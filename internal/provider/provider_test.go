package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderMetadata(t *testing.T) {
	t.Parallel()

	p := New("test")()
	var resp provider.MetadataResponse

	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

	if resp.TypeName != "tensorflow" {
		t.Fatalf("expected provider type name tensorflow, got %q", resp.TypeName)
	}

	if resp.Version != "test" {
		t.Fatalf("expected provider version test, got %q", resp.Version)
	}
}

func TestProviderSchema(t *testing.T) {
	t.Parallel()

	p := New("test")()
	var resp provider.SchemaResponse

	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	if resp.Schema.Attributes != nil {
		t.Fatalf("expected no provider attributes, got %d", len(resp.Schema.Attributes))
	}
}

func TestProviderDataSources(t *testing.T) {
	t.Parallel()

	p := New("test")()
	dataSources := p.DataSources(context.Background())

	expectedCount := 10 + len(GeneratedDataSources())
	if len(dataSources) != expectedCount {
		t.Fatalf("expected %d data sources, got %d", expectedCount, len(dataSources))
	}

	var resp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

	got := make(map[string]struct{}, len(dataSources))
	for _, newDataSource := range dataSources {
		ds := newDataSource()
		var dsResp datasource.MetadataResponse
		ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: resp.TypeName}, &dsResp)
		got[dsResp.TypeName] = struct{}{}
	}

	for _, name := range []string{
		"tensorflow_literal",
		"tensorflow_ref",
		"tensorflow_attr",
		"tensorflow_call",
		"tensorflow_assign",
		"tensorflow_return",
		"tensorflow_with",
		"tensorflow_function",
		"tensorflow_raw_op",
		"tensorflow_program",
		"tensorflow_constant",
		"tensorflow_reshape",
		"tensorflow_math_reduce_sum",
	} {
		if _, ok := got[name]; !ok {
			t.Fatalf("expected %q data source to be registered; got %#v", name, got)
		}
	}
}
