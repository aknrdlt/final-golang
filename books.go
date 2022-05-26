package library

import "errors"

type Book struct {
	Id     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title" binding:"required"`
	Author string `json:"author" db:"author"`
}

type Client struct {
	Id       int
	ClientId int
	BookId   int
}

type UpdateBook struct {
	Title  *string `json:"title"`
	Author *string `json:"author"`
}

func (i UpdateBook) Validate() error {
	if i.Title == nil && i.Author == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
