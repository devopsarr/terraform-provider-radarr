package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	importListPlexResourceName   = "import_list_plex"
	importListPlexImplementation = "PlexImport"
	importListPlexConfigContract = "PlexListSettings"
	importListPlexType           = "plex"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListPlexResource{}
	_ resource.ResourceWithImportState = &ImportListPlexResource{}
)

func NewImportListPlexResource() resource.Resource {
	return &ImportListPlexResource{}
}

// ImportListPlexResource defines the import list implementation.
type ImportListPlexResource struct {
	client *radarr.APIClient
}

// ImportListPlex describes the import list data model.
type ImportListPlex struct {
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	Monitor             types.String `tfsdk:"monitor"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	RootFolderPath      types.String `tfsdk:"root_folder_path"`
	AccessToken         types.String `tfsdk:"access_token"`
	ListOrder           types.Int64  `tfsdk:"list_order"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	EnableAuto          types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd         types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListPlex) toImportList() *ImportList {
	return &ImportList{
		Tags:                i.Tags,
		Name:                i.Name,
		Monitor:             i.Monitor,
		MinimumAvailability: i.MinimumAvailability,
		RootFolderPath:      i.RootFolderPath,
		AccessToken:         i.AccessToken,
		ListOrder:           i.ListOrder,
		ID:                  i.ID,
		QualityProfileID:    i.QualityProfileID,
		Enabled:             i.Enabled,
		EnableAuto:          i.EnableAuto,
		SearchOnAdd:         i.SearchOnAdd,
		Implementation:      types.StringValue(importListPlexImplementation),
		ConfigContract:      types.StringValue(importListPlexConfigContract),
		ListType:            types.StringValue(importListPlexType),
	}
}

func (i *ImportListPlex) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.Monitor = importList.Monitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.AccessToken = importList.AccessToken
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListPlexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListPlexResourceName
}

func (r *ImportListPlexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Plex resource.\nFor more information refer to [Import List](https://wiki.servarr.com/radarr/settings#import-lists) and [Plex](https://wiki.servarr.com/radarr/supported#pleximport).",
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
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *ImportListPlexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListPlexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListPlex
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListPlexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListPlexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListPlex current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListPlexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListPlexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListPlex
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListPlexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListPlexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListPlex current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListPlexResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListPlexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListPlexResourceName+": "+req.ID)
}

func (i *ImportListPlex) write(ctx context.Context, importList *radarr.ImportListResource) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList)
	i.fromImportList(genericImportList)
}

func (i *ImportListPlex) read(ctx context.Context) *radarr.ImportListResource {
	return i.toImportList().read(ctx)
}
