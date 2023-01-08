package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
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

func main() {
	fmt.Println("Welcome to my Jukebox!")

	Port := os.Getenv("PORT")
	if Port == "" {
		panic("$PORT must be set")
	}

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	mux := http.NewServeMux()

	// Login
	mux.HandleFunc("/login/", LoginHandler)

	// Spotify login redirection
	mux.HandleFunc("/auth_callback/", LoginHandler)

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
