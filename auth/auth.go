package auth

import (
	"net/http"

	. "github.com/paulrahul/jukebox-server"
)

type Auth interface {
	// GetInstance() *Auth
	Login(w http.ResponseWriter, r *http.Request)
	GetUser() User
}
