package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"net/http"
	"net/url"

	"fmt"

	"github.com/ncatelli/gates/pkg/config"
	"github.com/ncatelli/gates/pkg/gate"
	"github.com/ncatelli/gates/pkg/models"
	"github.com/ncatelli/gates/pkg/outputter"
	"github.com/ncatelli/gates/pkg/router"
)

func startStateManagerService(g *gate.GateService) (chan<- models.MessageInput, chan<- bool) {
	quit := make(chan bool)
	messages := make(chan models.MessageInput)

	go func(g *gate.GateService, msgs <-chan models.MessageInput, quit <-chan bool) {
		for {
			select {
			case msg := <-msgs:
				ts, err := g.ReceiveInput(msg.Tick, msg.Path, msg.Input)
				if err != nil {
					msg.Resp <- models.GateResponse{Err: err, OutputReady: false, Output: false}
					continue
				}

				inputs, err := ts.ReturnInputsIfReady()
				if err != nil {
					msg.Resp <- models.GateResponse{Err: nil, OutputReady: false, Output: false}
					continue
				}

				output, err := g.Compute(msg.Tick, inputs)
				if err != nil {
					msg.Resp <- models.GateResponse{Err: err, OutputReady: false, Output: false}
					continue
				}

				msg.Resp <- models.GateResponse{Err: nil, OutputReady: true, Output: output}
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

func instantiateOutputterFromConfig(c *config.Config) outputter.Outputter {
	var o outputter.Outputter = nil

	switch c.OutputTy {
	case "stdout":
		o = &outputter.StdOutOutputter{}
	case "http":
		endpoints := make([]*url.URL, 0, len(c.OutputAddrs))
		for _, addr := range c.OutputAddrs {
			endpoints = append(endpoints, &addr)
		}

		o = &outputter.HttpOutputter{Endpoints: endpoints}
	}

	return o
}

func startHTTPServer(c *config.Config, g *gate.GateService, msgs chan<- models.MessageInput, wg *sync.WaitGroup) *http.Server {
	r, _ := router.New(g, msgs)
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
			log.Fatalf("unable to parse config params: %v", e)
		}

		g := instantiateGateFromConfig(&c)
		if g == nil {
			panic("invalid gate config")
		}

		o := instantiateOutputterFromConfig(&c)
		if o == nil {
			panic("invalid gate config")
		}

		gg := gate.NewGenericGate(g, o)
		inboundMsgs, stateQuitChan := startStateManagerService(gg)

		log.Printf("Starting server on %s\n", c.ListenAddr)
		log.Printf("Configured as %s gate\n", c.GateTy)

		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := startHTTPServer(&c, gg, inboundMsgs, httpServerExitDone)

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
