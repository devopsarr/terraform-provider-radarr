package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/tools"
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

const namingResourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NamingResource{}
	_ resource.ResourceWithImportState = &NamingResource{}
)

func NewNamingResource() resource.Resource {
	return &NamingResource{}
}

// NamingResource defines the naming implementation.
type NamingResource struct {
	client *radarr.APIClient
}

// Naming describes the naming data model.
type Naming struct {
	ColonReplacementFormat   types.String `tfsdk:"colon_replacement_format"`
	StandardMovieFormat      types.String `tfsdk:"standard_movie_format"`
	MovieFolderFormat        types.String `tfsdk:"movie_folder_format"`
	ID                       types.Int64  `tfsdk:"id"`
	IncludeQuality           types.Bool   `tfsdk:"include_quality"`
	RenameMovies             types.Bool   `tfsdk:"rename_movies"`
	ReplaceIllegalCharacters types.Bool   `tfsdk:"replace_illegal_characters"`
	ReplaceSpaces            types.Bool   `tfsdk:"replace_spaces"`
}

func (r *NamingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingResourceName
}

func (r *NamingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Naming resource.\nFor more information refer to [Naming](https://wiki.servarr.com/radarr/settings#community-naming-suggestions) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Naming ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"include_quality": schema.BoolAttribute{
				MarkdownDescription: "Include quality in file name.",
				Required:            true,
			},
			"rename_movies": schema.BoolAttribute{
				MarkdownDescription: "Radarr will use the existing file name if false.",
				Required:            true,
			},
			"replace_illegal_characters": schema.BoolAttribute{
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Required:            true,
			},
			"replace_spaces": schema.BoolAttribute{
				MarkdownDescription: "Replace spaces.",
				Required:            true,
			},
			"colon_replacement_format": schema.StringAttribute{
				MarkdownDescription: "Change how Radarr handles colon replacement. Valid values are: 'delete', 'dash', 'spaceDash', and 'spaceDashSpace'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("delete", "dash", "spaceDash", "spaceDashSpace"),
				},
			},
			"movie_folder_format": schema.StringAttribute{
				MarkdownDescription: "Movie folder format.",
				Required:            true,
			},
			"standard_movie_format": schema.StringAttribute{
				MarkdownDescription: "Standard movie formatss.",
				Required:            true,
			},
		},
	}
}

func (r *NamingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *radarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NamingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var naming *Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	if _, _, err := r.client.NamingConfigApi.GetConfigNaming(ctx).Execute(); err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to init %s, got error: %s", namingResourceName, err))

		return
	}

	// Build Create resource
	request := naming.read()
	request.SetId(1)

	// Create new Naming
	response, _, err := r.client.NamingConfigApi.UpdateConfigNaming(ctx, strconv.Itoa(int(request.GetId()))).NamingConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var naming *Naming

	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get naming current value
	response, _, err := r.client.NamingConfigApi.GetConfigNaming(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var naming *Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := naming.read()

	// Update Naming
	response, _, err := r.client.NamingConfigApi.UpdateConfigNaming(ctx, strconv.Itoa(int(request.GetId()))).NamingConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Naming cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+namingResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *NamingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+namingResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (n *Naming) write(naming *radarr.NamingConfigResource) {
	n.IncludeQuality = types.BoolValue(naming.GetIncludeQuality())
	n.RenameMovies = types.BoolValue(naming.GetRenameMovies())
	n.ReplaceIllegalCharacters = types.BoolValue(naming.GetReplaceIllegalCharacters())
	n.ReplaceSpaces = types.BoolValue(naming.GetReplaceSpaces())
	n.ID = types.Int64Value(int64(naming.GetId()))
	n.ColonReplacementFormat = types.StringValue(string(naming.GetColonReplacementFormat()))
	n.StandardMovieFormat = types.StringValue(naming.GetStandardMovieFormat())
	n.MovieFolderFormat = types.StringValue(naming.GetMovieFolderFormat())
}

func (n *Naming) read() *radarr.NamingConfigResource {
	naming := radarr.NewNamingConfigResource()
	naming.SetColonReplacementFormat(radarr.ColonReplacementFormat(n.ColonReplacementFormat.ValueString()))
	naming.SetId(int32(n.ID.ValueInt64()))
	naming.SetIncludeQuality(n.IncludeQuality.ValueBool())
	naming.SetMovieFolderFormat(n.MovieFolderFormat.ValueString())
	naming.SetRenameMovies(n.RenameMovies.ValueBool())
	naming.SetReplaceIllegalCharacters(n.ReplaceIllegalCharacters.ValueBool())
	naming.SetReplaceSpaces(n.ReplaceSpaces.ValueBool())
	naming.SetStandardMovieFormat(n.StandardMovieFormat.ValueString())

	return naming
}
