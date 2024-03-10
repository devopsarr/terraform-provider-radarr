package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const systemStatusDataSourceName = "system_status"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the system status implementation.
type SystemStatusDataSource struct {
	client *radarr.APIClient
	auth   context.Context
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

func (d *SystemStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + systemStatusDataSourceName
}

func (d *SystemStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:System -->\nSystem Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/radarr/system#status) documentation.",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Status ID.",
				Computed:            true,
			},
			"is_debug": schema.BoolAttribute{
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
			},
			"is_production": schema.BoolAttribute{
				MarkdownDescription: "Is production flag.",
				Computed:            true,
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
			},
			"is_user_interactive": schema.BoolAttribute{
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
			},
			"is_net_core": schema.BoolAttribute{
				MarkdownDescription: "Is net core flag.",
				Computed:            true,
			},
			"is_docker": schema.BoolAttribute{
				MarkdownDescription: "Is docker flag.",
				Computed:            true,
			},
			"is_linux": schema.BoolAttribute{
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
			},
			"is_osx": schema.BoolAttribute{
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
			},
			"is_windows": schema.BoolAttribute{
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
			},
			"migration_version": schema.Int64Attribute{
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version.",
				Computed:            true,
			},
			"startup_path": schema.StringAttribute{
				MarkdownDescription: "Startup path.",
				Computed:            true,
			},
			"app_data": schema.StringAttribute{
				MarkdownDescription: "App data folder.",
				Computed:            true,
			},
			"os_name": schema.StringAttribute{
				MarkdownDescription: "OS name.",
				Computed:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "Application name.",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch.",
				Computed:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"database_type": schema.StringAttribute{
				MarkdownDescription: "Database type.",
				Computed:            true,
			},
			"database_version": schema.StringAttribute{
				MarkdownDescription: "Database version.",
				Computed:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode.",
				Computed:            true,
			},
			"package_author": schema.StringAttribute{
				MarkdownDescription: "Package author.",
				Computed:            true,
			},
			"package_update_mechanism": schema.StringAttribute{
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
			},
			"package_version": schema.StringAttribute{
				MarkdownDescription: "Package version.",
				Computed:            true,
			},
			"runtime_name": schema.StringAttribute{
				MarkdownDescription: "Runtime name.",
				Computed:            true,
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "Runtime version.",
				Computed:            true,
			},
			"build_time": schema.StringAttribute{
				MarkdownDescription: "Build time.",
				Computed:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Start time.",
				Computed:            true,
			},
		},
	}
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *SystemStatusDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, _, err := d.client.SystemAPI.GetSystemStatus(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	status := SystemStatus{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &status)...)
}

func (s *SystemStatus) write(status *radarr.SystemResource) {
	s.IsDebug = types.BoolValue(status.GetIsDebug())
	s.IsProduction = types.BoolValue(status.GetIsProduction())
	s.IsAdmin = types.BoolValue(status.GetIsProduction())
	s.IsUserInteractive = types.BoolValue(status.GetIsUserInteractive())
	s.IsNetCore = types.BoolValue(status.GetIsNetCore())
	s.IsDocker = types.BoolValue(status.GetIsDocker())
	s.IsLinux = types.BoolValue(status.GetIsLinux())
	s.IsOsx = types.BoolValue(status.GetIsOsx())
	s.IsWindows = types.BoolValue(status.GetIsWindows())
	s.ID = types.Int64Value(int64(1))
	s.MigrationVersion = types.Int64Value(int64(status.GetMigrationVersion()))
	s.Version = types.StringValue(status.GetVersion())
	s.StartupPath = types.StringValue(status.GetStartupPath())
	s.AppData = types.StringValue(status.GetAppData())
	s.OsName = types.StringValue(status.GetOsName())
	s.Branch = types.StringValue(status.GetBranch())
	s.Authentication = types.StringValue(string(status.GetAuthentication()))
	s.URLBase = types.StringValue(status.GetUrlBase())
	s.RuntimeVersion = types.StringValue(status.GetRuntimeVersion())
	s.RuntimeName = types.StringValue(status.GetRuntimeName())
	s.AppName = types.StringValue(status.GetAppName())
	s.DatabaseType = types.StringValue(string(status.GetDatabaseType()))
	s.DatabaseVersion = types.StringValue(status.GetDatabaseVersion())
	s.InstanceName = types.StringValue(status.GetInstanceName())
	s.Mode = types.StringValue(string(status.GetMode()))
	s.PackageAuthor = types.StringValue(status.GetPackageAuthor())
	s.PackageUpdateMechanism = types.StringValue(string(status.GetPackageUpdateMechanism()))
	s.PackageVersion = types.StringValue(status.GetPackageVersion())
	s.BuildTime = types.StringValue(status.BuildTime.String())
	s.StartTime = types.StringValue(status.StartTime.String())
}
