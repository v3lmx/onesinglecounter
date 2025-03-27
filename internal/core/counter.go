package core

import (
	"sync"
	"time"

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

type CountM struct {
	Count uint
	sync.RWMutex
}

func Game(commands <-chan Command, count *CountM, cond *Cond) {
	// func Game(events <-chan Event, clientsChan <-chan Client, best chan<- uint64, requestBest chan<- struct{}, responseBest <-chan CurrentBest, cronBest <-chan CurrentBest) {
	log.Debug("handle game")
	// count := uint64(0)

	// go func() {
	// 	for {
	// 		currentBest, ok := <-cronBest
	// 		if !ok {
	// 			log.Error("error receiving from channel currentBest")
	// 		}
	// 		dispatch(clients, MessageBest+":"+currentBest.Format())
	// 	}
	// }()

	// tickTime := make(chan time.Duration, 1)
	// tickTime <- 5 * time.Millisecond
	// t := time.NewTicker(1000 * time.Millisecond)
	t := time.NewTicker(5 * time.Millisecond)
	// >>> t.Reset(new_duration)

	defer t.Stop()

	go func() {
		for {
			select {
			// case newTickTime := <-tickTime:
			// 	t = time.NewTicker(newTickTime)
			case <-t.C:
				log.Debug("Broadcasting")
				cond.Broadcast()
				// tick := <-tickTime
				// tickTime <- tick
				// count.RLock()
				// c := count.Count
				// count.RUnlock()
				// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(tick))
				// go dispatch(ctx, clients, MessageCurrent+":"+strconv.Itoa(int(c)), tickTime)
				// cancel()
			}
		}
	}()

	for {
		log.Debug("waiting for event")

		cmd, ok := <-commands
		if !ok {
			log.Error("error receiving from channel event")
		}
		log.Debugf("received event : %v", cmd)

		// var msg Message
		switch cmd {
		case CommandReset:
			count.Lock()
			count.Count = 0
			count.Unlock()
		case CommandIncrement:
			count.Lock()
			count.Count++
			count.Unlock()
		case CommandCurrent:
			log.Debug("would current")
			// count.RLock()
			// c := count.Count
			// count.RUnlock()
			// msg := MessageCurrent + ":" + strconv.Itoa(int(c))
			// dispatchSingle(clients, event.ClientDest, msg)
		// case CommandBest:
		// 	requestBest <- struct{}{}
		// 	currentBest := <-responseBest
		// 	dispatchSingle(clients, event.ClientDest, MessageBest+":"+currentBest.Format())
		default:
			log.Errorf("invalid event: %v", cmd)
		}
		log.Debug(cmd)
		// best <- count
	}
}
