package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootFoldersDataSource{}

func NewRootFoldersDataSource() datasource.DataSource {
	return &RootFoldersDataSource{}
}

// RootFoldersDataSource defines the root folders implementation.
type RootFoldersDataSource struct {
	client *radarr.Radarr
}

// RootFolders describes the root folders data model.
type RootFolders struct {
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID          types.String `tfsdk:"id"`
	RootFolders types.Set    `tfsdk:"root_folders"`
}

func (d *RootFoldersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_folders"
}

func (d *RootFoldersDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->List all available [Root Folders](../resources/root_folder).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"root_folders": {
				MarkdownDescription: "Root Folder list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"path": {
						MarkdownDescription: "Root Folder absolute path.",
						Computed:            true,
						Type:                types.StringType,
					},
					"accessible": {
						MarkdownDescription: "Access flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"id": {
						MarkdownDescription: "Root Folder ID.",
						Computed:            true,
						Type:                types.Int64Type,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"unmapped_folders": {
						MarkdownDescription: "List of folders with no associated series.",
						Computed:            true,
						Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
							"path": {
								MarkdownDescription: "Path of unmapped folder.",
								Computed:            true,
								Type:                types.StringType,
							},
							"name": {
								MarkdownDescription: "Name of unmapped folder.",
								Computed:            true,
								Type:                types.StringType,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (d *RootFoldersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RootFoldersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RootFolders

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get rootfolders current value
	response, err := d.client.GetRootFoldersContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read rootfolders, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read root folders")
	// Map response body to resource schema attribute
	rootFolders := *writeRootFolders(ctx, response)
	tfsdk.ValueFrom(ctx, rootFolders, data.RootFolders.Type(context.Background()), &data.RootFolders)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeRootFolders(ctx context.Context, folders []*radarr.RootFolder) *[]RootFolder {
	output := make([]RootFolder, len(folders))
	for i, f := range folders {
		output[i] = *writeRootFolder(ctx, f)
	}

	return &output
}