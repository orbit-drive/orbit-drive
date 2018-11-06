package sys

import (
	"log"
	"strings"

	"github.com/gen2brain/beeep"
)

const (
	notifierTitle = "Orbit Drive"
	appIcon       = ""
)

// Notify sends a notification message to the system foreground
// and logs the message to console.
func Notify(m ...string) {
	msg := strings.Join(m, "")
	beeep.Notify(notifierTitle, msg, appIcon)
	log.Println(msg)
}

// Alert sends an alert notification to the system foreground
// and logs the alert message.
func Alert(m ...string) {
	msg := strings.Join(m, "")
	beeep.Alert(notifierTitle, msg, appIcon)
	log.Println(msg)
}

// Fatal sends an alert notification to the system foreground,
// logs the alert message and system os exits.
func Fatal(m ...string) {
	msg := strings.Join(m, "")
	beeep.Alert(notifierTitle, msg, appIcon)
	log.Fatalln(msg)
}
