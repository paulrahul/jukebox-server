package songs

import (
	. "github.com/paulrahul/jukebox-server"
)

type SongsHelper interface {
	GetTrack(id string) (Track, error)

	GetTracks(ids []string) ([]Track, error)

	GetUserLikedTracks() ([]Track, error)

	// GetPlaylistTracks(playlistID string) ([]Track, error)
}
