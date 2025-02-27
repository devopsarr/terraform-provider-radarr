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
	downloadClientFreeboxResourceName   = "download_client_freebox"
	downloadClientFreeboxImplementation = "TorrentFreeboxDownload"
	downloadClientFreeboxConfigContract = "FreeboxDownloadSettings"
	downloadClientFreeboxProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientFreeboxResource{}
	_ resource.ResourceWithImportState = &DownloadClientFreeboxResource{}
)

func NewDownloadClientFreeboxResource() resource.Resource {
	return &DownloadClientFreeboxResource{}
}

// DownloadClientFreeboxResource defines the download client implementation.
type DownloadClientFreeboxResource struct {
	client *radarr.APIClient
	auth   context.Context
}

// DownloadClientFreebox describes the download client data model.
type DownloadClientFreebox struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	APIURL                   types.String `tfsdk:"api_url"`
	AppID                    types.String `tfsdk:"app_id"`
	AppToken                 types.String `tfsdk:"app_token"`
	Category                 types.String `tfsdk:"category"`
	DestinationDirectory     types.String `tfsdk:"destination_directory"`
	RecentPriority           types.Int64  `tfsdk:"recent_priority"`
	OlderPriority            types.Int64  `tfsdk:"older_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	AddPaused                types.Bool   `tfsdk:"add_paused"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientFreebox) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		APIURL:                   d.APIURL,
		AppID:                    d.AppID,
		AppToken:                 d.AppToken,
		Category:                 d.Category,
		DestinationDirectory:     d.DestinationDirectory,
		RecentPriority:           d.RecentPriority,
		OlderPriority:            d.OlderPriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		AddPaused:                d.AddPaused,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		Implementation:           types.StringValue(downloadClientFreeboxImplementation),
		ConfigContract:           types.StringValue(downloadClientFreeboxConfigContract),
		Protocol:                 types.StringValue(downloadClientFreeboxProtocol),
	}
}

func (d *DownloadClientFreebox) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.APIURL = client.APIURL
	d.AppID = client.AppID
	d.AppToken = client.AppToken
	d.Category = client.Category
	d.DestinationDirectory = client.DestinationDirectory
	d.RecentPriority = client.RecentPriority
	d.OlderPriority = client.OlderPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.AddPaused = client.AddPaused
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientFreeboxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientFreeboxResourceName
}

func (r *DownloadClientFreeboxResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->\nDownload Client Freebox resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients) and [Freebox](https://wiki.servarr.com/radarr/supported#torrentfreeboxdownload).",
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
				Required:            true,
			},
			"recent_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"older_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Required:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "API URL.",
				Required:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "App ID.",
				Required:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App Token.",
				Required:            true,
				Sensitive:           true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "category.",
				Optional:            true,
				Computed:            true,
			},
			"destination_directory": schema.StringAttribute{
				MarkdownDescription: "Movie directory.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientFreeboxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *DownloadClientFreeboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientFreebox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientFreebox
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(r.auth).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientFreeboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientFreeboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientFreeboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientFreebox

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientFreebox current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(r.auth, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientFreeboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientFreeboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientFreeboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientFreebox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientFreebox
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(r.auth, request.GetId()).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientFreeboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientFreeboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientFreeboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientFreebox current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientFreeboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientFreeboxResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientFreeboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientFreeboxResourceName+": "+req.ID)
}

func (d *DownloadClientFreebox) write(ctx context.Context, downloadClient *radarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientFreebox) read(ctx context.Context, diags *diag.Diagnostics) *radarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}
