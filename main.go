package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/raymondsquared/r-go-pact/src/server/handler"
)

var env string
var port int
var wait time.Duration

func init() {
	env = os.Getenv("ENV")
	flag.IntVar(&port, "portFlag", 8100, "the port number of the api")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
}

func main() {
	flag.Parse()

	r := mux.NewRouter()

	// Add your routes as needed

	// Routes
	r.Handle("/health-check", handler.HealthCheckHandler()).Methods("GET")

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),

		// Pass our instance of gorilla/mux in.
		Handler: r,

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	fmt.Printf("ENV: %s\n", env)
	fmt.Printf("App is running on: %s\n", srv.Addr)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
