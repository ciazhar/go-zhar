package controller

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/clickhouse_analytic_report/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	repo *repository.TransactionRepository
}

func NewTransactionController(repo *repository.TransactionRepository) *TransactionController {
	return &TransactionController{
		repo: repo,
	}
}

// Aggregate Total Amount for the Last Month
func (tc *TransactionController) AggregateTotalAmount(c *fiber.Ctx) error {
	ctx := context.Background()

	result, err := tc.repo.AggregateTotalAmountLastMonth(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to aggregate total amount",
			Data: logger.LogAndReturnError(ctx, err, "Failed to aggregate total amount", map[string]interface{}{
				"error": err.Error(),
			}),
		})
	}
	return c.JSON(result)
}

// Count Transactions Per Type for the Last Month
func (tc *TransactionController) CountTransactions(c *fiber.Ctx) error {
	ctx := context.Background()

	result, err := tc.repo.CountTransactionsPerTypeLastMonth(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to count transactions",
			Data: logger.LogAndReturnError(ctx, err, "Failed to count transactions", map[string]interface{}{
				"error": err.Error(),
			}),
		})
	}
	return c.JSON(result)
}

// Average Transaction Value Per User
func (tc *TransactionController) AverageTransactionValue(c *fiber.Ctx) error {
	ctx := context.Background()

	result, err := tc.repo.AverageTransactionValuePerUser(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to get average transaction value",
			Data: logger.LogAndReturnError(ctx, err, "Failed to get average transaction value", map[string]interface{}{
				"error": err.Error(),
			}),
		})
	}
	return c.JSON(result)
}

// Top 10 Users by Total Transaction Value
func (tc *TransactionController) TopUsers(c *fiber.Ctx) error {
	ctx := context.Background()

	result, err := tc.repo.TopUsersByTotalTransactionValue(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to get top users",
			Data: logger.LogAndReturnError(ctx, err, "Failed to get top users", map[string]interface{}{
				"error": err.Error(),
			}),
		})
	}
	return c.JSON(result)
}

// Get Transaction Summary
func (tc *TransactionController) GetTransactionSummary(c *fiber.Ctx) error {
	ctx := context.Background()

	result, err := tc.repo.GetTransactionSummary(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to get transaction summary",
			Data: logger.LogAndReturnError(ctx, err, "Failed to get transaction summary", map[string]interface{}{
				"error": err.Error(),
			}),
		})
	}
	return c.JSON(result)
}

func (tc *TransactionController) GetUsersWithMoreThanXTransactions(c *fiber.Ctx) error {
	// Assuming X is passed as a query parameter, e.g. /users-x-transactions?threshold=10
	threshold := c.QueryInt("threshold", 10)

	users, err := tc.repo.GetUsersWithMoreThanXTransactions(context.Background(), threshold)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (tc *TransactionController) GetSumOfTransactionAmountsPerDay(c *fiber.Ctx) error {
	// Assuming transaction type is passed as a query parameter, e.g. /purchase-sum?type=Purchase
	transactionType := c.Query("type", "Purchase")

	sums, err := tc.repo.GetSumOfTransactionAmountsPerDay(context.Background(), transactionType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sums)
}

func (tc *TransactionController) GetTotalRefundsProcessed(c *fiber.Ctx) error {
	total, err := tc.repo.GetTotalRefundsProcessed(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(total)
}

func (tc *TransactionController) GetPeakHourForTransactions(c *fiber.Ctx) error {
	peakHour, err := tc.repo.GetPeakHourForTransactions(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"peak_hour": peakHour})
}

func (tc *TransactionController) GetUsersWithBothPurchasesAndRefunds(c *fiber.Ctx) error {
	users, err := tc.repo.GetUsersWithBothPurchasesAndRefunds(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}
