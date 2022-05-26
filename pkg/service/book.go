package service

import (
	library "github.com/aknrdlt/final-golang"
	"github.com/aknrdlt/final-golang/pkg/repository"
)

type BookService struct {
	repo repository.Book
}

func NewBookService(repo repository.Book) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(userId int, book library.Book) (int, error) {
	return s.repo.Create(userId, book)
}

func (s *BookService) GetAll(userId int) ([]library.Book, error) {
	return s.repo.GetAll(userId)
}

func (s *BookService) GetById(userId, bookId int) (library.Book, error) {
	return s.repo.GetById(userId, bookId)
}

func (s *BookService) Delete(userId, bookId int) error {
	return s.repo.Delete(userId, bookId)
}

func (s *BookService) Update(userId, bookId int, input library.UpdateBook) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, bookId, input)
}
