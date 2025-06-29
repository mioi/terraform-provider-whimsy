package provider

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

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

	// Generate deterministic name based on triggers
	generatedName, err := d.generateDeterministicName(data.Triggers)
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

// generateDeterministicName creates a deterministic name based on triggers
func (d *WhimsyDataSource) generateDeterministicName(triggers types.Map) (string, error) {
	// Create a seed based on triggers + data source type
	seedString := d.typeName

	if !triggers.IsNull() {
		// Convert triggers to a deterministic string
		triggerElements := triggers.Elements()
		for key, value := range triggerElements {
			if strValue, ok := value.(types.String); ok {
				seedString += fmt.Sprintf("%s=%s;", key, strValue.ValueString())
			}
		}
	}

	// Create deterministic seed from string
	hash := sha256.Sum256([]byte(seedString))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))

	// Use seed for deterministic selection
	return d.generateFromSeed(seed)
}

// generateFromSeed generates a name using a deterministic seed
func (d *WhimsyDataSource) generateFromSeed(seed int64) (string, error) {
	// Get all words for this data source type
	var wordList []string
	switch d.typeName {
	case "plant":
		wordList = whimsy.Plants()
	case "animal":
		wordList = whimsy.Animals()
	case "color":
		wordList = whimsy.Colors()
	default:
		return "", fmt.Errorf("unknown data source type: %s", d.typeName)
	}

	if len(wordList) == 0 {
		return "", fmt.Errorf("no words available for %s", d.typeName)
	}

	// Create seeded random generator
	rng := rand.New(rand.NewSource(seed))
	index := rng.Intn(len(wordList))

	return wordList[index], nil
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
