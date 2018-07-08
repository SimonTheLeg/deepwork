package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	affectedApps := []string{"Mail", "Rocket.Chat+", "Calendar"}

	var action func(name string) error

	flag.Parse()

	desStage := flag.Arg(0)

	switch desStage {
	case "on":
		action = CloseApp
	case "off":
		action = OpenApp
	default:
		fmt.Println("Usage: deepwork [on,off]")
		os.Exit(1)
	}

	for _, app := range affectedApps {
		err := action(app)
		if err != nil {
			log.Printf("Could not close app %s: %v", app, err)
		}
	}
}
