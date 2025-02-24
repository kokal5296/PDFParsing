package models

import "time"

// UserFile represents the relationship between users and files
type UserFile struct {
	UserID     int       `json:"user_id"`
	FileID     int       `json:"file_id"`
	Filename   string    `json:"filename"`
	UploadDate time.Time `json:"upload_date"`
	Status     string    `json:"status"`
}
