package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	Insert(context context.Context, book *model.Book) error
	UpdateQuantity(context context.Context, id string, amount int) error
}

type bookRepository struct {
	conn *mongo.Collection
}

func (b bookRepository) Insert(context context.Context, book *model.Book) error {
	one, err := b.conn.InsertOne(context, book)
	if err != nil {
		return err
	}
	book.ID = one.InsertedID.(primitive.ObjectID)
	return nil
}

func (b bookRepository) UpdateQuantity(context context.Context, id string, amount int) error {

	bookID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = b.conn.UpdateOne(context,
		bson.M{"_id": bookID, "quantity": bson.M{"$gte": amount}},
		bson.M{"$inc": bson.M{"quantity": -amount}},
	)
	return err
}

func NewBookRepository(conn *mongo.Database) BookRepository {

	collection := conn.Collection("book")

	return &bookRepository{
		conn: collection,
	}

}
