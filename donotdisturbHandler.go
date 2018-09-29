package main

import (
	"fmt"
	"os/exec"
	"time"
)

func restartNotificationCenter() error {
	cmd := exec.Command("bash", "-c", "killall NotificationCenter")

	_, err := cmd.Output()

	if err != nil {
		return err
	}

	return nil
}

// TurnDoNotDisturbOn turns on Do Not Disturb on Mac OS X
func TurnDoNotDisturbOn() error {

	// Enable Do Not Disturb Mode
	cmd := exec.Command("bash", "-c", "defaults -currentHost write ~/Library/Preferences/ByHost/com.apple.notificationcenterui doNotDisturb -boolean true")

	_, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	// Set Time forDo Not Disturb Mode
	t := time.Now()
	tfmt := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d +0000", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	cmd = exec.Command("bash", "-c", fmt.Sprintf("defaults -currentHost write ~/Library/Preferences/ByHost/com.apple.notificationcenterui doNotDisturbDate -date '%s'", tfmt))

	_, err = cmd.CombinedOutput()

	if err != nil {
		return err
	}

	// Restart Notification Center
	err = restartNotificationCenter()

	if err != nil {
		return err
	}

	return nil
}

// TurnDoNotDisturbOff turns off Do Not Disturb on Mac OS X
func TurnDoNotDisturbOff() error {

	// Disable Do Not Disturb Mode
	cmd := exec.Command("bash", "-c", "defaults -currentHost write ~/Library/Preferences/ByHost/com.apple.notificationcenterui doNotDisturb -boolean false")

	_, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	// Restart Notification Center
	err = restartNotificationCenter()

	if err != nil {
		return err
	}

	return nil
}
