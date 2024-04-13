package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func RunTransactionWithRetry(sctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
	for {
		err := txnFn(sctx) // Performs transaction.
		if err == nil {
			return nil
		}

		log.Println("Transaction aborted. Caught exception during transaction.")

		// If transient error, retry the whole transaction
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
			log.Println("TransientTransactionError, retrying transaction...")
			continue
		}
		return err
	}
}

func CommitWithRetry(sctx mongo.SessionContext) error {
	for {
		err := sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			return nil
		case mongo.CommandError:
			// Can retry commit
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				log.Println("UnknownTransactionCommitResult, retrying commit operation...")
				continue
			}
			return e
		default:
			return e
		}
	}
}
