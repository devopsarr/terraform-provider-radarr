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
	"golang.org/x/exp/slices"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

const notificationResourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationResource{}
var _ resource.ResourceWithImportState = &NotificationResource{}

var (
	notificationBoolFields        = []string{"alwaysUpdate", "cleanLibrary", "directMessage", "notify", "requireEncryption", "sendSilently", "useSsl", "updateLibrary", "useEuEndpoint"}
	notificationStringFields      = []string{"accessToken", "accessTokenSecret", "apiKey", "aPIKey", "appToken", "arguments", "author", "authToken", "authUser", "avatar", "cc", "bcc", "botToken", "channel", "chatId", "consumerKey", "consumerSecret", "deviceNames", "displayTime", "expires", "from", "host", "icon", "instanceName", "mention", "password", "path", "refreshToken", "senderDomain", "senderId", "server", "signIn", "sound", "to", "token", "url", "userKey", "username", "webHookUrl", "serverUrl", "userName", "clickUrl", "mapFrom", "mapTo", "key", "event"}
	notificationIntFields         = []string{"port", "grabFields", "importFields", "priority", "retry", "expire", "method"}
	notificationStringSliceFields = []string{"recipients", "topics", "tags", "channelTags", "devices"}
	notificationIntSliceFields    = []string{"deviceIds"}
)

func NewNotificationResource() resource.Resource {
	return &NotificationResource{}
}

// NotificationResource defines the notification implementation.
type NotificationResource struct {
	client *radarr.Radarr
}

// Notification describes the notification data model.
type Notification struct {
	OnGrab                      types.Bool   `tfsdk:"on_grab"`
	OnDownload                  types.Bool   `tfsdk:"on_download"`
	OnUpgrade                   types.Bool   `tfsdk:"on_upgrade"`
	OnRename                    types.Bool   `tfsdk:"on_rename"`
	OnMovieAdded                types.Bool   `tfsdk:"on_movie_added"`
	OnMovieDelete               types.Bool   `tfsdk:"on_movie_delete"`
	OnMovieFileDelete           types.Bool   `tfsdk:"on_movie_file_delete"`
	OnMovieFileDeleteForUpgrade types.Bool   `tfsdk:"on_movie_file_delete_for_upgrade"`
	OnHealthIssue               types.Bool   `tfsdk:"on_health_issue"`
	OnApplicationUpdate         types.Bool   `tfsdk:"on_application_update"`
	IncludeHealthWarnings       types.Bool   `tfsdk:"include_health_warnings"`
	Name                        types.String `tfsdk:"name"`
	Implementation              types.String `tfsdk:"implementation"`
	ConfigContract              types.String `tfsdk:"config_contract"`
	Tags                        types.Set    `tfsdk:"tags"`
	ID                          types.Int64  `tfsdk:"id"`
	Path                        types.String `tfsdk:"path"`
	Arguments                   types.String `tfsdk:"arguments"`
	UpdateLibrary               types.Bool   `tfsdk:"update_library"`
	UseEuEndpoint               types.Bool   `tfsdk:"use_eu_endpoint"`
	Notify                      types.Bool   `tfsdk:"notify"`
	UseSSL                      types.Bool   `tfsdk:"use_ssl"`
	SendSilently                types.Bool   `tfsdk:"send_silently"`
	RequireEncryption           types.Bool   `tfsdk:"require_encryption"`
	DirectMessage               types.Bool   `tfsdk:"direct_message"`
	CleanLibrary                types.Bool   `tfsdk:"clean_library"`
	AlwaysUpdate                types.Bool   `tfsdk:"always_update"`
	RefreshToken                types.String `tfsdk:"refresh_token"`
	WebHookURL                  types.String `tfsdk:"web_hook_url"`
	Username                    types.String `tfsdk:"username"`
	UserKey                     types.String `tfsdk:"user_key"`
	Mention                     types.String `tfsdk:"mention"`
	Avatar                      types.String `tfsdk:"avatar"`
	URL                         types.String `tfsdk:"url"`
	Token                       types.String `tfsdk:"token"`
	To                          types.String `tfsdk:"to"`
	Sound                       types.String `tfsdk:"sound"`
	Bcc                         types.String `tfsdk:"bcc"`
	SignIn                      types.String `tfsdk:"sign_in"`
	Server                      types.String `tfsdk:"server"`
	SenderID                    types.String `tfsdk:"sender_id"`
	BotToken                    types.String `tfsdk:"bot_token"`
	SenderDomain                types.String `tfsdk:"sender_domain"`
	Icon                        types.String `tfsdk:"icon"`
	Host                        types.String `tfsdk:"host"`
	From                        types.String `tfsdk:"from"`
	Expires                     types.String `tfsdk:"expires"`
	AccessToken                 types.String `tfsdk:"access_token"`
	AccessTokenSecret           types.String `tfsdk:"access_token_secret"`
	APIKey                      types.String `tfsdk:"api_key"`
	AppToken                    types.String `tfsdk:"app_token"`
	Author                      types.String `tfsdk:"author"`
	AuthToken                   types.String `tfsdk:"auth_token"`
	AuthUser                    types.String `tfsdk:"auth_user"`
	Password                    types.String `tfsdk:"password"`
	Cc                          types.String `tfsdk:"cc"`
	Channel                     types.String `tfsdk:"channel"`
	ChatID                      types.String `tfsdk:"chat_id"`
	ConsumerKey                 types.String `tfsdk:"consumer_key"`
	ConsumerSecret              types.String `tfsdk:"consumer_secret"`
	DeviceNames                 types.String `tfsdk:"device_names"`
	InstanceName                types.String `tfsdk:"instance_name"`
	ServerURL                   types.String `tfsdk:"server_url"`
	ClickURL                    types.String `tfsdk:"click_url"`
	MapFrom                     types.String `tfsdk:"map_from"`
	MapTo                       types.String `tfsdk:"map_to"`
	DisplayTime                 types.String `tfsdk:"display_time"`
	Key                         types.String `tfsdk:"key"`
	Event                       types.String `tfsdk:"event"`
	Priority                    types.Int64  `tfsdk:"priority"`
	Port                        types.Int64  `tfsdk:"port"`
	Method                      types.Int64  `tfsdk:"method"`
	Retry                       types.Int64  `tfsdk:"retry"`
	Expire                      types.Int64  `tfsdk:"expire"`
	ImportFields                types.Int64  `tfsdk:"import_fields"`
	GrabFields                  types.Int64  `tfsdk:"grab_fields"`
	Recipients                  types.Set    `tfsdk:"recipients"`
	Devices                     types.Set    `tfsdk:"devices"`
	DeviceIds                   types.Set    `tfsdk:"device_ids"`
	Topics                      types.Set    `tfsdk:"topics"`
	ChannelTags                 types.Set    `tfsdk:"channel_tags"`
	FieldTags                   types.Set    `tfsdk:"field_tags"`
}

func (r *NotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationResourceName
}

func (r *NotificationResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect).",
		Attributes: map[string]tfsdk.Attribute{
			"on_grab": {
				MarkdownDescription: "On grab flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_download": {
				MarkdownDescription: "On download flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_upgrade": {
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_rename": {
				MarkdownDescription: "On rename flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_movie_added": {
				MarkdownDescription: "On movie added flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_movie_delete": {
				MarkdownDescription: "On movie delete flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_movie_file_delete": {
				MarkdownDescription: "On movie file delete flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_movie_file_delete_for_upgrade": {
				MarkdownDescription: "On movie file delete for upgrade flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_health_issue": {
				MarkdownDescription: "On health issue flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_application_update": {
				MarkdownDescription: "On application update flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"include_health_warnings": {
				MarkdownDescription: "Include health warnings.",
				Required:            true,
				Type:                types.BoolType,
			},
			"config_contract": {
				MarkdownDescription: "Notification configuration template.",
				Required:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Notification implementation name.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Notification name.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"always_update": {
				MarkdownDescription: "Always update flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"clean_library": {
				MarkdownDescription: "Clean library flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"direct_message": {
				MarkdownDescription: "Direct message flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"notify": {
				MarkdownDescription: "Notify flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"require_encryption": {
				MarkdownDescription: "Require encryption flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"send_silently": {
				MarkdownDescription: "Add silently flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"update_library": {
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_eu_endpoint": {
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_ssl": {
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"port": {
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"grab_fields": {
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart, `10` CustomFormats, `11` CustomFormatScore.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
				},
			},
			"import_fields": {
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}),
				},
			},
			"method": {
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{1, 2}),
				},
			},
			"priority": {
				MarkdownDescription: "Priority.", // TODO: add values in description
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{-2, -1, 0, 1, 2, 3, 4, 5, 7, 8}),
				},
			},
			"retry": {
				MarkdownDescription: "Retry.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"expire": {
				MarkdownDescription: "Expire.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"access_token": {
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"access_token_secret": {
				MarkdownDescription: "Access token secret.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"app_token": {
				MarkdownDescription: "App token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"arguments": {
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"author": {
				MarkdownDescription: "Author.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"auth_token": {
				MarkdownDescription: "Auth token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"auth_user": {
				MarkdownDescription: "Auth user.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"avatar": {
				MarkdownDescription: "Avatar.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"instance_name": {
				MarkdownDescription: "Instance name.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"bcc": {
				MarkdownDescription: "Bcc.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"bot_token": {
				MarkdownDescription: "Bot token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"cc": {
				MarkdownDescription: "Cc.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"channel": {
				MarkdownDescription: "Channel.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"chat_id": {
				MarkdownDescription: "Chat ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"consumer_key": {
				MarkdownDescription: "Consumer key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"consumer_secret": {
				MarkdownDescription: "Consumer secret.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"device_names": {
				MarkdownDescription: "Device names.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"display_time": {
				MarkdownDescription: "Display time.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"expires": {
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"from": {
				MarkdownDescription: "From.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"host": {
				MarkdownDescription: "Host.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"icon": {
				MarkdownDescription: "Icon.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"mention": {
				MarkdownDescription: "Mention.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "password.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"path": {
				MarkdownDescription: "Path.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"refresh_token": {
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sender_domain": {
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sender_id": {
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"server": {
				MarkdownDescription: "server.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sign_in": {
				MarkdownDescription: "Sign in.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sound": {
				MarkdownDescription: "Sound.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"to": {
				MarkdownDescription: "To.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"token": {
				MarkdownDescription: "Token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"url": {
				MarkdownDescription: "URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"user_key": {
				MarkdownDescription: "User key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"web_hook_url": {
				MarkdownDescription: "Web hook url.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"server_url": {
				MarkdownDescription: "Server url.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"click_url": {
				MarkdownDescription: "Click URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"map_from": {
				MarkdownDescription: "Map From.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"map_to": {
				MarkdownDescription: "Map To.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"key": {
				MarkdownDescription: "Key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"event": {
				MarkdownDescription: "Event.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"device_ids": {
				MarkdownDescription: "Device IDs.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"channel_tags": {
				MarkdownDescription: "Channel tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"devices": {
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"topics": {
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"field_tags": {
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"recipients": {
				MarkdownDescription: "Recipients.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
		},
	}, nil
}

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Notification
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Notification current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Notification
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Notification current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *Notification) write(ctx context.Context, notification *radarr.NotificationOutput) {
	n.OnGrab = types.BoolValue(notification.OnGrab)
	n.OnDownload = types.BoolValue(notification.OnDownload)
	n.OnUpgrade = types.BoolValue(notification.OnUpgrade)
	n.OnRename = types.BoolValue(notification.OnRename)
	n.OnMovieAdded = types.BoolValue(notification.OnMovieAdded)
	n.OnMovieDelete = types.BoolValue(notification.OnMovieDelete)
	n.OnMovieFileDelete = types.BoolValue(notification.OnMovieFileDelete)
	n.OnMovieFileDeleteForUpgrade = types.BoolValue(notification.OnMovieFileDeleteForUpgrade)
	n.OnHealthIssue = types.BoolValue(notification.OnHealthIssue)
	n.OnApplicationUpdate = types.BoolValue(notification.OnApplicationUpdate)
	n.IncludeHealthWarnings = types.BoolValue(notification.IncludeHealthWarnings)
	n.ID = types.Int64Value(notification.ID)
	n.Name = types.StringValue(notification.Name)
	n.Implementation = types.StringValue(notification.Implementation)
	n.ConfigContract = types.StringValue(notification.ConfigContract)
	n.Tags = types.SetValueMust(types.Int64Type, nil)
	n.ChannelTags = types.SetValueMust(types.StringType, nil)
	n.DeviceIds = types.SetValueMust(types.Int64Type, nil)
	n.Topics = types.SetValueMust(types.StringType, nil)
	n.Devices = types.SetValueMust(types.StringType, nil)
	n.Recipients = types.SetValueMust(types.StringType, nil)
	n.FieldTags = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, notification.Tags, n.Tags.Type(ctx), &n.Tags)
	n.writeFields(ctx, notification.Fields)
}

func (n *Notification) writeFields(ctx context.Context, fields []*starr.FieldOutput) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(notificationStringFields, f.Name) {
			tools.WriteStringField(f, n)

			continue
		}

		if slices.Contains(notificationBoolFields, f.Name) {
			tools.WriteBoolField(f, n)

			continue
		}

		if slices.Contains(notificationIntFields, f.Name) {
			tools.WriteIntField(f, n)

			continue
		}

		if slices.Contains(notificationStringSliceFields, f.Name) {
			tools.WriteStringSliceField(ctx, f, n)
		}

		if slices.Contains(notificationIntSliceFields, f.Name) {
			tools.WriteIntSliceField(ctx, f, n)
		}
	}
}

func (n *Notification) read(ctx context.Context) *radarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &radarr.NotificationInput{
		OnGrab:                      n.OnGrab.ValueBool(),
		OnDownload:                  n.OnDownload.ValueBool(),
		OnUpgrade:                   n.OnUpgrade.ValueBool(),
		OnRename:                    n.OnRename.ValueBool(),
		OnMovieAdded:                n.OnMovieAdded.ValueBool(),
		OnMovieDelete:               n.OnMovieDelete.ValueBool(),
		OnMovieFileDelete:           n.OnMovieFileDelete.ValueBool(),
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade.ValueBool(),
		OnHealthIssue:               n.OnHealthIssue.ValueBool(),
		OnApplicationUpdate:         n.OnApplicationUpdate.ValueBool(),
		IncludeHealthWarnings:       n.IncludeHealthWarnings.ValueBool(),
		ID:                          n.ID.ValueInt64(),
		Name:                        n.Name.ValueString(),
		Implementation:              n.Implementation.ValueString(),
		ConfigContract:              n.ConfigContract.ValueString(),
		Tags:                        tags,
		Fields:                      n.readFields(ctx),
	}
}

func (n *Notification) readFields(ctx context.Context) []*starr.FieldInput {
	var output []*starr.FieldInput

	for _, b := range notificationBoolFields {
		if field := tools.ReadBoolField(b, n); field != nil {
			output = append(output, field)
		}
	}

	for _, i := range notificationIntFields {
		if field := tools.ReadIntField(i, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range notificationStringFields {
		if field := tools.ReadStringField(s, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range notificationStringSliceFields {
		if field := tools.ReadStringSliceField(ctx, s, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range notificationIntSliceFields {
		if field := tools.ReadIntSliceField(ctx, s, n); field != nil {
			output = append(output, field)
		}
	}

	return output
}
