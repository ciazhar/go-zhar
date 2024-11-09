package purchase

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/repository"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/purchase/controller"
	repository3 "github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/purchase/repository"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/purchase/service"
	repository2 "github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/transaction/repository"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(router fiber.Router, database *mongo.Database) {
	bookRepository := repository.NewBookRepository(database)
	transactionRepository := repository2.NewTransactionRepository(database)
	purchaseRepository := repository3.NewPurchaseRepository(database.Client(), bookRepository, transactionRepository)
	purchaseService := service.NewPurchaseService(purchaseRepository)
	purchaseController := controller.NewPurchaseController(purchaseService)

	router.Post("/purchase", purchaseController.Purchase)
}
