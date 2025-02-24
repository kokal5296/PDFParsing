package models

// Queue represents a file waiting for processing
type Queue struct {
	ID      int    `json:"id"`
	FileID  int    `json:"file_id"`
	PDFFile []byte `json:"pdf_file"`
}
