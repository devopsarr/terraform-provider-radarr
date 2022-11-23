package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

const customFormatResourceName = "custom_format"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CustomFormatResource{}
var _ resource.ResourceWithImportState = &CustomFormatResource{}

var (
	customFormatStringFields = []string{"value"}
	customFormatIntFields    = []string{"min", "max"}
)

func NewCustomFormatResource() resource.Resource {
	return &CustomFormatResource{}
}

// CustomFormatResource defines the custom format implementation.
type CustomFormatResource struct {
	client *radarr.Radarr
}

// CustomFormat describes the custom format data model.
type CustomFormat struct {
	Specifications                  types.Set    `tfsdk:"specifications"`
	Name                            types.String `tfsdk:"name"`
	ID                              types.Int64  `tfsdk:"id"`
	IncludeCustomFormatWhenRenaming types.Bool   `tfsdk:"include_custom_format_when_renaming"`
}

type Specification struct {
	Name           types.String `tfsdk:"name"`
	Implementation types.String `tfsdk:"implementation"`
	Value          types.String `tfsdk:"value"`
	Min            types.Int64  `tfsdk:"min"`
	Max            types.Int64  `tfsdk:"max"`
	Negate         types.Bool   `tfsdk:"negate"`
	Required       types.Bool   `tfsdk:"required"`
}

func (r *CustomFormatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatResourceName
}

func (r *CustomFormatResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Custom Format resource.\nFor more information refer to [Custom Format](https://wiki.servarr.com/radarr/settings#custom-formats).",
		Attributes: map[string]tfsdk.Attribute{
			"include_custom_format_when_renaming": {
				MarkdownDescription: "Include custom format when renaming flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"name": {
				MarkdownDescription: "Custom Format name.",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Custom Format ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"specifications": {
				MarkdownDescription: "Specifications.",
				Required:            true,
				Attributes:          tfsdk.SetNestedAttributes(r.getSpecificationSchema().Attributes),
			},
		},
	}, nil
}

func (r CustomFormatResource) getSpecificationSchema() tfsdk.Schema {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"negate": {
				MarkdownDescription: "Negate flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"required": {
				MarkdownDescription: "Required flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"name": {
				MarkdownDescription: "Specification name.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Implementation.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			// Field values
			"value": {
				MarkdownDescription: "Value.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"min": {
				MarkdownDescription: "Min.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"max": {
				MarkdownDescription: "Max.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
		},
	}
}

func (r *CustomFormatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *CustomFormatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *CustomFormat

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new CustomFormat
	request := client.read(ctx)

	response, err := r.client.AddCustomFormatContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+customFormatResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client CustomFormat

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get CustomFormat current value
	response, err := r.client.GetCustomFormatContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *CustomFormat

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update CustomFormat
	request := client.read(ctx)

	response, err := r.client.UpdateCustomFormatContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+customFormatResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *CustomFormat

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete CustomFormat current value
	err := r.client.DeleteCustomFormatContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+customFormatResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *CustomFormatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+customFormatResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (c *CustomFormat) write(ctx context.Context, customFormat *radarr.CustomFormatOutput) {
	c.ID = types.Int64Value(customFormat.ID)
	c.Name = types.StringValue(customFormat.Name)
	c.IncludeCustomFormatWhenRenaming = types.BoolValue(customFormat.IncludeCFWhenRenaming)
	c.Specifications = types.SetValueMust(CustomFormatResource{}.getSpecificationSchema().Type(), nil)

	specs := make([]Specification, len(customFormat.Specifications))
	for n, c := range customFormat.Specifications {
		specs[n].write(c)
	}

	tfsdk.ValueFrom(ctx, specs, c.Specifications.Type(ctx), &c.Specifications)
}

func (s *Specification) write(spec *radarr.CustomFormatOutputSpec) {
	s.Implementation = types.StringValue(spec.Implementation)
	s.Name = types.StringValue(spec.Name)
	s.Negate = types.BoolValue(spec.Negate)
	s.Required = types.BoolValue(spec.Required)
	s.writeFields(spec.Fields)
}

func (s *Specification) writeFields(fields []*starr.FieldOutput) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(customFormatStringFields, f.Name) {
			tools.WriteStringField(f, s)

			continue
		}

		if slices.Contains(customFormatIntFields, f.Name) {
			tools.WriteIntField(f, s)

			continue
		}
	}
}

func (c *CustomFormat) read(ctx context.Context) *radarr.CustomFormatInput {
	specifications := make([]Specification, len(c.Specifications.Elements()))
	tfsdk.ValueAs(ctx, c.Specifications, &specifications)
	specs := make([]*radarr.CustomFormatInputSpec, len(specifications))

	for n, d := range specifications {
		specs[n] = d.read()
	}

	return &radarr.CustomFormatInput{
		ID:                    c.ID.ValueInt64(),
		Name:                  c.Name.ValueString(),
		IncludeCFWhenRenaming: c.IncludeCustomFormatWhenRenaming.ValueBool(),
		Specifications:        specs,
	}
}

func (s *Specification) read() *radarr.CustomFormatInputSpec {
	return &radarr.CustomFormatInputSpec{
		Name:           s.Name.ValueString(),
		Implementation: s.Implementation.ValueString(),
		Negate:         s.Negate.ValueBool(),
		Required:       s.Required.ValueBool(),
		Fields:         s.readFields(),
	}
}

func (s *Specification) readFields() []*starr.FieldInput {
	var output []*starr.FieldInput

	for _, i := range customFormatIntFields {
		if field := tools.ReadIntField(i, s); field != nil {
			output = append(output, field)
		}
	}

	for _, str := range customFormatStringFields {
		if field := tools.ReadStringField(str, s); field != nil {
			output = append(output, field)
		}
	}

	return output
}
