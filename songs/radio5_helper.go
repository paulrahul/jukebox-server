package songs

import (
	"encoding/json"

	. "github.com/paulrahul/jukebox-server"
	"github.com/paulrahul/jukebox-server/auth"
)

var radio5SongsHelper *Radio5Helper

type Radio5Helper struct {
	Client *auth.Radio5Client
}

func GetRadio5Helper() *Radio5Helper {
	if radio5SongsHelper == nil {
		radio5SongsHelper = &Radio5Helper{&auth.GetRadio5Auth().Client}
	}

	return radio5SongsHelper
}

func convertRawRadio5TrackToJBTrack(r5Track string) (Track, error) {
	var rawData map[string]interface{}

	err := json.Unmarshal([]byte(r5Track), &rawData)
	if err != nil {
		return Track{}, err
	}

	return convertRadio5TrackToJBTrack(rawData)
}

func convertRadio5TrackToJBTrack(r5Track map[string]interface{}) (Track, error) {
	ret := Track{
		r5Track["_id"].(string),
		r5Track["title"].(string),
		[]string{r5Track["artist"].(string)},
	}

	return ret, nil
}

func (r Radio5Helper) GetTrack(id string) (Track, error) {
	if r.Client == nil {
		panic("No client injected into Radio5Helper")
	}

	_, resp, err := r.Client.Get("track/" + id)
	if err != nil {
		return Track{}, err
	}

	return convertRawRadio5TrackToJBTrack(resp)
}

func (r Radio5Helper) GetTracks(ids []string) ([]Track, error) {
	num := len(ids)
	ret := make([]Track, num)

	for i, v := range ids {
		track, err := r.GetTrack(v)
		if err != nil {
			return nil, err
		}

		ret[i] = track
	}

	return ret, nil
}

func (r Radio5Helper) GetUserLikedTracks() ([]Track, error) {
	if r.Client == nil {
		panic("No client injected into Radio5Helper")
	}

	_, resp, err := r.Client.Get(
		"contributor/likes/" + r.Client.Credentials.ContributorID + "?page=1&size=50")
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	err = json.Unmarshal([]byte(resp), &data)
	if err != nil {
		return nil, err
	}

	num := len(data)
	ret := make([]Track, num)
	for i, v := range data {
		ret[i], err = convertRadio5TrackToJBTrack(v)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
