package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const autoTagsDataSourceName = "auto_tags"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagsDataSource{}

func NewAutoTagsDataSource() datasource.DataSource {
	return &AutoTagsDataSource{}
}

// AutoTagsDataSource defines the download clients implementation.
type AutoTagsDataSource struct {
	client *radarr.APIClient
	auth   context.Context
}

// AutoTags describes the download clients data model.
type AutoTags struct {
	AutoTags types.Set    `tfsdk:"auto_tags"`
	ID       types.String `tfsdk:"id"`
}

func (d *AutoTagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagsDataSourceName
}

func (d *AutoTagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Tags -->\nList all available [Auto Tags](../resources/auto_tag).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"auto_tags": schema.SetNestedAttribute{
				MarkdownDescription: "Auto Tag list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"remove_tags_automatically": schema.BoolAttribute{
							MarkdownDescription: "Remove tags automatically flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Auto Tag name.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *AutoTagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *AutoTagsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get download clients current value
	response, _, err := d.client.AutoTaggingAPI.ListAutoTagging(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, autoTagsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagsDataSourceName)
	// Map response body to resource schema attribute
	autoTags := make([]AutoTag, len(response))
	for i, a := range response {
		autoTags[i].write(ctx, &a, &resp.Diagnostics)
	}

	autoList, diags := types.SetValueFrom(ctx, AutoTag{}.getType(), autoTags)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, AutoTags{AutoTags: autoList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
