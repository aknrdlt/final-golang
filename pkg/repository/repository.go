package repository

import (
	library "github.com/aknrdlt/final-golang"
	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user library.User) (int, error)
	Get(username, password string) (library.User, error)
}

type Book interface {
	Create(userId int, book library.Book) (int, error)
	GetAll(userId int) ([]library.Book, error)
	GetById(userId, bookId int) (library.Book, error)
	Delete(userId, bookId int) error
	Update(userId, bookId int, input library.UpdateBook) error
}

type Repository struct {
	User
	Book
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewAuthPostgres(db),
		Book: NewBookPostgres(db),
	}
}
