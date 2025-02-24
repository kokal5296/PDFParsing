package models

type Parser struct {
	ParsedFile   string `json:"parsed_file"`
	ParsedStatus string `json:"parsed_status"`
	ParsedError  string `json:"parsed_errors"`
}
