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

const systemStatusDataSourceName = "system_status"

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
	AppName                types.String `tfsdk:"app_name"`
	Version                types.String `tfsdk:"version"`
	Branch                 types.String `tfsdk:"branch"`
	OsName                 types.String `tfsdk:"os_name"`
	PackageVersion         types.String `tfsdk:"package_version"`
	PackageUpdateMechanism types.String `tfsdk:"package_update_mechanism"`
	PackageAuthor          types.String `tfsdk:"package_author"`
	Mode                   types.String `tfsdk:"mode"`
	InstanceName           types.String `tfsdk:"instance_name"`
	BuildTime              types.String `tfsdk:"build_time"`
	RuntimeName            types.String `tfsdk:"runtime_name"`
	DatabaseType           types.String `tfsdk:"database_type"`
	StartupPath            types.String `tfsdk:"startup_path"`
	AppData                types.String `tfsdk:"app_data"`
	StartTime              types.String `tfsdk:"start_time"`
	DatabaseVersion        types.String `tfsdk:"database_version"`
	Authentication         types.String `tfsdk:"authentication"`
	URLBase                types.String `tfsdk:"url_base"`
	RuntimeVersion         types.String `tfsdk:"runtime_version"`
	MigrationVersion       types.Int64  `tfsdk:"migration_version"`
	ID                     types.Int64  `tfsdk:"id"`
	IsWindows              types.Bool   `tfsdk:"is_windows"`
	IsDebug                types.Bool   `tfsdk:"is_debug"`
	IsAdmin                types.Bool   `tfsdk:"is_admin"`
	IsProduction           types.Bool   `tfsdk:"is_production"`
	IsOsx                  types.Bool   `tfsdk:"is_osx"`
	IsLinux                types.Bool   `tfsdk:"is_linux"`
	IsDocker               types.Bool   `tfsdk:"is_docker"`
	IsNetCore              types.Bool   `tfsdk:"is_net_core"`
	IsUserInteractive      types.Bool   `tfsdk:"is_user_interactive"`
}

func (d *SystemStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + systemStatusDataSourceName
}

func (d *SystemStatusDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Status -->System Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/radarr/system#status) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				MarkdownDescription: "Status ID.",
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
			"is_docker": {
				MarkdownDescription: "Is docker flag.",
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
			"app_name": {
				MarkdownDescription: "Application name.",
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
			"url_base": {
				MarkdownDescription: "Base URL.",
				Computed:            true,
				Type:                types.StringType,
			},
			"database_type": {
				MarkdownDescription: "Database type.",
				Computed:            true,
				Type:                types.StringType,
			},
			"database_version": {
				MarkdownDescription: "Database version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"instance_name": {
				MarkdownDescription: "Instance name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"mode": {
				MarkdownDescription: "Mode.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_author": {
				MarkdownDescription: "Package author.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_update_mechanism": {
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_version": {
				MarkdownDescription: "Package version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_name": {
				MarkdownDescription: "Runtime name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_version": {
				MarkdownDescription: "Runtime version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"build_time": {
				MarkdownDescription: "Build time.",
				Computed:            true,
				Type:                types.StringType,
			},
			"start_time": {
				MarkdownDescription: "Start time.",
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
			tools.UnexpectedDataSourceConfigureType,
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
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	status := SystemStatus{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &status)...)
}

func (s *SystemStatus) write(status *radarr.SystemStatus) {
	s.IsDebug = types.BoolValue(status.IsDebug)
	s.IsProduction = types.BoolValue(status.IsProduction)
	s.IsAdmin = types.BoolValue(status.IsProduction)
	s.IsUserInteractive = types.BoolValue(status.IsUserInteractive)
	s.IsNetCore = types.BoolValue(status.IsNetCore)
	s.IsDocker = types.BoolValue(status.IsDocker)
	s.IsLinux = types.BoolValue(status.IsLinux)
	s.IsOsx = types.BoolValue(status.IsOsx)
	s.IsWindows = types.BoolValue(status.IsWindows)
	s.ID = types.Int64Value(int64(1))
	s.MigrationVersion = types.Int64Value(status.MigrationVersion)
	s.Version = types.StringValue(status.Version)
	s.StartupPath = types.StringValue(status.StartupPath)
	s.AppData = types.StringValue(status.AppData)
	s.OsName = types.StringValue(status.OsName)
	s.Branch = types.StringValue(status.Branch)
	s.Authentication = types.StringValue(status.Authentication)
	s.URLBase = types.StringValue(status.URLBase)
	s.RuntimeVersion = types.StringValue(status.RuntimeVersion)
	s.RuntimeName = types.StringValue(status.RuntimeName)
	s.AppName = types.StringValue(status.AppName)
	s.DatabaseType = types.StringValue(status.DatabaseType)
	s.DatabaseVersion = types.StringValue(status.DatabaseVersion)
	s.InstanceName = types.StringValue(status.InstanceName)
	s.Mode = types.StringValue(status.Mode)
	s.PackageAuthor = types.StringValue(status.PackageAuthor)
	s.PackageUpdateMechanism = types.StringValue(status.PackageUpdateMechanism)
	s.PackageVersion = types.StringValue(status.PackageVersion)
	s.BuildTime = types.StringValue(status.BuildTime.String())
	s.StartTime = types.StringValue(status.StartTime.String())
}
