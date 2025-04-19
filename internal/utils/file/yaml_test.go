package file

import (
	"reflect"
	"testing"
)

type sample struct {
	Name string
	Age  int
}

func TestSetYamlAndUnmarshalYamlTo(t *testing.T) {
	in := sample{Name: "Alice", Age: 30}
	f := &File{}
	f.SetYaml(in)

	var out sample
	err := f.UnmarshalYamlTo(&out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("UnmarshalYamlTo() = %+v, want %+v", out, in)
	}
}

func TestUnmarshalYamlToWithError(t *testing.T) {
	f := &File{}
	f.errors = append(f.errors, errFake{})
	var out sample
	err := f.UnmarshalYamlTo(&out)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
