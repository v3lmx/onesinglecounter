package core

import "github.com/charmbracelet/log"

func HandleGame(events chan string) {
	log.Debug("handle game")

	for {
		event := <-events

		log.Debug(event)
	}
}

