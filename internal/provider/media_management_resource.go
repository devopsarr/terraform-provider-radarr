package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MediaManagementResource{}
var _ resource.ResourceWithImportState = &MediaManagementResource{}

func NewMediaManagementResource() resource.Resource {
	return &MediaManagementResource{}
}

// MediaManagementResource defines the media management implementation.
type MediaManagementResource struct {
	client *radarr.Radarr
}

// MediaManagement describes the media management data model.
type MediaManagement struct {
	ChmodFolder                             types.String `tfsdk:"chmod_folder"`
	RescanAfterRefresh                      types.String `tfsdk:"rescan_after_refresh"`
	RecycleBin                              types.String `tfsdk:"recycle_bin"`
	FileDate                                types.String `tfsdk:"file_date"`
	ExtraFileExtensions                     types.String `tfsdk:"extra_file_extensions"`
	DownloadPropersAndRepacks               types.String `tfsdk:"download_propers_and_repacks"`
	ChownGroup                              types.String `tfsdk:"chown_group"`
	ID                                      types.Int64  `tfsdk:"id"`
	MinimumFreeSpaceWhenImporting           types.Int64  `tfsdk:"minimum_free_space_when_importing"`
	RecycleBinCleanupDays                   types.Int64  `tfsdk:"recycle_bin_cleanup_days"`
	SetPermissionsLinux                     types.Bool   `tfsdk:"set_permissions_linux"`
	SkipFreeSpaceCheckWhenImporting         types.Bool   `tfsdk:"skip_free_space_check_when_importing"`
	AutoRenameFolders                       types.Bool   `tfsdk:"auto_rename_folders"`
	PathsDefaultStatic                      types.Bool   `tfsdk:"paths_default_static"`
	ImportExtraFiles                        types.Bool   `tfsdk:"import_extra_files"`
	EnableMediaInfo                         types.Bool   `tfsdk:"enable_media_info"`
	DeleteEmptyFolders                      types.Bool   `tfsdk:"delete_empty_folders"`
	CreateEmptyMovieFolders                 types.Bool   `tfsdk:"create_empty_movie_folders"`
	CopyUsingHardlinks                      types.Bool   `tfsdk:"copy_using_hardlinks"`
	AutoUnmonitorPreviouslyDownloadedMovies types.Bool   `tfsdk:"auto_unmonitor_previously_downloaded_movies"`
}

func (r *MediaManagementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_media_management"
}

func (r *MediaManagementResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Media Management resource.\nFor more information refer to [Naming](https://wiki.servarr.com/radarr/settings#file-management) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Media Management ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"auto_rename_folders": {
				MarkdownDescription: "Auto rename folders.",
				Required:            true,
				Type:                types.BoolType,
			},
			"auto_unmonitor_previously_downloaded_movies": {
				MarkdownDescription: "Auto unmonitor previously downloaded movies.",
				Required:            true,
				Type:                types.BoolType,
			},
			"copy_using_hardlinks": {
				MarkdownDescription: "Use hardlinks instead of copy.",
				Required:            true,
				Type:                types.BoolType,
			},
			"create_empty_movie_folders": {
				MarkdownDescription: "Create empty movies directories.",
				Required:            true,
				Type:                types.BoolType,
			},
			"delete_empty_folders": {
				MarkdownDescription: "Delete empty movies directories.",
				Required:            true,
				Type:                types.BoolType,
			},
			"enable_media_info": {
				MarkdownDescription: "Scan files details.",
				Required:            true,
				Type:                types.BoolType,
			},
			"import_extra_files": {
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Required:            true,
				Type:                types.BoolType,
			},
			"paths_default_static": {
				MarkdownDescription: "Path default static.",
				Required:            true,
				Type:                types.BoolType,
			},
			"set_permissions_linux": {
				MarkdownDescription: "Set permission for imported files.",
				Required:            true,
				Type:                types.BoolType,
			},
			"skip_free_space_check_when_importing": {
				MarkdownDescription: "Skip free space check before importing.",
				Required:            true,
				Type:                types.BoolType,
			},
			"minimum_free_space_when_importing": {
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"recycle_bin_cleanup_days": {
				MarkdownDescription: "Recyle bin days of retention.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"chmod_folder": {
				MarkdownDescription: "Permission in linux format.",
				Required:            true,
				Type:                types.StringType,
			},
			"chown_group": {
				MarkdownDescription: "Group used for permission.",
				Required:            true,
				Type:                types.StringType,
			},
			"download_propers_and_repacks": {
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					tools.StringMatch([]string{"preferAndUpgrade", "doNotUpgrade", "doNotPrefer"}),
				},
			},
			"extra_file_extensions": {
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Required:            true,
				Type:                types.StringType,
			},
			"file_date": {
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'cinemas, and 'release'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					tools.StringMatch([]string{"none", "cinemas", "release"}),
				},
			},
			"recycle_bin": {
				MarkdownDescription: "Recycle bin absolute path.",
				Required:            true,
				Type:                types.StringType,
			},
			"rescan_after_refresh": {
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					tools.StringMatch([]string{"always", "afterManual", "never"}),
				},
			},
		},
	}, nil
}

func (r *MediaManagementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MediaManagementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan MediaManagement

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readMediaManagement(&plan)
	data.ID = 1

	// Create new MediaManagement
	response, err := r.client.UpdateMediaManagementContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create mediamanagement, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created media_management: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeMediaManagement(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *MediaManagementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state MediaManagement

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get mediamanagement current value
	response, err := r.client.GetMediaManagementContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read mediamanagements, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read media_management: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeMediaManagement(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *MediaManagementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan MediaManagement

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readMediaManagement(&plan)

	// Update MediaManagement
	response, err := r.client.UpdateMediaManagementContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update mediamanagement, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated media_management: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeMediaManagement(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *MediaManagementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Mediamanagement cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled media_management: 1")
	resp.State.RemoveResource(ctx)
}

func (r *MediaManagementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported media_management: 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func writeMediaManagement(mediaMgt *radarr.MediaManagement) *MediaManagement {
	return &MediaManagement{
		AutoRenameFolders:                       types.BoolValue(mediaMgt.AutoRenameFolders),
		AutoUnmonitorPreviouslyDownloadedMovies: types.BoolValue(mediaMgt.AutoUnmonitorPreviouslyDownloadedMovies),
		CopyUsingHardlinks:                      types.BoolValue(mediaMgt.CopyUsingHardlinks),
		CreateEmptyMovieFolders:                 types.BoolValue(mediaMgt.CreateEmptyMovieFolders),
		DeleteEmptyFolders:                      types.BoolValue(mediaMgt.DeleteEmptyFolders),
		EnableMediaInfo:                         types.BoolValue(mediaMgt.EnableMediaInfo),
		ImportExtraFiles:                        types.BoolValue(mediaMgt.ImportExtraFiles),
		PathsDefaultStatic:                      types.BoolValue(mediaMgt.PathsDefaultStatic),
		SetPermissionsLinux:                     types.BoolValue(mediaMgt.SetPermissionsLinux),
		SkipFreeSpaceCheckWhenImporting:         types.BoolValue(mediaMgt.SkipFreeSpaceCheckWhenImporting),
		ID:                                      types.Int64Value(mediaMgt.ID),
		MinimumFreeSpaceWhenImporting:           types.Int64Value(mediaMgt.MinimumFreeSpaceWhenImporting),
		RecycleBinCleanupDays:                   types.Int64Value(mediaMgt.RecycleBinCleanupDays),
		ChmodFolder:                             types.StringValue(mediaMgt.ChmodFolder),
		ChownGroup:                              types.StringValue(mediaMgt.ChownGroup),
		DownloadPropersAndRepacks:               types.StringValue(mediaMgt.DownloadPropersAndRepacks),
		ExtraFileExtensions:                     types.StringValue(mediaMgt.ExtraFileExtensions),
		FileDate:                                types.StringValue(mediaMgt.FileDate),
		RecycleBin:                              types.StringValue(mediaMgt.RecycleBin),
		RescanAfterRefresh:                      types.StringValue(mediaMgt.RescanAfterRefresh),
	}
}

func readMediaManagement(mediaMgt *MediaManagement) *radarr.MediaManagement {
	return &radarr.MediaManagement{
		AutoRenameFolders:                       mediaMgt.AutoRenameFolders.ValueBool(),
		AutoUnmonitorPreviouslyDownloadedMovies: mediaMgt.AutoUnmonitorPreviouslyDownloadedMovies.ValueBool(),
		CopyUsingHardlinks:                      mediaMgt.CopyUsingHardlinks.ValueBool(),
		CreateEmptyMovieFolders:                 mediaMgt.CreateEmptyMovieFolders.ValueBool(),
		DeleteEmptyFolders:                      mediaMgt.DeleteEmptyFolders.ValueBool(),
		EnableMediaInfo:                         mediaMgt.EnableMediaInfo.ValueBool(),
		ImportExtraFiles:                        mediaMgt.ImportExtraFiles.ValueBool(),
		PathsDefaultStatic:                      mediaMgt.PathsDefaultStatic.ValueBool(),
		SetPermissionsLinux:                     mediaMgt.SetPermissionsLinux.ValueBool(),
		SkipFreeSpaceCheckWhenImporting:         mediaMgt.SkipFreeSpaceCheckWhenImporting.ValueBool(),
		ID:                                      mediaMgt.ID.ValueInt64(),
		MinimumFreeSpaceWhenImporting:           mediaMgt.MinimumFreeSpaceWhenImporting.ValueInt64(),
		RecycleBinCleanupDays:                   mediaMgt.RecycleBinCleanupDays.ValueInt64(),
		ChmodFolder:                             mediaMgt.ChmodFolder.ValueString(),
		ChownGroup:                              mediaMgt.ChownGroup.ValueString(),
		DownloadPropersAndRepacks:               mediaMgt.DownloadPropersAndRepacks.ValueString(),
		ExtraFileExtensions:                     mediaMgt.ExtraFileExtensions.ValueString(),
		FileDate:                                mediaMgt.FileDate.ValueString(),
		RecycleBin:                              mediaMgt.RecycleBin.ValueString(),
		RescanAfterRefresh:                      mediaMgt.RescanAfterRefresh.ValueString(),
	}
}
