package provider

import (
	"os"
	"silas.com/ssf-terraform/apisix-client/api"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestApisixUpstreamResource(t *testing.T) {
	os.Setenv(api.ApisixHost, "http://172.18.21.239:9180")
	os.Setenv(api.ApisixKey, "api-key")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "common" {
    id = "common"
    type = "roundbin"
    nodes = [["127.0.0.1", "80", "1"]]
    retries = 3
    name = "common"
    desc = "Common upstream for all services, forward requests to ingress"
    pass_host = "pass"
    upstream_host = "invalid"
 }
`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apisix_upstream.common", "id", "common"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "nodes.0.0", "127.0.0.1"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "retries", "3"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "name", "common"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "desc", "Common upstream for all services, forward requests to ingress"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "common" {
    id = "common"
    type = "roundbin"
    nodes = [["127.0.0.1", "80", "1"]]
    retries = 1
    name = "common"
    desc = "Common upstream for all services, forward requests to ingress"
    pass_host = "rewrite"
    upstream_host = "127.0.0.2:80"
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apisix_upstream.common", "id", "common"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "nodes.0.0", "127.0.0.1"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "retries", "1"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "name", "common"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "desc", "Common upstream for all services, forward requests to ingress"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "pass_host", "rewrite"),
					resource.TestCheckResourceAttr("apisix_upstream.common", "upstream_host", "127.0.0.2:80"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
