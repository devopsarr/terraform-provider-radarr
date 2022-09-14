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
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the system status implementation.
type SystemStatusDataSource struct {
	client *radarr.Radarr
}

// SystemStatus describes the system status data model.
type SystemStatus struct {
	IsDebug           types.Bool `tfsdk:"is_debug"`
	IsProduction      types.Bool `tfsdk:"is_production"`
	IsAdmin           types.Bool `tfsdk:"is_admin"`
	IsUserInteractive types.Bool `tfsdk:"is_user_interactive"`
	IsMono            types.Bool `tfsdk:"is_mono"`
	IsNetCore         types.Bool `tfsdk:"is_net_core"`
	IsLinux           types.Bool `tfsdk:"is_linux"`
	IsOsx             types.Bool `tfsdk:"is_osx"`
	IsWindows         types.Bool `tfsdk:"is_windows"`
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID               types.Int64  `tfsdk:"id"`
	MigrationVersion types.Int64  `tfsdk:"migration_version"`
	Version          types.String `tfsdk:"version"`
	StartupPath      types.String `tfsdk:"startup_path"`
	AppData          types.String `tfsdk:"app_data"`
	OsName           types.String `tfsdk:"os_name"`
	OsVersion        types.String `tfsdk:"os_version"`
	Branch           types.String `tfsdk:"branch"`
	Authentication   types.String `tfsdk:"authentication"`
	SqliteVersion    types.String `tfsdk:"sqlite_version"`
	URLBase          types.String `tfsdk:"url_base"`
	RuntimeVersion   types.String `tfsdk:"runtime_version"`
	RuntimeName      types.String `tfsdk:"runtime_name"`
	BuildTime        types.String `tfsdk:"build_time"`
}

func (d *SystemStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_status"
}

func (d *SystemStatusDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[subcategory:Status]: #\nSystem Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/radarr/system#status) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"is_debug": {
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_production": {
				MarkdownDescription: "Is production flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_admin": {
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_user_interactive": {
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_net_core": {
				MarkdownDescription: "Is net core flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_mono": {
				MarkdownDescription: "Is mono flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_linux": {
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_osx": {
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_windows": {
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"migration_version": {
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"version": {
				MarkdownDescription: "Version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"startup_path": {
				MarkdownDescription: "Startup path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"app_data": {
				MarkdownDescription: "App data folder.",
				Computed:            true,
				Type:                types.StringType,
			},
			"os_name": {
				MarkdownDescription: "OS name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"os_version": {
				MarkdownDescription: "OS version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"branch": {
				MarkdownDescription: "Branch.",
				Computed:            true,
				Type:                types.StringType,
			},
			"authentication": {
				MarkdownDescription: "Authentication.",
				Computed:            true,
				Type:                types.StringType,
			},
			"sqlite_version": {
				MarkdownDescription: "SQLite version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"url_base": {
				MarkdownDescription: "Base URL.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_version": {
				MarkdownDescription: "Runtime version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_name": {
				MarkdownDescription: "Runtime name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"build_time": {
				MarkdownDescription: "Build time.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, err := d.client.GetSystemStatusContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read system status, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read system status")

	result := writeSystemStatus(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func writeSystemStatus(status *radarr.SystemStatus) *SystemStatus {
	return &SystemStatus{
		IsDebug:           types.Bool{Value: status.IsDebug},
		IsProduction:      types.Bool{Value: status.IsProduction},
		IsAdmin:           types.Bool{Value: status.IsProduction},
		IsUserInteractive: types.Bool{Value: status.IsUserInteractive},
		IsMono:            types.Bool{Value: status.IsMono},
		IsNetCore:         types.Bool{Value: status.IsNetCore},
		IsLinux:           types.Bool{Value: status.IsLinux},
		IsOsx:             types.Bool{Value: status.IsOsx},
		IsWindows:         types.Bool{Value: status.IsWindows},
		ID:                types.Int64{Value: int64(1)},
		MigrationVersion:  types.Int64{Value: int64(status.MigrationVersion)},
		Version:           types.String{Value: status.Version},
		StartupPath:       types.String{Value: status.StartupPath},
		AppData:           types.String{Value: status.AppData},
		OsName:            types.String{Value: status.OsName},
		OsVersion:         types.String{Value: status.OsVersion},
		Branch:            types.String{Value: status.Branch},
		Authentication:    types.String{Value: status.Authentication},
		SqliteVersion:     types.String{Value: status.SqliteVersion},
		URLBase:           types.String{Value: status.URLBase},
		RuntimeVersion:    types.String{Value: status.RuntimeVersion},
		RuntimeName:       types.String{Value: status.RuntimeName},
		BuildTime:         types.String{Value: status.BuildTime.String()},
	}
}
