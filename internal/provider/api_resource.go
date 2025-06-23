package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &apiResource{}
	_ resource.ResourceWithConfigure = &apiResource{}
)

func NewApiResource() resource.Resource {
	return &apiResource{}
}

type apiResource struct {
	providerData *apiResourceProviderModel
}

type apiResourceModel struct {
	ID          types.String `tfsdk:"id"`
	EnpointPath types.String `tfsdk:"endpoint_path"`
	Payload     types.String `tfsdk:"payload"`
	Response    types.String `tfsdk:"response"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *apiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg, ok := req.ProviderData.(*apiResourceProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *provider.NttdataTerraformRestApiProviderModel. got %T.", req.ProviderData),
		)
		return
	}

	r.providerData = providerCfg
}

// Create implements resource.Resource.
func (r *apiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "resource Create func called")

	var data apiResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "error in data")
		return
	}

	fullUrl := strings.TrimSuffix(r.providerData.BaseURL.ValueString(), "/") + "/" + strings.TrimPrefix(data.EnpointPath.ValueString(), "/")
	tflog.Debug(ctx, "api resource fullUrl: "+fullUrl)

	body := strings.NewReader(data.Payload.ValueString())
	httpReq, err := http.NewRequest("GET", fullUrl, body)
	if err != nil {
		resp.Diagnostics.AddError("error creating request", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if !r.providerData.AuthToken.IsNull() {
		httpReq.Header.Set("Authorization", "Bearer "+r.providerData.AuthToken.ValueString())
	}

	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("error sending request", err.Error())
		return
	}
	defer httpResp.Body.Close()

	respBody, _ := io.ReadAll(httpResp.Body)
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		resp.Diagnostics.AddError("error sending request", "status code: "+httpResp.Status+" body: "+string(respBody))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("req-%d", time.Now().Unix()))
	data.Response = types.StringValue(string(respBody))

	resp.State.Set(ctx, &data)

	tflog.Debug(ctx, "resource Create func finished with payload: "+string(respBody))
}

// Delete implements resource.Resource.
func (r *apiResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("apiResource Delete unimplemented")
}

// Metadata implements resource.Resource.
func (r *apiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "resource Metadata func called")
	resp.TypeName = req.ProviderTypeName + "_apiresource"
}

// Read implements resource.Resource.
func (r *apiResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("apiResource Read unimplemented")
}

// Schema implements resource.Resource.
func (r *apiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "resource Schema func called")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"endpoint_path": schema.StringAttribute{
				Required: true,
			},
			"payload": schema.StringAttribute{
				Required: true,
			},
			"response": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Update implements resource.Resource.
func (r *apiResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
