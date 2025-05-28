package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/websocket"

	"github.com/v3lmx/counter/internal/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// todo:security check origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnect(mux *http.ServeMux, commands chan<- core.Command, count *atomic.Uint64, best *core.CurrentBest, tickClock *core.Cond, bestClock *core.Cond) {
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Get /connect")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("error upgrading connection: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())

		msg := make(chan string)

		go handleEvents(ctx, cancel, conn, commands)
		go handleCount(ctx, cancel, msg, count, tickClock)
		go handleBest(ctx, cancel, msg, best, bestClock)

		// Wait for any func to finish
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-msg:
				err := conn.WriteMessage(websocket.TextMessage, []byte(m))
				if err != nil {
					slog.Warn("Error writing message: ", err.Error())
					return
				}
			}
		}
	})
}

func handleEvents(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, commands chan<- core.Command) {
	defer cancel()
	for {
		if ctx.Err() != nil {
			return
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message: " + err.Error())
			return
		}
		slog.Debug(fmt.Sprintf("msg: %b", msg))

		switch string(msg) {
		case core.MessageReset:
			commands <- core.CommandReset
		case core.MessageIncrement:
			commands <- core.CommandIncrement
		case core.MessageCurrent:
			commands <- core.CommandCurrent
		case core.MessageBest:
			commands <- core.CommandBest
		default:
			slog.Error("Invalid command: " + string(msg))
			continue
		}
	}
}

func handleCount(ctx context.Context, cancel context.CancelFunc, msg chan<- string, count *atomic.Uint64, cond *core.Cond) {
	defer cancel()
	for {
		if ctx.Err() != nil {
			return
		}

		cond.L.Lock()
		// wait for broadcast (tick)
		// todo: maybe local tick ?? -> how to update it dynamically
		//		\-> update tick every tick from RWlock value like count?
		cond.Wait()

		c := count.Load()

		cond.L.Unlock()

		msg <- "current:" + strconv.Itoa(int(c))
	}
}

func handleBest(ctx context.Context, cancel context.CancelFunc, msg chan<- string, best *core.CurrentBest, cond *core.Cond) {
	defer cancel()
	msg <- core.MessageBest
	for {
		if ctx.Err() != nil {
			return
		}

		cond.L.Lock()
		cond.Wait()

		var b core.Best
		best.RLock()
		b = best.Copy()
		best.RUnlock()

		cond.L.Unlock()

		msg <- b.Format()
	}
}
