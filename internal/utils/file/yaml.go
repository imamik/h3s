package file

import (
	"gopkg.in/yaml.v3"
)

func (f *File) SetYaml(yml interface{}) *File {
	data, err := yaml.Marshal(yml)
	f.SetBytes(data)
	if err != nil {
		f.errors = append(f.errors, err)
	}
	return f
}

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
