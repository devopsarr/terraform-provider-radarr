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
	"golift.io/starr/radarr"
)

const (
	notificationCustomScriptResourceName   = "notification_custom_script"
	NotificationCustomScriptImplementation = "CustomScript"
	NotificationCustomScriptConfigContrat  = "CustomScriptSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationCustomScriptResource{}
var _ resource.ResourceWithImportState = &NotificationCustomScriptResource{}

func NewNotificationCustomScriptResource() resource.Resource {
	return &NotificationCustomScriptResource{}
}

// NotificationCustomScriptResource defines the notification implementation.
type NotificationCustomScriptResource struct {
	client *radarr.Radarr
}

// NotificationCustomScript describes the notification data model.
type NotificationCustomScript struct {
	Tags                        types.Set    `tfsdk:"tags"`
	Arguments                   types.String `tfsdk:"arguments"`
	Path                        types.String `tfsdk:"path"`
	Name                        types.String `tfsdk:"name"`
	ID                          types.Int64  `tfsdk:"id"`
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

func (n NotificationCustomScript) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		Path:                        n.Path,
		Arguments:                   n.Arguments,
		Name:                        n.Name,
		ID:                          n.ID,
		OnGrab:                      n.OnGrab,
		OnMovieAdded:                n.OnMovieAdded,
		OnMovieDelete:               n.OnMovieDelete,
		IncludeHealthWarnings:       n.IncludeHealthWarnings,
		OnApplicationUpdate:         n.OnApplicationUpdate,
		OnHealthIssue:               n.OnHealthIssue,
		OnMovieFileDelete:           n.OnMovieFileDelete,
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade,
		OnRename:                    n.OnRename,
		OnUpgrade:                   n.OnUpgrade,
		OnDownload:                  n.OnDownload,
	}
}

func (n *NotificationCustomScript) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Path = notification.Path
	n.Arguments = notification.Arguments
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

func (r *NotificationCustomScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationCustomScriptResourceName
}

func (r *NotificationCustomScriptResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Custom Script resource.\nFor more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Custom Script](https://wiki.servarr.com/radarr/supported#customscript).",
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
			"name": {
				MarkdownDescription: "NotificationCustomScript name.",
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
			"arguments": {
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"path": {
				MarkdownDescription: "Path.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *NotificationCustomScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationCustomScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationCustomScript
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationCustomScript current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationCustomScript
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationCustomScript current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationCustomScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationCustomScriptResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationCustomScript) write(ctx context.Context, notification *radarr.NotificationOutput) {
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

func (n *NotificationCustomScript) read(ctx context.Context) *radarr.NotificationInput {
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
		ConfigContract:              NotificationCustomScriptConfigContrat,
		Implementation:              NotificationCustomScriptImplementation,
		ID:                          n.ID.ValueInt64(),
		Name:                        n.Name.ValueString(),
		Tags:                        tags,
		Fields:                      n.toNotification().readFields(ctx),
	}
}
