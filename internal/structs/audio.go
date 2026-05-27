package structs

type AudioPlayRequest struct {
	AudioFile   string `json:"audiofile"`
	AudioBase64 string `json:"audiobase64"`
}
