package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/hashstructure/v2"
)

const (
	customFormatConditionEditionDataSourceName = "custom_format_condition_edition"
	customFormatConditionEditionImplementation = "EditionSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatConditionEditionDataSource{}

func NewCustomFormatConditionEditionDataSource() datasource.DataSource {
	return &CustomFormatConditionEditionDataSource{}
}

// CustomFormatConditionEditionDataSource defines the custom format condition edition implementation.
type CustomFormatConditionEditionDataSource struct {
	client *radarr.APIClient
}

func (d *CustomFormatConditionEditionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatConditionEditionDataSourceName
}

func (d *CustomFormatConditionEditionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles --> Custom Format Condition Edition data source.\nFor more information refer to [Custom Format Conditions](https://wiki.servarr.com/radarr/settings#conditions).",
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate flag.",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Computed flag.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Implementation.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specification name.",
				Required:            true,
			},
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom format condition edition ID.",
				Computed:            true,
			},
			// Field values
			"value": schema.StringAttribute{
				MarkdownDescription: "Edition RegEx.",
				Required:            true,
			},
		},
	}
}

func (d *CustomFormatConditionEditionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *CustomFormatConditionEditionDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormatConditionValue

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, customFormatConditionEditionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatConditionEditionDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), customFormatConditionEditionImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
