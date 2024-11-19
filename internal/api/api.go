package api

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	// "github.com/v3lmx/counter/internal/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// todo:security check origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnect(mux *http.ServeMux, events chan<- string, responses chan<- chan string, logger *log.Logger) {
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Get /connect")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("error upgrading connection: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make(chan string)
		responses <- response

		go handleEvents(conn, events, logger)
		handleResponses(conn, response, logger)

		logger.Debug("end")
	})
}

func handleEvents(conn *websocket.Conn, events chan<- string, logger *log.Logger) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error : %v", err)
			return
		}
		logger.Debugf("msg: %s", msg)

		events <- string(msg)
		logger.Debug("sent msg")
	}
}

func handleResponses(conn *websocket.Conn, response <-chan string, logger *log.Logger) {
	for {
		r := <-response
		logger.Debugf("resp: %v", r)

		err := conn.WriteMessage(websocket.TextMessage, []byte(r))
		if err != nil {
			logger.Errorf("Error : %v", err)
			return
		}
	}
}
