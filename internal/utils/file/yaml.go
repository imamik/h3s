package file

import (
	"gopkg.in/yaml.v3"
)

// SetYaml sets the file content to the YAML representation of the given interface
func (f *File) SetYaml(yml interface{}) *File {
	data, err := yaml.Marshal(yml)
	f.SetBytes(data)
	if err != nil {
		f.errors = append(f.errors, err)
	}
	return f
}

// UnmarshalYamlTo unmarshals the file content to the given interface
func (f *File) UnmarshalYamlTo(out interface{}) error {
	err := f.Error()
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f.data, out)
	if err != nil {
		f.errors = append(f.errors, err)
	}
	return f.Error()
}
