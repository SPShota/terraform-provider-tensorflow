package provider

import (
	"context"

	"github.com/SPShota/terraform-provider-tensorflow/internal/python"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &rawOpDataSource{}

func NewRawOpDataSource() datasource.DataSource {
	return &rawOpDataSource{}
}

type rawOpDataSource struct{}

type rawOpDataSourceModel struct {
	Op         types.String `tfsdk:"op"`
	Args       types.List   `tfsdk:"args"`
	Kwargs     types.Map    `tfsdk:"kwargs"`
	Expression types.String `tfsdk:"expression"`
	Statement  types.String `tfsdk:"statement"`
}

func (d *rawOpDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_raw_op"
}

func (d *rawOpDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python call to `tf.raw_ops.<op>`. Use this for raw TensorFlow ops that do not have dedicated generated wrappers.",
		Attributes: map[string]datasourceschema.Attribute{
			"op": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Raw op name under `tf.raw_ops`, such as `AddV2` or `Identity`.",
			},
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

func (d *rawOpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data rawOpDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := python.ValidateIdentifier(data.Op.ValueString()); err != nil {
		resp.Diagnostics.AddError("Invalid raw op", err.Error())
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

	function, err := python.RawExpression("tf.raw_ops." + data.Op.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid raw op function", err.Error())
		return
	}

	expr, err := python.Call(function, args, kwargs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid raw op call", err.Error())
		return
	}

	statement, err := python.ExpressionStatement(expr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid raw op statement", err.Error())
		return
	}

	data.Expression = types.StringValue(expr.Code())
	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
