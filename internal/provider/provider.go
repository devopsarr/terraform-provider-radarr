package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &RadarrProvider{}

// RadarrProvider defines the provider implementation.
type RadarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Radarr describes the provider data model.
type Radarr struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *RadarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "radarr"
	resp.Version = p.version
}

func (p *RadarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Radarr provider is used to interact with any [Radarr](https://radarr.video/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Radarr authentication. Can be specified via the `RADARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Radarr URL with protocol and port (e.g. `https://test.radarr.tv:7878`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `RADARR_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *RadarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Radarr

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract URL
	APIURL := data.URL.ValueString()
	if APIURL == "" {
		APIURL = os.Getenv("RADARR_URL")
	}

	parsedAPIURL, err := url.Parse(APIURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find valid URL",
			"URL cannot parsed",
		)

		return
	}

	// Extract key
	key := data.APIKey.ValueString()
	if key == "" {
		key = os.Getenv("RADARR_API_KEY")
	}

	if key == "" {
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Set context for API calls
	auth := context.WithValue(
		context.Background(),
		radarr.ContextAPIKeys,
		map[string]radarr.APIKey{
			"X-Api-Key": {Key: key},
		},
	)
	auth = context.WithValue(auth, radarr.ContextServerVariables, map[string]string{
		"protocol": parsedAPIURL.Scheme,
		"hostpath": parsedAPIURL.Host,
	})
	resp.DataSourceData = auth
	resp.ResourceData = auth
}

func (p *RadarrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Download Clients
		NewDownloadClientConfigResource,
		NewDownloadClientResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientFloodResource,
		NewDownloadClientFreeboxResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,
		NewRemotePathMappingResource,

		// Indexers
		NewIndexerConfigResource,
		NewIndexerResource,
		NewIndexerFilelistResource,
		NewIndexerIptorrentsResource,
		NewIndexerHdbitsResource,
		NewIndexerNewznabResource,
		NewIndexerNyaaResource,
		NewIndexerPassThePopcornResource,
		NewIndexerTorrentPotatoResource,
		NewIndexerTorrentRssResource,
		NewIndexerTorznabResource,

		// Import Lists
		NewImportListResource,
		NewImportListCustomResource,
		NewImportListCouchPotatoResource,
		NewImportListIMDBResource,
		NewImportListPlexResource,
		NewImportListRadarrResource,
		NewImportListRSSResource,
		NewImportListStevenluResource,
		NewImportListStevenlu2Resource,
		NewImportListTMDBCompanyResource,
		NewImportListTMDBKeywordResource,
		NewImportListTMDBListResource,
		NewImportListTMDBPersonResource,
		NewImportListTMDBPopularResource,
		NewImportListTMDBUserResource,
		NewImportListTraktListResource,
		NewImportListTraktPopularResource,
		NewImportListTraktUserResource,
		NewImportListConfigResource,
		NewImportListExclusionResource,

		// Media Management
		NewMediaManagementResource,
		NewNamingResource,
		NewRootFolderResource,

		// Metadata
		NewMetadataEmbyResource,
		NewMetadataKodiResource,
		NewMetadataRoksboxResource,
		NewMetadataWdtvResource,
		NewMetadataResource,
		NewMetadataConfigResource,

		// Movies
		NewMovieResource,

		// Notifications
		NewNotificationResource,
		NewNotificationAppriseResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationEmbyResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKodiResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationNtfyResource,
		NewNotificationPlexResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSimplepushResource,
		NewNotificationSlackResource,
		NewNotificationSynologyResource,
		NewNotificationTelegramResource,
		NewNotificationTraktResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// Profiles
		NewCustomFormatResource,
		NewDelayProfileResource,
		NewQualityProfileResource,
		NewQualityDefinitionResource,

		// System
		NewHostResource,

		// Tags
		NewTagResource,
		NewAutoTagResource,
	}
}

func (p *RadarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Download Clients
		NewDownloadClientConfigDataSource,
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,
		NewRemotePathMappingDataSource,
		NewRemotePathMappingsDataSource,

		// Indexers
		NewIndexerConfigDataSource,
		NewIndexerDataSource,
		NewIndexersDataSource,

		// Import Lists
		NewImportListConfigDataSource,
		NewImportListExclusionDataSource,
		NewImportListExclusionsDataSource,

		// Media Management
		NewMediaManagementDataSource,
		NewNamingDataSource,
		NewRootFolderDataSource,
		NewRootFoldersDataSource,

		// Metadata
		NewMetadataDataSource,
		NewMetadataConsumersDataSource,
		NewMetadataConfigDataSource,

		// Movies
		NewMovieDataSource,
		NewMoviesDataSource,

		// Notifications
		NewImportListDataSource,
		NewImportListsDataSource,
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// Profiles
		NewCustomFormatDataSource,
		NewCustomFormatsDataSource,
		NewDelayProfileDataSource,
		NewDelayProfilesDataSource,
		NewLanguageDataSource,
		NewLanguagesDataSource,
		NewQualityProfileDataSource,
		NewQualityProfilesDataSource,
		NewQualityDefinitionDataSource,
		NewQualityDefinitionsDataSource,
		NewCustomFormatConditionDataSource,
		NewCustomFormatConditionEditionDataSource,
		NewCustomFormatConditionIndexerFlagDataSource,
		NewCustomFormatConditionLanguageDataSource,
		NewCustomFormatConditionQualityModifierDataSource,
		NewCustomFormatConditionReleaseGroupDataSource,
		NewCustomFormatConditionReleaseTitleDataSource,
		NewCustomFormatConditionResolutionDataSource,
		NewCustomFormatConditionSizeDataSource,
		NewCustomFormatConditionSourceDataSource,
		NewQualityDataSource,

		// System
		NewSystemStatusDataSource,
		NewHostDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
		NewAutoTagDataSource,
		NewAutoTagsDataSource,
		NewAutoTagConditionDataSource,
		NewAutoTagConditionGenresDataSource,
		NewAutoTagConditionMonitoredDataSource,
		NewAutoTagConditionRootFolderDataSource,
		NewAutoTagConditionYearDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RadarrProvider{
			version: version,
		}
	}
}

// ResourceConfigure is a helper function to set the client for a specific resource.
func resourceConfigure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) (context.Context, *radarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(context.Context)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected context.Context, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData, radarr.NewAPIClient(radarr.NewConfiguration())
}

func dataSourceConfigure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (context.Context, *radarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(context.Context)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected context.Context, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData, radarr.NewAPIClient(radarr.NewConfiguration())
}
