package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"silas.com/ssf-terraform/apisix-client/api"
	"silas.com/ssf-terraform/apisix-client/model"
	"strconv"
	"strings"
)

const InvalidUpstreamHost = "invalid"
const RewriteUpstreamHost = "rewrite"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UpstreamResource{}
var _ resource.ResourceWithImportState = &UpstreamResource{}

func NewUpstreamResource() resource.Resource {
	return &UpstreamResource{}
}

type UpstreamResource struct {
	client *api.ApisixClient
}

type UpstreamResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Type         types.String `tfsdk:"type"`
	Nodes        [][]string   `tfsdk:"nodes"`
	Retries      types.Int32  `tfsdk:"retries"`
	Name         types.String `tfsdk:"name"`
	Desc         types.String `tfsdk:"desc"`
	PassHost     types.String `tfsdk:"pass_host"`
	UpstreamHost types.String `tfsdk:"upstream_host"`
}

func (r *UpstreamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_upstream"
}

func (r *UpstreamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "upstream resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway upstream ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway upstream type",
				Optional:            true,
			},
			"nodes": schema.ListAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "Apisix gateway upstream nodes",
				Optional:            true,
			},
			"retries": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Apisix gateway upstream retries",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway upstream name",
				Optional:            true,
			},
			"desc": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway upstream desc",
				Optional:            true,
			},
			"pass_host": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway upstream pass host",
				Optional:            true,
			},
			"upstream_host": schema.StringAttribute{
				MarkdownDescription: "Apisix gateway upstream upstream host",
				Optional:            true,
			},
		},
	}
}

func (r *UpstreamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func arrayToMap(arrays [][]string) map[string]int {
	mappings := make(map[string]int)
	for _, array := range arrays {
		host := array[0]
		port := array[1]
		priority, _ := strconv.Atoi(array[2])
		mappings[host+":"+port] = priority
	}
	return mappings
}

func mapToArray(mappings map[string]int) [][]string {
	arrays := make([][]string, 0, len(mappings))
	for k, v := range mappings {
		array := make([]string, 0, 3)
		host := strings.Split(k, ":")[0]
		port := strings.Split(k, ":")[1]
		priority := strconv.Itoa(v)
		array = append(array, host)
		array = append(array, port)
		array = append(array, priority)
		arrays = append(arrays, array)
	}
	return arrays
}

func (r *UpstreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UpstreamResourceModel
	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	upstream := &model.Upstream{
		ID:           data.ID.ValueString(),
		Type:         data.Type.ValueString(),
		Nodes:        arrayToMap(data.Nodes),
		Retries:      int(data.Retries.ValueInt32()),
		Name:         data.Name.ValueString(),
		Desc:         data.Desc.ValueString(),
		PassHost:     data.PassHost.ValueString(),
		UpstreamHost: data.UpstreamHost.ValueString(),
	}

	createdUpstream, err := r.client.CreateUpstreams(upstream)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating upstream",
			"Could not create upstream, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the plan with create API response and save it to .tfstate file
	// API response reflects the latest state of the route
	data.ID = types.StringValue(createdUpstream.ID)
	data.Nodes = mapToArray(createdUpstream.Nodes)
	data.Retries = types.Int32Value(int32(createdUpstream.Retries))
	data.Name = types.StringValue(createdUpstream.Name)
	data.Desc = types.StringValue(createdUpstream.Desc)
	data.PassHost = types.StringValue(createdUpstream.PassHost)
	data.UpstreamHost = types.StringValue(createdUpstream.UpstreamHost)
	if createdUpstream.PassHost != RewriteUpstreamHost {
		data.UpstreamHost = types.StringValue(InvalidUpstreamHost)
	}

	tflog.Trace(ctx, "created a resource "+data.ID.ValueString())
	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *UpstreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UpstreamResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fetchedUpstream, err := r.client.GetUpstreamById(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting upstream",
			"Could not get upstream, unexpected error: "+err.Error(),
		)
		return
	}

	data.ID = types.StringValue(fetchedUpstream.ID)
	data.Nodes = mapToArray(fetchedUpstream.Nodes)
	data.Retries = types.Int32Value(int32(fetchedUpstream.Retries))
	data.Name = types.StringValue(fetchedUpstream.Name)
	data.Desc = types.StringValue(fetchedUpstream.Desc)
	data.PassHost = types.StringValue(fetchedUpstream.PassHost)
	data.UpstreamHost = types.StringValue(fetchedUpstream.UpstreamHost)
	if fetchedUpstream.PassHost != RewriteUpstreamHost {
		data.UpstreamHost = types.StringValue(InvalidUpstreamHost)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UpstreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UpstreamResourceModel
	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	upstream := &model.Upstream{
		ID:           data.ID.ValueString(),
		Type:         data.Type.ValueString(),
		Nodes:        arrayToMap(data.Nodes),
		Retries:      int(data.Retries.ValueInt32()),
		Name:         data.Name.ValueString(),
		Desc:         data.Desc.ValueString(),
		PassHost:     data.PassHost.ValueString(),
		UpstreamHost: data.UpstreamHost.ValueString(),
	}

	createdUpstream, err := r.client.UpdateUpstream(upstream)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating upstream",
			"Could not create upstream, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the plan with create API response and save it to .tfstate file
	// API response reflects the latest state of the route
	data.ID = types.StringValue(createdUpstream.ID)
	data.Nodes = mapToArray(createdUpstream.Nodes)
	data.Retries = types.Int32Value(int32(createdUpstream.Retries))
	data.Name = types.StringValue(createdUpstream.Name)
	data.Desc = types.StringValue(createdUpstream.Desc)
	data.PassHost = types.StringValue(createdUpstream.PassHost)
	data.UpstreamHost = types.StringValue(createdUpstream.UpstreamHost)
	if createdUpstream.PassHost != RewriteUpstreamHost {
		data.UpstreamHost = types.StringValue(InvalidUpstreamHost)
	}

	tflog.Trace(ctx, "created a resource "+data.ID.ValueString())
	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *UpstreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UpstreamResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUpstreamById(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting upstream",
			"Could not delete upstream, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *UpstreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
