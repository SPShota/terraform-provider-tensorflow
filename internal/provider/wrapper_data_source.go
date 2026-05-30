package provider

import (
	"context"

	"github.com/SPShota/terraform-provider-tensorflow/internal/python"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &wrapperDataSource{}

type WrapperDataSourceSpec struct {
	TypeNameSuffix string
	Function       string
	DocURL         string
}

func NewWrapperDataSource(spec WrapperDataSourceSpec) datasource.DataSource {
	return &wrapperDataSource{spec: spec}
}

type wrapperDataSource struct {
	spec WrapperDataSourceSpec
}

func (d *wrapperDataSource) WrapperSpec() WrapperDataSourceSpec {
	return d.spec
}

type wrapperDataSourceModel struct {
	Args       types.List   `tfsdk:"args"`
	Kwargs     types.Map    `tfsdk:"kwargs"`
	Expression types.String `tfsdk:"expression"`
	Statement  types.String `tfsdk:"statement"`
}

func (d *wrapperDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.spec.TypeNameSuffix
}

func (d *wrapperDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	description := "Generated wrapper for `" + d.spec.Function + "`."
	if d.spec.DocURL != "" {
		description += "\n\nDocumentation: " + d.spec.DocURL
	}

	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: description,
		Attributes: map[string]datasourceschema.Attribute{
			"args": datasourceschema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Positional argument expressions.",
			},
			"kwargs": datasourceschema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Keyword argument expressions.",
			},
			"expression": expressionAttribute(),
			"statement":  statementAttribute(),
		},
	}
}

func (d *wrapperDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data wrapperDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := expressionList(ctx, data.Args)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	kwargs, diags := keywordArguments(ctx, data.Kwargs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	function, err := python.RawExpression(d.spec.Function)
	if err != nil {
		resp.Diagnostics.AddError("Invalid wrapper function", err.Error())
		return
	}

	expr, err := python.Call(function, args, kwargs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid wrapper call", err.Error())
		return
	}

	statement, err := python.ExpressionStatement(expr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid wrapper statement", err.Error())
		return
	}

	data.Expression = types.StringValue(expr.Code())
	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
