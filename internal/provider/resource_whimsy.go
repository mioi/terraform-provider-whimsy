package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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