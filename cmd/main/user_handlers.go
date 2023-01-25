package main

import (
	"fmt"
	"net/http"

	"github.com/paulrahul/jukebox-server/songs"
	log "github.com/sirupsen/logrus"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Debug("/user called.")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		w.WriteHeader(http.StatusOK)
		return
	}

	// TODO
	// Here, for now, get the Spotify user playlists only.
	// For that, get the Spotify Playlists helper, inject the spotify client (I guess)
	// and then call the respective method.
	// user, err := user.GetSpotifyUserHelper().GetCurrentUser()
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, err)
	// 	return
	// }

	// playlists, err := playlist.GetSpotifyPlaylistHelper().GetUserPlaylists(user.ID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, err)
	// 	return
	// }

	// fmt.Fprintln(w, playlists)

	// Get Spotify liked songs.
	// tracks, err := songs.GetSpotifySongsHelper().GetUserLikedTracks()
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, err)
	// 	return
	// }
	// fmt.Fprintln(w, tracks)

	// Get Radio5 liked songs.
	tracks, err := songs.GetRadio5Helper().GetUserLikedTracks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, tracks)

	w.WriteHeader(http.StatusOK)
}
