package api

import (
	"net/http"
	"strconv"
	"sync"

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

func HandleConnect(mux *http.ServeMux, commands chan<- core.Command, count *core.CountM, best *core.CurrentBest, tickClock *core.Cond, bestClock *core.Cond) {
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Get /connect")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("error upgrading connection: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var wg sync.WaitGroup

		go handleEvents(conn, commands, &wg)
		go handleCount(conn, count, tickClock, &wg)
		go handleBest(conn, best, bestClock, &wg)

		// When either finishes, we have an error and we must cleanup
		wg.Add(1)
		wg.Wait()

		log.Debug("end")
	})
}

func handleEvents(conn *websocket.Conn, commands chan<- core.Command, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Error : %v", err)
			return
		}
		log.Debugf("msg: %b", msg)

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
			log.Errorf("Invalid command: %s", msg)
			continue
		}
	}
}

func handleCount(conn *websocket.Conn, count *core.CountM, cond *core.Cond, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		cond.L.Lock()
		// wait for broadcast (tick)
		// todo: maybe local tick ?? -> how to update it dynamically
		//		\-> update tick every tick from RWlock value like count?
		cond.Wait()

		count.RLock()
		c := count.Count
		count.RUnlock()

		cond.L.Unlock()

		err := conn.WriteMessage(websocket.TextMessage, []byte("current:"+strconv.Itoa(int(c))))
		if err != nil {
			log.Errorf("Error : %v", err)
			return
		}
	}
}

func handleBest(conn *websocket.Conn, best *core.CurrentBest, cond *core.Cond, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		cond.L.Lock()
		cond.Wait()

		var b core.CurrentBest
		best.RLock()
		b.AllTime = best.AllTime
		b.Year = best.Year
		b.Month = best.Month
		b.Week = best.Week
		b.Day = best.Day
		b.Hour = best.Hour
		b.Minute = best.Minute
		best.RUnlock()

		cond.L.Unlock()

		err := conn.WriteMessage(websocket.TextMessage, []byte(b.Format()))
		if err != nil {
			log.Errorf("Error : %v", err)
			return
		}
	}
}
