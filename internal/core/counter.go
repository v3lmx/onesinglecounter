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
	CommandBest
)

const (
	MessageIncrement = "inc"
	MessageReset     = "res"
	MessageCurrent   = "current"
	MessageBest      = "best"
)

type Event struct {
	Cmd        Command
	ClientDest uuid.UUID
}

type Client struct {
	Id   uuid.UUID
	C    chan string
	Done bool
}

func NewClient() Client {
	return Client{
		Id:   uuid.New(),
		C:    make(chan string),
		Done: false,
	}
}

func Game(events <-chan Event, clientsChan <-chan Client, best chan<- uint64, requestBest chan<- struct{}, responseBest <-chan CurrentBest, cronBest <-chan CurrentBest) {
	log.Debug("handle game")
	count := uint64(0)
	clients := make(map[uuid.UUID]chan string, 0)

	go func() {
		for {
			currentBest, ok := <-cronBest
			if !ok {
				log.Error("error receiving from channel currentBest")
			}
			dispatch(clients, MessageBest+":"+currentBest.Format())
		}
	}()

	for {
		log.Debug("waiting for event")
		select {
		case event, ok := <-events:
			if !ok {
				log.Error("error receiving from channel event")
			}
			log.Debugf("received event : %v", event)

			// var msg Message
			switch event.Cmd {
			case CommandReset:
				count = 0
				dispatch(clients, MessageReset)
			case CommandIncrement:
				count += 1
				dispatch(clients, MessageIncrement)
			case CommandCurrent:
				msg := MessageCurrent + ":" + strconv.Itoa(int(count))
				dispatchSingle(clients, event.ClientDest, msg)
			case CommandBest:
				requestBest <- struct{}{}
				currentBest := <-responseBest
				dispatchSingle(clients, event.ClientDest, MessageBest+":"+currentBest.Format())
			default:
				log.Errorf("invalid event: %v", event)
			}
			log.Debug(event)
			best <- count
		case r := <-clientsChan:
			if r.Done {
				delete(clients, r.Id)
				continue
			}
			clients[r.Id] = r.C
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
