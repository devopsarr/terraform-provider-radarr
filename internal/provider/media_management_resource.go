package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const mediaManagementResourceName = "media_management"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MediaManagementResource{}
	_ resource.ResourceWithImportState = &MediaManagementResource{}
)

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
	resp.TypeName = req.ProviderTypeName + "_" + mediaManagementResourceName
}

func (r *MediaManagementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Media Management resource.\nFor more information refer to [Naming](https://wiki.servarr.com/radarr/settings#file-management) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Media Management ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"auto_rename_folders": schema.BoolAttribute{
				MarkdownDescription: "Auto rename folders.",
				Required:            true,
			},
			"auto_unmonitor_previously_downloaded_movies": schema.BoolAttribute{
				MarkdownDescription: "Auto unmonitor previously downloaded movies.",
				Required:            true,
			},
			"copy_using_hardlinks": schema.BoolAttribute{
				MarkdownDescription: "Use hardlinks instead of copy.",
				Required:            true,
			},
			"create_empty_movie_folders": schema.BoolAttribute{
				MarkdownDescription: "Create empty movies directories.",
				Required:            true,
			},
			"delete_empty_folders": schema.BoolAttribute{
				MarkdownDescription: "Delete empty movies directories.",
				Required:            true,
			},
			"enable_media_info": schema.BoolAttribute{
				MarkdownDescription: "Scan files details.",
				Required:            true,
			},
			"import_extra_files": schema.BoolAttribute{
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Required:            true,
			},
			"paths_default_static": schema.BoolAttribute{
				MarkdownDescription: "Path default static.",
				Required:            true,
			},
			"set_permissions_linux": schema.BoolAttribute{
				MarkdownDescription: "Set permission for imported files.",
				Required:            true,
			},
			"skip_free_space_check_when_importing": schema.BoolAttribute{
				MarkdownDescription: "Skip free space check before importing.",
				Required:            true,
			},
			"minimum_free_space_when_importing": schema.Int64Attribute{
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Required:            true,
			},
			"recycle_bin_cleanup_days": schema.Int64Attribute{
				MarkdownDescription: "Recyle bin days of retention.",
				Required:            true,
			},
			"chmod_folder": schema.StringAttribute{
				MarkdownDescription: "Permission in linux format.",
				Required:            true,
			},
			"chown_group": schema.StringAttribute{
				MarkdownDescription: "Group used for permission.",
				Required:            true,
			},
			"download_propers_and_repacks": schema.StringAttribute{
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("preferAndUpgrade", "doNotUpgrade", "doNotPrefer"),
				},
			},
			"extra_file_extensions": schema.StringAttribute{
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Required:            true,
			},
			"file_date": schema.StringAttribute{
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'cinemas, and 'release'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "cinemas", "release"),
				},
			},
			"recycle_bin": schema.StringAttribute{
				MarkdownDescription: "Recycle bin absolute path.",
				Required:            true,
			},
			"rescan_after_refresh": schema.StringAttribute{
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("always", "afterManual", "never"),
				},
			},
		},
	}
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
	var management *MediaManagement

	resp.Diagnostics.Append(req.Plan.Get(ctx, &management)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := management.read()
	data.ID = 1

	// Create new MediaManagement
	response, err := r.client.UpdateMediaManagementContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create mediamanagement, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created media_management: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	management.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &management)...)
}

func (r *MediaManagementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var management *MediaManagement

	resp.Diagnostics.Append(req.State.Get(ctx, &management)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get mediamanagement current value
	response, err := r.client.GetMediaManagementContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", mediaManagementResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+mediaManagementResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	management.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &management)...)
}

func (r *MediaManagementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var management *MediaManagement

	resp.Diagnostics.Append(req.Plan.Get(ctx, &management)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := management.read()

	// Update MediaManagement
	response, err := r.client.UpdateMediaManagementContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", mediaManagementResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+mediaManagementResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	management.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &management)...)
}

func (r *MediaManagementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Mediamanagement cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+mediaManagementResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *MediaManagementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+mediaManagementResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (m *MediaManagement) write(mediaMgt *radarr.MediaManagement) {
	m.AutoRenameFolders = types.BoolValue(mediaMgt.AutoRenameFolders)
	m.AutoUnmonitorPreviouslyDownloadedMovies = types.BoolValue(mediaMgt.AutoUnmonitorPreviouslyDownloadedMovies)
	m.CopyUsingHardlinks = types.BoolValue(mediaMgt.CopyUsingHardlinks)
	m.CreateEmptyMovieFolders = types.BoolValue(mediaMgt.CreateEmptyMovieFolders)
	m.DeleteEmptyFolders = types.BoolValue(mediaMgt.DeleteEmptyFolders)
	m.EnableMediaInfo = types.BoolValue(mediaMgt.EnableMediaInfo)
	m.ImportExtraFiles = types.BoolValue(mediaMgt.ImportExtraFiles)
	m.PathsDefaultStatic = types.BoolValue(mediaMgt.PathsDefaultStatic)
	m.SetPermissionsLinux = types.BoolValue(mediaMgt.SetPermissionsLinux)
	m.SkipFreeSpaceCheckWhenImporting = types.BoolValue(mediaMgt.SkipFreeSpaceCheckWhenImporting)
	m.ID = types.Int64Value(mediaMgt.ID)
	m.MinimumFreeSpaceWhenImporting = types.Int64Value(mediaMgt.MinimumFreeSpaceWhenImporting)
	m.RecycleBinCleanupDays = types.Int64Value(mediaMgt.RecycleBinCleanupDays)
	m.ChmodFolder = types.StringValue(mediaMgt.ChmodFolder)
	m.ChownGroup = types.StringValue(mediaMgt.ChownGroup)
	m.DownloadPropersAndRepacks = types.StringValue(mediaMgt.DownloadPropersAndRepacks)
	m.ExtraFileExtensions = types.StringValue(mediaMgt.ExtraFileExtensions)
	m.FileDate = types.StringValue(mediaMgt.FileDate)
	m.RecycleBin = types.StringValue(mediaMgt.RecycleBin)
	m.RescanAfterRefresh = types.StringValue(mediaMgt.RescanAfterRefresh)
}

func (m *MediaManagement) read() *radarr.MediaManagement {
	return &radarr.MediaManagement{
		AutoRenameFolders:                       m.AutoRenameFolders.ValueBool(),
		AutoUnmonitorPreviouslyDownloadedMovies: m.AutoUnmonitorPreviouslyDownloadedMovies.ValueBool(),
		CopyUsingHardlinks:                      m.CopyUsingHardlinks.ValueBool(),
		CreateEmptyMovieFolders:                 m.CreateEmptyMovieFolders.ValueBool(),
		DeleteEmptyFolders:                      m.DeleteEmptyFolders.ValueBool(),
		EnableMediaInfo:                         m.EnableMediaInfo.ValueBool(),
		ImportExtraFiles:                        m.ImportExtraFiles.ValueBool(),
		PathsDefaultStatic:                      m.PathsDefaultStatic.ValueBool(),
		SetPermissionsLinux:                     m.SetPermissionsLinux.ValueBool(),
		SkipFreeSpaceCheckWhenImporting:         m.SkipFreeSpaceCheckWhenImporting.ValueBool(),
		ID:                                      m.ID.ValueInt64(),
		MinimumFreeSpaceWhenImporting:           m.MinimumFreeSpaceWhenImporting.ValueInt64(),
		RecycleBinCleanupDays:                   m.RecycleBinCleanupDays.ValueInt64(),
		ChmodFolder:                             m.ChmodFolder.ValueString(),
		ChownGroup:                              m.ChownGroup.ValueString(),
		DownloadPropersAndRepacks:               m.DownloadPropersAndRepacks.ValueString(),
		ExtraFileExtensions:                     m.ExtraFileExtensions.ValueString(),
		FileDate:                                m.FileDate.ValueString(),
		RecycleBin:                              m.RecycleBin.ValueString(),
		RescanAfterRefresh:                      m.RescanAfterRefresh.ValueString(),
	}
}
