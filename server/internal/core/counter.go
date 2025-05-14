package core

import (
	"sync"
	"sync/atomic"
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

func Game(commands <-chan Command, count *atomic.Uint64, cond *Cond) {
	log.Debug("handle game")
	t := time.NewTicker(10 * time.Millisecond)

	defer t.Stop()

	go func() {
		for range t.C {
			cond.Broadcast()
		}
	}()

	for {
		cmd, ok := <-commands
		if !ok {
			log.Error("error receiving from channel event")
		}
		log.Debugf("received event : %v", cmd)

		switch cmd {
		case CommandReset:
			count.Store(0)
		case CommandIncrement:
			count.Add(1)
		case CommandCurrent:
			log.Debug("would current")
			// count.RLock()
			// c := count.Count
			// count.RUnlock()
			// msg := MessageCurrent + ":" + strconv.Itoa(int(c))
			// dispatchSingle(clients, event.ClientDest, msg)

		// >>>>>>>>>>>>>>>>>>> TODO
		// update best every minute with cron

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
