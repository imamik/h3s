package releases

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const k3sReleasesURL = "https://api.github.com/repos/k3s-io/k3s/releases"

func getAllReleases() ([]Release, error) {
	resp, err := http.Get(k3sReleasesURL)
	if err != nil {
		fmt.Println("Error fetching releases:", err)
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

func GetLatestRelease() (Release, error) {
	var releases, err = getAllReleases()
	return releases[len(releases)-1], err
}
