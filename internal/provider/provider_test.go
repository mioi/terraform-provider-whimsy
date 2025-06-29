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

func TestDataSources(t *testing.T) {
	// Test that data sources can be instantiated
	plantDS := NewPlantDataSource()
	animalDS := NewAnimalDataSource()
	colorDS := NewColorDataSource()

	if plantDS == nil {
		t.Error("Plant data source should not be nil")
	}
	if animalDS == nil {
		t.Error("Animal data source should not be nil")
	}
	if colorDS == nil {
		t.Error("Color data source should not be nil")
	}
}
