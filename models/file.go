package models

import "time"

// File represents a PDF file in the system
type File struct {
	ID         int       `json:"id"`
	FileName   string    `json:"file_name"`
	FileHash   string    `json:"file_hash"`
	ParsedFile string    `json:"-"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
