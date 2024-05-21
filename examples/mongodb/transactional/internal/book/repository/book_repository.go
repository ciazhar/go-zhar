package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository struct {
	conn *mongo.Collection
}

func (b *BookRepository) Insert(context context.Context, book *model.Book) (err error) {
	one, err := b.conn.InsertOne(context, book)
	if err != nil {
		return
	}
	book.ID = one.InsertedID.(primitive.ObjectID)
	return
}

func (b *BookRepository) UpdateQuantity(context context.Context, id string, amount int) (err error) {

	bookID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	_, err = b.conn.UpdateOne(context,
		bson.M{"_id": bookID, "quantity": bson.M{"$gte": amount}},
		bson.M{"$inc": bson.M{"quantity": -amount}},
	)
	return
}

func NewBookRepository(conn *mongo.Database) *BookRepository {

	collection := conn.Collection("book")

	return &BookRepository{
		conn: collection,
	}

}
