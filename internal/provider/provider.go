// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"silas.com/ssf-terraform/apisix-client/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ApisixGatewayProvider satisfies various provider interfaces.
var _ provider.Provider = &ApisixGatewayProvider{}
var _ provider.ProviderWithFunctions = &ApisixGatewayProvider{}
var _ provider.ProviderWithEphemeralResources = &ApisixGatewayProvider{}

// ApisixGatewayProvider defines the provider implementation.
type ApisixGatewayProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ApisixGatewayProviderModel describes the provider data model.
type ApisixGatewayProviderModel struct {
	Env types.String `tfsdk:"env"`
}

func (p *ApisixGatewayProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "apisix"
	resp.Version = p.version
}

func (p *ApisixGatewayProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"env": schema.StringAttribute{
				MarkdownDescription: "apisix gateway running env, like dev,uat",
				Optional:            true,
			},
		},
	}
}

func (p *ApisixGatewayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ApisixGatewayProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, ok := os.LookupEnv(api.ApisixHost)
	if !ok || host == "" {
		resp.Diagnostics.AddError(
			"Env 'APISIX_HOST' not set",
			"User must set env 'APISIX_HOST', it represent the addr of apisix gateway.",
		)
	}
	key, ok := os.LookupEnv(api.ApisixKey)
	if !ok || key == "" {
		resp.Diagnostics.AddError(
			"Env 'APISIX_KEY' not set",
			"User must set env 'APISIX_KEY', it contains the authentication info of apisix gateway.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "apisix_gateway_env", data.Env)
	ctx = tflog.SetField(ctx, "apisix_gateway_host", host)
	ctx = tflog.SetField(ctx, "apisix_gateway_key", key)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "apisix_gateway_key")

	client := api.NewApisixClient()
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Apisix Client", map[string]any{"success": true})
}

func (p *ApisixGatewayProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRouteResource,
	}
}

func (p *ApisixGatewayProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *ApisixGatewayProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ApisixGatewayProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ApisixGatewayProvider{
			version: version,
		}
	}
}
