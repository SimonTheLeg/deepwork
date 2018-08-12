package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

var UserHomeDir string

func init() {
	// Determine Users home directory
	var err error
	UserHomeDir, err = homedir.Dir()
	if err != nil {
		log.Fatalf("Could not determine Users home directory: %v", err)
	}
}

func TestParseConfig(t *testing.T) {
	difFolderPerm := []struct {
		name           string
		folderPerm     os.FileMode
		configLocation string
		expResp        string
	}{
		{
			"no perm, directory exists",
			0000,
			"/tmp/noPermissionsHere/config.json",
			"Could not access config file: open /tmp/noPermissionsHere/config.json: permission denied",
		},
		{
			"execute only, directory doesn't exist yet",
			0100,
			"/tmp/executeOnly/config.json",
			"Could not write default config: open /tmp/executeOnly/config.json: permission denied",
		},
		{
			"write only",
			0200,
			"/tmp/writeOnly/config.json",
			"Could not access config file: open /tmp/writeOnly/config.json: permission denied",
		},
		{
			"write & execute",
			0300,
			"/tmp/writeAndExecute/config.json",
			"no error",
		}, {
			"read only",
			0400,
			"/tmp/readOnly/config.json",
			"Could not access config file: open /tmp/readOnly/config.json: permission denied",
		}, {
			"read & execute",
			0500,
			"/tmp/readAndExecute/config.json",
			"Could not write default config: open /tmp/readAndExecute/config.json: permission denied",
		}, {
			"read & write",
			0600,
			"/tmp/readAndWrite/config.json",
			"Could not access config file: open /tmp/readAndWrite/config.json: permission denied",
		}, {
			"all permissions",
			0700,
			"/tmp/allPermissions/config.json",
			"no error",
		},
	}

	for _, tc := range difFolderPerm {
		t.Run(tc.name, func(t *testing.T) {
			confDir := path.Dir(tc.configLocation)

			if err := os.Mkdir(confDir, tc.folderPerm); err != nil {
				t.Fatalf("Could not create dir: %s: %v", confDir, err)
			}

			_, err := parseConfig(tc.configLocation)

			// Check if the response from parseConfig, is the expected one
			if err != nil {
				if err.Error() != tc.expResp {
					t.Errorf("Expected: '%v' got: '%v'", tc.expResp, err.Error())
				}
			} else {
				if tc.expResp != "no error" {
					t.Errorf("Got no error, but expected '%v'", tc.expResp)
				}
			}

			// Make sure in the successful cases the file is actually there
			if tc.expResp == "no error" {
				cont, err := ioutil.ReadFile(tc.configLocation)
				if err != nil {
					t.Errorf("Expected a file at %s, but could not read it due to: %v", tc.configLocation, err)
				}
				defaultConfig := []byte(`{"affectedApps":["Mail","Calendar"]}`)
				if bytes.Compare(cont, defaultConfig) != 0 {
					t.Errorf("%s should have the following content: '%v', but has '%v'", tc.configLocation, defaultConfig, cont)
				}
			}

			// Clean up
			os.Chmod(confDir, 0700)
			if err := os.RemoveAll(confDir); err != nil {
				t.Fatalf("Could not clean up '%s': %v; Please clean up this directory manually, as it can corrupt further testing", confDir, err)
			}
		})
	}
}

func TestDetermineAction(t *testing.T) {
	tt := []struct {
		name      string
		command   string
		expResult func(name string) error
	}{
		{"Test action On", "on", CloseApp},
		{"Test action Off", "off", OpenApp},
		{"Non Valid command", "not-valid-command", nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := determineAction(tc.command)
			if res != tc.expResult {
				t.Errorf("Expected a pointer to function: '%v', got '%v'", tc.expResult, &res)
			}
		})
	}
}
