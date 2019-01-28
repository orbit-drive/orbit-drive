package sys

import (
	"strings"

	"github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"
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
	log.Info(msg)
}

// Alert sends an alert notification to the system foreground
// and logs the alert message.
func Alert(m ...string) {
	msg := strings.Join(m, "")
	beeep.Alert(notifierTitle, msg, appIcon)
	log.Warn(msg)
}

// Fatal sends an alert notification to the system foreground,
// logs the alert message and system os exits.
func Fatal(m ...string) {
	msg := strings.Join(m, "")
	beeep.Alert(notifierTitle, msg, appIcon)
	log.Fatal(msg)
}
