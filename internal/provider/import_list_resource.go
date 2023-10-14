package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

const importListResourceName = "import_list"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListResource{}
	_ resource.ResourceWithImportState = &ImportListResource{}
)

var importListFields = helpers.Fields{
	Bools:             []string{"onlyActive", "personCast", "personCastDirector", "personCastProducer", "personCastSound", "personCastWriting"},
	Ints:              []string{"port", "source", "minScore", "tMDbListType", "limit", "traktListType", "languageCode", "userListType"},
	IntsExceptions:    []string{"filterCriteria.languageCode", "listType"},
	Strings:           []string{"baseUrl", "urlBase", "link", "apiKey", "url", "accessToken", "refreshToken", "expires", "companyId", "keywordId", "listId", "personId", "accountId", "authUser", "username", "listname", "traktAdditionalParameters", "tmdbCertification", "genres", "years", "rating", "minVoteAverage", "minVotes", "certification", "includeGenreIds", "excludeGenreIds"},
	StringsExceptions: []string{"filterCriteria.certification", "filterCriteria.minVoteAverage", "filterCriteria.minVotes", "filterCriteria.includeGenreIds", "filterCriteria.excludeGenreIds"},
	IntSlices:         []string{"profileIds", "tagIds"},
}

func NewImportListResource() resource.Resource {
	return &ImportListResource{}
}

// ImportListResource defines the import list implementation.
type ImportListResource struct {
	client *radarr.APIClient
}

// ImportList describes the import list data model.
type ImportList struct {
	ProfileIds                types.Set    `tfsdk:"profile_ids"`
	TagIds                    types.Set    `tfsdk:"tag_ids"`
	Tags                      types.Set    `tfsdk:"tags"`
	Name                      types.String `tfsdk:"name"`
	ConfigContract            types.String `tfsdk:"config_contract"`
	Implementation            types.String `tfsdk:"implementation"`
	Monitor                   types.String `tfsdk:"monitor"`
	MinimumAvailability       types.String `tfsdk:"minimum_availability"`
	RootFolderPath            types.String `tfsdk:"root_folder_path"`
	ListType                  types.String `tfsdk:"list_type"`
	TraktAdditionalParameters types.String `tfsdk:"trakt_additional_parameters"`
	Certification             types.String `tfsdk:"certification"`
	Genres                    types.String `tfsdk:"genres"`
	Years                     types.String `tfsdk:"years"`
	Rating                    types.String `tfsdk:"rating"`
	MinVoteAverage            types.String `tfsdk:"min_vote_average"`
	MinVotes                  types.String `tfsdk:"min_votes"`
	TMDBCertification         types.String `tfsdk:"tmdb_certification"`
	IncludeGenreIds           types.String `tfsdk:"include_genre_ids"`
	ExcludeGenreIds           types.String `tfsdk:"exclude_genre_ids"`
	AuthUser                  types.String `tfsdk:"auth_user"`
	Username                  types.String `tfsdk:"username"`
	Listname                  types.String `tfsdk:"listname"`
	KeywordID                 types.String `tfsdk:"keyword_id"`
	CompanyID                 types.String `tfsdk:"company_id"`
	ListID                    types.String `tfsdk:"list_id"`
	PersonID                  types.String `tfsdk:"person_id"`
	AccountID                 types.String `tfsdk:"account_id"`
	AccessToken               types.String `tfsdk:"access_token"`
	RefreshToken              types.String `tfsdk:"refresh_token"`
	Expires                   types.String `tfsdk:"expires"`
	BaseURL                   types.String `tfsdk:"base_url"`
	URLBase                   types.String `tfsdk:"url_base"`
	URL                       types.String `tfsdk:"url"`
	Link                      types.String `tfsdk:"link"`
	APIKey                    types.String `tfsdk:"api_key"`
	ListOrder                 types.Int64  `tfsdk:"list_order"`
	ID                        types.Int64  `tfsdk:"id"`
	QualityProfileID          types.Int64  `tfsdk:"quality_profile_id"`
	Port                      types.Int64  `tfsdk:"port"`
	Source                    types.Int64  `tfsdk:"source"`
	MinScore                  types.Int64  `tfsdk:"min_score"`
	TMDBListType              types.Int64  `tfsdk:"tmdb_list_type"`
	UserListType              types.Int64  `tfsdk:"user_list_type"`
	Limit                     types.Int64  `tfsdk:"limit"`
	TraktListType             types.Int64  `tfsdk:"trakt_list_type"`
	LanguageCode              types.Int64  `tfsdk:"language_code"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	EnableAuto                types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd               types.Bool   `tfsdk:"search_on_add"`
	OnlyActive                types.Bool   `tfsdk:"only_active"`
	PersonCast                types.Bool   `tfsdk:"cast"`
	PersonCastDirector        types.Bool   `tfsdk:"cast_director"`
	PersonCastProducer        types.Bool   `tfsdk:"cast_producer"`
	PersonCastSound           types.Bool   `tfsdk:"cast_sound"`
	PersonCastWriting         types.Bool   `tfsdk:"cast_writing"`
}

func (i ImportList) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tag_ids":                     types.SetType{}.WithElementType(types.Int64Type),
			"tags":                        types.SetType{}.WithElementType(types.Int64Type),
			"profile_ids":                 types.SetType{}.WithElementType(types.Int64Type),
			"name":                        types.StringType,
			"config_contract":             types.StringType,
			"implementation":              types.StringType,
			"monitor":                     types.StringType,
			"minimum_availability":        types.StringType,
			"root_folder_path":            types.StringType,
			"list_type":                   types.StringType,
			"trakt_additional_parameters": types.StringType,
			"certification":               types.StringType,
			"genres":                      types.StringType,
			"years":                       types.StringType,
			"rating":                      types.StringType,
			"min_vote_average":            types.StringType,
			"min_votes":                   types.StringType,
			"tmdb_certification":          types.StringType,
			"include_genre_ids":           types.StringType,
			"exclude_genre_ids":           types.StringType,
			"auth_user":                   types.StringType,
			"username":                    types.StringType,
			"listname":                    types.StringType,
			"keyword_id":                  types.StringType,
			"company_id":                  types.StringType,
			"list_id":                     types.StringType,
			"person_id":                   types.StringType,
			"account_id":                  types.StringType,
			"access_token":                types.StringType,
			"refresh_token":               types.StringType,
			"expires":                     types.StringType,
			"base_url":                    types.StringType,
			"url_base":                    types.StringType,
			"url":                         types.StringType,
			"link":                        types.StringType,
			"api_key":                     types.StringType,
			"list_order":                  types.Int64Type,
			"id":                          types.Int64Type,
			"quality_profile_id":          types.Int64Type,
			"port":                        types.Int64Type,
			"source":                      types.Int64Type,
			"min_score":                   types.Int64Type,
			"tmdb_list_type":              types.Int64Type,
			"user_list_type":              types.Int64Type,
			"limit":                       types.Int64Type,
			"trakt_list_type":             types.Int64Type,
			"language_code":               types.Int64Type,
			"enabled":                     types.BoolType,
			"enable_auto":                 types.BoolType,
			"search_on_add":               types.BoolType,
			"only_active":                 types.BoolType,
			"cast":                        types.BoolType,
			"cast_director":               types.BoolType,
			"cast_producer":               types.BoolType,
			"cast_sound":                  types.BoolType,
			"cast_writing":                types.BoolType,
		})
}

func (r *ImportListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListResourceName
}

func (r *ImportListResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Generic Import List resource. When possible use a specific resource instead.\nFor more information refer to [Import List](https://wiki.servarr.com/radarr/settings#import-lists).",
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
			"implementation": schema.StringAttribute{
				MarkdownDescription: "ImportList implementation name.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "ImportList configuration template.",
				Required:            true,
			},
			"list_type": schema.StringAttribute{
				MarkdownDescription: "List type.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("program", "tmdb", "trakt", "plex", "other", "advanced"),
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
			"only_active": schema.BoolAttribute{
				MarkdownDescription: "Only active.",
				Optional:            true,
				Computed:            true,
			},
			"cast": schema.BoolAttribute{
				MarkdownDescription: "Include cast.",
				Optional:            true,
				Computed:            true,
			},
			"cast_director": schema.BoolAttribute{
				MarkdownDescription: "Include cast director.",
				Optional:            true,
				Computed:            true,
			},
			"cast_producer": schema.BoolAttribute{
				MarkdownDescription: "Include cast producer.",
				Optional:            true,
				Computed:            true,
			},
			"cast_sound": schema.BoolAttribute{
				MarkdownDescription: "Include cast sound.",
				Optional:            true,
				Computed:            true,
			},
			"cast_writing": schema.BoolAttribute{
				MarkdownDescription: "Include cast writing.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"source": schema.Int64Attribute{
				MarkdownDescription: "Source.",
				Optional:            true,
				Computed:            true,
			},
			"min_score": schema.Int64Attribute{
				MarkdownDescription: "Min score.",
				Optional:            true,
				Computed:            true,
			},
			"tmdb_list_type": schema.Int64Attribute{
				MarkdownDescription: "TMDB list type.",
				Optional:            true,
				Computed:            true,
			},
			"user_list_type": schema.Int64Attribute{
				MarkdownDescription: "User list type.",
				Optional:            true,
				Computed:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "limit.",
				Optional:            true,
				Computed:            true,
			},
			"trakt_list_type": schema.Int64Attribute{
				MarkdownDescription: "Trakt list type.",
				Optional:            true,
				Computed:            true,
			},
			"language_code": schema.Int64Attribute{
				MarkdownDescription: "Language code.",
				Optional:            true,
				Computed:            true,
			},
			"listname": schema.StringAttribute{
				MarkdownDescription: "List name.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Optional:            true,
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"company_id": schema.StringAttribute{
				MarkdownDescription: "Company ID.",
				Optional:            true,
				Computed:            true,
			},
			"keyword_id": schema.StringAttribute{
				MarkdownDescription: "Keyword ID.",
				Optional:            true,
				Computed:            true,
			},
			"list_id": schema.StringAttribute{
				MarkdownDescription: "List ID.",
				Optional:            true,
				Computed:            true,
			},
			"person_id": schema.StringAttribute{
				MarkdownDescription: "Person ID.",
				Optional:            true,
				Computed:            true,
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Account ID.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Optional:            true,
				Computed:            true,
			},
			"link": schema.StringAttribute{
				MarkdownDescription: "Link.",
				Optional:            true,
				Computed:            true,
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
			"profile_ids": schema.SetAttribute{
				MarkdownDescription: "Profile IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ImportListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportList

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportList
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.writeSensitive(importList)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportList

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportList current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.writeSensitive(importList)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportList

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportList
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.writeSensitive(importList)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportList current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListResourceName+": "+req.ID)
}

func (i *ImportList) write(ctx context.Context, importList *radarr.ImportListResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	i.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	diags.Append(localDiag...)

	i.Enabled = types.BoolValue(importList.GetEnabled())
	i.EnableAuto = types.BoolValue(importList.GetEnableAuto())
	i.SearchOnAdd = types.BoolValue(importList.GetSearchOnAdd())
	i.QualityProfileID = types.Int64Value(int64(importList.GetQualityProfileId()))
	i.ID = types.Int64Value(int64(importList.GetId()))
	i.ListOrder = types.Int64Value(int64(importList.GetListOrder()))
	i.ConfigContract = types.StringValue(importList.GetConfigContract())
	i.Implementation = types.StringValue(importList.GetImplementation())
	i.Monitor = types.StringValue(string(importList.GetMonitor()))
	i.MinimumAvailability = types.StringValue(string(importList.GetMinimumAvailability()))
	i.RootFolderPath = types.StringValue(importList.GetRootFolderPath())
	i.ListType = types.StringValue(string(importList.GetListType()))
	i.Name = types.StringValue(importList.GetName())
	i.ProfileIds = types.SetValueMust(types.Int64Type, nil)
	i.TagIds = types.SetValueMust(types.Int64Type, nil)
	helpers.WriteFields(ctx, i, importList.GetFields(), importListFields)
}

func (i *ImportList) read(ctx context.Context, diags *diag.Diagnostics) *radarr.ImportListResource {
	list := radarr.NewImportListResource()
	list.SetEnabled(i.Enabled.ValueBool())
	list.SetEnableAuto(i.EnableAuto.ValueBool())
	list.SetSearchOnAdd(i.SearchOnAdd.ValueBool())
	list.SetQualityProfileId(int32(i.QualityProfileID.ValueInt64()))
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetListOrder(int32(i.ListOrder.ValueInt64()))
	list.SetMonitor(radarr.MonitorTypes(i.Monitor.ValueString()))
	list.SetRootFolderPath(i.RootFolderPath.ValueString())
	list.SetMinimumAvailability(radarr.MovieStatusType(i.MinimumAvailability.ValueString()))
	list.SetListType(radarr.ImportListType(i.ListType.ValueString()))
	list.SetConfigContract(i.ConfigContract.ValueString())
	list.SetImplementation(i.Implementation.ValueString())
	list.SetName(i.Name.ValueString())
	diags.Append(i.Tags.ElementsAs(ctx, &list.Tags, true)...)
	list.SetFields(helpers.ReadFields(ctx, i, importListFields))

	return list
}

// writeSensitive copy sensitive data from another resource.
func (i *ImportList) writeSensitive(importList *ImportList) {
	if !importList.APIKey.IsUnknown() {
		i.APIKey = importList.APIKey
	}
}
