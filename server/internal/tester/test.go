package tester

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MessageIncrement = "inc"
	MessageReset     = "res"
)

// Tester mimics a client to send increment messages to the counter server.
func Tester(ctx context.Context, wg *sync.WaitGroup, url string) {
	defer wg.Done()
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic("dial:" + err.Error())
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
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
				fmt.Printf("error write: %v", err)
				return
			}
		case <-ctx.Done():
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
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
