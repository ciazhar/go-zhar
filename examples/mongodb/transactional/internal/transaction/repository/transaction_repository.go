package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/transaction/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository interface {
	Insert(context context.Context, transaction *model.Transaction) error
}

type transactionRepository struct {
	conn *mongo.Collection
}

func (t transactionRepository) Insert(context context.Context, transaction *model.Transaction) error {
	_, err := t.conn.InsertOne(context, transaction)
	return err
}

func NewTransactionRepository(conn *mongo.Database) TransactionRepository {
	return transactionRepository{
		conn: conn.Collection("transaction"),
	}
}
