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
	customFormatConditionQualityModifierDataSourceName = "custom_format_condition_quality_modifier"
	customFormatConditionQualityModifierImplementation = "QualityModifierSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatConditionQualityModifierDataSource{}

func NewCustomFormatConditionQualityModifierDataSource() datasource.DataSource {
	return &CustomFormatConditionQualityModifierDataSource{}
}

// CustomFormatConditionQualityModifierDataSource defines the custom_format_condition_quality_modifier implementation.
type CustomFormatConditionQualityModifierDataSource struct {
	client *radarr.APIClient
}

func (d *CustomFormatConditionQualityModifierDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatConditionQualityModifierDataSourceName
}

func (d *CustomFormatConditionQualityModifierDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles --> Custom Format Condition Quality Modifier data source.\nFor more information refer to [Custom Format Conditions](https://wiki.servarr.com/radarr/settings#conditions).",
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate modifier.",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Computed modifier.",
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
				MarkdownDescription: "Custom format condition quality modifier ID.",
				Computed:            true,
			},
			// Field values
			"value": schema.StringAttribute{
				MarkdownDescription: "Quality modifier ID. `0` NONE, `1` REGIONAL, `2` SCREENER, `3` RAWHD, `4` BRDISK, `5` REMUX.",
				Required:            true,
			},
		},
	}
}

func (d *CustomFormatConditionQualityModifierDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *CustomFormatConditionQualityModifierDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormatConditionValue

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, customFormatConditionQualityModifierDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatConditionQualityModifierDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), customFormatConditionQualityModifierImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
