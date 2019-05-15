package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// OpenApp opens an application using AppleScript
func OpenApp(appName string, reschan chan string, errchan chan error) func() {
	return func() {
		script := fmt.Sprintf(`
tell application "%s"
	activate
end tell
`, appName)

		err := executeAppleScript(script)

		if err != nil {
			errchan <- fmt.Errorf("Could not open app '%s': '%v'", appName, err)
			return
		}

		reschan <- fmt.Sprintf("Successfully opened app: '%s'", appName)
		return
	}
}

// CloseApp closes an application using AppleScript
func CloseApp(appName string, reschan chan string, errchan chan error) func() {
	return func() {
		script := fmt.Sprintf(`
tell application "%s"
	quit
end tell
`, appName)

		err := executeAppleScript(script)

		if err != nil {
			errchan <- fmt.Errorf("Could not close app '%s': '%v'", appName, err)
			return
		}
		reschan <- fmt.Sprintf("Successfully closed app: '%s'", appName)
		return
	}
}

// executeAppleScript takes in a fully parsed Apple-Script and executes the command using osascript
func executeAppleScript(command string) error {
	cmd := exec.Command("osascript", "-e", command)

	output, err := cmd.CombinedOutput()
	prettyOutput := strings.Replace(string(output), "\n", "", -1)

	if err != nil {
		return fmt.Errorf("Could not execute script: %s ; %v", prettyOutput, err)
	}

	return nil
}
