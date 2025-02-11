package tests

import (
	"awesomeProject3/funcs"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func FilterRecipesByLevel(Recipes []funcs.Recipe, RecipeType string) []funcs.Recipe {
	var filtered []funcs.Recipe
	for _, r := range Recipes {
		if r.Level == RecipeType {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
func TestFilterRecipesByLevel(t *testing.T) {
	Recipes := []funcs.Recipe{
		{Name: "Mushroom Risotto 1", Level: "master"},
		{Name: "Vegetable Stir Fry 3", Level: "amateur"},
		{Name: "Vegetable Stir Fry 4", Level: "beginner"},
	}
	RecipeLevel := "master"

	filtered := FilterRecipesByLevel(Recipes, RecipeLevel)
	fmt.Println("TESTING UNIT TEST")
	time.Sleep(1 * time.Second)
	expected := []funcs.Recipe{{Name: "Mushroom Risotto 1", Level: "master"}}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Incorrect filtered Recipes. Expected: %v, Got: %v", expected, filtered)
	}
}
