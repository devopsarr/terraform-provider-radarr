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
	notificationBoxcarResourceName   = "notification_boxcar"
	notificationBoxcarImplementation = "Boxcar"
	notificationBoxcarConfigContract = "BoxcarSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationBoxcarResource{}
	_ resource.ResourceWithImportState = &NotificationBoxcarResource{}
)

func NewNotificationBoxcarResource() resource.Resource {
	return &NotificationBoxcarResource{}
}

// NotificationBoxcarResource defines the notification implementation.
type NotificationBoxcarResource struct {
	client *radarr.Radarr
}

// NotificationBoxcar describes the notification data model.
type NotificationBoxcar struct {
	Tags                        types.Set    `tfsdk:"tags"`
	Token                       types.String `tfsdk:"token"`
	Name                        types.String `tfsdk:"name"`
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

func (n NotificationBoxcar) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		Token:                       n.Token,
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

func (n *NotificationBoxcar) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Token = notification.Token
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

func (r *NotificationBoxcarResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationBoxcarResourceName
}

func (r *NotificationBoxcarResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Boxcar resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Boxcar](https://wiki.servarr.com/radarr/supported#boxcar).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_movie_added": schema.BoolAttribute{
				MarkdownDescription: "On movie added flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_movie_delete": schema.BoolAttribute{
				MarkdownDescription: "On movie delete flag.",
				Required:            true,
			},
			"on_movie_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_movie_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete for upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationBoxcar name.",
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
			"token": schema.StringAttribute{
				MarkdownDescription: "Token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationBoxcarResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationBoxcarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationBoxcar

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationBoxcar
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationBoxcarResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationBoxcarResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationBoxcarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationBoxcar

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationBoxcar current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationBoxcarResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationBoxcarResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationBoxcarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationBoxcar

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationBoxcar
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationBoxcarResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationBoxcarResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationBoxcarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationBoxcar

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationBoxcar current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationBoxcarResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationBoxcarResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationBoxcarResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationBoxcarResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationBoxcar) write(ctx context.Context, notification *radarr.NotificationOutput) {
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

func (n *NotificationBoxcar) read(ctx context.Context) *radarr.NotificationInput {
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
		ConfigContract:              notificationBoxcarConfigContract,
		Implementation:              notificationBoxcarImplementation,
		ID:                          n.ID.ValueInt64(),
		Name:                        n.Name.ValueString(),
		Tags:                        tags,
		Fields:                      n.toNotification().readFields(ctx),
	}
}
