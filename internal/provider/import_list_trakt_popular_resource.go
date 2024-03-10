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
	importListTraktPopularResourceName   = "import_list_trakt_popular"
	importListTraktPopularImplementation = "TraktPopularImport"
	importListTraktPopularConfigContract = "TraktPopularSettings"
	importListTraktPopularType           = "trakt"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListTraktPopularResource{}
	_ resource.ResourceWithImportState = &ImportListTraktPopularResource{}
)

func NewImportListTraktPopularResource() resource.Resource {
	return &ImportListTraktPopularResource{}
}

// ImportListTraktPopularResource defines the import list implementation.
type ImportListTraktPopularResource struct {
	client *radarr.APIClient
	auth   context.Context
}

// ImportListTraktPopular describes the import list data model.
type ImportListTraktPopular struct {
	Tags                      types.Set    `tfsdk:"tags"`
	Name                      types.String `tfsdk:"name"`
	Monitor                   types.String `tfsdk:"monitor"`
	MinimumAvailability       types.String `tfsdk:"minimum_availability"`
	RootFolderPath            types.String `tfsdk:"root_folder_path"`
	AuthUser                  types.String `tfsdk:"auth_user"`
	TraktAdditionalParameters types.String `tfsdk:"trakt_additional_parameters"`
	AccessToken               types.String `tfsdk:"access_token"`
	RefreshToken              types.String `tfsdk:"refresh_token"`
	Expires                   types.String `tfsdk:"expires"`
	Certification             types.String `tfsdk:"certification"`
	Genres                    types.String `tfsdk:"genres"`
	Years                     types.String `tfsdk:"years"`
	Rating                    types.String `tfsdk:"rating"`
	TraktListType             types.Int64  `tfsdk:"trakt_list_type"`
	Limit                     types.Int64  `tfsdk:"limit"`
	ListOrder                 types.Int64  `tfsdk:"list_order"`
	ID                        types.Int64  `tfsdk:"id"`
	QualityProfileID          types.Int64  `tfsdk:"quality_profile_id"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	EnableAuto                types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd               types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListTraktPopular) toImportList() *ImportList {
	return &ImportList{
		Tags:                      i.Tags,
		Name:                      i.Name,
		Monitor:                   i.Monitor,
		MinimumAvailability:       i.MinimumAvailability,
		RootFolderPath:            i.RootFolderPath,
		ListOrder:                 i.ListOrder,
		RefreshToken:              i.RefreshToken,
		AccessToken:               i.AccessToken,
		Expires:                   i.Expires,
		AuthUser:                  i.AuthUser,
		TraktAdditionalParameters: i.TraktAdditionalParameters,
		Certification:             i.Certification,
		Genres:                    i.Genres,
		Years:                     i.Years,
		Rating:                    i.Rating,
		TraktListType:             i.TraktListType,
		Limit:                     i.Limit,
		ID:                        i.ID,
		QualityProfileID:          i.QualityProfileID,
		Enabled:                   i.Enabled,
		EnableAuto:                i.EnableAuto,
		SearchOnAdd:               i.SearchOnAdd,
		Implementation:            types.StringValue(importListTraktPopularImplementation),
		ConfigContract:            types.StringValue(importListTraktPopularConfigContract),
		ListType:                  types.StringValue(importListTraktPopularType),
	}
}

func (i *ImportListTraktPopular) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.Monitor = importList.Monitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.RefreshToken = importList.RefreshToken
	i.AccessToken = importList.AccessToken
	i.Expires = importList.Expires
	i.AuthUser = importList.AuthUser
	i.TraktAdditionalParameters = importList.TraktAdditionalParameters
	i.Certification = importList.Certification
	i.Genres = importList.Genres
	i.Years = importList.Years
	i.Rating = importList.Rating
	i.TraktListType = importList.TraktListType
	i.Limit = importList.Limit
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListTraktPopularResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListTraktPopularResourceName
}

func (r *ImportListTraktPopularResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Trakt Popular resource.\nFor more information refer to [Import List](https://wiki.servarr.com/radarr/settings#import-lists) and [Trakt Popular](https://wiki.servarr.com/radarr/supported#traktpopularimport).",
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
			"limit": schema.Int64Attribute{
				MarkdownDescription: "limit.",
				Required:            true,
			},
			"trakt_list_type": schema.Int64Attribute{
				MarkdownDescription: "Trakt list type.`0` Trending, `1` Popular, `2` Anticipated, `3` BoxOffice, `4` TopWatchedByWeek, `5` TopWatchedByMonth, `6` TopWatchedByYear, `7` TopWatchedByAllTime, `8` RecommendedByWeek, `9` RecommendedByMonth, `10` RecommendedByYear, `10` RecommendedByAllTime.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11),
				},
			},
			"certification": schema.StringAttribute{
				MarkdownDescription: "Certification.",
				Optional:            true,
				Computed:            true,
			},
			"genres": schema.StringAttribute{
				MarkdownDescription: "Genres.",
				Optional:            true,
				Computed:            true,
			},
			"years": schema.StringAttribute{
				MarkdownDescription: "Years.",
				Optional:            true,
				Computed:            true,
			},
			"rating": schema.StringAttribute{
				MarkdownDescription: "Rating.",
				Optional:            true,
				Computed:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Required:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"trakt_additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Trakt additional parameters.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *ImportListTraktPopularResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *ImportListTraktPopularResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListTraktPopular

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListTraktPopular
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.CreateImportList(r.auth).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListTraktPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListTraktPopularResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktPopularResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListTraktPopular

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListTraktPopular current value
	response, _, err := r.client.ImportListAPI.GetImportListById(r.auth, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTraktPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListTraktPopularResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktPopularResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListTraktPopular

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListTraktPopular
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.UpdateImportList(r.auth, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListTraktPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListTraktPopularResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktPopularResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListTraktPopular current value
	_, err := r.client.ImportListAPI.DeleteImportList(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListTraktPopularResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListTraktPopularResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListTraktPopularResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListTraktPopularResourceName+": "+req.ID)
}

func (i *ImportListTraktPopular) write(ctx context.Context, importList *radarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListTraktPopular) read(ctx context.Context, diags *diag.Diagnostics) *radarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}
