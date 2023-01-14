package playlist

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"

	. "github.com/paulrahul/jukebox-server"
	"github.com/paulrahul/jukebox-server/auth"
)

var spotifyPlaylistHelper *SpotifyHelper

type SpotifyHelper struct {
	Client *spotify.Client
}

func GetSpotifyPlaylistHelper() *SpotifyHelper {
	log.Debug("GetSpotifyPlaylisHelper called.")

	if spotifyPlaylistHelper == nil {
		spotifyPlaylistHelper = &SpotifyHelper{auth.GetSpotifyAuth().Client}
	}

	return spotifyPlaylistHelper
}

func convertSpotifyPlaylistToJBPlaylist(spotifyPlaylist *spotify.SimplePlaylist) Playlist {
	return Playlist{string(spotifyPlaylist.ID), spotifyPlaylist.Name, nil}
}

func (s SpotifyHelper) GetPlaylist(id string) (Playlist, error) {
	spotifyPlaylist, err := s.Client.GetPlaylist(context.Background(), spotify.ID(id))

	if err != nil {
		return Playlist{}, err
	}

	return convertSpotifyPlaylistToJBPlaylist(&spotifyPlaylist.SimplePlaylist), err
}

func (s SpotifyHelper) GetUserPlaylists(userID string) ([]Playlist, error) {
	playlists, err := s.Client.GetPlaylistsForUser(context.Background(), userID)

	if err != nil {
		return nil, err
	}

	num := len(playlists.Playlists)
	ret := make([]Playlist, num)
	for i := 0; i < num; i++ {
		ret[i] = convertSpotifyPlaylistToJBPlaylist(&playlists.Playlists[i])
	}

	return ret, nil
}
