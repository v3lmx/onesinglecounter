package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/v3lmx/counter/internal/tester"
)

var addr = flag.String("addr", "localhost:10001", "http service address")
var numClients = flag.Int("numClients", 1, "Number of clients to spawn")

func main() {
	flag.Parse()
	log.SetFlags(0)

	ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		cancel()
	}()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	var wg sync.WaitGroup
	wg.Add(*numClients)
	for range *numClients {
		go tester.Tester(ctx, &wg, u.String())
	}

	wg.Wait()
}
