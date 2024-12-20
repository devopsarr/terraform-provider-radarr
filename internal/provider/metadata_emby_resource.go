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
	metadataEmbyResourceName   = "metadata_emby"
	metadataEmbyImplementation = "MediaBrowserMetadata"
	metadataEmbyConfigContract = "MediaBrowserMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataEmbyResource{}
	_ resource.ResourceWithImportState = &MetadataEmbyResource{}
)

func NewMetadataEmbyResource() resource.Resource {
	return &MetadataEmbyResource{}
}

// MetadataEmbyResource defines the Emby metadata implementation.
type MetadataEmbyResource struct {
	client *radarr.APIClient
	auth   context.Context
}

// MetadataEmby describes the Emby metadata data model.
type MetadataEmby struct {
	Tags          types.Set    `tfsdk:"tags"`
	Name          types.String `tfsdk:"name"`
	ID            types.Int64  `tfsdk:"id"`
	Enable        types.Bool   `tfsdk:"enable"`
	MovieMetadata types.Bool   `tfsdk:"movie_metadata"`
}

func (m MetadataEmby) toMetadata() *Metadata {
	return &Metadata{
		Tags:           m.Tags,
		Name:           m.Name,
		ID:             m.ID,
		Enable:         m.Enable,
		MovieMetadata:  m.MovieMetadata,
		Implementation: types.StringValue(metadataEmbyImplementation),
		ConfigContract: types.StringValue(metadataEmbyConfigContract),
	}
}

func (m *MetadataEmby) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.MovieMetadata = metadata.MovieMetadata
}

func (r *MetadataEmbyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataEmbyResourceName
}

func (r *MetadataEmbyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->\nMetadata Emby resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/radarr/settings#metadata) and [Emby](https://wiki.servarr.com/radarr/supported#mediabrowsermetadata).",
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
		},
	}
}

func (r *MetadataEmbyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *MetadataEmbyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataEmby

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataEmby
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.CreateMetadata(r.auth).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataEmbyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataEmbyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataEmby

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataEmby current value
	response, _, err := r.client.MetadataAPI.GetMetadataById(r.auth, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataEmbyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataEmbyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataEmby

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataEmby
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.UpdateMetadata(r.auth, request.GetId()).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataEmbyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataEmbyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataEmby current value
	_, err := r.client.MetadataAPI.DeleteMetadata(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataEmbyResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataEmbyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataEmbyResourceName+": "+req.ID)
}

func (m *MetadataEmby) write(ctx context.Context, metadata *radarr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataEmby) read(ctx context.Context, diags *diag.Diagnostics) *radarr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}
