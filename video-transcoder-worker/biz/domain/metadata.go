package domain

type VideoMetadataMessage struct {
	VideoURL  string `json:"video_url"`
	Thumbnail string `json:"thumbnail"`
}
