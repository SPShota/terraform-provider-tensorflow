package provider

import (
	"context"

	"github.com/SPShota/terraform-provider-tensorflow/internal/python"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &returnDataSource{}
	_ datasource.DataSource = &withDataSource{}
	_ datasource.DataSource = &functionDataSource{}
)

func NewReturnDataSource() datasource.DataSource {
	return &returnDataSource{}
}

func NewWithDataSource() datasource.DataSource {
	return &withDataSource{}
}

func NewFunctionDataSource() datasource.DataSource {
	return &functionDataSource{}
}

type returnDataSource struct{}

type returnDataSourceModel struct {
	Value     types.String `tfsdk:"value"`
	Statement types.String `tfsdk:"statement"`
}

func (d *returnDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_return"
}

func (d *returnDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python return statement.",
		Attributes: map[string]datasourceschema.Attribute{
			"value": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Return value expression.",
			},
			"statement": statementAttribute(),
		},
	}
}

func (d *returnDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data returnDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := python.RawExpression(data.Value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid return value", err.Error())
		return
	}

	statement, err := python.Return(value)
	if err != nil {
		resp.Diagnostics.AddError("Invalid return statement", err.Error())
		return
	}

	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type withDataSource struct{}

type withDataSourceModel struct {
	Context    types.String `tfsdk:"context"`
	Alias      types.String `tfsdk:"alias"`
	Statements types.List   `tfsdk:"statements"`
	Statement  types.String `tfsdk:"statement"`
}

func (d *withDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_with"
}

func (d *withDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python `with` block statement.",
		Attributes: map[string]datasourceschema.Attribute{
			"context": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Context manager expression.",
			},
			"alias": datasourceschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional alias name for `as ...`.",
			},
			"statements": datasourceschema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Statements in the with block body. Empty bodies render `pass`.",
			},
			"statement": statementAttribute(),
		},
	}
}

func (d *withDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data withDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contextExpr, err := python.RawExpression(data.Context.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid with context", err.Error())
		return
	}

	statements, ok := statementsFromList(ctx, data.Statements, resp)
	if !ok {
		return
	}

	alias := ""
	if !data.Alias.IsNull() && !data.Alias.IsUnknown() {
		alias = data.Alias.ValueString()
	}

	statement, err := python.With(contextExpr, alias, statements)
	if err != nil {
		resp.Diagnostics.AddError("Invalid with statement", err.Error())
		return
	}

	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type functionDataSource struct{}

type functionDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Args       types.List   `tfsdk:"args"`
	Decorators types.List   `tfsdk:"decorators"`
	Statements types.List   `tfsdk:"statements"`
	Expression types.String `tfsdk:"expression"`
	Statement  types.String `tfsdk:"statement"`
}

func (d *functionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}

func (d *functionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python function definition statement.",
		Attributes: map[string]datasourceschema.Attribute{
			"name": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Function name.",
			},
			"args": datasourceschema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Function argument names.",
			},
			"decorators": datasourceschema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Decorator expressions without the leading `@`.",
			},
			"statements": datasourceschema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Statements in the function body. Empty bodies render `pass`.",
			},
			"expression": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reference expression for the defined function.",
			},
			"statement": statementAttribute(),
		},
	}
}

func (d *functionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data functionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, ok := stringList(ctx, data.Args, resp)
	if !ok {
		return
	}

	rawDecorators, ok := stringList(ctx, data.Decorators, resp)
	if !ok {
		return
	}
	decorators := make([]python.Expression, 0, len(rawDecorators))
	for _, rawDecorator := range rawDecorators {
		decorator, err := python.RawExpression(rawDecorator)
		if err != nil {
			resp.Diagnostics.AddError("Invalid function decorator", err.Error())
			return
		}
		decorators = append(decorators, decorator)
	}

	statements, ok := statementsFromList(ctx, data.Statements, resp)
	if !ok {
		return
	}

	statement, err := python.FunctionDef(data.Name.ValueString(), args, decorators, statements)
	if err != nil {
		resp.Diagnostics.AddError("Invalid function statement", err.Error())
		return
	}

	ref, err := python.Reference(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid function reference", err.Error())
		return
	}

	data.Expression = types.StringValue(ref.Code())
	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func statementsFromList(ctx context.Context, values types.List, resp *datasource.ReadResponse) ([]python.Statement, bool) {
	rawStatements, ok := stringList(ctx, values, resp)
	if !ok {
		return nil, false
	}

	statements := make([]python.Statement, 0, len(rawStatements))
	for _, rawStatement := range rawStatements {
		statement, err := python.RawStatement(rawStatement)
		if err != nil {
			resp.Diagnostics.AddError("Invalid block statement", err.Error())
			return nil, false
		}
		statements = append(statements, statement)
	}

	return statements, true
}

func stringList(ctx context.Context, values types.List, resp *datasource.ReadResponse) ([]string, bool) {
	if values.IsNull() || values.IsUnknown() {
		return nil, true
	}

	var rawValues []string
	resp.Diagnostics.Append(values.ElementsAs(ctx, &rawValues, false)...)
	if resp.Diagnostics.HasError() {
		return nil, false
	}

	return rawValues, true
}
