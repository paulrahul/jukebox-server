package user

import (
	. "github.com/paulrahul/jukebox-server"
)

type UserHelper interface {
	GetUserInfo(id string) (User, error)
	GetCurrentUser() (User, error)
}
