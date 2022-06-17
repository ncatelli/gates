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

	"github.com/gorilla/mux"
	"github.com/ncatelli/gates/pkg/config"
)

// healthHandler takes a GET request and returns a 200 response to simulate a
// health check.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "Ok"}`)
}

func startHTTPServer(c *config.Config, wg *sync.WaitGroup) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")
	srv := &http.Server{
		Addr:    c.ListenAddr,
		Handler: router,
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

		log.Printf("Starting server on %s\n", c.ListenAddr)

		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := startHTTPServer(&c, httpServerExitDone)

		// blocks for shutdown. If a SIGHUP happens it will gracefully
		// restart the server.
		<-sigs

		log.Println("reloading configuration...")

		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}

		// wait for goroutine started in startHttpServer() to stop
		httpServerExitDone.Wait()
	}
}
