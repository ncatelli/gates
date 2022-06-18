package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"net/http"

	"fmt"

	"github.com/ncatelli/gates/pkg/config"
	"github.com/ncatelli/gates/pkg/gate"
	"github.com/ncatelli/gates/pkg/router"
)

type MessageInput struct {
	resp  chan error
	tick  uint
	path  rune
	input gate.IO
}

func startStateManagerService(g *gate.GenericGate) (chan<- MessageInput, chan<- bool) {
	quit := make(chan bool)
	messages := make(chan MessageInput)

	go func(g *gate.GenericGate, msgs <-chan MessageInput, quit <-chan bool) {
		for {
			select {
			case msg := <-msgs:
				ts, err := g.ReceiveInput(msg.tick, msg.path, msg.input)
				if err != nil {
					msg.resp <- err
				}

				inputs, err := ts.ReturnInputsIfReady()
				if err != nil {
					msg.resp <- err
				}

				_, err = g.Compute(msg.tick, inputs)
				if err != nil {
					msg.resp <- err
				}
			case <-quit:
				return
			}
		}
	}(g, messages, quit)

	return messages, quit

}

// healthHandler takes a GET request and returns a 200 response to simulate a
// health check.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "Ok"}`)
}

func instantiateGateFromConfig(c *config.Config) gate.Gate {
	var g gate.Gate = nil

	switch c.GateTy {
	case "not":
		g = &gate.Not{}
	case "and":
		g = &gate.And{}
	}

	return g
}

func startHTTPServer(c *config.Config, g gate.Gate, wg *sync.WaitGroup) *http.Server {
	r, _ := router.New(g)
	r.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")
	srv := &http.Server{
		Addr:    c.ListenAddr,
		Handler: r,
	}

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	for {
		c, e := config.New()
		if e != nil {
			log.Fatal("unable to parse config params")
		}

		g := instantiateGateFromConfig(&c)
		if g == nil {
			panic("invalid gate config")
		}
		gg := gate.NewGenericGate(g)
		_, stateQuitChan := startStateManagerService(gg)

		log.Printf("Starting server on %s\n", c.ListenAddr)
		log.Printf("Configured as %s gate\n", c.GateTy)

		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := startHTTPServer(&c, g, httpServerExitDone)

		// blocks for shutdown. If a SIGHUP happens it will gracefully
		// restart the server.
		<-sigs

		log.Println("reloading configuration...")

		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}

		stateQuitChan <- true

		// wait for goroutine started in startHttpServer() to stop
		httpServerExitDone.Wait()
	}
}
