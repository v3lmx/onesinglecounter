package api

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/v3lmx/counter/internal/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// todo:security check origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnect(mux *http.ServeMux, logger *log.Logger) {
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Get /connect")

		// w.Write("ee")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("error upgrading connection: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		events := make(chan string)

		go handleEvents(conn, events)

		core.HandleGame(events)

		logger.Debug("end")
	})
}

func handleEvents(conn *websocket.Conn, events chan string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Error : %v", err)
		}

		events <- string(msg)
	}
}
