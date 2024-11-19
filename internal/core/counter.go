package core

import (
	"strconv"

	"github.com/charmbracelet/log"
)

func Game(events <-chan string, responses <-chan chan string, logger *log.Logger) {
	logger.Debug("handle game")
	count := 0
	res := make([]chan string, 0)

	for {
		logger.Debug("waiting for event")
		select {
		case event, ok := <-events:
			if !ok {
				logger.Error("error receiving from channel")
			}
			logger.Debugf("received event : %v", event)

			var msg string
			switch event {
			case "increment":
				count += 1
				msg = "increment"
			case "reset":
				count = 0
				msg = "reset"
			case "current":
				msg = strconv.Itoa(count)
			default:
				logger.Warnf("invalid event: %v", event)
			}
			dispatch(res, msg)
			logger.Debug(event)
		case r := <-responses:
			res = append(res, r)
		}

	}
}

func dispatch(responses []chan string, msg string) {
	for _, r := range responses {
		r <- msg
	}
}
