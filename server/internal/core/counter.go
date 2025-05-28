package core

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

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

func Game(commands <-chan Command, count *atomic.Uint64, cond *Cond, counterTickTime time.Duration) {
	slog.Debug("handle game")
	t := time.NewTicker(counterTickTime)

	defer t.Stop()

	go func() {
		for range t.C {
			cond.Broadcast()
		}
	}()

	for {
		cmd, ok := <-commands
		if !ok {
			slog.Error("error receiving from channel event")
		}
		slog.Debug("received event: ", "cmd", cmd)

		switch cmd {
		case CommandReset:
			count.Store(0)
		case CommandIncrement:
			count.Add(1)
		default:
			slog.Error("invalid event: %v", "cmd", cmd)
		}
		slog.Debug("cmd:", "cmd", cmd)
	}
}
