package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationSimplepushResourceName   = "notification_simplepush"
	notificationSimplepushImplementation = "Simplepush"
	notificationSimplepushConfigContract = "SimplepushSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSimplepushResource{}
	_ resource.ResourceWithImportState = &NotificationSimplepushResource{}
)

func NewNotificationSimplepushResource() resource.Resource {
	return &NotificationSimplepushResource{}
}

// NotificationSimplepushResource defines the notification implementation.
type NotificationSimplepushResource struct {
	client *radarr.APIClient
}

// NotificationSimplepush describes the notification data model.
type NotificationSimplepush struct {
	Tags                        types.Set    `tfsdk:"tags"`
	Event                       types.String `tfsdk:"event"`
	Name                        types.String `tfsdk:"name"`
	Key                         types.String `tfsdk:"key"`
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

func (n NotificationSimplepush) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		Event:                       n.Event,
		Key:                         n.Key,
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

func (n *NotificationSimplepush) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Event = notification.Event
	n.Key = notification.Key
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

func (r *NotificationSimplepushResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSimplepushResourceName
}

func (r *NotificationSimplepushResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Simplepush resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Simplepush](https://wiki.servarr.com/radarr/supported#simplepush).",
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
				MarkdownDescription: "NotificationSimplepush name.",
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
			"event": schema.StringAttribute{
				MarkdownDescription: "Event.",
				Optional:            true,
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationSimplepushResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *radarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationSimplepushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSimplepush
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSimplepush current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSimplepush
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSimplepush current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSimplepushResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSimplepushResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationSimplepushResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationSimplepush) write(ctx context.Context, notification *radarr.NotificationResource) {
	genericNotification := Notification{
		OnGrab:                      types.BoolValue(notification.GetOnGrab()),
		OnDownload:                  types.BoolValue(notification.GetOnDownload()),
		OnUpgrade:                   types.BoolValue(notification.GetOnUpgrade()),
		OnMovieAdded:                types.BoolValue(notification.GetOnMovieAdded()),
		OnMovieDelete:               types.BoolValue(notification.GetOnMovieDelete()),
		OnMovieFileDelete:           types.BoolValue(notification.GetOnMovieFileDelete()),
		OnMovieFileDeleteForUpgrade: types.BoolValue(notification.GetOnMovieFileDeleteForUpgrade()),
		OnHealthIssue:               types.BoolValue(notification.GetOnHealthIssue()),
		OnApplicationUpdate:         types.BoolValue(notification.GetOnApplicationUpdate()),
		IncludeHealthWarnings:       types.BoolValue(notification.GetIncludeHealthWarnings()),
		ID:                          types.Int64Value(int64(notification.GetId())),
		Name:                        types.StringValue(notification.GetName()),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationSimplepush) read(ctx context.Context) *radarr.NotificationResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	notification := radarr.NewNotificationResource()
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetOnDownload(n.OnDownload.ValueBool())
	notification.SetOnUpgrade(n.OnUpgrade.ValueBool())
	notification.SetOnMovieAdded(n.OnMovieAdded.ValueBool())
	notification.SetOnMovieDelete(n.OnMovieDelete.ValueBool())
	notification.SetOnMovieFileDelete(n.OnMovieFileDelete.ValueBool())
	notification.SetOnMovieFileDeleteForUpgrade(n.OnMovieFileDeleteForUpgrade.ValueBool())
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnApplicationUpdate(n.OnApplicationUpdate.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetConfigContract(notificationSimplepushConfigContract)
	notification.SetImplementation(notificationSimplepushImplementation)
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetTags(tags)
	notification.SetFields(n.toNotification().readFields(ctx))

	return notification
}
