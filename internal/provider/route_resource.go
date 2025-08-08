// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"silas.com/ssf-terraform/apisix-client/api"
	"silas.com/ssf-terraform/apisix-client/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RouteResource{}
var _ resource.ResourceWithImportState = &RouteResource{}

func NewRouteResource() resource.Resource {
	return &RouteResource{}
}

// RouteResource defines the resource implementation.
type RouteResource struct {
	client *api.ApisixClient
}

// RouteResourceModel describes the resource data model.
type RouteResourceModel struct {
	ID         *string                `tfsdk:"id"`
	Uris       []string               `tfsdk:"uris"`
	UpstreamId *string                `tfsdk:"upstream_id"`
	Plugins    *Plugins               `tfsdk:"plugins"`
	Name       *string                `tfsdk:"name"`
	Desc       *string                `tfsdk:"desc"`
	Hosts      []string               `tfsdk:"hosts"`
	Methods    []string               `tfsdk:"methods"`
	Priority   *int                   `tfsdk:"priority"`
	Vars       [][]interface{}        `tfsdk:"vars"`
	Labels     map[string]interface{} `tfsdk:"labels"`
	Timeout    *Timeout               `tfsdk:"timeout"`
	Status     *int                   `tfsdk:"status"`
}

type Plugins struct {
	OpenIdConnectPlugin *OpenIdConnectPlugin `json:"openid_connect"`
}

type OpenIdConnectPlugin struct {
	ClientId       string   `tfsdk:"client_id"`
	ClientSecret   string   `tfsdk:"client_secret"`
	Discovery      string   `tfsdk:"discovery"`
	RequiredScopes []string `tfsdk:"required_scopes"`
}

type Timeout struct {
	Connect int `json:"connect"`
	Send    int `json:"send"`
	Read    int `json:"read"`
}

func (r *RouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (r *RouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway route ID",
				Required:            true,
			},
			"uris": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Apisix gateway route URIs",
				Optional:            true,
			},
			"upstream_id": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway route upstream ID",
				Optional:            true,
			},
			"plugins": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"openid_connect": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								MarkdownDescription: "Client ID",
								Optional:            true,
							},
							"client_secret": schema.StringAttribute{
								MarkdownDescription: "Client Secret",
								Optional:            true,
							},
							"discovery": schema.StringAttribute{
								MarkdownDescription: "Discovery endpoint",
								Optional:            true,
							},
							"required_scopes": schema.ListAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "Required scopes",
								Required:            true,
							},
						},
						MarkdownDescription: "openid_connect auth plugin",
						Optional:            true,
					},
				},
				MarkdownDescription: "Apisix gateway route plugins",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway route name",
				Optional:            true,
			},
			"desc": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway route desc",
				Optional:            true,
			},
			"hosts": schema.ListAttribute{
				MarkdownDescription: "Apisix gateway route hosts",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"methods": schema.ListAttribute{
				MarkdownDescription: "Apisix gateway route methods",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"priority": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Apisix gateway route priority",
			},
			"vars": schema.ListAttribute{
				ElementType: types.ListType{
					ElemType: types.ObjectType{},
				},
				Optional:            true,
				MarkdownDescription: "Apisix gateway route vars",
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: "Apisix gateway route labels",
				Optional:            true,
				ElementType:         types.ObjectType{},
			},
			"timeout": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"connect": schema.Int64Attribute{
						MarkdownDescription: "Connect timeout",
						Optional:            true,
					},
					"send": schema.Int64Attribute{
						MarkdownDescription: "Send timeout",
						Optional:            true,
					},
					"read": schema.Int64Attribute{
						MarkdownDescription: "Read timeout",
						Optional:            true,
					},
				},
				Optional:            true,
				MarkdownDescription: "Apisix gateway route timeout",
			},
			"status": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Apisix gateway route status",
			},
		},
	}
}

func (r *RouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.ApisixClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ApisixClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *RouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RouteResourceModel
	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	authPlugin := model.OpenIdConnectPlugin{
		ClientId:              data.Plugins.OpenIdConnectPlugin.ClientId,
		ClientSecret:          data.Plugins.OpenIdConnectPlugin.ClientSecret,
		Discovery:             data.Plugins.OpenIdConnectPlugin.Discovery,
		RequiredScopes:        data.Plugins.OpenIdConnectPlugin.RequiredScopes,
		BearerOnly:            true,
		UseJwks:               true,
		JwkExpiresIn:          600,
		AudienceRequired:      true,
		Audience:              "aud",
		AudienceMatchClientId: true,
		Realm:                 "silas-apisix-gateway",
	}
	route := &model.Route{
		ID:         data.ID,
		Uris:       data.Uris,
		UpstreamId: data.UpstreamId,
		Plugins: &model.Plugins{
			OpenIdConnectPlugin: &authPlugin,
		},
		Name:     data.Name,
		Desc:     data.Desc,
		Hosts:    data.Hosts,
		Methods:  data.Methods,
		Priority: data.Priority,
		Vars:     data.Vars,
		Labels:   data.Labels,
		Timeout: &model.Timeout{
			Connect: data.Timeout.Connect,
			Send:    data.Timeout.Send,
			Read:    data.Timeout.Read,
		},
		Status: data.Status,
	}

	createdRoute, err := r.client.CreateRoute(route)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating route",
			"Could not create route, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the plan with create API response and save it to .tfstate file
	// API response reflects the latest state of the route
	data.ID = createdRoute.ID
	data.Uris = createdRoute.Uris
	data.UpstreamId = createdRoute.UpstreamId
	data.Plugins = &Plugins{
		OpenIdConnectPlugin: &OpenIdConnectPlugin{
			ClientId:       createdRoute.Plugins.OpenIdConnectPlugin.ClientId,
			ClientSecret:   createdRoute.Plugins.OpenIdConnectPlugin.ClientSecret,
			Discovery:      createdRoute.Plugins.OpenIdConnectPlugin.Discovery,
			RequiredScopes: createdRoute.Plugins.OpenIdConnectPlugin.RequiredScopes,
		},
	}
	data.Name = createdRoute.Name
	data.Desc = createdRoute.Desc
	data.Hosts = createdRoute.Hosts
	data.Methods = createdRoute.Methods
	data.Priority = createdRoute.Priority
	data.Vars = createdRoute.Vars
	data.Labels = createdRoute.Labels
	data.Timeout = &Timeout{
		Connect: createdRoute.Timeout.Connect,
		Send:    createdRoute.Timeout.Send,
		Read:    createdRoute.Timeout.Read,
	}
	data.Status = createdRoute.Status

	tflog.Trace(ctx, "created a resource "+data.ID)
	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *RouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.client.GetRouteById(data.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting route",
			"Could not get route, unexpected error: "+err.Error(),
		)
		return
	}

	data.ID = route.ID
	data.Uris = route.Uris
	data.UpstreamId = route.UpstreamId
	data.Plugins = &Plugins{
		OpenIdConnectPlugin: &OpenIdConnectPlugin{
			ClientId:       route.Plugins.OpenIdConnectPlugin.ClientId,
			ClientSecret:   route.Plugins.OpenIdConnectPlugin.ClientSecret,
			Discovery:      route.Plugins.OpenIdConnectPlugin.Discovery,
			RequiredScopes: route.Plugins.OpenIdConnectPlugin.RequiredScopes,
		},
	}
	data.Name = route.Name
	data.Desc = route.Desc
	data.Hosts = route.Hosts
	data.Methods = route.Methods
	data.Priority = route.Priority
	data.Vars = route.Vars
	data.Labels = route.Labels
	data.Timeout = &Timeout{
		Connect: route.Timeout.Connect,
		Send:    route.Timeout.Send,
		Read:    route.Timeout.Read,
	}
	data.Status = route.Status

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RouteResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	authPlugin := model.OpenIdConnectPlugin{
		ClientId:       data.Plugins.OpenIdConnectPlugin.ClientId,
		ClientSecret:   data.Plugins.OpenIdConnectPlugin.ClientSecret,
		Discovery:      data.Plugins.OpenIdConnectPlugin.Discovery,
		RequiredScopes: data.Plugins.OpenIdConnectPlugin.RequiredScopes,
	}
	route := &model.Route{
		ID:         data.ID,
		Uris:       data.Uris,
		UpstreamId: data.UpstreamId,
		Plugins: &model.Plugins{
			OpenIdConnectPlugin: &authPlugin,
		},
		Name:     data.Name,
		Desc:     data.Desc,
		Hosts:    data.Hosts,
		Methods:  data.Methods,
		Priority: data.Priority,
		Vars:     data.Vars,
		Labels:   data.Labels,
		Timeout: &model.Timeout{
			Connect: data.Timeout.Connect,
			Send:    data.Timeout.Send,
			Read:    data.Timeout.Read,
		},
		Status: data.Status,
	}

	updatedRoute, err := r.client.UpdateRoute(route)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating route",
			"Could not create route, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the plan with create API response and save it to .tfstate file
	// API response reflects the latest state of the route
	data.ID = updatedRoute.ID
	data.Uris = updatedRoute.Uris
	data.UpstreamId = updatedRoute.UpstreamId
	data.Plugins = &Plugins{
		OpenIdConnectPlugin: &OpenIdConnectPlugin{
			ClientId:       updatedRoute.Plugins.OpenIdConnectPlugin.ClientId,
			ClientSecret:   updatedRoute.Plugins.OpenIdConnectPlugin.ClientSecret,
			Discovery:      updatedRoute.Plugins.OpenIdConnectPlugin.Discovery,
			RequiredScopes: updatedRoute.Plugins.OpenIdConnectPlugin.RequiredScopes,
		},
	}
	data.Name = updatedRoute.Name
	data.Desc = updatedRoute.Desc
	data.Hosts = updatedRoute.Hosts
	data.Methods = updatedRoute.Methods
	data.Priority = updatedRoute.Priority
	data.Vars = updatedRoute.Vars
	data.Labels = updatedRoute.Labels
	data.Timeout = &Timeout{
		Connect: updatedRoute.Timeout.Connect,
		Send:    updatedRoute.Timeout.Send,
		Read:    updatedRoute.Timeout.Read,
	}
	data.Status = updatedRoute.Status

	tflog.Trace(ctx, "created a resource "+data.ID)
	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *RouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRouteById(data.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting route",
			"Could not delete route, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *RouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
