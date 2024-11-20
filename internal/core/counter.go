package core

import (
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type Command int

const (
	CommandIncrement Command = iota
	CommandReset
	CommandCurrent
)

type Event struct {
	Cmd  Command
	Dest uuid.UUID
}

type Response struct {
	Id   uuid.UUID
	C    chan string
	Done bool
}

func NewResponse() Response {
	return Response{
		Id:   uuid.New(),
		C:    make(chan string),
		Done: false,
	}
}

func Game(events <-chan Event, responses <-chan Response, best chan<- uint) {
	log.Debug("handle game")
	count := 0
	res := make(map[uuid.UUID]chan string, 0)

	for {
		log.Debug("waiting for event")
		select {
		case event, ok := <-events:
			if !ok {
				log.Error("error receiving from channel")
			}
			log.Debugf("received event : %v", event)

			var msg string
			switch event.Cmd {
			case CommandIncrement:
				count += 1
				msg = "increment"
				dispatch(res, msg)
			case CommandReset:
				count = 0
				msg = "reset"
				dispatch(res, msg)
			case CommandCurrent:
				msg = "current:" + strconv.Itoa(count)
				dispatchSingle(res, event.Dest, msg)
			default:
				log.Warnf("invalid event: %v", event)
			}
			log.Debug(event)
		case r := <-responses:
			if r.Done {
				delete(res, r.Id)
				continue
			}
			res[r.Id] = r.C
		}

	}
}

func dispatch(responses map[uuid.UUID]chan string, msg string) {
	for _, r := range responses {
		r <- msg
	}
}

func dispatchSingle(responses map[uuid.UUID]chan string, dest uuid.UUID, msg string) {
	c, ok := responses[dest]
	if !ok {
		log.Errorf("Could not send message to %v", dest)
	}

	c <- msg
}
