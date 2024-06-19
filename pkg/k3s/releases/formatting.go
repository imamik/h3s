package releases

import (
	"fmt"
	"strings"
)

// Type returns "Release Candidate" if Prerelease is true, otherwise "Stable"
var rc = "Release Candidate"
var stable = "Stable"

func (r Release) Type() string {
	if r.Prerelease {
		return rc
	}
	return stable
}

func (r Release) FormattedDate() string {
	return r.PublishedAt.Format("2006-01-02")
}

func PrintReleases(releases []Release) {
	// Calculate column widths
	maxNameLength := len("Name")
	maxTypeLength := max(len("Type"), len(rc), len(stable)) + 5
	for _, r := range releases {
		if len(r.Name) > maxNameLength {
			maxNameLength = len(r.Name) + 5
		}
	}

	// Print the header
	fmt.Printf("%-*s %-*s %-10s\n", maxNameLength, "Name", maxTypeLength, "Type", "Date")
	fmt.Println(strings.Repeat("-", maxNameLength+maxTypeLength+11))

	// Print each release in the formatted table
	for _, v := range releases {
		fmt.Printf("%-*s %-*s %-10s\n", maxNameLength, v.Name, maxTypeLength, v.Type(), v.FormattedDate())
	}
}
