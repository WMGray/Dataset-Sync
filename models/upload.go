package models

type UploadDetails struct {
	ImageName    string `json:"image_name"`
	DatasetName  string `json:"dataset_name"`
	ImagePath    string `json:"image_path"`
	ImageSize    string `json:"image_size"`
	UploadTime   string `json:"upload_time"`
	UploadStatus string `json:"upload_status"`
}
