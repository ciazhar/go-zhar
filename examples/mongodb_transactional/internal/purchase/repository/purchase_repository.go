package repository

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/repository"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/transaction/model"
	repository2 "github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/transaction/repository"
	mongo2 "github.com/ciazhar/go-start-small/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type PurchaseRepository struct {
	client                *mongo.Client
	bookRepository        *repository.BookRepository
	transactionRepository *repository2.TransactionRepository
}

// PurchaseWithAutomaticTransaction handles purchasing with automatic transaction.
//
// context - context.Context, transaction *model.Transaction
// error
func (p *PurchaseRepository) PurchaseWithAutomaticTransaction(context context.Context, transaction *model.Transaction) (err error) {

	callback := func(sessionCtx mongo.SessionContext) (res interface{}, err error) {
		if err = p.bookRepository.UpdateQuantity(sessionCtx, transaction.BookID.Hex(), transaction.Amount); err != nil {
			return
		}

		if err = p.transactionRepository.Insert(sessionCtx, transaction); err != nil {
			return
		}

		//TO DEBUG TRX UNCOMMENT THIS LINE AND COMMENT ABOVE
		//if err == nil {
		//	return nil, errors.New("UnknownTransactionCommitResult")
		//}

		return
	}

	session, err := p.client.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(context)

	if _, err = session.WithTransaction(context, callback); err != nil {
		return
	}

	return
}

// PurchaseWithManualTransaction performs a manual transaction for a purchase.
//
// It takes a context and a transaction model as parameters and returns an error.
func (p *PurchaseRepository) PurchaseWithManualTransaction(context context.Context, transaction *model.Transaction) (err error) {

	callback := func(sessionCtx mongo.SessionContext) (err error) {
		if err = sessionCtx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		); err != nil {
			return
		}

		err = p.bookRepository.UpdateQuantity(sessionCtx, transaction.BookID.Hex(), transaction.Amount)
		if err != nil {
			err = sessionCtx.AbortTransaction(sessionCtx)
			if err != nil {
				return
			}
			return
		}

		err = p.transactionRepository.Insert(sessionCtx, transaction)
		if err != nil {
			err = sessionCtx.AbortTransaction(sessionCtx)
			if err != nil {
				return
			}
			return
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
	bookRepository *repository.BookRepository,
	transactionRepository *repository2.TransactionRepository,
) *PurchaseRepository {
	return &PurchaseRepository{
		client:                client,
		bookRepository:        bookRepository,
		transactionRepository: transactionRepository,
	}

}
