package provider

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mioi/whimsy"
)

var _ resource.Resource = &WhimsyResource{}

type WhimsyResource struct {
	typeName    string
	description string
	generator   func() (string, error)
}

type WhimsyResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Triggers types.Map    `tfsdk:"triggers"`
}

func NewWhimsyResource(typeName, description string, generator func() (string, error)) resource.Resource {
	return &WhimsyResource{
		typeName:    typeName,
		description: description,
		generator:   generator,
	}
}

func (r *WhimsyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *WhimsyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random " + r.description + " that persists until triggers change.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this resource instance.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The generated " + r.description + ".",
				Computed:            true,
			},
			"triggers": schema.MapAttribute{
				MarkdownDescription: "Arbitrary map of values that, when changed, will trigger regeneration of the " + r.description + ".",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *WhimsyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *WhimsyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WhimsyResourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating "+r.typeName+" resource")

	// Generate random name using whimsy library
	generatedName, err := r.generator()
	if err != nil {
		resp.Diagnostics.AddError("Random Generation Error", "Unable to generate random "+r.description+": "+err.Error())
		return
	}

	data.Name = types.StringValue(generatedName)
	data.Id = types.StringValue(generatedName)

	tflog.Debug(ctx, "Created "+r.typeName+" resource", map[string]any{
		"name": generatedName,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WhimsyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WhimsyResourceModel

	// Read prior state
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading "+r.typeName+" resource", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Resource exists as-is, no changes needed
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WhimsyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WhimsyResourceModel
	var state WhimsyResourceModel

	// Read configuration and state
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating "+r.typeName+" resource")

	// Check if triggers changed
	if !data.Triggers.Equal(state.Triggers) {
		// Regenerate name when triggers change
		generatedName, err := r.generator()
		if err != nil {
			resp.Diagnostics.AddError("Random Generation Error", "Unable to generate random "+r.description+": "+err.Error())
			return
		}
		
		data.Name = types.StringValue(generatedName)
		data.Id = types.StringValue(generatedName)
		
		tflog.Debug(ctx, "Regenerated "+r.typeName+" due to trigger change", map[string]any{
			"new_name": generatedName,
		})
	} else {
		// Keep existing name if triggers haven't changed
		data.Name = state.Name
		data.Id = state.Id
		
		tflog.Debug(ctx, "Keeping existing "+r.typeName+" name", map[string]any{
			"name": state.Name.ValueString(),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WhimsyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: removing a random name from state doesn't require any external cleanup
	tflog.Debug(ctx, "Deleted "+r.typeName+" resource")
}

// Factory functions for each type
func NewPlantResource() resource.Resource {
	return NewWhimsyResource("plant", "plant name", whimsy.RandomPlant)
}

func NewAnimalResource() resource.Resource {
	return NewWhimsyResource("animal", "animal name", whimsy.RandomAnimal)
}

func NewColorResource() resource.Resource {
	return NewWhimsyResource("color", "color name", whimsy.RandomColor)
}

// WhimsyNameResource is a specialized resource for combining multiple parts
var _ resource.Resource = &WhimsyNameResource{}

type WhimsyNameResource struct{}

type WhimsyNameResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Parts     types.List   `tfsdk:"parts"`
	Delimiter types.String `tfsdk:"delimiter"`
	Random    types.Bool   `tfsdk:"random"`
	Triggers  types.Map    `tfsdk:"triggers"`
}

func NewWhimsyNameResource() resource.Resource {
	return &WhimsyNameResource{}
}

func (r *WhimsyNameResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

func (r *WhimsyNameResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random name by combining multiple parts (plants, animals, colors) with a delimiter.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this resource instance.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The generated combined name.",
				Computed:            true,
			},
			"parts": schema.ListAttribute{
				MarkdownDescription: "List of name parts to combine. Valid values: 'plant', 'animal', 'color'. Default: ['color', 'animal'].",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("color"), types.StringValue("animal")})),
			},
			"delimiter": schema.StringAttribute{
				MarkdownDescription: "String used to separate each part. Default: '-'.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("-"),
			},
			"random": schema.BoolAttribute{
				MarkdownDescription: "If true, ignores the order of parts and randomizes their arrangement. Default: false.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"triggers": schema.MapAttribute{
				MarkdownDescription: "Arbitrary map of values that, when changed, will trigger regeneration of the name.",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *WhimsyNameResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *WhimsyNameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WhimsyNameResourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating whimsy_name resource")

	// Generate the combined name
	generatedName, err := r.generateCombinedName(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Random Generation Error", "Unable to generate random name: "+err.Error())
		return
	}

	data.Name = types.StringValue(generatedName)
	data.Id = types.StringValue(generatedName)

	// Set the computed values from defaults if not provided
	if data.Parts.IsNull() {
		data.Parts = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("color"), types.StringValue("animal")})
	}
	if data.Delimiter.IsNull() {
		data.Delimiter = types.StringValue("-")
	}
	if data.Random.IsNull() {
		data.Random = types.BoolValue(false)
	}

	tflog.Debug(ctx, "Created whimsy_name resource", map[string]any{
		"name": generatedName,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WhimsyNameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WhimsyNameResourceModel

	// Read prior state
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading whimsy_name resource", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Resource exists as-is, no changes needed
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WhimsyNameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WhimsyNameResourceModel
	var state WhimsyNameResourceModel

	// Read configuration and state
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating whimsy_name resource")

	// Check if any configuration changed that would require regeneration
	shouldRegenerate := !data.Parts.Equal(state.Parts) ||
		!data.Delimiter.Equal(state.Delimiter) ||
		!data.Random.Equal(state.Random) ||
		!data.Triggers.Equal(state.Triggers)

	if shouldRegenerate {
		// Regenerate name when configuration or triggers change
		generatedName, err := r.generateCombinedName(ctx, data)
		if err != nil {
			resp.Diagnostics.AddError("Random Generation Error", "Unable to generate random name: "+err.Error())
			return
		}

		data.Name = types.StringValue(generatedName)
		data.Id = types.StringValue(generatedName)

		tflog.Debug(ctx, "Regenerated whimsy_name due to configuration change", map[string]any{
			"new_name": generatedName,
		})
	} else {
		// Keep existing name if nothing changed
		data.Name = state.Name
		data.Id = state.Id

		tflog.Debug(ctx, "Keeping existing whimsy_name", map[string]any{
			"name": state.Name.ValueString(),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WhimsyNameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: removing a random name from state doesn't require any external cleanup
	tflog.Debug(ctx, "Deleted whimsy_name resource")
}

// generateCombinedName creates a name by combining multiple parts
func (r *WhimsyNameResource) generateCombinedName(ctx context.Context, data WhimsyNameResourceModel) (string, error) {
	// Extract parts from the list, use defaults if null
	var parts []string
	if data.Parts.IsNull() {
		parts = []string{"color", "animal"}
	} else {
		for _, element := range data.Parts.Elements() {
			if strVal, ok := element.(types.String); ok {
				parts = append(parts, strVal.ValueString())
			}
		}
	}

	// Get delimiter, use default if null
	delimiter := "-"
	if !data.Delimiter.IsNull() {
		delimiter = data.Delimiter.ValueString()
	}

	// Get random flag, use default if null
	randomOrder := false
	if !data.Random.IsNull() {
		randomOrder = data.Random.ValueBool()
	}

	// Validate parts list is not empty
	if len(parts) == 0 {
		return "", fmt.Errorf("parts list cannot be empty: must contain at least one of 'plant', 'animal', 'color'")
	}

	// Validate parts
	validParts := map[string]bool{"plant": true, "animal": true, "color": true}
	for _, part := range parts {
		if !validParts[part] {
			return "", fmt.Errorf("invalid part '%s': must be one of 'plant', 'animal', 'color'", part)
		}
	}

	// If random is true, shuffle the parts
	if randomOrder {
		shuffledParts := make([]string, len(parts))
		copy(shuffledParts, parts)
		
		// Fisher-Yates shuffle using crypto/rand
		for i := len(shuffledParts) - 1; i > 0; i-- {
			j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
			if err != nil {
				return "", fmt.Errorf("failed to generate random number: %w", err)
			}
			shuffledParts[i], shuffledParts[int(j.Int64())] = shuffledParts[int(j.Int64())], shuffledParts[i]
		}
		parts = shuffledParts
	}

	// Generate a name for each part
	var nameParts []string
	for _, part := range parts {
		var name string
		var err error

		switch part {
		case "plant":
			name, err = whimsy.RandomPlant()
		case "animal":
			name, err = whimsy.RandomAnimal()
		case "color":
			name, err = whimsy.RandomColor()
		default:
			return "", fmt.Errorf("unknown part type: %s", part)
		}

		if err != nil {
			return "", fmt.Errorf("failed to generate %s name: %w", part, err)
		}

		nameParts = append(nameParts, name)
	}

	// Join with delimiter
	combinedName := strings.Join(nameParts, delimiter)

	tflog.Debug(ctx, "Generated combined name", map[string]any{
		"parts":      parts,
		"nameParts":  nameParts,
		"delimiter":  delimiter,
		"random":     randomOrder,
		"result":     combinedName,
	})

	return combinedName, nil
}