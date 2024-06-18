package k3s

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const k3sReleasesURL = "https://api.github.com/repos/k3s-io/k3s/releases"

type Release struct {
	Name        string    `json:"name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
}

// Type returns "Release Candidate" if Prerelease is true, otherwise "Stable"
var rc = "Release Candidate"
var stable = "Stable"

func (r Release) Type() string {
	if r.Prerelease {
		return rc
	}
	return stable
}

// FormattedDate returns the PublishedAt date formatted as YYYY-MM-DD
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

func GetK3sReleases() ([]Release, error) {
	resp, err := http.Get(k3sReleasesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch releases")
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

func GetK3sLatestRelease() (Release, error) {
	var releases, err = GetK3sReleases()
	return releases[len(releases)-1], err
}
