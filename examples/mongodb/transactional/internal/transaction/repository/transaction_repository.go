package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/transaction/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	conn *mongo.Collection
}

func (t TransactionRepository) Insert(context context.Context, transaction *model.Transaction) (err error) {
	_, err = t.conn.InsertOne(context, transaction)
	return
}

func NewTransactionRepository(conn *mongo.Database) *TransactionRepository {
	return &TransactionRepository{
		conn: conn.Collection("transaction"),
	}
}
