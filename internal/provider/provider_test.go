package provider

import (
	"testing"

	"github.com/mioi/whimsy"
)

func TestProvider(t *testing.T) {
	// Simple test to verify provider can be instantiated
	provider := New("test")()
	if provider == nil {
		t.Error("Provider should not be nil")
	}
}

func validateNameList(t *testing.T, names []string, nameType string) {
	if len(names) < 200 {
		t.Errorf("Expected at least 200 %s names, got %d", nameType, len(names))
	}

	for _, name := range names {
		if len(name) > 6 {
			t.Errorf("%s name '%s' is longer than 6 characters", nameType, name)
		}
		for _, char := range name {
			if char < 'a' || char > 'z' {
				t.Errorf("%s name '%s' contains invalid character '%c'", nameType, name, char)
			}
		}
	}

	// Check alphabetical order
	for i := 1; i < len(names); i++ {
		if names[i-1] >= names[i] {
			t.Errorf("%s names are not in alphabetical order: '%s' should come after '%s'", nameType, names[i-1], names[i])
		}
	}

	// Check for duplicates
	for i := 1; i < len(names); i++ {
		if names[i-1] == names[i] {
			t.Errorf("%s names have duplicates: '%s' found multiple times", nameType, names[i-1])
		}
	}
}

func TestPlantNames(t *testing.T) {
	validateNameList(t, whimsy.Plants(), "Plant")
}

func TestAnimalNames(t *testing.T) {
	validateNameList(t, whimsy.Animals(), "Animal")
}

func TestColorNames(t *testing.T) {
	validateNameList(t, whimsy.Colors(), "Color")
}

func TestResources(t *testing.T) {
	// Test that resources can be instantiated
	plantRes := NewPlantResource()
	animalRes := NewAnimalResource()
	colorRes := NewColorResource()

	if plantRes == nil {
		t.Error("Plant resource should not be nil")
	}
	if animalRes == nil {
		t.Error("Animal resource should not be nil")
	}
	if colorRes == nil {
		t.Error("Color resource should not be nil")
	}
}

func TestResourceGeneration(t *testing.T) {
	// Test that resources can generate valid names
	tests := []struct {
		name      string
		generator func() (string, error)
		wordList  []string
	}{
		{"plant", whimsy.RandomPlant, whimsy.Plants()},
		{"animal", whimsy.RandomAnimal, whimsy.Animals()},
		{"color", whimsy.RandomColor, whimsy.Colors()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Generate multiple names to test randomness
			generatedNames := make(map[string]bool)
			for i := 0; i < 10; i++ {
				name, err := test.generator()
				if err != nil {
					t.Errorf("%s generator returned error: %v", test.name, err)
					continue
				}

				if name == "" {
					t.Errorf("%s generator returned empty string", test.name)
					continue
				}

				// Verify name is in the word list
				found := false
				for _, word := range test.wordList {
					if word == name {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s generator returned '%s' which is not in the %s list", test.name, name, test.name)
				}

				generatedNames[name] = true
			}

			// Should have at least some variety (not all the same name)
			if len(generatedNames) < 2 {
				t.Errorf("%s generator should produce some variety, got %d unique names from 10 generations", test.name, len(generatedNames))
			}
		})
	}
}
