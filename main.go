package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"net/http"
	"net/url"

	"fmt"

	"github.com/ncatelli/gates/pkg/gate"
	"github.com/ncatelli/gates/pkg/models"
	"github.com/ncatelli/gates/pkg/outputter"
	"github.com/ncatelli/gates/pkg/router"
)

const (
	subCommandUsage string = `Usage of ./gates:
  Available Subcommands:
    not
    and
    or
    xor
    nand
    nor
`
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

func instantiateOutputter(ot outputter.OutputTy, endpoints []url.URL) outputter.Outputter {
	var o outputter.Outputter = nil

	switch ot {
	case outputter.StdOut:
		o = &outputter.StdOutOutputter{}
	case outputter.HTTP:
		o = &outputter.HttpOutputter{Endpoints: endpoints}
	}

	return o
}

func startGateHTTPServer(listenAddr string, g *gate.GateService, msgs chan<- models.MessageInput, wg *sync.WaitGroup) *http.Server {
	r, _ := router.New(g, msgs)
	r.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")
	srv := &http.Server{
		Addr:    listenAddr,
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

func parseCommands() (func(*sync.WaitGroup) (*http.Server, chan<- bool), error) {
	if len(os.Args) < 2 {
		fmt.Print(subCommandUsage)
		return nil, nil
	}

	var g gate.Gate = nil
	cmd := os.Args[1]
	if cmd == "-h" || cmd == "--help" || cmd == "-help" {
		fmt.Print(subCommandUsage)
		return nil, nil
	}

	/// switch case to handle for gate types
	isGate := true
	switch cmd {
	case "not":
		g = &gate.Not{}
	case "and":
		g = &gate.And{}
	case "or":
		g = &gate.Or{}
	case "xor":
		g = &gate.Xor{}
	case "nand":
		g = &gate.Nand{}
	case "nor":
		g = &gate.Nor{}
	default:
		isGate = false
	}

	if isGate {
		var outputAddr string
		var listenAddr string

		// setup the flags
		gateCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		gateCmd.StringVar(&listenAddr, "listen-addr", "127.0.0.1:8080", "The server address gates binds to.")
		gateCmd.StringVar(&outputAddr, "output-addrs", "", `An optional comma-separated list of address for the http outputter to send a computed
output to. If empty the stdout outputter is used.`)
		gateCmd.Parse(os.Args[2:])

		// configure outputter
		preparsedAddrs := strings.Split(outputAddr, ",")
		endpoints := make([]url.URL, 0, len(preparsedAddrs))
		for _, addrStr := range preparsedAddrs {
			if len(addrStr) == 0 {
				break
			}

			addr, err := url.Parse(addrStr)
			if err != nil {
				return nil, err
			} else if addr == nil {
				continue
			}

			endpoints = append(endpoints, *addr)
		}

		ot := outputter.StdOut
		if len(endpoints) > 0 {
			ot = outputter.HTTP
		}
		o := instantiateOutputter(ot, endpoints)

		// instantiate the gate and start a listener
		gg := gate.NewGenericGate(g, o)
		inboundMsgs, stateQuitChan := startStateManagerService(gg)

		log.Printf("Starting server on %s\n", listenAddr)
		log.Printf("Configured as %s gate\n", cmd)

		return func(httpServerExitDone *sync.WaitGroup) (*http.Server, chan<- bool) {
			srv := startGateHTTPServer(listenAddr, gg, inboundMsgs, httpServerExitDone)
			return srv, stateQuitChan
		}, nil
	}

	return nil, fmt.Errorf("unimplemented mode: %s", cmd)
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	for {
		runner, err := parseCommands()
		if err != nil {
			log.Fatal(err)
		} else if runner == nil && err == nil {
			// early termination case
			os.Exit(0)
		}

		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)

		srv, quitChan := runner(httpServerExitDone)

		// blocks for shutdown. If a SIGHUP happens it will gracefully
		// restart the server.
		<-sigs

		log.Println("reloading configuration...")

		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}

		quitChan <- true

		// wait for goroutine started in startHttpServer() to stop
		httpServerExitDone.Wait()
	}
}
