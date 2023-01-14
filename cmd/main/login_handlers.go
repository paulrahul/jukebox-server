package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/paulrahul/jukebox-server/auth"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Debug("/login or /auth_callback called.")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		w.WriteHeader(http.StatusOK)
		return
	}

	if strings.HasSuffix(r.URL.Path, "auth_callback") {
		// Since this is Spotify specific.
		auth.GetSpotifyAuth().RedirectHandler(w, r)
	}

	// platform := r.PostFormValue("platform")
	platform := r.URL.Query().Get("platform")

	authInstance := getAuthInstance(platform)
	if authInstance == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Platform %s not supported for /login/\n", platform)
		return
	}

	authInstance.Login(w, r)

	// w.WriteHeader(http.StatusOK)
}

func getAuthInstance(platform string) auth.Auth {
	switch platform {
	case "spotify":
		return auth.GetSpotifyAuth()
	case "radio5":
	}

	return nil
}
