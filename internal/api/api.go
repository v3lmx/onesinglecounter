package api

import (
	"net/http"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/v3lmx/counter/internal/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// todo:security check origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnect(mux *http.ServeMux, events chan<- core.Event, clients chan<- core.Client) {
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Get /connect")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("error upgrading connection: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		client := core.NewClient()
		clients <- client

		var wg sync.WaitGroup

		go handleEvents(conn, events, client.Id, &wg)
		go handleResponses(conn, client.C, &wg)

		// When either finishes, we have an error and we must cleanup
		wg.Add(1)
		wg.Wait()

		client.Done = true
		clients <- client

		log.Debug("end")
	})
}

func handleEvents(conn *websocket.Conn, events chan<- core.Event, dest uuid.UUID, wg *sync.WaitGroup) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Error : %v", err)
			wg.Done()
			return
		}
		log.Debugf("msg: %b", msg)

		var cmd core.Command
		switch string(msg) {
		case core.MessageReset:
			cmd = core.CommandReset
		case core.MessageIncrement:
			cmd = core.CommandIncrement
		case core.MessageCurrent:
			cmd = core.CommandCurrent
		case core.MessageBest:
			cmd = core.CommandBest
		default:
			log.Errorf("Invalid command: %s", msg)
			continue
		}

		events <- core.Event{Cmd: cmd, ClientDest: dest}
		log.Debug("sent msg")
	}
}

func handleResponses(conn *websocket.Conn, response <-chan string, wg *sync.WaitGroup) {
	for {
		r := <-response
		log.Debugf("resp: %v", r)

		err := conn.WriteMessage(websocket.TextMessage, []byte(r))
		if err != nil {
			log.Errorf("Error : %v", err)
			wg.Done()
			return
		}
	}
}
