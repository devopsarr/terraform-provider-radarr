package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	metadataRoksboxResourceName   = "metadata_roksbox"
	metadataRoksboxImplementation = "RoksboxMetadata"
	metadataRoksboxConfigContract = "RoksboxMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataRoksboxResource{}
	_ resource.ResourceWithImportState = &MetadataRoksboxResource{}
)

func NewMetadataRoksboxResource() resource.Resource {
	return &MetadataRoksboxResource{}
}

// MetadataRoksboxResource defines the Roksbox metadata implementation.
type MetadataRoksboxResource struct {
	client *radarr.APIClient
	auth   context.Context
}

// MetadataRoksbox describes the Roksbox metadata data model.
type MetadataRoksbox struct {
	Tags          types.Set    `tfsdk:"tags"`
	Name          types.String `tfsdk:"name"`
	ID            types.Int64  `tfsdk:"id"`
	Enable        types.Bool   `tfsdk:"enable"`
	MovieMetadata types.Bool   `tfsdk:"movie_metadata"`
	MovieImages   types.Bool   `tfsdk:"movie_images"`
}

func (m MetadataRoksbox) toMetadata() *Metadata {
	return &Metadata{
		Tags:           m.Tags,
		Name:           m.Name,
		ID:             m.ID,
		Enable:         m.Enable,
		MovieMetadata:  m.MovieMetadata,
		MovieImages:    m.MovieImages,
		Implementation: types.StringValue(metadataRoksboxImplementation),
		ConfigContract: types.StringValue(metadataRoksboxConfigContract),
	}
}

func (m *MetadataRoksbox) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.MovieMetadata = metadata.MovieMetadata
	m.MovieImages = metadata.MovieImages
}

func (r *MetadataRoksboxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataRoksboxResourceName
}

func (r *MetadataRoksboxResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->\nMetadata Roksbox resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/radarr/settings#metadata) and [ROKSBOX](https://wiki.servarr.com/radarr/supported#roksboxmetadata).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"movie_metadata": schema.BoolAttribute{
				MarkdownDescription: "Movie metadata flag.",
				Required:            true,
			},
			"movie_images": schema.BoolAttribute{
				MarkdownDescription: "Movie images flag.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataRoksboxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *MetadataRoksboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataRoksbox
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.CreateMetadata(r.auth).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataRoksbox current value
	response, _, err := r.client.MetadataAPI.GetMetadataById(r.auth, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataRoksbox
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.UpdateMetadata(r.auth, request.GetId()).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataRoksbox current value
	_, err := r.client.MetadataAPI.DeleteMetadata(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataRoksboxResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataRoksboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataRoksboxResourceName+": "+req.ID)
}

func (m *MetadataRoksbox) write(ctx context.Context, metadata *radarr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataRoksbox) read(ctx context.Context, diags *diag.Diagnostics) *radarr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}
