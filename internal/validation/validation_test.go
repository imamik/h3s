package validation

import (
	"testing"
)

type BasicStruct struct {
	Name  string `validate:"required,min=4,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18,max=99"`
}

type NestedStruct struct {
	Inner BasicStruct `validate:"required"`
}

func TestValidateStruct_Valid(t *testing.T) {
	b := BasicStruct{Name: "Alice", Email: "alice@example.com", Age: 30}
	err := ValidateStruct(b)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateStruct_MissingRequired(t *testing.T) {
	b := BasicStruct{Email: "alice@example.com", Age: 30}
	err := ValidateStruct(b)
	if err == nil {
		t.Error("expected error for missing required field 'Name', got nil")
	}
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	b := BasicStruct{Name: "Alice", Email: "not-an-email", Age: 30}
	err := ValidateStruct(b)
	if err == nil {
		t.Error("expected error for invalid email, got nil")
	}
}

func TestValidateStruct_BoundaryValues(t *testing.T) {
	b := BasicStruct{Name: "Al", Email: "alice@example.com", Age: 17}
	err := ValidateStruct(b)
	if err == nil {
		t.Error("expected error for short name and age < 18, got nil")
	}
	b = BasicStruct{Name: "Aliceeeeeee", Email: "alice@example.com", Age: 100}
	err = ValidateStruct(b)
	if err == nil {
		t.Error("expected error for long name and age > 99, got nil")
	}
}

func TestValidateStruct_NestedStruct(t *testing.T) {
	n := NestedStruct{Inner: BasicStruct{Name: "Bobby", Email: "bob@example.com", Age: 25}}
	err := ValidateStruct(n)
	if err != nil {
		t.Errorf("expected no error for valid nested struct, got: %v", err)
	}
}

func TestValidateStruct_NestedStruct_Invalid(t *testing.T) {
	n := NestedStruct{Inner: BasicStruct{Name: "", Email: "bad", Age: 10}}
	err := ValidateStruct(n)
	if err == nil {
		t.Error("expected error for invalid nested struct, got nil")
	}
}
