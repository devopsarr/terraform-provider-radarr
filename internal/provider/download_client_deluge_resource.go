package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientDelugeResourceName   = "download_client_deluge"
	downloadClientDelugeImplementation = "Deluge"
	downloadClientDelugeConfigContract = "DelugeSettings"
	downloadClientDelugeProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientDelugeResource{}
	_ resource.ResourceWithImportState = &DownloadClientDelugeResource{}
)

func NewDownloadClientDelugeResource() resource.Resource {
	return &DownloadClientDelugeResource{}
}

// DownloadClientDelugeResource defines the download client implementation.
type DownloadClientDelugeResource struct {
	client *radarr.APIClient
	auth   context.Context
}

// DownloadClientDeluge describes the download client data model.
type DownloadClientDeluge struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Password                 types.String `tfsdk:"password"`
	MovieCategory            types.String `tfsdk:"movie_category"`
	MovieImportedCategory    types.String `tfsdk:"movie_imported_category"`
	RecentMoviePriority      types.Int64  `tfsdk:"recent_movie_priority"`
	OlderMoviePriority       types.Int64  `tfsdk:"older_movie_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	AddPaused                types.Bool   `tfsdk:"add_paused"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientDeluge) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Password:                 d.Password,
		MovieCategory:            d.MovieCategory,
		MovieImportedCategory:    d.MovieImportedCategory,
		RecentMoviePriority:      d.RecentMoviePriority,
		OlderMoviePriority:       d.OlderMoviePriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		AddPaused:                d.AddPaused,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		Implementation:           types.StringValue(downloadClientDelugeImplementation),
		ConfigContract:           types.StringValue(downloadClientDelugeConfigContract),
		Protocol:                 types.StringValue(downloadClientDelugeProtocol),
	}
}

func (d *DownloadClientDeluge) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Password = client.Password
	d.MovieCategory = client.MovieCategory
	d.MovieImportedCategory = client.MovieImportedCategory
	d.RecentMoviePriority = client.RecentMoviePriority
	d.OlderMoviePriority = client.OlderMoviePriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.AddPaused = client.AddPaused
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientDelugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientDelugeResourceName
}

func (r *DownloadClientDelugeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Deluge resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients) and [Deluge](https://wiki.servarr.com/radarr/supported#deluge).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_completed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_failed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"add_paused": schema.BoolAttribute{
				MarkdownDescription: "Add paused flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"older_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"movie_category": schema.StringAttribute{
				MarkdownDescription: "Movie category.",
				Optional:            true,
				Computed:            true,
			},
			"movie_imported_category": schema.StringAttribute{
				MarkdownDescription: "Movie imported category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientDelugeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *DownloadClientDelugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientDeluge

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientDeluge
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(r.auth).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientDelugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientDeluge

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientDeluge current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(r.auth, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientDelugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientDeluge

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientDeluge
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(r.auth, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientDelugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientDeluge current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientDelugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientDelugeResourceName+": "+req.ID)
}

func (d *DownloadClientDeluge) write(ctx context.Context, downloadClient *radarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientDeluge) read(ctx context.Context, diags *diag.Diagnostics) *radarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}
