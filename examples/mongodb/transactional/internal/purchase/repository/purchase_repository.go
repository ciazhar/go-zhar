package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/repository"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/transaction/model"
	repository2 "github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/transaction/repository"
	mongo2 "github.com/ciazhar/go-zhar/pkg/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type PurchaseRepository interface {
	PurchaseWithAutomaticTransaction(context context.Context, transaction *model.Transaction) error
	PurchaseWithManualTransaction(context context.Context, transaction *model.Transaction) error
}

type purchaseRepository struct {
	client                *mongo.Client
	bookRepository        repository.BookRepository
	transactionRepository repository2.TransactionRepository
}

// PurchaseWithAutomaticTransaction handles purchasing with automatic transaction.
//
// context - context.Context, transaction *model.Transaction
// error
func (p purchaseRepository) PurchaseWithAutomaticTransaction(context context.Context, transaction *model.Transaction) error {

	callback := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		err := p.bookRepository.UpdateQuantity(sessionCtx, transaction.BookID.Hex(), transaction.Amount)
		if err != nil {
			return nil, err
		}

		err = p.transactionRepository.Insert(sessionCtx, transaction)
		if err != nil {
			return nil, err
		}

		//TO DEBUG TRX UNCOMMENT THIS LINE AND COMMENT ABOVE
		//if err == nil {
		//	return nil, errors.New("UnknownTransactionCommitResult")
		//}

		return nil, nil
	}

	session, err := p.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context)

	_, err = session.WithTransaction(context, callback)
	if err != nil {
		return err
	}

	return nil
}

// PurchaseWithManualTransaction performs a manual transaction for a purchase.
//
// It takes a context and a transaction model as parameters and returns an error.
func (p purchaseRepository) PurchaseWithManualTransaction(context context.Context, transaction *model.Transaction) error {

	callback := func(sessionCtx mongo.SessionContext) error {
		err := sessionCtx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}

		err = p.bookRepository.UpdateQuantity(sessionCtx, transaction.BookID.Hex(), transaction.Amount)
		if err != nil {
			sessionCtx.AbortTransaction(sessionCtx)
			return err
		}

		err = p.transactionRepository.Insert(sessionCtx, transaction)
		if err != nil {
			sessionCtx.AbortTransaction(sessionCtx)
			return err
		}

		//TO DEBUG TRX UNCOMMENT THIS LINE AND COMMENT ABOVE
		//if err == nil {
		//	sessionCtx.AbortTransaction(sessionCtx)
		//	return errors.New("UnknownTransactionCommitResult")
		//}

		return mongo2.CommitWithRetry(sessionCtx)
	}

	return p.client.UseSessionWithOptions(
		context, options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return mongo2.RunTransactionWithRetry(sctx, callback)
		},
	)
}

func NewPurchaseRepository(
	client *mongo.Client,
	bookRepository repository.BookRepository,
	transactionRepository repository2.TransactionRepository,
) PurchaseRepository {
	return &purchaseRepository{
		client:                client,
		bookRepository:        bookRepository,
		transactionRepository: transactionRepository,
	}

}
