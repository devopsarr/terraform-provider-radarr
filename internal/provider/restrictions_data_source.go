package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const restrictionsDataSourceName = "restrictions"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RestrictionsDataSource{}

func NewRestrictionsDataSource() datasource.DataSource {
	return &RestrictionsDataSource{}
}

// RestrictionsDataSource defines the restrictions implementation.
type RestrictionsDataSource struct {
	client *radarr.Radarr
}

// Restrictions describes the restrictions data model.
type Restrictions struct {
	Restrictions types.Set    `tfsdk:"restrictions"`
	ID           types.String `tfsdk:"id"`
}

func (d *RestrictionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + restrictionsDataSourceName
}

func (d *RestrictionsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->List all available [Restrictions](../resources/restriction).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"restrictions": {
				MarkdownDescription: "Restriction list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"required": {
						MarkdownDescription: "Required.",
						Computed:            true,
						Type:                types.StringType,
					},
					"ignored": {
						MarkdownDescription: "Ignored.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tags": {
						MarkdownDescription: "List of associated tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
					"id": {
						MarkdownDescription: "Restriction ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
				}),
			},
		},
	}, nil
}

func (d *RestrictionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RestrictionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Restrictions

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get restrictions current value
	response, err := d.client.GetRestrictionsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", restrictionsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+restrictionsDataSourceName)
	// Map response body to resource schema attribute
	restrictions := make([]Restriction, len(response))
	for i, p := range response {
		restrictions[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, restrictions, data.Restrictions.Type(context.Background()), &data.Restrictions)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}