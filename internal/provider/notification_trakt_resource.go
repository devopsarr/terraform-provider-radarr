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
	notificationTraktResourceName   = "notification_trakt"
	NotificationTraktImplementation = "Trakt"
	NotificationTraktConfigContrat  = "TraktSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationTraktResource{}
var _ resource.ResourceWithImportState = &NotificationTraktResource{}

func NewNotificationTraktResource() resource.Resource {
	return &NotificationTraktResource{}
}

// NotificationTraktResource defines the notification implementation.
type NotificationTraktResource struct {
	client *radarr.Radarr
}

// NotificationTrakt describes the notification data model.
type NotificationTrakt struct {
	Tags                        types.Set    `tfsdk:"tags"`
	AuthUser                    types.String `tfsdk:"auth_user"`
	AccessToken                 types.String `tfsdk:"access_token"`
	RefreshToken                types.String `tfsdk:"refresh_token"`
	Expires                     types.String `tfsdk:"expires"`
	Name                        types.String `tfsdk:"name"`
	ID                          types.Int64  `tfsdk:"id"`
	OnMovieFileDeleteForUpgrade types.Bool   `tfsdk:"on_movie_file_delete_for_upgrade"`
	OnMovieFileDelete           types.Bool   `tfsdk:"on_movie_file_delete"`
	OnMovieAdded                types.Bool   `tfsdk:"on_movie_added"`
	IncludeHealthWarnings       types.Bool   `tfsdk:"include_health_warnings"`
	OnMovieDelete               types.Bool   `tfsdk:"on_movie_delete"`
	OnUpgrade                   types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                  types.Bool   `tfsdk:"on_download"`
}

func (n NotificationTrakt) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		AuthUser:                    n.AuthUser,
		Name:                        n.Name,
		AccessToken:                 n.AccessToken,
		RefreshToken:                n.RefreshToken,
		ID:                          n.ID,
		Expires:                     n.Expires,
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade,
		OnMovieAdded:                n.OnMovieAdded,
		OnMovieFileDelete:           n.OnMovieFileDelete,
		IncludeHealthWarnings:       n.IncludeHealthWarnings,
		OnMovieDelete:               n.OnMovieDelete,
		OnUpgrade:                   n.OnUpgrade,
		OnDownload:                  n.OnDownload,
	}
}

func (n *NotificationTrakt) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.AuthUser = notification.AuthUser
	n.Name = notification.Name
	n.AccessToken = notification.AccessToken
	n.RefreshToken = notification.RefreshToken
	n.Expires = notification.Expires
	n.ID = notification.ID
	n.OnMovieFileDeleteForUpgrade = notification.OnMovieFileDeleteForUpgrade
	n.OnMovieFileDelete = notification.OnMovieFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnMovieAdded = notification.OnMovieAdded
	n.OnMovieDelete = notification.OnMovieDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationTraktResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationTraktResourceName
}

func (r *NotificationTraktResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Trakt resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Trakt](https://wiki.servarr.com/radarr/supported#trakt).",
		Attributes: map[string]schema.Attribute{
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
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationTrakt name.",
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
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access Token.",
				Required:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Access Token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Required:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "expires.",
				Computed:            true,
			},
		},
	}
}

func (r *NotificationTraktResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationTraktResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationTrakt
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationTraktResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTraktResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationTrakt current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationTraktResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTraktResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationTrakt
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationTraktResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTraktResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationTrakt current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationTraktResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationTraktResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationTraktResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationTrakt) write(ctx context.Context, notification *radarr.NotificationOutput) {
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

func (n *NotificationTrakt) read(ctx context.Context) *radarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &radarr.NotificationInput{
		OnDownload:                  n.OnDownload.ValueBool(),
		OnUpgrade:                   n.OnUpgrade.ValueBool(),
		OnMovieAdded:                n.OnMovieAdded.ValueBool(),
		OnMovieDelete:               n.OnMovieDelete.ValueBool(),
		OnMovieFileDelete:           n.OnMovieFileDelete.ValueBool(),
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade.ValueBool(),
		IncludeHealthWarnings:       n.IncludeHealthWarnings.ValueBool(),
		ConfigContract:              NotificationTraktConfigContrat,
		Implementation:              NotificationTraktImplementation,
		ID:                          n.ID.ValueInt64(),
		Name:                        n.Name.ValueString(),
		Tags:                        tags,
		Fields:                      n.toNotification().readFields(ctx),
	}
}