package user

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"

	. "github.com/paulrahul/jukebox-server"
	"github.com/paulrahul/jukebox-server/auth"
)

var spotifyUserHelper *SpotifyHelper

type SpotifyHelper struct {
	Client *spotify.Client
}

func GetSpotifyUserHelper() *SpotifyHelper {
	log.Debug("GetSpotifyUserHelper called.")

	if spotifyUserHelper == nil {
		spotifyUserHelper = &SpotifyHelper{auth.GetSpotifyAuth().Client}
	}

	return spotifyUserHelper
}

func convertSpotifyUserToJBUser(spotifyUser *spotify.User) User {
	return User{spotifyUser.ID, spotifyUser.DisplayName, "spotify"}
}

func (s SpotifyHelper) GetUserInfo(id string) (User, error) {
	spotifyUser, err := s.Client.GetUsersPublicProfile(context.Background(), spotify.ID(id))

	if err != nil {
		return User{}, err
	}

	return convertSpotifyUserToJBUser(spotifyUser), nil
}

func (s SpotifyHelper) GetCurrentUser() (User, error) {
	currUser, err := s.Client.CurrentUser(context.Background())

	if err != nil {
		return User{}, nil
	}

	return s.GetUserInfo(currUser.ID)
}
