package version

import "fmt"

// Version represents a semantic version number
type Version struct {
	Major int // Major version number - changes if incompatible API changes are made
	Minor int // Minor version number - changes if new functionality is added in a backwards-compatible manner
	Patch int // Patch version number - changes if backwards-compatible bug fixes are made
}

// New creates a new Version
func New(major, minor, patch int) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// String returns the string representation of the Version
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
