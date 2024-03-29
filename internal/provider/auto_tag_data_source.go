package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const autoTagDataSourceName = "auto_tag"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagDataSource{}

func NewAutoTagDataSource() datasource.DataSource {
	return &AutoTagDataSource{}
}

// AutoTagDataSource defines the auto_tag implementation.
type AutoTagDataSource struct {
	client *radarr.APIClient
	auth   context.Context
}

func (d *AutoTagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagDataSourceName
}

func (d *AutoTagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Tags -->\n\nSingle [Auto Tag](../resources/auto_tag).",
		Attributes: map[string]schema.Attribute{
			"remove_tags_automatically": schema.BoolAttribute{
				MarkdownDescription: "Remove tags automatically flag.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Auto Tag name.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Auto Tag ID.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"specifications": schema.SetNestedAttribute{
				MarkdownDescription: "Specifications.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"negate": schema.BoolAttribute{
							MarkdownDescription: "Negate flag.",
							Computed:            true,
						},
						"required": schema.BoolAttribute{
							MarkdownDescription: "Computed flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Specification name.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Implementation.",
							Computed:            true,
						},
						// Field values
						"value": schema.StringAttribute{
							MarkdownDescription: "Value.",
							Computed:            true,
						},
						"min": schema.Int64Attribute{
							MarkdownDescription: "Min.",
							Computed:            true,
						},
						"max": schema.Int64Attribute{
							MarkdownDescription: "Max.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AutoTagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *AutoTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AutoTag

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get autoTag current value
	response, _, err := d.client.AutoTaggingAPI.ListAutoTagging(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, autoTagDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+autoTagDataSourceName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (c *AutoTag) find(ctx context.Context, name string, autoTags []radarr.AutoTaggingResource, diags *diag.Diagnostics) {
	for _, i := range autoTags {
		if i.GetName() == name {
			c.write(ctx, &i, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(autoTagDataSourceName, "name", name))
}
