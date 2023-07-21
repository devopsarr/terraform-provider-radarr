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

const importListConfigResourceName = "import_list_config"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListConfigResource{}
	_ resource.ResourceWithImportState = &ImportListConfigResource{}
)

func NewImportListConfigResource() resource.Resource {
	return &ImportListConfigResource{}
}

// ImportListConfigResource defines the import list config implementation.
type ImportListConfigResource struct {
	client *radarr.APIClient
}

// ImportListConfig describes the import list config data model.
type ImportListConfig struct {
	SyncLevel    types.String `tfsdk:"sync_level"`
	SyncInterval types.Int64  `tfsdk:"sync_interval"`
	ID           types.Int64  `tfsdk:"id"`
}

func (r *ImportListConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListConfigResourceName
}

func (r *ImportListConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Config resource.\nFor more information refer to [Import List](https://wiki.servarr.com/radarr/settings#completed-download-handling) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List Config ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sync_interval": schema.Int64Attribute{
				MarkdownDescription: "List Update Interval.",
				Required:            true,
			},
			"sync_level": schema.StringAttribute{
				MarkdownDescription: "Clean library level.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "logOnly", "keepAndUnmonitor", "removeAndKeep", "removeAndDelete"),
				},
			},
		},
	}
}

func (r *ImportListConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *ImportListConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := config.read()
	request.SetId(1)

	// Create new ImportListConfig
	response, _, err := r.client.ImportListConfigApi.UpdateImportListConfig(ctx, strconv.Itoa(int(request.GetId()))).ImportListConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *ImportListConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var config *ImportListConfig

	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListConfig current value
	response, _, err := r.client.ImportListConfigApi.GetImportListConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *ImportListConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var config *ImportListConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := config.read()

	// Update ImportListConfig
	response, _, err := r.client.ImportListConfigApi.UpdateImportListConfig(ctx, strconv.Itoa(int(request.GetId()))).ImportListConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (r *ImportListConfigResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// ImportListConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+importListConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *ImportListConfigResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+importListConfigResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *ImportListConfig) write(importListConfig *radarr.ImportListConfigResource) {
	c.ID = types.Int64Value(int64(importListConfig.GetId()))
	c.SyncInterval = types.Int64Value(int64(importListConfig.GetImportListSyncInterval()))
	c.SyncLevel = types.StringValue(importListConfig.GetListSyncLevel())
}

func (c *ImportListConfig) read() *radarr.ImportListConfigResource {
	config := radarr.NewImportListConfigResource()
	config.SetId(int32(c.ID.ValueInt64()))
	config.SetListSyncLevel(c.SyncLevel.ValueString())
	config.SetImportListSyncInterval(int32(c.SyncInterval.ValueInt64()))

	return config
}
