package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MediaManagementDataSource{}

func NewMediaManagementDataSource() datasource.DataSource {
	return &MediaManagementDataSource{}
}

// MediaManagementDataSource defines the media management implementation.
type MediaManagementDataSource struct {
	client *radarr.Radarr
}

func (d *MediaManagementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_media_management"
}

func (d *MediaManagementDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->[Media Management](../resources/media_management).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"auto_rename_folders": {
				MarkdownDescription: "Auto rename folders.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"auto_unmonitor_previously_downloaded_movies": {
				MarkdownDescription: "Auto unmonitor previously downloaded movies.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"copy_using_hardlinks": {
				MarkdownDescription: "Use hardlinks instead of copy.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"create_empty_movie_folders": {
				MarkdownDescription: "Create empty movies directories.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"delete_empty_folders": {
				MarkdownDescription: "Delete empty movies directories.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_media_info": {
				MarkdownDescription: "Scan files details.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"import_extra_files": {
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"paths_default_static": {
				MarkdownDescription: "Path default static.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"set_permissions_linux": {
				MarkdownDescription: "Set permission for imported files.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"skip_free_space_check_when_importing": {
				MarkdownDescription: "Skip free space check before importing.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"minimum_free_space_when_importing": {
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"recycle_bin_cleanup_days": {
				MarkdownDescription: "Recyle bin days of retention.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"chmod_folder": {
				MarkdownDescription: "Permission in linux format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"chown_group": {
				MarkdownDescription: "Group used for permission.",
				Computed:            true,
				Type:                types.StringType,
			},
			"download_propers_and_repacks": {
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Computed:            true,
				Type:                types.StringType,
			},
			"extra_file_extensions": {
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Computed:            true,
				Type:                types.StringType,
			},
			"file_date": {
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'localAirDate, and 'utcAirDate'.",
				Computed:            true,
				Type:                types.StringType,
			},
			"recycle_bin": {
				MarkdownDescription: "Recycle bin absolute path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"rescan_after_refresh": {
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *MediaManagementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MediaManagementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, err := d.client.GetMediaManagementContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read indexer cofig, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read madia_management")

	result := writeMediaManagement(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
