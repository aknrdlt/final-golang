package models

// User schema of the user table
type Book struct {
	ID     int64  `json:"id"`
	Title  string `json:"name"`
	Author string `json:"location"`
	Year   int64  `json:"age"`
}
