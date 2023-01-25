package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	. "github.com/paulrahul/jukebox-server"
)

func serverShutDownHandler(done chan bool) func() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return func() {
		sig := <-sigs
		log.Info()
		log.Info(sig)
		done <- true
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status: OK")
}

func main() {
	fmt.Println("Welcome to my Jukebox!")
	log.SetLevel(log.DebugLevel)

	Host = os.Getenv("JUKEBOX_HOST")
	if Host == "" {
		panic("$JUKEBOX_HOST must be set")
	}

	Port = os.Getenv("JUKEBOX_PORT")
	if Port == "" {
		panic("$JUKEBOX_PORT must be set")
	}

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	mux := http.NewServeMux()

	// Index
	mux.HandleFunc("/", (statusHandler))

	// Login
	mux.HandleFunc("/login", LoginHandler)

	// Radio5 login redirection
	mux.HandleFunc("/radio5_auth_callback", LoginHandler)

	// Spotify login redirection
	mux.HandleFunc("/auth_callback", LoginHandler)

	// User
	mux.HandleFunc("/user", UserHandler)

	go func() {
		err := http.ListenAndServe(":"+Port, corsWrapper.Handler(mux))
		if err != nil {
			panic(err)
		}
	}()

	done := make(chan bool, 1)
	serverShutDownHandler(done)()
	log.Info("Server Awaiting shutdown signal")
	<-done
	log.Info("Server Shutting down..")
}
