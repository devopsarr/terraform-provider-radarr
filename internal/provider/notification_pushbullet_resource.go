package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const (
	notificationPushbulletResourceName   = "notification_pushbullet"
	NotificationPushbulletImplementation = "PushBullet"
	NotificationPushbulletConfigContrat  = "PushBulletSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationPushbulletResource{}
var _ resource.ResourceWithImportState = &NotificationPushbulletResource{}

func NewNotificationPushbulletResource() resource.Resource {
	return &NotificationPushbulletResource{}
}

// NotificationPushbulletResource defines the notification implementation.
type NotificationPushbulletResource struct {
	client *radarr.Radarr
}

// NotificationPushbullet describes the notification data model.
type NotificationPushbullet struct {
	Tags                        types.Set    `tfsdk:"tags"`
	DeviceIds                   types.Set    `tfsdk:"device_ids"`
	ChannelTags                 types.Set    `tfsdk:"channel_tags"`
	SenderID                    types.String `tfsdk:"sender_id"`
	Name                        types.String `tfsdk:"name"`
	APIKey                      types.String `tfsdk:"api_key"`
	ID                          types.Int64  `tfsdk:"id"`
	OnGrab                      types.Bool   `tfsdk:"on_grab"`
	OnMovieFileDeleteForUpgrade types.Bool   `tfsdk:"on_movie_file_delete_for_upgrade"`
	OnMovieFileDelete           types.Bool   `tfsdk:"on_movie_file_delete"`
	OnMovieAdded                types.Bool   `tfsdk:"on_movie_added"`
	IncludeHealthWarnings       types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate         types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue               types.Bool   `tfsdk:"on_health_issue"`
	OnMovieDelete               types.Bool   `tfsdk:"on_movie_delete"`
	OnUpgrade                   types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                  types.Bool   `tfsdk:"on_download"`
}

func (n NotificationPushbullet) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		DeviceIds:                   n.DeviceIds,
		ChannelTags:                 n.ChannelTags,
		SenderID:                    n.SenderID,
		APIKey:                      n.APIKey,
		Name:                        n.Name,
		ID:                          n.ID,
		OnGrab:                      n.OnGrab,
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade,
		OnMovieAdded:                n.OnMovieAdded,
		OnMovieFileDelete:           n.OnMovieFileDelete,
		IncludeHealthWarnings:       n.IncludeHealthWarnings,
		OnApplicationUpdate:         n.OnApplicationUpdate,
		OnHealthIssue:               n.OnHealthIssue,
		OnMovieDelete:               n.OnMovieDelete,
		OnUpgrade:                   n.OnUpgrade,
		OnDownload:                  n.OnDownload,
	}
}

func (n *NotificationPushbullet) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.DeviceIds = notification.DeviceIds
	n.ChannelTags = notification.ChannelTags
	n.SenderID = notification.SenderID
	n.APIKey = notification.APIKey
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnMovieFileDeleteForUpgrade = notification.OnMovieFileDeleteForUpgrade
	n.OnMovieFileDelete = notification.OnMovieFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnMovieAdded = notification.OnMovieAdded
	n.OnMovieDelete = notification.OnMovieDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationPushbulletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationPushbulletResourceName
}

func (r *NotificationPushbulletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Pushbullet resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Pushbullet](https://wiki.servarr.com/radarr/supported#pushbullet).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_movie_added": schema.BoolAttribute{
				MarkdownDescription: "On movie added flag.",
				Required:            true,
			},
			"on_movie_delete": schema.BoolAttribute{
				MarkdownDescription: "On movie delete flag.",
				Required:            true,
			},
			"on_movie_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete flag.",
				Required:            true,
			},
			"on_movie_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete for upgrade flag.",
				Required:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Required:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationPushbullet name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"sender_id": schema.StringAttribute{
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "List of devices IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"channel_tags": schema.SetAttribute{
				MarkdownDescription: "List of channel tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationPushbulletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationPushbulletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationPushbullet
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationPushbulletResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationPushbulletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationPushbullet current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationPushbulletResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationPushbulletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationPushbullet
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationPushbulletResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationPushbulletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationPushbullet current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationPushbulletResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationPushbulletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationPushbulletResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationPushbullet) write(ctx context.Context, notification *radarr.NotificationOutput) {
	genericNotification := Notification{
		OnGrab:                      types.BoolValue(notification.OnGrab),
		OnDownload:                  types.BoolValue(notification.OnDownload),
		OnUpgrade:                   types.BoolValue(notification.OnUpgrade),
		OnMovieAdded:                types.BoolValue(notification.OnMovieAdded),
		OnMovieDelete:               types.BoolValue(notification.OnMovieDelete),
		OnMovieFileDelete:           types.BoolValue(notification.OnMovieFileDelete),
		OnMovieFileDeleteForUpgrade: types.BoolValue(notification.OnMovieFileDeleteForUpgrade),
		OnHealthIssue:               types.BoolValue(notification.OnHealthIssue),
		OnApplicationUpdate:         types.BoolValue(notification.OnApplicationUpdate),
		IncludeHealthWarnings:       types.BoolValue(notification.IncludeHealthWarnings),
		ID:                          types.Int64Value(notification.ID),
		Name:                        types.StringValue(notification.Name),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationPushbullet) read(ctx context.Context) *radarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &radarr.NotificationInput{
		OnGrab:                      n.OnGrab.ValueBool(),
		OnDownload:                  n.OnDownload.ValueBool(),
		OnUpgrade:                   n.OnUpgrade.ValueBool(),
		OnMovieAdded:                n.OnMovieAdded.ValueBool(),
		OnMovieDelete:               n.OnMovieDelete.ValueBool(),
		OnMovieFileDelete:           n.OnMovieFileDelete.ValueBool(),
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade.ValueBool(),
		OnHealthIssue:               n.OnHealthIssue.ValueBool(),
		OnApplicationUpdate:         n.OnApplicationUpdate.ValueBool(),
		IncludeHealthWarnings:       n.IncludeHealthWarnings.ValueBool(),
		ConfigContract:              NotificationPushbulletConfigContrat,
		Implementation:              NotificationPushbulletImplementation,
		ID:                          n.ID.ValueInt64(),
		Name:                        n.Name.ValueString(),
		Tags:                        tags,
		Fields:                      n.toNotification().readFields(ctx),
	}
}
