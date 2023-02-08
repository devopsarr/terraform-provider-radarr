package provider

import (
	"os"
	"testing"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"radarr": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("RADARR_URL"); v == "" {
		t.Skip("RADARR_URL must be set for acceptance tests")
	}

	if v := os.Getenv("RADARR_API_KEY"); v == "" {
		t.Skip("RADARR_API_KEY must be set for acceptance tests")
	}
}

func testAccAPIClient() *radarr.APIClient {
	config := radarr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", os.Getenv("RADARR_API_KEY"))
	config.Servers[0].URL = os.Getenv("RADARR_URL")

	return radarr.NewAPIClient(config)
}

const testUnauthorizedProvider = `
provider "radarr" {
	url = "http://localhost:7878"
	api_key = "ErrorAPIKey"
  }
`
