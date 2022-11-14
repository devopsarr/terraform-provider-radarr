package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const notificationsDataSourceName = "notifications"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NotificationsDataSource{}

func NewNotificationsDataSource() datasource.DataSource {
	return &NotificationsDataSource{}
}

// NotificationsDataSource defines the notifications implementation.
type NotificationsDataSource struct {
	client *radarr.Radarr
}

// Notifications describes the notifications data model.
type Notifications struct {
	Notifications types.Set    `tfsdk:"notifications"`
	ID            types.String `tfsdk:"id"`
}

func (d *NotificationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationsDataSourceName
}

func (d *NotificationsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Notifications -->List all available [Notifications](../resources/notification).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"notifications": {
				MarkdownDescription: "Notification list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"on_grab": {
						MarkdownDescription: "On grab flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_download": {
						MarkdownDescription: "On download flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_upgrade": {
						MarkdownDescription: "On upgrade flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_rename": {
						MarkdownDescription: "On rename flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_movie_added": {
						MarkdownDescription: "On movie added flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_movie_delete": {
						MarkdownDescription: "On movie delete flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_movie_file_delete": {
						MarkdownDescription: "On movie file delete flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_movie_file_delete_for_upgrade": {
						MarkdownDescription: "On movie file delete for upgrade flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_health_issue": {
						MarkdownDescription: "On health issue flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"on_application_update": {
						MarkdownDescription: "On application update flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"include_health_warnings": {
						MarkdownDescription: "Include health warnings.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"config_contract": {
						MarkdownDescription: "Notification configuration template.",
						Computed:            true,
						Type:                types.StringType,
					},
					"implementation": {
						MarkdownDescription: "Notification implementation name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"name": {
						MarkdownDescription: "Notification name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tags": {
						MarkdownDescription: "List of associated tags.",
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
						Computed:            true,
						Type:                types.BoolType,
					},
					"clean_library": {
						MarkdownDescription: "Clean library flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"direct_message": {
						MarkdownDescription: "Direct message flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"notify": {
						MarkdownDescription: "Notify flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"require_encryption": {
						MarkdownDescription: "Require encryption flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"send_silently": {
						MarkdownDescription: "Add silently flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"update_library": {
						MarkdownDescription: "Update library flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"use_eu_endpoint": {
						MarkdownDescription: "Use EU endpoint flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"use_ssl": {
						MarkdownDescription: "Use SSL flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"port": {
						MarkdownDescription: "Port.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"grab_fields": {
						MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart, `10` CustomFormats, `11` CustomFormatScore.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"import_fields": {
						MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"method": {
						MarkdownDescription: "Method. `1` POST, `2` PUT.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"priority": {
						MarkdownDescription: "Priority.", // TODO: add values in description
						Computed:            true,
						Type:                types.Int64Type,
					},
					"retry": {
						MarkdownDescription: "Retry.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"expire": {
						MarkdownDescription: "Expire.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"access_token": {
						MarkdownDescription: "Access token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"access_token_secret": {
						MarkdownDescription: "Access token secret.",
						Computed:            true,
						Type:                types.StringType,
					},
					"api_key": {
						MarkdownDescription: "API key.",
						Computed:            true,
						Type:                types.StringType,
					},
					"app_token": {
						MarkdownDescription: "App token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"arguments": {
						MarkdownDescription: "Arguments.",
						Computed:            true,
						Type:                types.StringType,
					},
					"author": {
						MarkdownDescription: "Author.",
						Computed:            true,
						Type:                types.StringType,
					},
					"auth_token": {
						MarkdownDescription: "Auth token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"auth_user": {
						MarkdownDescription: "Auth user.",
						Computed:            true,
						Type:                types.StringType,
					},
					"avatar": {
						MarkdownDescription: "Avatar.",
						Computed:            true,
						Type:                types.StringType,
					},
					"instance_name": {
						MarkdownDescription: "Instance name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"bcc": {
						MarkdownDescription: "Bcc.",
						Computed:            true,
						Type:                types.StringType,
					},
					"bot_token": {
						MarkdownDescription: "Bot token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"cc": {
						MarkdownDescription: "Cc.",
						Computed:            true,
						Type:                types.StringType,
					},
					"channel": {
						MarkdownDescription: "Channel.",
						Computed:            true,
						Type:                types.StringType,
					},
					"chat_id": {
						MarkdownDescription: "Chat ID.",
						Computed:            true,
						Type:                types.StringType,
					},
					"consumer_key": {
						MarkdownDescription: "Consumer key.",
						Computed:            true,
						Type:                types.StringType,
					},
					"consumer_secret": {
						MarkdownDescription: "Consumer secret.",
						Computed:            true,
						Type:                types.StringType,
					},
					"device_names": {
						MarkdownDescription: "Device names.",
						Computed:            true,
						Type:                types.StringType,
					},
					"display_time": {
						MarkdownDescription: "Display time.",
						Computed:            true,
						Type:                types.StringType,
					},
					"expires": {
						MarkdownDescription: "Expires.",
						Computed:            true,
						Type:                types.StringType,
					},
					"from": {
						MarkdownDescription: "From.",
						Computed:            true,
						Type:                types.StringType,
					},
					"host": {
						MarkdownDescription: "Host.",
						Computed:            true,
						Type:                types.StringType,
					},
					"icon": {
						MarkdownDescription: "Icon.",
						Computed:            true,
						Type:                types.StringType,
					},
					"mention": {
						MarkdownDescription: "Mention.",
						Computed:            true,
						Type:                types.StringType,
					},
					"password": {
						MarkdownDescription: "password.",
						Computed:            true,
						Type:                types.StringType,
					},
					"path": {
						MarkdownDescription: "Path.",
						Computed:            true,
						Type:                types.StringType,
					},
					"refresh_token": {
						MarkdownDescription: "Refresh token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"sender_domain": {
						MarkdownDescription: "Sender domain.",
						Computed:            true,
						Type:                types.StringType,
					},
					"sender_id": {
						MarkdownDescription: "Sender ID.",
						Computed:            true,
						Type:                types.StringType,
					},
					"server": {
						MarkdownDescription: "server.",
						Computed:            true,
						Type:                types.StringType,
					},
					"sign_in": {
						MarkdownDescription: "Sign in.",
						Computed:            true,
						Type:                types.StringType,
					},
					"sound": {
						MarkdownDescription: "Sound.",
						Computed:            true,
						Type:                types.StringType,
					},
					"to": {
						MarkdownDescription: "To.",
						Computed:            true,
						Type:                types.StringType,
					},
					"token": {
						MarkdownDescription: "Token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"url": {
						MarkdownDescription: "URL.",
						Computed:            true,
						Type:                types.StringType,
					},
					"user_key": {
						MarkdownDescription: "User key.",
						Computed:            true,
						Type:                types.StringType,
					},
					"username": {
						MarkdownDescription: "Username.",
						Computed:            true,
						Type:                types.StringType,
					},
					"web_hook_url": {
						MarkdownDescription: "Web hook url.",
						Computed:            true,
						Type:                types.StringType,
					},
					"server_url": {
						MarkdownDescription: "Server url.",
						Computed:            true,
						Type:                types.StringType,
					},
					"click_url": {
						MarkdownDescription: "Click URL.",
						Computed:            true,
						Type:                types.StringType,
					},
					"map_from": {
						MarkdownDescription: "Map From.",
						Computed:            true,
						Type:                types.StringType,
					},
					"map_to": {
						MarkdownDescription: "Map To.",
						Computed:            true,
						Type:                types.StringType,
					},
					"key": {
						MarkdownDescription: "Key.",
						Computed:            true,
						Type:                types.StringType,
					},
					"event": {
						MarkdownDescription: "Event.",
						Computed:            true,
						Type:                types.StringType,
					},
					"device_ids": {
						MarkdownDescription: "Device IDs.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
					"channel_tags": {
						MarkdownDescription: "Channel tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
					"devices": {
						MarkdownDescription: "Devices.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
					"topics": {
						MarkdownDescription: "Devices.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
					"field_tags": {
						MarkdownDescription: "Devices.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
					"recipients": {
						MarkdownDescription: "Recipients.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
				}),
			},
		},
	}, nil
}

func (d *NotificationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NotificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Notifications

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get notifications current value
	response, err := d.client.GetNotificationsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationsDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]Notification, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.Notifications.Type(context.Background()), &data.Notifications)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
