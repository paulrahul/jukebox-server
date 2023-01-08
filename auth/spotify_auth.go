package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	. "github.com/paulrahul/jukebox-server"
)

type SpotifyAuth struct {
	InitTimeMilliSecs int64
	User              User
}

var spotifyAuthInstance *SpotifyAuth

var REDIRECT_URL string
var STATE string

var auth *spotifyauth.Authenticator
var client *spotify.Client

func Init() {
	REDIRECT_URL = "http://localhost:" + Port + "/auth_callback"
	STATE = "jukebox-server"
}

func GetSpotifyAuth() *SpotifyAuth {
	log.Debug("GetSpotifyAuth called.")

	if spotifyAuthInstance == nil {
		Init()
		spotifyAuthInstance = &SpotifyAuth{}
	}

	return spotifyAuthInstance
}

func (s SpotifyAuth) Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("SpotifyAuth.Login called.")

	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(REDIRECT_URL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))

	// get the user to this URL - how you do that is up to you
	// you should specify a unique state string to identify the session
	url := auth.AuthURL(STATE)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// the user will eventually be redirected back to your redirect URL
// typically you'll have a handler set up like the following:
func (s SpotifyAuth) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("SpotifyAuth.RedirectHandler called.")

	// use the same state string here that you used to generate the URL
	token, err := auth.Token(r.Context(), STATE, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	// create a client using the specified token
	client = spotify.New(auth.Client(r.Context(), token))

	// the client can now be used to make authenticated requests
}

func (s SpotifyAuth) GetUser() User {
	return s.User
}
