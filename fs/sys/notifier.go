package sys

import (
	"log"
	"strings"

	"github.com/gen2brain/beeep"
)

const (
	NOTIFIER_TITLE = "Orbit Drive"
	APP_ICON       = ""
)

func Notify(m ...string) {
	msg := strings.Join(m, "")
	log.Println(msg)
	beeep.Notify(NOTIFIER_TITLE, msg, APP_ICON)
}

func Alert(m ...string) {
	msg := strings.Join(m, "")
	log.Println(msg)
	beeep.Alert(NOTIFIER_TITLE, msg, APP_ICON)
}
