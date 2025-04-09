package main

import (
	"flag"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"

	netpprof "net/http/pprof"
	// "syscall"

	"github.com/charmbracelet/log"
	"github.com/v3lmx/counter/internal/api"
	"github.com/v3lmx/counter/internal/core"
)

func checkCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// origin := r.Header.Get("Origin")
		// if slices.Contains(originAllowlist, origin) {
		// 	w.Header().Set("Access-Control-Allow-Origin", origin)
		// 	w.Header().Add("Vary", "Origin")
		// }
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// var count = core.CountM{Count: 0}
var count atomic.Uint64
var best = core.CurrentBest{}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	/////////////////

	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.ErrorLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", netpprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", netpprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", netpprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", netpprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	mux.Handle("/debug/pprof/goroutine", netpprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", netpprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", netpprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", netpprof.Handler("block"))

	commands := make(chan core.Command)

	var m1, m2 sync.Mutex
	tickClock := core.NewCond(&m1)
	bestClock := core.NewCond(&m2)

	go core.Game(commands, &count, &tickClock)
	// go core.Game(events, clients, count, requestBest, responseBest, cronBest)
	go core.Best(&count, &best, &tickClock, &bestClock)

	api.HandleConnect(mux, commands, &count, &best, &tickClock, &bestClock)

	log.Info("starting server on port 10001")
	log.Error(http.ListenAndServe(":10001", checkCORS(mux)))

	/////////////////

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		// Lookup("allocs") creates a profile similar to go test -memprofile.
		// Alternatively, use Lookup("heap") for a profile
		// that has inuse_space as the default index.
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
