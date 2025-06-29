package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &WhimsyProvider{}

type WhimsyProvider struct {
	version string
}

type WhimsyProviderModel struct{}

func (p *WhimsyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "whimsy"
	resp.Version = p.version
}

func (p *WhimsyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Whimsy provider generates memorable random names using plants, animals, and colors.",
	}
}

func (p *WhimsyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *WhimsyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPlantResource,
		NewAnimalResource,
		NewColorResource,
	}
}

func (p *WhimsyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WhimsyProvider{
			version: version,
		}
	}
}