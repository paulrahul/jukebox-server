package jukebox

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Platform    string `jsong:"platform"`
}

type Track struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type Playlist struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Track []Track `json:"track_list"`
}
