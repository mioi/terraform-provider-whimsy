package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mioi/whimsy"
)

var _ datasource.DataSource = &WhimsyDataSource{}

type WhimsyDataSource struct {
	typeName    string
	description string
	generator   func() (string, error)
}

type WhimsyDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Triggers types.Map    `tfsdk:"triggers"`
}

func NewWhimsyDataSource(typeName, description string, generator func() (string, error)) datasource.DataSource {
	return &WhimsyDataSource{
		typeName:    typeName,
		description: description,
		generator:   generator,
	}
}

func (d *WhimsyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.typeName
}

func (d *WhimsyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random " + d.description + ".",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this data source instance.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The generated " + d.description + ".",
				Computed:            true,
			},
			"triggers": schema.MapAttribute{
				MarkdownDescription: "Arbitrary map of values that, when changed, will trigger regeneration of the " + d.description + ".",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (d *WhimsyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *WhimsyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WhimsyDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading "+d.typeName+" data source")

	// Generate random name
	generatedName, err := d.generator()
	if err != nil {
		resp.Diagnostics.AddError("Random Generation Error", "Unable to generate random "+d.description+": "+err.Error())
		return
	}
	data.Name = types.StringValue(generatedName)
	data.Id = types.StringValue(generatedName)

	tflog.Debug(ctx, "Generated "+d.typeName+" name", map[string]any{
		"name": generatedName,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Factory functions for each type
func NewPlantDataSource() datasource.DataSource {
	return NewWhimsyDataSource("plant", "plant name", whimsy.RandomPlant)
}

func NewAnimalDataSource() datasource.DataSource {
	return NewWhimsyDataSource("animal", "animal name", whimsy.RandomAnimal)
}

func NewColorDataSource() datasource.DataSource {
	return NewWhimsyDataSource("color", "color name", whimsy.RandomColor)
}
