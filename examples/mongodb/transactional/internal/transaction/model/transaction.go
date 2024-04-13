package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	BookID primitive.ObjectID `bson:"book_id" json:"book_id"`
	Amount int                `bson:"amount" json:"amount"`
}
