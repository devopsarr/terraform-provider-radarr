package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const customFormatDataSourceName = "custom_format"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatDataSource{}

func NewCustomFormatDataSource() datasource.DataSource {
	return &CustomFormatDataSource{}
}

// CustomFormatDataSource defines the custom_format implementation.
type CustomFormatDataSource struct {
	client *radarr.Radarr
}

func (d *CustomFormatDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatDataSourceName
}

func (d *CustomFormatDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Download Client](../resources/custom_format).",
		Attributes: map[string]tfsdk.Attribute{
			"include_custom_format_when_renaming": {
				MarkdownDescription: "Include custom format when renaming flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"name": {
				MarkdownDescription: "Custom Format name.",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Custom Format ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"specifications": {
				MarkdownDescription: "Specifications.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"negate": {
						MarkdownDescription: "Negate flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"required": {
						MarkdownDescription: "Computed flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"name": {
						MarkdownDescription: "Specification name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"implementation": {
						MarkdownDescription: "Implementation.",
						Computed:            true,
						Type:                types.StringType,
					},
					// Field values
					"value": {
						MarkdownDescription: "Value.",
						Computed:            true,
						Type:                types.StringType,
					},
					"min": {
						MarkdownDescription: "Min.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"max": {
						MarkdownDescription: "Max.",
						Computed:            true,
						Type:                types.Int64Type,
					},
				}),
			},
		},
	}, nil
}

func (d *CustomFormatDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomFormatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormat

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get customFormat current value
	response, err := d.client.GetCustomFormatsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatDataSourceName, err))

		return
	}

	customFormat, err := findCustomFormat(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", customFormatDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatDataSourceName)
	data.write(ctx, customFormat)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findCustomFormat(name string, customFormats []*radarr.CustomFormatOutput) (*radarr.CustomFormatOutput, error) {
	for _, i := range customFormats {
		if i.Name == name {
			return i, nil
		}
	}

	return nil, tools.ErrDataNotFoundError(customFormatDataSourceName, "name", name)
}
