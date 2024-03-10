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
	autoTagConditionYearDataSourceName = "auto_tag_condition_year"
	autoTagConditionYearImplementation = "YearSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagConditionYearDataSource{}

func NewAutoTagConditionYearDataSource() datasource.DataSource {
	return &AutoTagConditionYearDataSource{}
}

// AutoTagConditionYearDataSource defines the auto_tag_condition_series type implementation.
type AutoTagConditionYearDataSource struct {
	client *radarr.APIClient
	auth   context.Context
}

func (d *AutoTagConditionYearDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagConditionYearDataSourceName
}

func (d *AutoTagConditionYearDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Tags -->\n Auto Tag Condition Series Type data source.\nFor more intagion refer to [Auto Tag Conditions](https://wiki.servarr.com/radarr/settings#conditions).",
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
				MarkdownDescription: "Auto tag condition series type ID.",
				Computed:            true,
			},
			// Field values
			"min": schema.Int64Attribute{
				MarkdownDescription: "Minimum Year.",
				Required:            true,
			},
			"max": schema.Int64Attribute{
				MarkdownDescription: "Maximum Year.",
				Required:            true,
			},
		},
	}
}

func (d *AutoTagConditionYearDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *AutoTagConditionYearDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AutoTagCondition

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, autoTagConditionYearDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagConditionYearDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), autoTagConditionYearImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
