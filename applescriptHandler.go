package main

import (
	"fmt"
	"os/exec"
	"strings"
)



// OpenApp opens an application using AppleScript
func OpenApp(appName string) error {
	script := fmt.Sprintf(`
tell application "%s"
	activate
end tell
`, appName)

	err := ExecuteAppleScript(script)

	return fmt.Errorf("Could not open app '%s': '%v'", appName, err)
}

// CloseApp closes an application using AppleScript
func CloseApp(appName string) error {
	script := fmt.Sprintf(`
tell application "%s"
	quit
end tell
`, appName)

	err := ExecuteAppleScript(script)

	return fmt.Errorf("Could not close app '%s': '%v'", appName, err)
}

// ExecuteAppleScript takes in a fully parsed Apple-Script executes the command using osascript
func ExecuteAppleScript(command string) error {
	cmd := exec.Command("osascript", "-e", command)

	output, err := cmd.CombinedOutput()
	prettyOutput := strings.Replace(string(output), "\n", "", -1)

	if err != nil {
		return fmt.Errorf("Could not execute script: %s ; %v", prettyOutput, err)
	}

	return nil
}
