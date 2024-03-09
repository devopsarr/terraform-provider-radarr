package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListTMDBPopularResourceName   = "import_list_tmdb_popular"
	importListTMDBPopularImplementation = "TMDbPopularImport"
	importListTMDBPopularConfigContract = "TMDbPopularSettings"
	importListTMDBPopularType           = "tmdb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListTMDBPopularResource{}
	_ resource.ResourceWithImportState = &ImportListTMDBPopularResource{}
)

func NewImportListTMDBPopularResource() resource.Resource {
	return &ImportListTMDBPopularResource{}
}

// ImportListTMDBPopularResource defines the import list implementation.
type ImportListTMDBPopularResource struct {
	client *radarr.APIClient
}

// ImportListTMDBPopular describes the import list data model.
type ImportListTMDBPopular struct {
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	Monitor             types.String `tfsdk:"monitor"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	RootFolderPath      types.String `tfsdk:"root_folder_path"`
	MinVoteAverage      types.String `tfsdk:"min_vote_average"`
	MinVotes            types.String `tfsdk:"min_votes"`
	TMDBCertification   types.String `tfsdk:"tmdb_certification"`
	IncludeGenreIDs     types.String `tfsdk:"include_genre_ids"`
	ExcludeGenreIDs     types.String `tfsdk:"exclude_genre_ids"`
	LanguageCode        types.Int64  `tfsdk:"language_code"`
	TMDBListType        types.Int64  `tfsdk:"tmdb_list_type"`
	ListOrder           types.Int64  `tfsdk:"list_order"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	EnableAuto          types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd         types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListTMDBPopular) toImportList() *ImportList {
	return &ImportList{
		Tags:                i.Tags,
		Name:                i.Name,
		Monitor:             i.Monitor,
		MinimumAvailability: i.MinimumAvailability,
		RootFolderPath:      i.RootFolderPath,
		TMDBListType:        i.TMDBListType,
		MinVoteAverage:      i.MinVoteAverage,
		MinVotes:            i.MinVotes,
		TMDBCertification:   i.TMDBCertification,
		IncludeGenreIDs:     i.IncludeGenreIDs,
		ExcludeGenreIDs:     i.ExcludeGenreIDs,
		LanguageCode:        i.LanguageCode,
		ListOrder:           i.ListOrder,
		ID:                  i.ID,
		QualityProfileID:    i.QualityProfileID,
		Enabled:             i.Enabled,
		EnableAuto:          i.EnableAuto,
		SearchOnAdd:         i.SearchOnAdd,
		Implementation:      types.StringValue(importListTMDBPopularImplementation),
		ConfigContract:      types.StringValue(importListTMDBPopularConfigContract),
		ListType:            types.StringValue(importListTMDBPopularType),
	}
}

func (i *ImportListTMDBPopular) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.Monitor = importList.Monitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.TMDBListType = importList.TMDBListType
	i.MinVoteAverage = importList.MinVoteAverage
	i.MinVotes = importList.MinVotes
	i.TMDBCertification = importList.TMDBCertification
	i.IncludeGenreIDs = importList.IncludeGenreIDs
	i.ExcludeGenreIDs = importList.ExcludeGenreIDs
	i.LanguageCode = importList.LanguageCode
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListTMDBPopularResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListTMDBPopularResourceName
}

func (r *ImportListTMDBPopularResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List TMDB Popular resource.\nFor more information refer to [Import List](https://wiki.servarr.com/radarr/settings#import-lists) and [TMDB Popular](https://wiki.servarr.com/radarr/supported#tmdbpopularimport).",
		Attributes: map[string]schema.Attribute{
			"enable_auto": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Optional:            true,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled flag.",
				Optional:            true,
				Computed:            true,
			},
			"search_on_add": schema.BoolAttribute{
				MarkdownDescription: "Search on add flag.",
				Optional:            true,
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"list_order": schema.Int64Attribute{
				MarkdownDescription: "List order.",
				Optional:            true,
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Required:            true,
			},
			"monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("movieOnly", "movieAndCollection", "none"),
				},
			},
			"minimum_availability": schema.StringAttribute{
				MarkdownDescription: "Minimum availability.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("tba", "announced", "inCinemas", "released", "deleted"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"tmdb_list_type": schema.Int64Attribute{
				MarkdownDescription: "TMDB list type. `1` Theaters, `2` Popular, `3` Top, `4` Upcoming.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2, 3, 4),
				},
			},
			"language_code": schema.Int64Attribute{
				MarkdownDescription: "Language code.",
				Required:            true,
			},
			"min_vote_average": schema.StringAttribute{
				MarkdownDescription: "Min vote average.",
				Optional:            true,
				Computed:            true,
			},
			"min_votes": schema.StringAttribute{
				MarkdownDescription: "Min votes.",
				Optional:            true,
				Computed:            true,
			},
			"tmdb_certification": schema.StringAttribute{
				MarkdownDescription: "Certification.",
				Optional:            true,
				Computed:            true,
			},
			"include_genre_ids": schema.StringAttribute{
				MarkdownDescription: "Include genre IDs.",
				Optional:            true,
				Computed:            true,
			},
			"exclude_genre_ids": schema.StringAttribute{
				MarkdownDescription: "Exclude genre IDs.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *ImportListTMDBPopularResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListTMDBPopularResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListTMDBPopular

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListTMDBPopular
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListTMDBPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListTMDBPopularResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBPopularResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListTMDBPopular

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListTMDBPopular current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTMDBPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListTMDBPopularResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBPopularResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListTMDBPopular

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListTMDBPopular
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListTMDBPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListTMDBPopularResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBPopularResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListTMDBPopular current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListTMDBPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListTMDBPopularResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListTMDBPopularResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListTMDBPopularResourceName+": "+req.ID)
}

func (i *ImportListTMDBPopular) write(ctx context.Context, importList *radarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListTMDBPopular) read(ctx context.Context, diags *diag.Diagnostics) *radarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}
