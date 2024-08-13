package k3s

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	k3sReleasesURL = "https://api.github.com/repos/k3s-io/k3s/releases"
	rc             = "Release Candidate"
	stable         = "Stable"
)

type Release struct {
	Name        string    `json:"name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
}

func getAllReleases() ([]Release, error) {
	res, err := http.Get(k3sReleasesURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch releases")
	}

	var releases []Release
	if err := json.NewDecoder(res.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

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
