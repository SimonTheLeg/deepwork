package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var configLocation string

type config struct {
	AffectedApps []string `json:"affectedApps"`
}

func main() {
	// Get Users homedirectory
	configLocation, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Could not determine users home directory: %v", err)
	}
	configLocation += "/.deepwork/config.json"

	// Parse Configuration
	config, err := parseConfig(configLocation)

	if err != nil {
		log.Fatalf("Could not parse config file: %v", err)
	}

	// Parse Command Line Flags
	flag.Parse()
	desStage := flag.Arg(0)

	// Determine desired action
	var action func(name string) error
	switch desStage {
	case "on":
		action = CloseApp
	case "off":
		action = OpenApp
	default:
		fmt.Println("Usage: deepwork [on,off]")
		os.Exit(1)
	}

	// Execute action
	for _, app := range config.AffectedApps {
		err := action(app)
		if err != nil {
			log.Printf("Could not close app %s: %v", app, err)
		}
	}
}

func parseConfig(configLocation string) (config, error) {
	var conf config
	confDir := path.Dir(configLocation)

	// Try to open config
	jsonFile, err := os.Open(configLocation)
	defer jsonFile.Close()

	// Check if there is a config file at the specified location, if not create a default config
	if os.IsNotExist(err) {
		defaultConfig := []byte(`{"affectedApps":["Mail","Calendar"]}`)

		// Create required directories if necessary
		if err = os.MkdirAll(confDir, 0744); err != nil {
			return config{}, fmt.Errorf("Could not create required directories for config: %v", err)
		}
		// Write default config
		if err = ioutil.WriteFile(configLocation, defaultConfig, 0644); err != nil {
			return config{}, fmt.Errorf("Could not write default config: %v", err)
		}
		// Call itself again to parse newly created conf
		return parseConfig(configLocation)
	}

	// Otherwise (e.g. no permissions on conf file or it's dir), return the error
	if err != nil {
		return config{}, fmt.Errorf("Could not access config file: %v", err)
	}

	// Read in the config
	confByte, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return config{}, fmt.Errorf("Could not read config: %v", err)
	}

	err = json.Unmarshal(confByte, &conf)
	if err != nil {
		return config{}, fmt.Errorf("Could not parse config: %v", err)
	}

	return conf, nil
}
