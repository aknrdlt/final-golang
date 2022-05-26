package repository

import (
	"fmt"
	library "github.com/aknrdlt/final-golang"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type BookPostgres struct {
	db *sqlx.DB
}

func NewBookPostgres(db *sqlx.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) Create(userId int, book library.Book) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, author) VALUES ($1, $2) RETURNING id", booksTable)
	row := tx.QueryRow(createListQuery, book.Title, book.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, book_id) VALUES ($1, $2)", clientsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *BookPostgres) GetAll(userId int) ([]library.Book, error) {
	var books []library.Book

	query := fmt.Sprintf("SELECT b.id, b.title, b.author FROM %s b INNER JOIN %s ul on b.id = ul.book_id WHERE ul.user_id = $1",
		booksTable, createTable)
	err := r.db.Select(&books, query, userId)

	return books, err
}

func (r *BookPostgres) GetById(userId, listId int) (library.Book, error) {
	var list library.Book

	query := fmt.Sprintf(`SELECT b.id, b.title, b.author FROM %s b
								INNER JOIN %s ul on tl.id = ul.book_id WHERE ul.user_id = $1 AND ul.book_id = $2`,
		booksTable, clientsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *BookPostgres) Delete(userId, bookId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.book_id AND ul.user_id=$1 AND ul.book_id=$2",
		booksTable, clientsTable)
	_, err := r.db.Exec(query, userId, bookId)

	return err
}

func (r *BookPostgres) Update(userId, bookId int, input library.UpdateBook) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, *input.Author)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.book_id AND ul.book_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, bookId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
