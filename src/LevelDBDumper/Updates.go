package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-version"
)

func checkUpdate(ver string) (bool, string) {
	currentVersion, _ := version.NewSemver(ver)

	if currentVersion.Prerelease() != "" {
		latestVersion, tag := checkUpdatePreleaseStream(ver)
		return currentVersion.LessThan(latestVersion), tag
	} else {
		latestVersion, tag := checkUpdateNormalReleaseStream(ver)
		return currentVersion.LessThan(latestVersion), tag
	}
}

func checkUpdatePreleaseStream(ver string) (*version.Version, string) {
	url := "https://api.github.com/repos/mdawsonuk/LevelDBDumper/releases"

	resp, err := http.Get(url)
	checkError(err)
	if resp == nil {
		retVer, _ := version.NewSemver(ver)
		return retVer, ver
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var results interface{}

	json.Unmarshal(body, &results)

	switch results := results.(type) {
	case []map[string]interface{}:
		// Drop the v from the tag
		tag := fmt.Sprintf("%s", results[0]["tag_name"])[1:]

		latestVersion, _ := version.NewSemver(tag)

		return latestVersion, tag
	}

	retVer, _ := version.NewSemver(ver)
	return retVer, ver
}

func checkUpdateNormalReleaseStream(ver string) (*version.Version, string) {
	url := "https://api.github.com/repos/mdawsonuk/LevelDBDumper/releases/latest"

	resp, err := http.Get(url)
	checkError(err)
	if resp == nil {
		retVer, _ := version.NewSemver(ver)
		return retVer, ver
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var results interface{}

	json.Unmarshal(body, &results)

	switch results := results.(type) {
	case []map[string]interface{}:
		// Drop the v from the tag
		tag := fmt.Sprintf("%s", results[0]["tag_name"])[1:]

		latestVersion, _ := version.NewSemver(tag)

		return latestVersion, tag
	}

	retVer, _ := version.NewSemver(ver)
	return retVer, ver
}
