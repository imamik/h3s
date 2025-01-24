// Package version contains version information for the application
package version

// Variables containing build information, set at build time
var (
	Version   = "dev"
	Commit    = "unknown"
	GoVersion = "unknown"
)

// BuildInfo contains version information
type BuildInfo struct {
	Version   string
	Commit    string
	GoVersion string
}

// GetBuildInfo returns the current build information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		Commit:    Commit,
		GoVersion: GoVersion,
	}
}
