package service

import (
	library "github.com/aknrdlt/final-golang/"
	"github.com/aknrdlt/final-golang/pkg/repository"
)

type Authorization interface {
	CreateUser(user library.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Book interface {
	Create(userId int, book library.Book) (int, error)
	GetAll(userId int) ([]library.Book, error)
	GetById(userId, bookId int) (library.Book, error)
	Delete(userId, bookId int) error
	Update(userId, bookId int, input library.Book) error
}

type Service struct {
	Authorization
	Book
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.User),
		Book:          NewBookService(repos.Book),
	}
}
