package auth

import (
	"net/http"

	. "github.com/paulrahul/jukebox-server"
)

type Auth interface {
	Login(w http.ResponseWriter, r *http.Request)
	GetUser() User
}
