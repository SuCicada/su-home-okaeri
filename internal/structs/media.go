package structs

type MediaStatus struct {
	Artist      *string `json:"artist,omitempty"`
	Title       *string `json:"title,omitempty"`
	AlbumArtist *string `json:"album_artist,omitempty"`
	AlbumTitle  *string `json:"album_title,omitempty"`
	TrackNumber int     `json:"trackNumber"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
	IsPlaying   bool    `json:"is_playing"`
	Volume      int     `json:"volume"`
	IsMute      bool    `json:"is_mute"`
}
