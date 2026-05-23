package structs

type MediaStatus struct {
	Artist      *string `json:"artist,omitempty"`
	Title       *string `json:"title,omitempty"`
	AlbumArtist *string `json:"albumArtist,omitempty"`
	AlbumTitle  *string `json:"albumTitle,omitempty"`
	TrackNumber int     `json:"trackNumber"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
	IsPlaying   bool    `json:"isPlaying"`
	Volume      int     `json:"volume"`
	IsMute      bool    `json:"isMute"`
}

type MediaCommand struct {
	Command string         `json:"command"`
	Args    map[string]any `json:"args"`
}
