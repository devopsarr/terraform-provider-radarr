package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	client *radarr.APIClient
}

// Restriction describes the remote path restriction data model.
type Restriction struct {
	Tags     types.Set    `tfsdk:"tags"`
	Required types.String `tfsdk:"required"`
	Ignored  types.String `tfsdk:"ignored"`
	ID       types.Int64  `tfsdk:"id"`
}

func (r Restriction) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"id":       types.Int64Type,
			"required": types.StringType,
			"ignored":  types.StringType,
			"tags":     types.SetType{}.WithElementType(types.Int64Type),
		})
}

func (r *RestrictionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *RestrictionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + restrictionName
}

func (r *RestrictionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *RestrictionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var restriction *Restriction

	resp.Diagnostics.Append(req.Plan.Get(ctx, &restriction)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Restriction
	request := restriction.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.RestrictionApi.CreateRestriction(ctx).RestrictionResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, restrictionName, err))

		return
	}

	tflog.Trace(ctx, "created "+restrictionName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	restriction.write(ctx, response, &resp.Diagnostics)
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
	response, _, err := r.client.RestrictionApi.GetRestrictionById(ctx, int32(restriction.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, restrictionName, err))

		return
	}

	tflog.Trace(ctx, "read "+restrictionName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	restriction.write(ctx, response, &resp.Diagnostics)
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
	request := restriction.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.RestrictionApi.UpdateRestriction(ctx, strconv.Itoa(int(request.GetId()))).RestrictionResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, restrictionName, err))

		return
	}

	tflog.Trace(ctx, "updated "+restrictionName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	restriction.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &restriction)...)
}

func (r *RestrictionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete restriction current value
	_, err := r.client.RestrictionApi.DeleteRestriction(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, restrictionName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+restrictionName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *RestrictionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+restrictionName+": "+req.ID)
}

func (r *Restriction) write(ctx context.Context, restriction *radarr.RestrictionResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	r.ID = types.Int64Value(int64(restriction.GetId()))
	r.Ignored = types.StringValue(restriction.GetIgnored())
	r.Required = types.StringValue(restriction.GetRequired())
	r.Tags, tempDiag = types.SetValueFrom(ctx, types.Int64Type, restriction.GetTags())
	diags.Append(tempDiag...)
}

func (r *Restriction) read(ctx context.Context, diags *diag.Diagnostics) *radarr.RestrictionResource {
	restriction := radarr.NewRestrictionResource()
	restriction.SetId(int32(r.ID.ValueInt64()))
	restriction.SetIgnored(r.Ignored.ValueString())
	restriction.SetRequired(r.Required.ValueString())
	diags.Append(r.Tags.ElementsAs(ctx, &restriction.Tags, true)...)

	return restriction
}
