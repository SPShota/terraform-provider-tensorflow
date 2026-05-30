package provider

import (
	"context"
	"testing"

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
