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

const metadataConfigResourceName = "metadata_config"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataConfigResource{}
	_ resource.ResourceWithImportState = &MetadataConfigResource{}
)

func NewMetadataConfigResource() resource.Resource {
	return &MetadataConfigResource{}
}

// MetadataConfigResource defines the metadata config implementation.
type MetadataConfigResource struct {
	client *radarr.APIClient
}

// MetadataConfig describes the metadata config data model.
type MetadataConfig struct {
	CertificationCountry types.String `tfsdk:"certification_country"`
	ID                   types.Int64  `tfsdk:"id"`
}

func (r *MetadataConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataConfigResourceName
}

func (r *MetadataConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Config resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/radarr/settings#options) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata Config ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"certification_country": schema.StringAttribute{
				MarkdownDescription: "Certification Country.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("au", "br", "ca", "fr", "de", "gb", "it", "es", "us", "nz"),
				},
			},
		},
	}
}

func (r *MetadataConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *MetadataConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := config.read()
	request.SetId(1)

	// Create new MetadataConfig
	response, _, err := r.client.MetadataConfigApi.UpdateMetadataConfig(ctx, strconv.Itoa(int(request.GetId()))).MetadataConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *MetadataConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var config *MetadataConfig

	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get metadataConfig current value
	response, _, err := r.client.MetadataConfigApi.GetMetadataConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *MetadataConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var config *MetadataConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := config.read()

	// Update MetadataConfig
	response, _, err := r.client.MetadataConfigApi.UpdateMetadataConfig(ctx, strconv.Itoa(int(request.GetId()))).MetadataConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *MetadataConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// MetadataConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+metadataConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *MetadataConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+metadataConfigResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *MetadataConfig) write(metadataConfig *radarr.MetadataConfigResource) {
	c.ID = types.Int64Value(int64(metadataConfig.GetId()))
	c.CertificationCountry = types.StringValue(string(metadataConfig.GetCertificationCountry()))
}

func (c *MetadataConfig) read() *radarr.MetadataConfigResource {
	config := radarr.NewMetadataConfigResource()
	config.SetCertificationCountry(radarr.TMDbCountryCode(c.CertificationCountry.ValueString()))
	config.SetId(int32(c.ID.ValueInt64()))

	return config
}
