package mongo

import (
	"errors"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

func RunTransactionWithRetry(ctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
	for {
		err := txnFn(ctx) // Performs transaction.
		if err == nil {
			return nil
		}

		logger.LogError(ctx, err, "Transaction failed", nil)

		// If transient error, retry the whole transaction
		var cmdErr mongo.CommandError
		if errors.As(err, &cmdErr) && cmdErr.HasErrorLabel("TransientTransactionError") {
			logger.LogError(ctx, err, "TransientTransactionError, retrying transaction...", nil)
			continue
		}
		return err
	}
}

func CommitWithRetry(ctx mongo.SessionContext) error {
	for {
		err := ctx.CommitTransaction(ctx)
		switch e := err.(type) {
		case nil:
			return nil
		case mongo.CommandError:
			// Can retry commit
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				logger.LogInfo(ctx, "UnknownTransactionCommitResult, retrying commit operation...", nil)
				continue
			}
			return e
		default:
			return e
		}
	}
}
