package playlist

import (
	. "github.com/paulrahul/jukebox-server"
)

type PlaylistHelper interface {
	GetPlaylist(id string) (Playlist, error)

	GetUserPlaylists(userID string) ([]Playlist, error)
}
