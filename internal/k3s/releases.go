package k3s

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// k3sReleasesURL is the URL to fetch k3s releases from
	k3sReleasesURL = "https://api.github.com/repos/k3s-io/k3s/releases"
	// rc is the string representation of a release candidate
	rc = "Release Candidate"
	// stable is the string representation of a stable release
	stable = "Stable"
)

// Release represents a k3s release as returned by the GitHub API
type Release struct {
	PublishedAt time.Time `json:"published_at"`
	Name        string    `json:"name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

// getAllReleases fetches all k3s releases from the GitHub API
func getAllReleases() ([]Release, error) {
	// Fetch releases from the GitHub API
	res, err := http.Get(k3sReleasesURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Correct: defer closing until after decoding

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases: status %d", res.StatusCode)
	}

	// Decode the response body into a slice of Release structs
	var releases []Release
	if err := json.NewDecoder(res.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

// GetFilteredReleases fetches all k3s releases from the GitHub API and then filters them based on the provided flags
func GetFilteredReleases(prerelease, stable bool, limit int) ([]Release, error) {
	releases, err := getAllReleases()
	if err != nil {
		return nil, err
	}

	filteredReleases := filterReleases(releases, prerelease, stable)

	// Apply limit
	if limit > 0 && limit < len(filteredReleases) {
		filteredReleases = filteredReleases[:limit]
	}

	return filteredReleases, nil
}

func filterReleases(releases []Release, prerelease, stable bool) []Release {
	var filteredReleases []Release
	for _, version := range releases {
		switch {
		case stable && !version.Prerelease && !version.Draft:
			filteredReleases = append(filteredReleases, version)
		case prerelease && version.Prerelease && !version.Draft:
			filteredReleases = append(filteredReleases, version)
		case !stable && !prerelease:
			filteredReleases = append(filteredReleases, version)
		}
	}
	return filteredReleases
}

// Type returns "Release Candidate" if Prerelease is true, otherwise "Stable"
func (r Release) Type() string {
	if r.Prerelease {
		return rc
	}
	return stable
}

// FormattedDate returns the date of the release in the format YYYY-MM-DD
func (r Release) FormattedDate() string {
	return r.PublishedAt.Format("2006-01-02")
}

// PrintReleases prints the releases in a formatted table
func PrintReleases(releases []Release) {
	// Calculate column widths
	maxNameLength := len("Name")
	maxTypeLength := maxInt(len("Type"), maxInt(len(rc), len(stable))) + 5
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

// maxInt returns the larger of two integers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
