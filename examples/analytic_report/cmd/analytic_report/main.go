package main

import (
	"log"

	"github.com/ciazhar/go-start-small/examples/analytic_report/internal/controller"
	"github.com/ciazhar/go-start-small/examples/analytic_report/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/clickhouse"
	"github.com/gofiber/fiber/v2"
)

func main() {
    // Initialize Fiber app
    app := fiber.New()

    // ClickHouse connection setup

	conn := clickhouse.Init(
		"localhost:9000",
		"default",
		"default",
		"",
		true,
	)

    repo := repository.NewTransactionRepository(conn)
    tc := controller.NewTransactionController(repo)

    // Define routes
    app.Get("/aggregate-amount", tc.AggregateTotalAmount)
	app.Get("/count-transactions", tc.CountTransactions)
	app.Get("/average-transaction-value", tc.AverageTransactionValue)
	app.Get("/top-users", tc.TopUsers)
    app.Get("/daily-volume", tc.GetTransactionSummary)
	app.Get("/users-x-transactions", tc.GetUsersWithMoreThanXTransactions)
	app.Get("/purchase-sum", tc.GetSumOfTransactionAmountsPerDay)
	app.Get("/total-refunds", tc.GetTotalRefundsProcessed)
	app.Get("/peak-hour", tc.GetPeakHourForTransactions)
	app.Get("/both-purchases-refunds", tc.GetUsersWithBothPurchasesAndRefunds)

    // Start the server
    log.Fatal(app.Listen(":3000"))
}
