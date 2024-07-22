package models

type Image struct {
	ID          int    `db:"id" json:"id"`
	Path        string `db:"path" json:"path"`
	Dimensions  string `db:"dimensions" json:"dimensions"`
	CameraModel string `db:"camera_model" json:"camera_model"`
	Location    string `db:"location" json:"location"`
	UploadedAt  string `db:"uploaded_at" json:"uploaded_at"`
}
