package main

import (
	"testing"

	"github.com/hashicorp/go-version"
)

func TestCheckUpdate(t *testing.T) {
	update, ver := checkUpdate("1.0.0")

	if !update {
		t.Errorf("checkUpdate failed, actual: %t, expected: true", update)
	}
	if ver == "1.0.0" {
		t.Errorf("checkUpdate failed, actual: %s, expected: not 1.0.0", ver)
	}
}

func TestCheckUpdatePrerelease(t *testing.T) {
	update, ver := checkUpdate("1.0.0-alpha.1")

	if !update {
		t.Errorf("checkUpdate failed, actual: %t, expected: true", update)
	}
	if ver == "1.0.0" {
		t.Errorf("checkUpdate failed, actual: %s, expected: not 1.0.0-alpha.1", ver)
	}
}

func TestCheckUpdatePreleaseStream(t *testing.T) {
	currentVersion, _ := version.NewSemver("1.0.0-alpha.1")
	version, _ := checkUpdatePreleaseStream("1.0.0-alpha.1")
	if currentVersion == version {
		t.Errorf("checkUpdatePreleaseStream failed, actual: %s, expected: not 1.0.0-alpha.1", version)
	}
}

func TestCheckUpdateNormalReleaseStream(t *testing.T) {
	currentVersion, _ := version.NewSemver("1.0.0")
	version, _ := checkUpdateNormalReleaseStream("1.0.0")
	if currentVersion == version {
		t.Errorf("checkUpdatePreleaseStream failed, actual: %s, expected: not 1.0.0", version)
	}
}
