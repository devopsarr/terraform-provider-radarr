package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const restrictionsDataSourceName = "restrictions"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RestrictionsDataSource{}

func NewRestrictionsDataSource() datasource.DataSource {
	return &RestrictionsDataSource{}
}

// RestrictionsDataSource defines the restrictions implementation.
type RestrictionsDataSource struct {
	client *radarr.APIClient
}

// Restrictions describes the restrictions data model.
type Restrictions struct {
	Restrictions types.Set    `tfsdk:"restrictions"`
	ID           types.String `tfsdk:"id"`
}

func (d *RestrictionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + restrictionsDataSourceName
}

func (d *RestrictionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->List all available [Restrictions](../resources/restriction).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"restrictions": schema.SetNestedAttribute{
				MarkdownDescription: "Restriction list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"required": schema.StringAttribute{
							MarkdownDescription: "Required.",
							Computed:            true,
						},
						"ignored": schema.StringAttribute{
							MarkdownDescription: "Ignored.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Restriction ID.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RestrictionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *RestrictionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Restrictions

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get restrictions current value
	response, _, err := d.client.RestrictionApi.ListRestriction(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, restrictionsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+restrictionsDataSourceName)
	// Map response body to resource schema attribute
	restrictions := make([]Restriction, len(response))
	for i, p := range response {
		restrictions[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, restrictions, data.Restrictions.Type(ctx), &data.Restrictions)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
