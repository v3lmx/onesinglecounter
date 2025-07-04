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
)

const (
	MessageIncrement = "inc"
	MessageReset     = "res"
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

		switch cmd {
		case CommandReset:
			count.Store(0)
		case CommandIncrement:
			count.Add(1)
		default:
			slog.Error("invalid event: %v", "cmd", cmd)
		}
	}
}
