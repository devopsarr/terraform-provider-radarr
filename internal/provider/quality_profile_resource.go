package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

const qualityProfileResourceName = "quality_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QualityProfileResource{}
	_ resource.ResourceWithImportState = &QualityProfileResource{}
)

func NewQualityProfileResource() resource.Resource {
	return &QualityProfileResource{}
}

// QualityProfileResource defines the quality profile implementation.
type QualityProfileResource struct {
	client *radarr.Radarr
}

// QualityProfile describes the quality profile data model.
type QualityProfile struct {
	QualityGroups     types.Set    `tfsdk:"quality_groups"`
	FormatItems       types.Set    `tfsdk:"format_items"`
	Name              types.String `tfsdk:"name"`
	Language          types.Object `tfsdk:"language"`
	ID                types.Int64  `tfsdk:"id"`
	Cutoff            types.Int64  `tfsdk:"cutoff"`
	MinFormatScore    types.Int64  `tfsdk:"min_format_score"`
	CutoffFormatScore types.Int64  `tfsdk:"cutoff_format_score"`
	UpgradeAllowed    types.Bool   `tfsdk:"upgrade_allowed"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	Qualities types.Set    `tfsdk:"qualities"`
	Name      types.String `tfsdk:"name"`
	ID        types.Int64  `tfsdk:"id"`
}

// Quality is part of QualityGroup.
type Quality struct {
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
}

// FormatItem is part of QualityProfile.
type FormatItem struct {
	Name   types.String `tfsdk:"name"`
	Format types.Int64  `tfsdk:"format"`
	Score  types.Int64  `tfsdk:"score"`
}

// Language is part of QualityProfile.
type Language struct {
	Name types.String `tfsdk:"name"`
	ID   types.Int64  `tfsdk:"id"`
}

func (r *QualityProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfileResourceName
}

func (r *QualityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Quality Profile resource.\nFor more information refer to [Quality Profile](https://wiki.servarr.com/radarr/settings#quality-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed flag.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff": schema.Int64Attribute{
				MarkdownDescription: "Quality ID to which cutoff.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff_format_score": schema.Int64Attribute{
				MarkdownDescription: "Cutoff format score.",
				Optional:            true,
				Computed:            true,
			},
			"min_format_score": schema.Int64Attribute{
				MarkdownDescription: "Min format score.",
				Optional:            true,
				Computed:            true,
			},
			"language": schema.SingleNestedAttribute{
				MarkdownDescription: "Language.",
				Required:            true,
				Attributes:          r.getLanguageSchema().Attributes,
			},
			"quality_groups": schema.SetNestedAttribute{
				MarkdownDescription: "Quality groups.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualityGroupSchema().Attributes,
				},
			},
			"format_items": schema.SetNestedAttribute{
				MarkdownDescription: "Format items.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getFormatItemsSchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualityGroupSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality group ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality group name.",
				Optional:            true,
				Computed:            true,
			},
			"qualities": schema.SetNestedAttribute{
				MarkdownDescription: "Qualities in group.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualitySchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualitySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Optional:            true,
				Computed:            true,
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Resolution.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality name.",
				Optional:            true,
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Source.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r QualityProfileResource) getFormatItemsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"format": schema.Int64Attribute{
				MarkdownDescription: "Format.",
				Optional:            true,
				Computed:            true,
			},
			"score": schema.Int64Attribute{
				MarkdownDescription: "Score.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r QualityProfileResource) getLanguageSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *QualityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *QualityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := profile.read(ctx)

	// Create new QualityProfile
	response, err := r.client.AddQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualityprofile current value
	response, err := r.client.GetQualityProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := profile.read(ctx)

	// Update QualityProfile
	response, err := r.client.UpdateQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	err := r.client.DeleteQualityProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+qualityProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *QualityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+qualityProfileResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (p *QualityProfile) write(ctx context.Context, profile *radarr.QualityProfile) {
	p.UpgradeAllowed = types.BoolValue(profile.UpgradeAllowed)
	p.ID = types.Int64Value(profile.ID)
	p.Name = types.StringValue(profile.Name)
	p.Cutoff = types.Int64Value(profile.Cutoff)
	p.CutoffFormatScore = types.Int64Value(profile.CutoffFormatScore)
	p.MinFormatScore = types.Int64Value(profile.MinFormatScore)
	p.QualityGroups = types.SetValueMust(QualityProfileResource{}.getQualityGroupSchema().Type(), nil)
	p.FormatItems = types.SetValueMust(QualityProfileResource{}.getFormatItemsSchema().Type(), nil)

	qualityGroups := make([]QualityGroup, len(profile.Qualities))
	for n, g := range profile.Qualities {
		qualityGroups[n].write(ctx, g)
	}

	formatItems := make([]FormatItem, len(profile.FormatItems))
	for n, f := range profile.FormatItems {
		formatItems[n].write(f)
	}

	language := Language{}
	language.write(profile.Language)
	tfsdk.ValueFrom(ctx, language, QualityProfileResource{}.getLanguageSchema().Type(), &p.Language)

	tfsdk.ValueFrom(ctx, qualityGroups, p.QualityGroups.Type(ctx), &p.QualityGroups)
}

func (q *QualityGroup) write(ctx context.Context, group *starr.Quality) {
	var (
		name      string
		id        int64
		qualities []Quality
	)

	if len(group.Items) == 0 {
		name = group.Quality.Name
		id = group.Quality.ID
		qualities = []Quality{{
			ID:         types.Int64Value(group.Quality.ID),
			Name:       types.StringValue(group.Quality.Name),
			Source:     types.StringValue(group.Quality.Source),
			Resolution: types.Int64Value(int64(group.Quality.Resolution)),
		}}
	} else {
		name = group.Name
		id = int64(group.ID)
		qualities = make([]Quality, len(group.Items))
		for m, q := range group.Items {
			qualities[m].write(q)
		}
	}

	q.Name = types.StringValue(name)
	q.ID = types.Int64Value(id)
	q.Qualities = types.SetValueMust(QualityProfileResource{}.getQualitySchema().Type(), nil)

	tfsdk.ValueFrom(ctx, qualities, q.Qualities.Type(ctx), &q.Qualities)
}

func (q *Quality) write(quality *starr.Quality) {
	q.ID = types.Int64Value(quality.Quality.ID)
	q.Name = types.StringValue(quality.Quality.Name)
	q.Source = types.StringValue(quality.Quality.Source)
	q.Resolution = types.Int64Value(int64(quality.Quality.Resolution))
}

func (f *FormatItem) write(format *starr.FormatItem) {
	f.Name = types.StringValue(format.Name)
	f.Format = types.Int64Value(format.Format)
	f.Score = types.Int64Value(format.Score)
}

func (l *Language) write(language *starr.Value) {
	l.Name = types.StringValue(language.Name)
	l.ID = types.Int64Value(language.ID)
}

func (p *QualityProfile) read(ctx context.Context) *radarr.QualityProfile {
	groups := make([]QualityGroup, len(p.QualityGroups.Elements()))
	tfsdk.ValueAs(ctx, p.QualityGroups, &groups)
	qualities := make([]*starr.Quality, len(groups))

	for n, g := range groups {
		q := make([]Quality, len(g.Qualities.Elements()))
		tfsdk.ValueAs(ctx, g.Qualities, &q)

		if len(q) == 0 {
			qualities[n] = &starr.Quality{
				Allowed: true,
				Quality: &starr.BaseQuality{
					ID:   g.ID.ValueInt64(),
					Name: g.Name.ValueString(),
				},
			}

			continue
		}

		items := make([]*starr.Quality, len(q))
		for m, q := range q {
			items[m] = q.read()
		}

		qualities[n] = &starr.Quality{
			Name:    g.Name.ValueString(),
			ID:      int(g.ID.ValueInt64()),
			Allowed: true,
			Items:   items,
		}
	}

	formats := make([]FormatItem, len(p.FormatItems.Elements()))
	tfsdk.ValueAs(ctx, p.FormatItems, &formats)
	formatItems := make([]*starr.FormatItem, len(formats))

	for n, f := range formats {
		formatItems[n] = f.read()
	}

	language := Language{}
	tfsdk.ValueAs(ctx, p.Language, &language)

	return &radarr.QualityProfile{
		UpgradeAllowed:    p.UpgradeAllowed.ValueBool(),
		ID:                p.ID.ValueInt64(),
		Cutoff:            p.Cutoff.ValueInt64(),
		Name:              p.Name.ValueString(),
		MinFormatScore:    p.MinFormatScore.ValueInt64(),
		CutoffFormatScore: p.Cutoff.ValueInt64(),
		Language:          language.read(),
		Qualities:         qualities,
		FormatItems:       formatItems,
	}
}

func (q *Quality) read() *starr.Quality {
	return &starr.Quality{
		Allowed: true,
		Quality: &starr.BaseQuality{
			Name:       q.Name.ValueString(),
			ID:         q.ID.ValueInt64(),
			Source:     q.Source.ValueString(),
			Resolution: int(q.Resolution.ValueInt64()),
		},
	}
}

func (f *FormatItem) read() *starr.FormatItem {
	return &starr.FormatItem{
		Format: f.Format.ValueInt64(),
		Name:   f.Name.ValueString(),
		Score:  f.Score.ValueInt64(),
	}
}

func (l *Language) read() *starr.Value {
	return &starr.Value{
		ID:   l.ID.ValueInt64(),
		Name: l.Name.ValueString(),
	}
}
