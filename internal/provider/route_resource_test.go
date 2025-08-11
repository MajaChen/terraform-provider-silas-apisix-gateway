// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"silas.com/ssf-terraform/apisix-client/api"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestApisixRouteResource(t *testing.T) {
	os.Setenv(api.ApisixHost, "http://172.18.21.239:9180")
	os.Setenv(api.ApisixKey, "api-key")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_route" "ssf-java-sdk-springboot3-demo-dynLoggingLevel" {
    id = "ssf-java-sdk-springboot3-demo-dynLoggingLevel"
    uris = ["/api/v1/demo/dynLoggingLevel"]
    upstream_id = "1"
    plugins = {
      openid_connect = {
        client_id = "client-id"
        discovery = "https://domain.authing.cn/oidc/.well-known/jwks.json"
        required_scopes = ["admin", "book"]    
       }
    }
    name = "ssf-java-sdk-springboot3-demo-dynLoggingLevel"
    desc = "ssf-java-sdk-springboot3-demo dynLoggingLevel"
    methods = ["GET", "POST"]
    priority = 10
    timeout = {
      connect = 10
      send = 10
      read = 10
    }
    status = 1
 }
`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "id", "ssf-java-sdk-springboot3-demo-dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "uris.0", "/api/v1/demo/dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "upstream_id", "1"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.client_id", "client-id"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.discovery", "https://domain.authing.cn/oidc/.well-known/jwks.json"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.required_scopes.0", "admin"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.required_scopes.1", "book"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "name", "ssf-java-sdk-springboot3-demo-dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "desc", "ssf-java-sdk-springboot3-demo dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "methods.0", "GET"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "priority", "10"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "timeout.send", "10"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "timeout.read", "10"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "status", "1"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_route" "ssf-java-sdk-springboot3-demo-dynLoggingLevel" {
    id = "ssf-java-sdk-springboot3-demo-dynLoggingLevel"
    uris = ["/api/v1/demo/dynLoggingLevel"]
    upstream_id = "1"
    plugins = {
      openid_connect = {
        client_id = "client-id"
        discovery = "https://domain.authing.cn/oidc/.well-known/jwks.json"
        required_scopes = ["admin", "book", "stuff"]    
       }
    }
    name = "ssf-java-sdk-springboot3-demo-dynLoggingLevel"
    desc = "ssf-java-sdk-springboot3-demo dynLoggingLevel"
    methods = ["GET", "POST"]
    priority = 20
    timeout = {
      connect = 10
      send = 30
      read = 30
    }
    status = 1
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "id", "ssf-java-sdk-springboot3-demo-dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "uris.0", "/api/v1/demo/dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "upstream_id", "1"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.client_id", "client-id"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.discovery", "https://domain.authing.cn/oidc/.well-known/jwks.json"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.required_scopes.0", "admin"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.required_scopes.1", "book"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "plugins.openid_connect.required_scopes.2", "stuff"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "name", "ssf-java-sdk-springboot3-demo-dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "desc", "ssf-java-sdk-springboot3-demo dynLoggingLevel"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "methods.0", "GET"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "priority", "20"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "timeout.send", "30"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "timeout.read", "30"),
					resource.TestCheckResourceAttr("apisix_route.ssf-java-sdk-springboot3-demo-dynLoggingLevel", "status", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
