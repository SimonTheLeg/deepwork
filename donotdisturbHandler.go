package main

import (
	"fmt"
	"os/exec"
)

func isNotificationCenterOn() (bool, error) {
	cmd := exec.Command("defaults", "-currentHost", "read", "~/Library/Preferences/ByHost/com.apple.notificationcenterui", "doNotDisturb")

	output, err := cmd.Output()

	if err != nil {
		return false, fmt.Errorf("Could not determine Notification Center State; %v", err)
	}

	outputstr := string(output)

	switch outputstr {
	case "0": // Aka Notification Center is off
		return false, nil
	case "1": // Aka Notification Center is on
		return true, nil
	}

	return false, fmt.Errorf("Could not determine Notification Center State: An unexpected error occurred")
}
