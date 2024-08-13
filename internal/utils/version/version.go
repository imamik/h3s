package version

import "fmt"

type Version struct {
	Major int // Major version number - breaking changes
	Minor int // Minor version number - new features
	Patch int // Patch version number - bug fixes
}

func Create(major, minor, patch int) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
