package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	sonicclient "terraform-provider-community-sonic/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &sonicProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sonicProvider{
			version: version,
		}
	}
}

// sonicProvider is the provider implementation.
type sonicProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *sonicProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sonic"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *sonicProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Terraform provider for Community SONiC " +
			"can be used to interact with a device running SONiC in order to manage resources.",
		MarkdownDescription: "The Terraform provider for Community SONiC " +
			"can be used to interact with a device running SONiC in order to manage resources.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "IP or FQDN of the SONiC device",
				Description:         "IP or FQDN of the SONiC device",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the SONiC device",
				Description:         "The username of the SONiC device",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the SONiC device",
				Description:         "The password of the SONiC device",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Boolean variable to specify whether to validate SSL certificate or not.",
				Description:         "Boolean variable to specify whether to validate SSL certificate or not.",
				Optional:            true,
			},
		},
	}
}

// Configure prepares an API client for data sources and resources.
func (p *sonicProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config sonicProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown SONiC API Host",
			"The provider cannot create the SONiC API client as there is an unknown configuration value for the SONiC API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SONIC_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown SONiC API Username",
			"The provider cannot create the SONiC API client as there is an unknown configuration value for the SONiC API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SONIC_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown SONiC API Password",
			"The provider cannot create the SONiC API client as there is an unknown configuration value for the SONiC API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SONIC_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("SONIC_HOST")
	username := os.Getenv("SONIC_USERNAME")
	password := os.Getenv("SONIC_PASSWORD")
	insecure := false

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing SONiC API Host",
			"The provider cannot create the SONiC API client as there is a missing or empty value for the SONiC API host. "+
				"Set the host value in the configuration or use the SONIC_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing SONiC API Username",
			"The provider cannot create the SONiC API client as there is a missing or empty value for the SONiC API username. "+
				"Set the username value in the configuration or use the SONIC_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing SONiC API Password",
			"The provider cannot create the SONiC API client as there is a missing or empty value for the SONiC API password. "+
				"Set the password value in the configuration or use the SONIC_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new SONiC client using the configuration values
	client, err := sonicclient.NewClient(&host, &username, &password, insecure)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SONiC API Client",
			"An unexpected error occurred when creating the SONiC API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"SONiC Client Error: "+err.Error(),
		)
		return
	}

	// Make the SONiC client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *sonicProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSonicDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *sonicProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// sonicProviderModel maps provider schema data to a Go type.
type sonicProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}
