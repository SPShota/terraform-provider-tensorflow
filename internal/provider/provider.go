package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &tensorflowProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &tensorflowProvider{
			version: version,
		}
	}
}

type tensorflowProvider struct {
	version string
}

func (p *tensorflowProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tensorflow"
	resp.Version = p.version
}

func (p *tensorflowProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *tensorflowProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *tensorflowProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{
		NewLiteralDataSource,
		NewRefDataSource,
		NewAttrDataSource,
		NewCallDataSource,
		NewAssignDataSource,
		NewReturnDataSource,
		NewWithDataSource,
		NewFunctionDataSource,
		NewRawOpDataSource,
		NewProgramDataSource,
	}

	return append(dataSources, GeneratedDataSources()...)
}

func (p *tensorflowProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
