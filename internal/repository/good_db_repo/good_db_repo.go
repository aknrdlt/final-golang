package good_db_repo

import (
	"context"
	"errors"

	"github.com/yeldisbayev/thrid-group/library/internal/domain"
	"github.com/yeldisbayev/thrid-group/library/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookMongoRepo struct {
	collection *mongo.Collection
}

func NewBookMongoRepo(collection *mongo.Collection) repository.BookRepo {
	if collection == nil {
		panic("empty collection argument")
	}

	return bookMongoRepo{
		collection: collection,
	}

}

func (repo bookMongoRepo) Store(
	ctx context.Context,
	book domain.Book,
) (*domain.Book, error) {
	result, err := repo.collection.InsertOne(ctx, mapToBookSchema(book))

	if err != nil {
		return nil, err
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		book.ID = objectID.Hex()
	} else {
		return nil, errors.New("can't get inserted id")
	}

	return &book, nil

}

func (repo bookMongoRepo) Get(
	ctx context.Context,
	id string,
) (*domain.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := primitive.M{"_id": objectID}

	result := repo.collection.FindOne(ctx, filter)

	if err := result.Err(); err != nil {
		return nil, err
	}

	bookSchema := BookSchema{}

	if err := result.Decode(&bookSchema); err != nil {
		return nil, err
	}

	book := mapFromBookSchema(bookSchema)

	return &book, nil

}

func (repo bookMongoRepo) Remove(
	ctx context.Context,
	id string,
) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := primitive.M{"_id": objectID}

	result, err := repo.collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("can't delete book, maybe wrong id")
	}

	return nil

}
