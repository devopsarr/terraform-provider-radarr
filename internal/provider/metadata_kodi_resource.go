package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
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
	metadataKodiResourceName   = "metadata_kodi"
	metadataKodiImplementation = "XbmcMetadata"
	metadataKodiConfigContract = "XbmcMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataKodiResource{}
	_ resource.ResourceWithImportState = &MetadataKodiResource{}
)

func NewMetadataKodiResource() resource.Resource {
	return &MetadataKodiResource{}
}

// MetadataKodiResource defines the Kodi metadata implementation.
type MetadataKodiResource struct {
	client *radarr.APIClient
}

// MetadataKodi describes the Kodi metadata data model.
type MetadataKodi struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	ID                    types.Int64  `tfsdk:"id"`
	MovieMetadataLanguage types.Int64  `tfsdk:"movie_metadata_language"`
	Enable                types.Bool   `tfsdk:"enable"`
	MovieMetadata         types.Bool   `tfsdk:"movie_metadata"`
	MovieMetadataURL      types.Bool   `tfsdk:"movie_metadata_url"`
	MovieImages           types.Bool   `tfsdk:"movie_images"`
	UseMovieNfo           types.Bool   `tfsdk:"use_movie_nfo"`
	AddCollectionName     types.Bool   `tfsdk:"add_collection_name"`
}

func (m MetadataKodi) toMetadata() *Metadata {
	return &Metadata{
		Tags:                  m.Tags,
		Name:                  m.Name,
		ID:                    m.ID,
		MovieMetadataLanguage: m.MovieMetadataLanguage,
		Enable:                m.Enable,
		MovieMetadata:         m.MovieMetadata,
		MovieMetadataURL:      m.MovieMetadataURL,
		MovieImages:           m.MovieImages,
		UseMovieNfo:           m.UseMovieNfo,
		AddCollectionName:     m.AddCollectionName,
	}
}

func (m *MetadataKodi) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.MovieMetadataLanguage = metadata.MovieMetadataLanguage
	m.Enable = metadata.Enable
	m.MovieMetadata = metadata.MovieMetadata
	m.MovieImages = metadata.MovieImages
	m.MovieMetadataURL = metadata.MovieMetadataURL
	m.UseMovieNfo = metadata.UseMovieNfo
	m.AddCollectionName = metadata.AddCollectionName
}

func (r *MetadataKodiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataKodiResourceName
}

func (r *MetadataKodiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Kodi resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/radarr/settings#metadata) and [KODI](https://wiki.servarr.com/radarr/supported#xbmcmetadata).",
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
			"add_collection_name": schema.BoolAttribute{
				MarkdownDescription: "Add collection name flag.",
				Required:            true,
			},
			"use_movie_nfo": schema.BoolAttribute{
				MarkdownDescription: "Use movie nfo flag.",
				Required:            true,
			},
			"movie_images": schema.BoolAttribute{
				MarkdownDescription: "Movie images flag.",
				Required:            true,
			},
			"movie_metadata": schema.BoolAttribute{
				MarkdownDescription: "Movie metafata flag.",
				Required:            true,
			},
			"movie_metadata_url": schema.BoolAttribute{
				MarkdownDescription: "Movie metadata URL flag.",
				Required:            true,
			},
			"movie_metadata_language": schema.Int64Attribute{
				MarkdownDescription: "Movie metadata language.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataKodiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataKodiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataKodi
	request := metadata.read(ctx)

	response, _, err := r.client.MetadataApi.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataKodi current value
	response, _, err := r.client.MetadataApi.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataKodi
	request := metadata.read(ctx)

	response, _, err := r.client.MetadataApi.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+metadataKodiResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataKodi current value
	_, err := r.client.MetadataApi.DeleteMetadata(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataKodiResourceName+": "+strconv.Itoa(int(metadata.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataKodiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataKodiResourceName+": "+req.ID)
}

func (m *MetadataKodi) write(ctx context.Context, metadata *radarr.MetadataResource) {
	genericMetadata := Metadata{
		Name:   types.StringValue(metadata.GetName()),
		ID:     types.Int64Value(int64(metadata.GetId())),
		Enable: types.BoolValue(metadata.GetEnable()),
	}
	genericMetadata.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, metadata.Tags)
	genericMetadata.writeFields(metadata.Fields)
	m.fromMetadata(&genericMetadata)
}

func (m *MetadataKodi) read(ctx context.Context) *radarr.MetadataResource {
	tags := make([]*int32, len(m.Tags.Elements()))
	tfsdk.ValueAs(ctx, m.Tags, &tags)

	metadata := radarr.NewMetadataResource()
	metadata.SetEnable(m.Enable.ValueBool())
	metadata.SetId(int32(m.ID.ValueInt64()))
	metadata.SetConfigContract(metadataKodiConfigContract)
	metadata.SetImplementation(metadataKodiImplementation)
	metadata.SetName(m.Name.ValueString())
	metadata.SetTags(tags)
	metadata.SetFields(m.toMetadata().readFields())

	return metadata
}
