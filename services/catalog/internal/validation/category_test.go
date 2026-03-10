package validation

import (
	"testing"
)

func TestValidateCategoryInput_Success(t *testing.T) {
	input := &CategoryInput{
		Name:     "Eletrônicos",
		Slug:     "eletronicos",
		ParentID: "",
	}
	if err := ValidateCategoryInput(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Slug != "eletronicos" {
		t.Errorf("expected normalized slug, got %s", input.Slug)
	}
}

func TestValidateCategoryInput_NameRequired(t *testing.T) {
	input := &CategoryInput{Name: "", Slug: "test"}
	if err := ValidateCategoryInput(input); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateCategoryInput_SlugRequired(t *testing.T) {
	input := &CategoryInput{Name: "Test", Slug: ""}
	if err := ValidateCategoryInput(input); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateCategoryInput_InvalidSlugFormat(t *testing.T) {
	input := &CategoryInput{Name: "Test", Slug: "Invalid_Slug!"}
	if err := ValidateCategoryInput(input); err == nil {
		t.Fatal("expected validation error for invalid slug")
	}
}
