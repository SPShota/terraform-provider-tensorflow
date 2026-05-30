package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/SPShota/terraform-provider-tensorflow/internal/python"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &literalDataSource{}
	_ datasource.DataSource = &refDataSource{}
	_ datasource.DataSource = &attrDataSource{}
	_ datasource.DataSource = &callDataSource{}
	_ datasource.DataSource = &assignDataSource{}
)

func NewLiteralDataSource() datasource.DataSource {
	return &literalDataSource{}
}

func NewRefDataSource() datasource.DataSource {
	return &refDataSource{}
}

func NewAttrDataSource() datasource.DataSource {
	return &attrDataSource{}
}

func NewCallDataSource() datasource.DataSource {
	return &callDataSource{}
}

func NewAssignDataSource() datasource.DataSource {
	return &assignDataSource{}
}

type expressionResult struct {
	Expression types.String `tfsdk:"expression"`
}

type statementResult struct {
	Statement types.String `tfsdk:"statement"`
}

type literalDataSource struct{}

type literalDataSourceModel struct {
	ValueJSON  types.String `tfsdk:"value_json"`
	Expression types.String `tfsdk:"expression"`
}

func (d *literalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_literal"
}

func (d *literalDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Converts a JSON value into a Python literal expression.",
		Attributes: map[string]datasourceschema.Attribute{
			"value_json": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "JSON value to convert to Python literal syntax. Use Terraform `jsonencode(...)` for HCL values.",
			},
			"expression": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Python literal expression.",
			},
		},
	}
}

func (d *literalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data literalDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	expr, err := literalExpressionFromJSON(data.ValueJSON.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid literal", err.Error())
		return
	}

	data.Expression = types.StringValue(expr.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type refDataSource struct{}

type refDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Expression types.String `tfsdk:"expression"`
}

func (d *refDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ref"
}

func (d *refDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python reference expression.",
		Attributes: map[string]datasourceschema.Attribute{
			"name": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Python identifier or dotted identifier, such as `tf.float32`.",
			},
			"expression": expressionAttribute(),
		},
	}
}

func (d *refDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data refDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	expr, err := python.Reference(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid reference", err.Error())
		return
	}

	data.Expression = types.StringValue(expr.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type attrDataSource struct{}

type attrDataSourceModel struct {
	Receiver   types.String `tfsdk:"receiver"`
	Name       types.String `tfsdk:"name"`
	Expression types.String `tfsdk:"expression"`
}

func (d *attrDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attr"
}

func (d *attrDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python attribute access expression.",
		Attributes: map[string]datasourceschema.Attribute{
			"receiver": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Receiver Python expression.",
			},
			"name": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Attribute name.",
			},
			"expression": expressionAttribute(),
		},
	}
}

func (d *attrDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data attrDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	receiver, err := python.RawExpression(data.Receiver.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid receiver", err.Error())
		return
	}

	expr, err := python.Attribute(receiver, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid attribute", err.Error())
		return
	}

	data.Expression = types.StringValue(expr.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type callDataSource struct{}

type callDataSourceModel struct {
	Function   types.String `tfsdk:"function"`
	Args       types.List   `tfsdk:"args"`
	Kwargs     types.Map    `tfsdk:"kwargs"`
	Expression types.String `tfsdk:"expression"`
	Statement  types.String `tfsdk:"statement"`
}

func (d *callDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_call"
}

func (d *callDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python function call expression.",
		Attributes: map[string]datasourceschema.Attribute{
			"function": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Python callable expression, such as `tf.constant`.",
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

func (d *callDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data callDataSourceModel
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

	function, err := python.RawExpression(data.Function.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid call function", err.Error())
		return
	}

	expr, err := python.Call(function, args, kwargs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid call", err.Error())
		return
	}

	statement, err := python.ExpressionStatement(expr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid call statement", err.Error())
		return
	}

	data.Expression = types.StringValue(expr.Code())
	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type assignDataSource struct{}

type assignDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Value      types.String `tfsdk:"value"`
	Expression types.String `tfsdk:"expression"`
	Statement  types.String `tfsdk:"statement"`
}

func (d *assignDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assign"
}

func (d *assignDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Creates a Python assignment statement.",
		Attributes: map[string]datasourceschema.Attribute{
			"name": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Target variable name.",
			},
			"value": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Value expression.",
			},
			"expression": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reference expression for the assigned variable.",
			},
			"statement": statementAttribute(),
		},
	}
}

func (d *assignDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data assignDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := python.RawExpression(data.Value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid assignment value", err.Error())
		return
	}

	statement, err := python.Assign(data.Name.ValueString(), value)
	if err != nil {
		resp.Diagnostics.AddError("Invalid assignment", err.Error())
		return
	}

	ref, err := python.Reference(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid assignment reference", err.Error())
		return
	}

	data.Expression = types.StringValue(ref.Code())
	data.Statement = types.StringValue(statement.Code())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func literalExpressionFromJSON(valueJSON string) (python.Expression, error) {
	decoder := json.NewDecoder(strings.NewReader(valueJSON))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return python.Expression{}, fmt.Errorf("decode JSON value: %w", err)
	}

	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return python.Expression{}, fmt.Errorf("value_json must contain exactly one JSON value")
		}
		return python.Expression{}, fmt.Errorf("value_json must contain exactly one JSON value")
	}

	return python.Literal(value)
}

func expressionList(ctx context.Context, values types.List) ([]python.Expression, diag.Diagnostics) {
	var diags diag.Diagnostics
	if values.IsNull() || values.IsUnknown() {
		return nil, diags
	}

	var rawValues []string
	diags.Append(values.ElementsAs(ctx, &rawValues, false)...)
	if diags.HasError() {
		return nil, diags
	}

	expressions := make([]python.Expression, 0, len(rawValues))
	for _, rawValue := range rawValues {
		expr, err := python.RawExpression(rawValue)
		if err != nil {
			diags.AddError("Invalid argument expression", err.Error())
			return nil, diags
		}
		expressions = append(expressions, expr)
	}

	return expressions, diags
}

func keywordArguments(ctx context.Context, values types.Map) ([]python.KeywordArgument, diag.Diagnostics) {
	var diags diag.Diagnostics
	if values.IsNull() || values.IsUnknown() {
		return nil, diags
	}

	var rawValues map[string]string
	diags.Append(values.ElementsAs(ctx, &rawValues, false)...)
	if diags.HasError() {
		return nil, diags
	}

	names := make([]string, 0, len(rawValues))
	for name := range rawValues {
		names = append(names, name)
	}
	sort.Strings(names)

	kwargs := make([]python.KeywordArgument, 0, len(names))
	for _, name := range names {
		expr, err := python.RawExpression(rawValues[name])
		if err != nil {
			diags.AddError("Invalid keyword argument expression", err.Error())
			return nil, diags
		}
		kwargs = append(kwargs, python.KeywordArgument{Name: name, Value: expr})
	}

	return kwargs, diags
}

func expressionAttribute() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Generated Python expression.",
	}
}

func statementAttribute() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Generated Python statement.",
	}
}
