package k3s

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	// Name is the name of the release
	Name string `json:"name"`
	// Prerelease is true if the release is a release candidate
	Prerelease bool `json:"prerelease"`
	// Draft is true if the release is a draft
	Draft bool `json:"draft"`
	// PublishedAt is the time the release was published
	PublishedAt time.Time `json:"published_at"`
}

// getAllReleases fetches all k3s releases from the GitHub API
func getAllReleases() ([]Release, error) {
	// Fetch releases from the GitHub API
	res, err := http.Get(k3sReleasesURL)
	if err != nil {
		return nil, err
	}

	// Close the response body when the function returns
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body")
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch releases")
	}

	// Decode the response body into a slice of Release structs
	var releases []Release
	if err := json.NewDecoder(res.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

// GetFilteredReleases fetches all k3s releases from the GitHub API and then filters them based on the provided flags
func GetFilteredReleases(prerelease bool, stable bool, limit int) ([]Release, error) {
	releases, err := getAllReleases()
	if err != nil {
		return nil, err
	}

	// Filter releases based on flags
	var filteredReleases []Release
	for _, version := range releases {
		if stable && !version.Prerelease && !version.Draft {
			filteredReleases = append(filteredReleases, version)
		} else if prerelease && version.Prerelease && !version.Draft {
			filteredReleases = append(filteredReleases, version)
		} else if !stable && !prerelease {
			filteredReleases = append(filteredReleases, version)
		}
	}

	// Apply limit
	if limit > 0 && limit < len(filteredReleases) {
		filteredReleases = filteredReleases[:limit]
	}

	return filteredReleases, nil
}

// Type returns "Release Candidate" if Prerelease is true, otherwise "Stable"
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
