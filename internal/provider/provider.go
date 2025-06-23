package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &apiResourceProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &apiResourceProvider{
			version: version,
		}
	}
}

// apiResourceProvider is the provider implementation.
type apiResourceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type apiResourceProviderModel struct {
	BaseURL   types.String `tfsdk:"base_url"`
	AuthToken types.String `tfsdk:"auth_token"`
}

// Metadata returns the provider type name.
func (p *apiResourceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nttdata-rest-api"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *apiResourceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "Base URL of the API",
				Required:    true,
			},
			"auth_token": schema.StringAttribute{
				Description: "Authentication token for the API",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a apiResourceProvider API client for data sources and resources.
func (p *apiResourceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config apiResourceProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = &config
	resp.ResourceData = &config
}

// DataSources defines the data sources implemented in the provider.
func (p *apiResourceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *apiResourceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiResource,
	}
}
