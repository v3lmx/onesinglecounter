package tester

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

const (
	MessageIncrement = "inc"
	MessageReset     = "res"
	MessageCurrent   = "current"
	MessageBest      = "best"
)

func Tester(ctx context.Context, wg *sync.WaitGroup, url string) {
	defer wg.Done()
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				// log.Errorf("read: %v", err)
				return
			}
			// log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(MessageIncrement))
			if err != nil {
				log.Errorf("write: %v", err)
				return
			}
		case <-ctx.Done():
			// log.Error("done ctx")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				// log.Errorf("write close: %v", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
