package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NamingResource{}
var _ resource.ResourceWithImportState = &NamingResource{}

func NewNamingResource() resource.Resource {
	return &NamingResource{}
}

// NamingResource defines the naming implementation.
type NamingResource struct {
	client *radarr.Radarr
}

// Naming describes the naming data model.
type Naming struct {
	IncludeQuality           types.Bool   `tfsdk:"include_quality"`
	RenameMovies             types.Bool   `tfsdk:"rename_movies"`
	ReplaceIllegalCharacters types.Bool   `tfsdk:"replace_illegal_characters"`
	ReplaceSpaces            types.Bool   `tfsdk:"replace_spaces"`
	ID                       types.Int64  `tfsdk:"id"`
	ColonReplacementFormat   types.String `tfsdk:"colon_replacement_format"`
	StandardMovieFormat      types.String `tfsdk:"standard_movie_format"`
	MovieFolderFormat        types.String `tfsdk:"movie_folder_format"`
}

func (r *NamingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_naming"
}

func (r *NamingResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Media Management]: #\nNaming resource.\nFor more information refer to [Naming](https://wiki.servarr.com/radarr/settings#community-naming-suggestions) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Naming ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"include_quality": {
				MarkdownDescription: "Include quality in file name.",
				Required:            true,
				Type:                types.BoolType,
			},
			"rename_movies": {
				MarkdownDescription: "Radarr will use the existing file name if false.",
				Required:            true,
				Type:                types.BoolType,
			},
			"replace_illegal_characters": {
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Required:            true,
				Type:                types.BoolType,
			},
			"replace_spaces": {
				MarkdownDescription: "Replace spaces.",
				Required:            true,
				Type:                types.BoolType,
			},
			"colon_replacement_format": {
				MarkdownDescription: "Change how Radarr handles colon replacement. Valid values are: 'delete', 'dash', 'spaceDash', and 'spaceDashSpace'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"delete", "dash", "spaceDash", "spaceDashSpace"}),
				},
			},
			"movie_folder_format": {
				MarkdownDescription: "Movie folder format.",
				Required:            true,
				Type:                types.StringType,
			},
			"standard_movie_format": {
				MarkdownDescription: "Standard movie formatss.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *NamingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NamingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	init, err := r.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to init naming, got error: %s", err))

		return
	}

	_, err = r.client.UpdateNamingContext(ctx, init)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to init naming, got error: %s", err))

		return
	}

	// Build Create resource
	data := readNaming(&plan)
	data.ID = 1

	// Create new Naming
	response, err := r.client.UpdateNamingContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create naming, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created naming: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *NamingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Naming

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get naming current value
	response, err := r.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read namings, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read naming: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *NamingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readNaming(&plan)

	// Update Naming
	response, err := r.client.UpdateNamingContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update naming, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated naming: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *NamingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Naming cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled naming: 1")
	resp.State.RemoveResource(ctx)
}

func (r *NamingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported naming: 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func writeNaming(naming *radarr.Naming) *Naming {
	return &Naming{
		IncludeQuality:           types.BoolValue(naming.IncludeQuality),
		RenameMovies:             types.BoolValue(naming.RenameMovies),
		ReplaceIllegalCharacters: types.BoolValue(naming.ReplaceIllegalCharacters),
		ReplaceSpaces:            types.BoolValue(naming.ReplaceSpaces),
		ID:                       types.Int64Value(naming.ID),
		ColonReplacementFormat:   types.StringValue(naming.ColonReplacementFormat),
		StandardMovieFormat:      types.StringValue(naming.StandardMovieFormat),
		MovieFolderFormat:        types.StringValue(naming.MovieFolderFormat),
	}
}

func readNaming(naming *Naming) *radarr.Naming {
	return &radarr.Naming{
		IncludeQuality:           naming.IncludeQuality.ValueBool(),
		RenameMovies:             naming.RenameMovies.ValueBool(),
		ReplaceIllegalCharacters: naming.ReplaceIllegalCharacters.ValueBool(),
		ReplaceSpaces:            naming.ReplaceSpaces.ValueBool(),
		ID:                       naming.ID.ValueInt64(),
		ColonReplacementFormat:   naming.ColonReplacementFormat.ValueString(),
		StandardMovieFormat:      naming.StandardMovieFormat.ValueString(),
		MovieFolderFormat:        naming.MovieFolderFormat.ValueString(),
	}
}
