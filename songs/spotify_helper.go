package songs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"

	. "github.com/paulrahul/jukebox-server"
	"github.com/paulrahul/jukebox-server/auth"
)

var spotifySongsHelper *SpotifyHelper

type SpotifyHelper struct {
	Client *spotify.Client
}

func GetSpotifySongsHelper() *SpotifyHelper {
	log.Debug("GetSpotifySongsHelper called.")

	if spotifySongsHelper == nil {
		spotifySongsHelper = &SpotifyHelper{auth.GetSpotifyAuth().Client}
	}

	return spotifySongsHelper
}

func convertSpotifyTrackToJBTrack(spotifyTrack *spotify.FullTrack) Track {
	ret := Track{string(spotifyTrack.ID), spotifyTrack.Name, make([]string, len(spotifyTrack.Artists))}
	for i, v := range spotifyTrack.Artists {
		ret.Artists[i] = v.Name
	}

	return ret
}

func (s SpotifyHelper) GetTrack(id string) (Track, error) {
	if s.Client == nil {
		panic("Client object not injected into SpotifyHelper")
	}

	fullTrack, err := s.Client.GetTrack(context.Background(), spotify.ID(id))

	if err != nil {
		return Track{}, err
	}

	return convertSpotifyTrackToJBTrack(fullTrack), nil
}

func (s SpotifyHelper) GetTracks(ids []string) ([]Track, error) {
	num := len(ids)
	ret := make([]Track, num)

	// Can send requests only for 50 tracks at a time.
	nIters := num / 50
	for i := 0; i < nIters; i++ {
		start := i * 50
		end := start + 50
		if end > num {
			end = num
		}

		spotifyIDs := make([]spotify.ID, end-start)
		for j, k := start, 0; j < end; j, k = j+1, k+1 {
			spotifyIDs[k] = spotify.ID(ids[j])
		}
		tracks, err := s.Client.GetTracks(context.Background(), spotifyIDs)
		if err != nil {
			return nil, err
		}

		for j, k := start, 0; j < end; j, k = j+1, k+1 {
			ret[j] = convertSpotifyTrackToJBTrack(tracks[k])
		}
	}

	return ret, nil
}

func (s SpotifyHelper) GetUserLikedTracks() ([]Track, error) {
	if s.Client == nil {
		panic("Client object not injected into SpotifyHelper")
	}

	tracksPage, err := s.Client.CurrentUsersTracks(context.Background())

	if err != nil {
		return nil, err
	}

	num := len(tracksPage.Tracks)
	ret := make([]Track, num)

	for i, v := range tracksPage.Tracks {
		ret[i] = convertSpotifyTrackToJBTrack(&v.FullTrack)
	}

	return ret, nil
}

// func (s SpotifyHelper) GetPlaylistTracks(playlistID string) ([]Track, error) {
// 	// Get the pla
// }
