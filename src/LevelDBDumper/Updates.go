package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-version"
)

func checkUpdate() (bool, string) {
	url := "https://api.github.com/repos/mdawsonuk/LevelDBDumper/releases/latest"

	resp, err := http.Get(url)
	checkError(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var results map[string]interface{}

	json.Unmarshal(body, &results)

	// Drop the v from the tag
	tag := fmt.Sprintf("%s", results["tag_name"])[1:]

	currentVersion, _ := version.NewSemver(VERSION)
	latestVersion, _ := version.NewSemver(tag)

	return currentVersion.LessThan(latestVersion), tag
}
