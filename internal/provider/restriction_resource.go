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

const restrictionName = "restriction"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RestrictionResource{}
	_ resource.ResourceWithImportState = &RestrictionResource{}
)

func NewRestrictionResource() resource.Resource {
	return &RestrictionResource{}
}

// RestrictionResource defines the remote path restriction implementation.
type RestrictionResource struct {
	client *radarr.Radarr
}

// Restriction describes the remote path restriction data model.
type Restriction struct {
	Tags     types.Set    `tfsdk:"tags"`
	Required types.String `tfsdk:"required"`
	Ignored  types.String `tfsdk:"ignored"`
	ID       types.Int64  `tfsdk:"id"`
}

func (r *RestrictionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Restriction resource.\nFor more information refer to [Restriction](https://wiki.servarr.com/radarr/settings#remote-path-restrictions) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Restriction ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"required": schema.StringAttribute{
				MarkdownDescription: "Required. Either one of 'required' or 'ignored' must be set.",
				Optional:            true,
				Computed:            true,
			},
			"ignored": schema.StringAttribute{
				MarkdownDescription: "Ignored. Either one of 'required' or 'ignored' must be set.",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *RestrictionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + restrictionName
}

func (r *RestrictionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RestrictionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var restriction *Restriction

	resp.Diagnostics.Append(req.Plan.Get(ctx, &restriction)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Restriction
	request := restriction.read(ctx)

	response, err := r.client.AddRestrictionContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", restrictionName, err))

		return
	}

	tflog.Trace(ctx, "created "+restrictionName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	restriction.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &restriction)...)
}

func (r *RestrictionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var restriction *Restriction

	resp.Diagnostics.Append(req.State.Get(ctx, &restriction)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get restriction current value
	response, err := r.client.GetRestrictionContext(ctx, restriction.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", restrictionName, err))

		return
	}

	tflog.Trace(ctx, "read "+restrictionName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	restriction.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &restriction)...)
}

func (r *RestrictionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var restriction *Restriction

	resp.Diagnostics.Append(req.Plan.Get(ctx, &restriction)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Restriction
	request := restriction.read(ctx)

	response, err := r.client.UpdateRestrictionContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", restrictionName, err))

		return
	}

	tflog.Trace(ctx, "updated "+restrictionName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	restriction.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &restriction)...)
}

func (r *RestrictionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *Restriction

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete restriction current value
	err := r.client.DeleteRestrictionContext(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", restrictionName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+restrictionName+": "+strconv.Itoa(int(state.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *RestrictionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+restrictionName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *Restriction) write(ctx context.Context, restriction *radarr.Restriction) {
	r.ID = types.Int64Value(restriction.ID)
	r.Ignored = types.StringValue(restriction.Ignored)
	r.Required = types.StringValue(restriction.Required)
	r.Tags = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, restriction.Tags, r.Tags.Type(ctx), &r.Tags)
}

func (r *Restriction) read(ctx context.Context) *radarr.Restriction {
	var tags []int

	tfsdk.ValueAs(ctx, r.Tags, &tags)

	return &radarr.Restriction{
		ID:       r.ID.ValueInt64(),
		Ignored:  r.Ignored.ValueString(),
		Required: r.Required.ValueString(),
		Tags:     tags,
	}
}
