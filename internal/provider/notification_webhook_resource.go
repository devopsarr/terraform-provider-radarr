package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const (
	notificationWebhookResourceName   = "notification_webhook"
	NotificationWebhookImplementation = "Webhook"
	NotificationWebhookConfigContrat  = "WebhookSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationWebhookResource{}
var _ resource.ResourceWithImportState = &NotificationWebhookResource{}

func NewNotificationWebhookResource() resource.Resource {
	return &NotificationWebhookResource{}
}

// NotificationWebhookResource defines the notification implementation.
type NotificationWebhookResource struct {
	client *radarr.Radarr
}

// NotificationWebhook describes the notification data model.
type NotificationWebhook struct {
	Tags                        types.Set    `tfsdk:"tags"`
	URL                         types.String `tfsdk:"url"`
	Name                        types.String `tfsdk:"name"`
	Username                    types.String `tfsdk:"username"`
	Password                    types.String `tfsdk:"password"`
	ID                          types.Int64  `tfsdk:"id"`
	Method                      types.Int64  `tfsdk:"method"`
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
}

func (n NotificationWebhook) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		URL:                         n.URL,
		Method:                      n.Method,
		Username:                    n.Username,
		Password:                    n.Password,
		Name:                        n.Name,
		ID:                          n.ID,
		OnGrab:                      n.OnGrab,
		OnMovieAdded:                n.OnMovieAdded,
		OnMovieFileDelete:           n.OnMovieFileDelete,
		OnMovieDelete:               n.OnMovieDelete,
		IncludeHealthWarnings:       n.IncludeHealthWarnings,
		OnApplicationUpdate:         n.OnApplicationUpdate,
		OnHealthIssue:               n.OnHealthIssue,
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade,
		OnRename:                    n.OnRename,
		OnUpgrade:                   n.OnUpgrade,
		OnDownload:                  n.OnDownload,
	}
}

func (n *NotificationWebhook) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.URL = notification.URL
	n.Method = notification.Method
	n.Username = notification.Username
	n.Password = notification.Password
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
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationWebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationWebhookResourceName
}

func (r *NotificationWebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Webhook resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Webhook](https://wiki.servarr.com/radarr/supported#webhook).",
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
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
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
				MarkdownDescription: "NotificationWebhook name.",
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
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password.",
				Optional:            true,
				Computed:            true,
			},
			"method": schema.Int64Attribute{
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
		},
	}
}

func (r *NotificationWebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationWebhook
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationWebhookResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationWebhook current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationWebhookResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationWebhook
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationWebhookResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationWebhook current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationWebhookResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationWebhookResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationWebhook) write(ctx context.Context, notification *radarr.NotificationOutput) {
	genericNotification := Notification{
		OnGrab:                      types.BoolValue(notification.OnGrab),
		OnDownload:                  types.BoolValue(notification.OnDownload),
		OnUpgrade:                   types.BoolValue(notification.OnUpgrade),
		OnRename:                    types.BoolValue(notification.OnRename),
		OnMovieAdded:                types.BoolValue(notification.OnMovieAdded),
		OnMovieDelete:               types.BoolValue(notification.OnMovieDelete),
		OnMovieFileDelete:           types.BoolValue(notification.OnMovieFileDelete),
		OnMovieFileDeleteForUpgrade: types.BoolValue(notification.OnMovieFileDeleteForUpgrade),
		OnHealthIssue:               types.BoolValue(notification.OnHealthIssue),
		OnApplicationUpdate:         types.BoolValue(notification.OnApplicationUpdate),
		IncludeHealthWarnings:       types.BoolValue(notification.IncludeHealthWarnings),
		ID:                          types.Int64Value(notification.ID),
		Name:                        types.StringValue(notification.Name),
		Tags:                        types.SetValueMust(types.Int64Type, nil),
	}
	tfsdk.ValueFrom(ctx, notification.Tags, genericNotification.Tags.Type(ctx), &genericNotification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationWebhook) read(ctx context.Context) *radarr.NotificationInput {
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
		ConfigContract:              NotificationWebhookConfigContrat,
		Implementation:              NotificationWebhookImplementation,
		ID:                          n.ID.ValueInt64(),
		Name:                        n.Name.ValueString(),
		Tags:                        tags,
		Fields:                      n.toNotification().readFields(ctx),
	}
}
