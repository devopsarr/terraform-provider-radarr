package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Value is generic ID/Name struct applied to a few places.
type Value struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Tag is the tag resource.
type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// Tags is a list of Tag.
type Tags struct {
	ID   types.String `tfsdk:"id"`
	Tags []Tag        `tfsdk:"tags"`
}

// DelayProfile is the delay_profile resource.
type DelayProfile struct {
	EnableUsenet           types.Bool    `tfsdk:"enable_usenet"`
	EnableTorrent          types.Bool    `tfsdk:"enable_torrent"`
	BypassIfHighestQuality types.Bool    `tfsdk:"bypass_if_highest_quality"`
	UsenetDelay            types.Int64   `tfsdk:"usenet_delay"`
	TorrentDelay           types.Int64   `tfsdk:"torrent_delay"`
	ID                     types.Int64   `tfsdk:"id"`
	Order                  types.Int64   `tfsdk:"order"`
	Tags                   []types.Int64 `tfsdk:"tags"`
	PreferredProtocol      types.String  `tfsdk:"preferred_protocol"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// DelayProfiles is a list of DelayProfile.
type DelayProfiles struct {
	ID            types.String   `tfsdk:"id"`
	DelayProfiles []DelayProfile `tfsdk:"delay_profiles"`
}

// QualityProfile is the quality_profile resource.
type QualityProfile struct {
	UpgradeAllowed    types.Bool     `tfsdk:"upgrade_allowed"`
	ID                types.Int64    `tfsdk:"id"`
	Cutoff            types.Int64    `tfsdk:"cutoff"`
	MinFormatScore    types.Int64    `tfsdk:"min_format_score"`
	CutoffFormatScore types.Int64    `tfsdk:"cutoff_format_score"`
	Name              types.String   `tfsdk:"name"`
	Language          types.String   `tfsdk:"language"`
	FormatItems       []FormatItem   `tfsdk:"format_items"`
	QualityGroups     []QualityGroup `tfsdk:"quality_groups"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// QualityProfiles is a list of QualityProfile.
type QualityProfiles struct {
	ID              types.String     `tfsdk:"id"`
	QualityProfiles []QualityProfile `tfsdk:"quality_profiles"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Qualities []Quality    `tfsdk:"qualities"`
}

//Quality is part of QualityGroup.
type Quality struct {
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
}

//FormatItems is part of QualityProfile.
type FormatItem struct {
	Format types.Int64  `tfsdk:"format"`
	Score  types.Int64  `tfsdk:"score"`
	Name   types.String `tfsdk:"name"`
}
