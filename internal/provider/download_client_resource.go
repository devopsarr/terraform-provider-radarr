package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

const downloadClientResourceName = "download_client"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientResource{}
	_ resource.ResourceWithImportState = &DownloadClientResource{}
)

var downloadClientFields = helpers.Fields{
	Bools:                  []string{"addPaused", "useSsl", "startOnAdd", "sequentialOrder", "firstAndLast", "addStopped", "saveMagnetFiles", "readOnly"},
	Ints:                   []string{"port", "recentMoviePriority", "olderMoviePriority", "recentPriority", "olderPriority", "initialState", "intialState"},
	Strings:                []string{"host", "apiKey", "urlBase", "rpcPath", "secretToken", "password", "username", "movieCategory", "movieImportedCategory", "movieDirectory", "destinationDirectory", "destination", "category", "nzbFolder", "strmFolder", "torrentFolder", "magnetFileExtension", "watchFolder", "apiUrl", "appId", "appToken"},
	StringSlices:           []string{"fieldTags", "postImportTags"},
	StringSlicesExceptions: []string{"tags"},
	IntSlices:              []string{"additionalTags"},
}

func NewDownloadClientResource() resource.Resource {
	return &DownloadClientResource{}
}

// DownloadClientResource defines the download client implementation.
type DownloadClientResource struct {
	client *radarr.APIClient
	auth   context.Context
}

// DownloadClient describes the download client data model.
type DownloadClient struct {
	Tags                     types.Set    `tfsdk:"tags"`
	PostImportTags           types.Set    `tfsdk:"post_import_tags"`
	FieldTags                types.Set    `tfsdk:"field_tags"`
	AdditionalTags           types.Set    `tfsdk:"additional_tags"`
	NzbFolder                types.String `tfsdk:"nzb_folder"`
	Category                 types.String `tfsdk:"category"`
	Implementation           types.String `tfsdk:"implementation"`
	Name                     types.String `tfsdk:"name"`
	Protocol                 types.String `tfsdk:"protocol"`
	MagnetFileExtension      types.String `tfsdk:"magnet_file_extension"`
	TorrentFolder            types.String `tfsdk:"torrent_folder"`
	StrmFolder               types.String `tfsdk:"strm_folder"`
	Host                     types.String `tfsdk:"host"`
	ConfigContract           types.String `tfsdk:"config_contract"`
	Destination              types.String `tfsdk:"destination"`
	MovieDirectory           types.String `tfsdk:"movie_directory"`
	Username                 types.String `tfsdk:"username"`
	MovieImportedCategory    types.String `tfsdk:"movie_imported_category"`
	MovieCategory            types.String `tfsdk:"movie_category"`
	Password                 types.String `tfsdk:"password"`
	SecretToken              types.String `tfsdk:"secret_token"`
	RPCPath                  types.String `tfsdk:"rpc_path"`
	URLBase                  types.String `tfsdk:"url_base"`
	APIURL                   types.String `tfsdk:"api_url"`
	APIKey                   types.String `tfsdk:"api_key"`
	AppID                    types.String `tfsdk:"app_id"`
	AppToken                 types.String `tfsdk:"app_token"`
	WatchFolder              types.String `tfsdk:"watch_folder"`
	DestinationDirectory     types.String `tfsdk:"destination_directory"`
	RecentPriority           types.Int64  `tfsdk:"recent_priority"`
	OlderPriority            types.Int64  `tfsdk:"older_priority"`
	RecentMoviePriority      types.Int64  `tfsdk:"recent_movie_priority"`
	IntialState              types.Int64  `tfsdk:"intial_state"`
	InitialState             types.Int64  `tfsdk:"initial_state"`
	OlderMoviePriority       types.Int64  `tfsdk:"older_movie_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	AddStopped               types.Bool   `tfsdk:"add_stopped"`
	SaveMagnetFiles          types.Bool   `tfsdk:"save_magnet_files"`
	ReadOnly                 types.Bool   `tfsdk:"read_only"`
	FirstAndLast             types.Bool   `tfsdk:"first_and_last"`
	SequentialOrder          types.Bool   `tfsdk:"sequential_order"`
	StartOnAdd               types.Bool   `tfsdk:"start_on_add"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	AddPaused                types.Bool   `tfsdk:"add_paused"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClient) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":                       types.SetType{}.WithElementType(types.Int64Type),
			"additional_tags":            types.SetType{}.WithElementType(types.Int64Type),
			"post_import_tags":           types.SetType{}.WithElementType(types.StringType),
			"field_tags":                 types.SetType{}.WithElementType(types.StringType),
			"nzb_folder":                 types.StringType,
			"category":                   types.StringType,
			"implementation":             types.StringType,
			"name":                       types.StringType,
			"protocol":                   types.StringType,
			"magnet_file_extension":      types.StringType,
			"torrent_folder":             types.StringType,
			"strm_folder":                types.StringType,
			"host":                       types.StringType,
			"config_contract":            types.StringType,
			"destination":                types.StringType,
			"movie_directory":            types.StringType,
			"username":                   types.StringType,
			"movie_imported_category":    types.StringType,
			"movie_category":             types.StringType,
			"password":                   types.StringType,
			"secret_token":               types.StringType,
			"rpc_path":                   types.StringType,
			"url_base":                   types.StringType,
			"api_url":                    types.StringType,
			"api_key":                    types.StringType,
			"app_id":                     types.StringType,
			"app_token":                  types.StringType,
			"watch_folder":               types.StringType,
			"destination_directory":      types.StringType,
			"recent_priority":            types.Int64Type,
			"older_priority":             types.Int64Type,
			"recent_movie_priority":      types.Int64Type,
			"intial_state":               types.Int64Type,
			"initial_state":              types.Int64Type,
			"older_movie_priority":       types.Int64Type,
			"priority":                   types.Int64Type,
			"port":                       types.Int64Type,
			"id":                         types.Int64Type,
			"add_stopped":                types.BoolType,
			"save_magnet_files":          types.BoolType,
			"read_only":                  types.BoolType,
			"first_and_last":             types.BoolType,
			"sequential_order":           types.BoolType,
			"start_on_add":               types.BoolType,
			"use_ssl":                    types.BoolType,
			"add_paused":                 types.BoolType,
			"enable":                     types.BoolType,
			"remove_failed_downloads":    types.BoolType,
			"remove_completed_downloads": types.BoolType,
		})
}

func (r *DownloadClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientResourceName
}

func (r *DownloadClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->\nGeneric Download Client resource. When possible use a specific resource instead.\nFor more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients).",
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
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "DownloadClient configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "DownloadClient implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("usenet", "torrent"),
				},
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
			"start_on_add": schema.BoolAttribute{
				MarkdownDescription: "Start on add flag.",
				Optional:            true,
				Computed:            true,
			},
			"sequential_order": schema.BoolAttribute{
				MarkdownDescription: "Sequential order flag.",
				Optional:            true,
				Computed:            true,
			},
			"first_and_last": schema.BoolAttribute{
				MarkdownDescription: "First and last flag.",
				Optional:            true,
				Computed:            true,
			},
			"add_stopped": schema.BoolAttribute{
				MarkdownDescription: "Add stopped flag.",
				Optional:            true,
				Computed:            true,
			},
			"save_magnet_files": schema.BoolAttribute{
				MarkdownDescription: "Save magnet files flag.",
				Optional:            true,
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Read only flag.",
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
			"initial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state. `0` Start, `1` ForceStart, `2` Pause.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"intial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"rpc_path": schema.StringAttribute{
				MarkdownDescription: "RPC path.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "API URL.",
				Optional:            true,
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "App ID.",
				Optional:            true,
				Computed:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App Token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"secret_token": schema.StringAttribute{
				MarkdownDescription: "Secret token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
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
			"movie_directory": schema.StringAttribute{
				MarkdownDescription: "Movie directory.",
				Optional:            true,
				Computed:            true,
			},
			"destination_directory": schema.StringAttribute{
				MarkdownDescription: "Movie directory.",
				Optional:            true,
				Computed:            true,
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: "Destination.",
				Optional:            true,
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
			},
			"nzb_folder": schema.StringAttribute{
				MarkdownDescription: "NZB folder.",
				Optional:            true,
				Computed:            true,
			},
			"strm_folder": schema.StringAttribute{
				MarkdownDescription: "STRM folder.",
				Optional:            true,
				Computed:            true,
			},
			"torrent_folder": schema.StringAttribute{
				MarkdownDescription: "Torrent folder.",
				Optional:            true,
				Computed:            true,
			},
			"magnet_file_extension": schema.StringAttribute{
				MarkdownDescription: "Magnet file extension.",
				Optional:            true,
				Computed:            true,
			},
			"watch_folder": schema.StringAttribute{
				MarkdownDescription: "Watch folder flag.",
				Optional:            true,
				Computed:            true,
			},
			"additional_tags": schema.SetAttribute{
				MarkdownDescription: "Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Field tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"post_import_tags": schema.SetAttribute{
				MarkdownDescription: "Post import tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *DownloadClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *DownloadClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClient
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(r.auth).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.writeSensitive(client)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client *DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClient current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(r.auth, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.writeSensitive(client)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClient
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(r.auth, request.GetId()).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.writeSensitive(client)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClient current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientResourceName+": "+req.ID)
}

func (d *DownloadClient) write(ctx context.Context, downloadClient *radarr.DownloadClientResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	d.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, downloadClient.Tags)
	diags.Append(localDiag...)

	d.Enable = types.BoolValue(downloadClient.GetEnable())
	d.RemoveCompletedDownloads = types.BoolValue(downloadClient.GetRemoveCompletedDownloads())
	d.RemoveFailedDownloads = types.BoolValue(downloadClient.GetRemoveFailedDownloads())
	d.Priority = types.Int64Value(int64(downloadClient.GetPriority()))
	d.ID = types.Int64Value(int64(downloadClient.GetId()))
	d.ConfigContract = types.StringValue(downloadClient.GetConfigContract())
	d.Implementation = types.StringValue(downloadClient.GetImplementation())
	d.Name = types.StringValue(downloadClient.GetName())
	d.Protocol = types.StringValue(string(downloadClient.GetProtocol()))
	d.AdditionalTags = types.SetValueMust(types.Int64Type, nil)
	d.FieldTags = types.SetValueMust(types.StringType, nil)
	d.PostImportTags = types.SetValueMust(types.StringType, nil)
	helpers.WriteFields(ctx, d, downloadClient.GetFields(), downloadClientFields)
}

func (d *DownloadClient) read(ctx context.Context, diags *diag.Diagnostics) *radarr.DownloadClientResource {
	client := radarr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetRemoveCompletedDownloads(d.RemoveCompletedDownloads.ValueBool())
	client.SetRemoveFailedDownloads(d.RemoveFailedDownloads.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(d.ConfigContract.ValueString())
	client.SetImplementation(d.Implementation.ValueString())
	client.SetName(d.Name.ValueString())
	client.SetProtocol(radarr.DownloadProtocol(d.Protocol.ValueString()))
	diags.Append(d.Tags.ElementsAs(ctx, &client.Tags, true)...)
	client.SetFields(helpers.ReadFields(ctx, d, downloadClientFields))

	return client
}

// writeSensitive copy sensitive data from another resource.
func (d *DownloadClient) writeSensitive(client *DownloadClient) {
	if !client.Password.IsUnknown() {
		d.Password = client.Password
	}

	if !client.APIKey.IsUnknown() {
		d.APIKey = client.APIKey
	}

	if !client.SecretToken.IsUnknown() {
		d.SecretToken = client.SecretToken
	}

	if !client.AppToken.IsUnknown() {
		d.AppToken = client.AppToken
	}
}
