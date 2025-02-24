package service

type FileStatus string

const (
	InQueue  FileStatus = "in_queue"
	Parsing  FileStatus = "parsing"
	Error    FileStatus = "error"
	Success  FileStatus = "success"
	Imported FileStatus = "imported"
)
