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
	importListRSSResourceName   = "import_list_rss"
	importListRSSImplementation = "RSSImport"
	importListRSSConfigContract = "RSSImportSettings"
	importListRSSType           = "advanced"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListRSSResource{}
	_ resource.ResourceWithImportState = &ImportListRSSResource{}
)

func NewImportListRSSResource() resource.Resource {
	return &ImportListRSSResource{}
}

// ImportListRSSResource defines the import list implementation.
type ImportListRSSResource struct {
	client *radarr.APIClient
}

// ImportListRSS describes the import list data model.
type ImportListRSS struct {
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	Monitor             types.String `tfsdk:"monitor"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	RootFolderPath      types.String `tfsdk:"root_folder_path"`
	Link                types.String `tfsdk:"link"`
	ListOrder           types.Int64  `tfsdk:"list_order"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	EnableAuto          types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd         types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListRSS) toImportList() *ImportList {
	return &ImportList{
		Tags:                i.Tags,
		Name:                i.Name,
		Monitor:             i.Monitor,
		MinimumAvailability: i.MinimumAvailability,
		RootFolderPath:      i.RootFolderPath,
		Link:                i.Link,
		ListOrder:           i.ListOrder,
		ID:                  i.ID,
		QualityProfileID:    i.QualityProfileID,
		Enabled:             i.Enabled,
		EnableAuto:          i.EnableAuto,
		SearchOnAdd:         i.SearchOnAdd,
		Implementation:      types.StringValue(importListRSSImplementation),
		ConfigContract:      types.StringValue(importListRSSConfigContract),
		ListType:            types.StringValue(importListRSSType),
	}
}

func (i *ImportListRSS) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.Monitor = importList.Monitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.Link = importList.Link
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListRSSResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListRSSResourceName
}

func (r *ImportListRSSResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List RSS resource.\nFor more information refer to [Import List](https://wiki.servarr.com/radarr/settings#import-lists) and [RSS](https://wiki.servarr.com/radarr/supported#rssimport).",
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
			"link": schema.StringAttribute{
				MarkdownDescription: "Link.",
				Required:            true,
			},
		},
	}
}

func (r *ImportListRSSResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListRSSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListRSS

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListRSS
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListRSSResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListRSSResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListRSSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListRSS

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListRSS current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListRSSResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListRSSResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListRSSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListRSS

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListRSS
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListRSSResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListRSSResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListRSSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListRSS

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListRSS current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListRSSResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListRSSResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListRSSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListRSSResourceName+": "+req.ID)
}

func (i *ImportListRSS) write(ctx context.Context, importList *radarr.ImportListResource) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList)
	i.fromImportList(genericImportList)
}

func (i *ImportListRSS) read(ctx context.Context) *radarr.ImportListResource {
	return i.toImportList().read(ctx)
}
