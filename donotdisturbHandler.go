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
func TurnDoNotDisturbOn(reschan chan string, errchan chan error) func() {
	return func() {
		// Enable Do Not Disturb Mode
		cmd := exec.Command("bash", "-c", "defaults -currentHost write ~/Library/Preferences/ByHost/com.apple.notificationcenterui doNotDisturb -boolean true")

		_, err := cmd.CombinedOutput()

		if err != nil {
			errchan <- fmt.Errorf("Could not turn Do Not Disturb on: %v", err)
			return
		}

		// Set Time forDo Not Disturb Mode
		t := time.Now()
		tfmt := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d +0000", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

		cmd = exec.Command("bash", "-c", fmt.Sprintf("defaults -currentHost write ~/Library/Preferences/ByHost/com.apple.notificationcenterui doNotDisturbDate -date '%s'", tfmt))

		_, err = cmd.CombinedOutput()

		if err != nil {
			errchan <- fmt.Errorf("Could not set time for turning Do Not Disturb on: %v", err)
			return
		}

		// Restart Notification Center
		err = restartNotificationCenter()

		if err != nil {
			errchan <- fmt.Errorf("Could not restart Notification Center: %v", err)
			return
		}

		reschan <- "Successfully turned Do Not Disturb On"
		return
	}
}

// TurnDoNotDisturbOff turns off Do Not Disturb on Mac OS X
func TurnDoNotDisturbOff(reschan chan string, errchan chan error) func() {
	return func() {
		// Disable Do Not Disturb Mode
		cmd := exec.Command("bash", "-c", "defaults -currentHost write ~/Library/Preferences/ByHost/com.apple.notificationcenterui doNotDisturb -boolean false")

		_, err := cmd.CombinedOutput()

		if err != nil {
			errchan <- fmt.Errorf("Could not turn Do Not Disturb off: %v", err)
			return
		}

		// Restart Notification Center
		err = restartNotificationCenter()

		if err != nil {
			errchan <- fmt.Errorf("Could not restart Notification Center: %v", err)
			return
		}

		reschan <- "Successfully turned Do Not Disturb Off"
		return
	}
}
