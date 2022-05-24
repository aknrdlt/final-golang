package good_db_repo

import (
	"time"

	"github.com/yeldisbayev/thrid-group/library/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookSchema struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	ReleaseDate time.Time          `bson:"releaseDate"`
	Author      string             `bson:"author"`
	Genre       string             `bson:"genre"`
}

func mapToBookSchema(book domain.Book) BookSchema {
	return BookSchema{
		ID:          primitive.NewObjectID(),
		Name:        book.Name,
		Description: book.Description,
		ReleaseDate: book.ReleaseDate,
		Author:      book.Author,
		Genre:       book.Genre,
	}

}

func mapFromBookSchema(bookSchema BookSchema) domain.Book {
	return domain.Book{
		ID:          bookSchema.ID.Hex(),
		Name:        bookSchema.Name,
		Description: bookSchema.Description,
		ReleaseDate: bookSchema.ReleaseDate,
		Author:      bookSchema.Author,
		Genre:       bookSchema.Genre,
	}

}
