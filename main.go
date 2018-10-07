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

type config struct {
	AffectedApps []string `json:"affectedApps"`
}

var version = "dev-build"
var defaultConfig = []byte(`{"affectedApps":["Mail","Calendar"]}`)

var curConfig config
var reschan = make(chan string)
var errchan = make(chan error)

func main() {
	// Get Users homedirectory to set the configLocation
	configLocation, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Could not determine users home directory: %v", err)
	}
	configLocation += "/.deepwork/config.json"

	// Parse Configuration
	curConfig, err = parseConfig(configLocation)

	if err != nil {
		log.Fatalf("Could not parse config file: %v", err)
	}

	// Parse Command Line Flags
	flag.Parse()
	desAction := flag.Arg(0)

	// Determine desired actions
	actions := determineActions(desAction)

	if actions == nil {
		os.Exit(0)
	}

	// Execute all actions in parallel
	for _, action := range actions {
		go action()
	}

	for i := 0; i < len(actions); i++ {
		select {
		case res := <-reschan:
			fmt.Println(res)
		case err := <-errchan:
			fmt.Println(err)
		}
	}
}

func determineActions(desAction string) []func() {
	var actions []func()

	switch desAction {
	case "on":
		// Handle Apps
		for _, app := range curConfig.AffectedApps {
			actions = append(actions, CloseApp(app, reschan, errchan))
		}
		// Handle Notification Center
		actions = append(actions, TurnDoNotDisturbOn(reschan, errchan))
	case "off":
		// Handle Apps
		for _, app := range curConfig.AffectedApps {
			actions = append(actions, OpenApp(app, reschan, errchan))
		}
		// Handle Notification Center
		actions = append(actions, TurnDoNotDisturbOff(reschan, errchan))
	case "version":
		fmt.Println(version)
		return nil
	default:
		fmt.Println("Usage: deepwork [on,off,version]")
		return nil
	}
	return actions
}

func parseConfig(configLocation string) (config, error) {
	var conf config
	confDir := path.Dir(configLocation)

	// Try to open config
	jsonFile, err := os.Open(configLocation)
	defer jsonFile.Close()

	// Check if there is a config file at the specified location, if not create a default config
	if os.IsNotExist(err) {
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
