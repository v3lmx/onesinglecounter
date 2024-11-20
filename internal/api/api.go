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

func HandleConnect(mux *http.ServeMux, events chan<- core.Event, responses chan<- core.Response) {
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Get /connect")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("error upgrading connection: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := core.NewResponse()
		responses <- response

		var wg sync.WaitGroup

		go handleEvents(conn, events, response.Id, &wg)
		go handleResponses(conn, response.C, &wg)

		// When either finishes, we have an error and we must cleanup
		wg.Add(1)
		wg.Wait()

		response.Done = true
		responses <- response

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
		log.Debugf("msg: %s", msg)

		var cmd core.Command
		switch string(msg) {
		case "increment":
			cmd = core.CommandIncrement
		case "reset":
			cmd = core.CommandReset
		case "current":
			cmd = core.CommandCurrent
		default:
			log.Errorf("Invalid command: %s", msg)
			continue
		}

		events <- core.Event{Cmd: cmd, Dest: dest}
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
