package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/clickhouse_analytic_report/internal/model"
)

type TransactionRepository struct {
	conn clickhouse.Conn
}

// NewTransactionRepository initializes a new TransactionRepository
func NewTransactionRepository(conn clickhouse.Conn) *TransactionRepository {
	return &TransactionRepository{conn: conn}
}

// AggregateTotalAmountLastMonth calculates the total amount for each transaction type in the last month.
func (r *TransactionRepository) AggregateTotalAmountLastMonth(ctx context.Context) ([]model.AggregateTotalAmountResponse, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            transaction_type AS Type, 
            CAST(SUM(amount) AS Float64) AS TotalAmount
        FROM 
            transactions
        WHERE 
            timestamp >= NOW() - INTERVAL 1 MONTH
        GROUP BY 
            transaction_type
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.AggregateTotalAmountResponse
	for rows.Next() {
		var resp model.AggregateTotalAmountResponse
		if err := rows.Scan(&resp.Type, &resp.TotalAmount); err != nil {
			return nil, err
		}

		results = append(results, resp)
	}

	return results, nil
}

// CountTransactionsPerTypeLastMonth counts transactions per type for the last month.
func (r *TransactionRepository) CountTransactionsPerTypeLastMonth(ctx context.Context) ([]model.CountTransactionsPerTypeResponse, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            transaction_type AS Type, 
            COUNT(*) AS TransactionCount
        FROM 
            transactions
        WHERE 
            timestamp >= NOW() - INTERVAL 1 MONTH
        GROUP BY 
            transaction_type
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.CountTransactionsPerTypeResponse
	for rows.Next() {
		var resp model.CountTransactionsPerTypeResponse
		if err := rows.Scan(&resp.Type, &resp.TransactionCount); err != nil {
			return nil, err
		}
		results = append(results, resp)
	}
	return results, nil
}

// AverageTransactionValuePerUser calculates the average transaction value for each user.
func (r *TransactionRepository) AverageTransactionValuePerUser(ctx context.Context) ([]model.AverageTransactionValuePerUserResponse, error) {
	query := `
        SELECT 
            user_id, 
            AVG(amount) AS avg_transaction_value
        FROM transactions
        GROUP BY user_id
        ORDER BY avg_transaction_value desc 
    `
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []model.AverageTransactionValuePerUserResponse
	for rows.Next() {
		var response model.AverageTransactionValuePerUserResponse
		if err := rows.Scan(&response.UserID, &response.AverageTransactionValue); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, response)
	}
	return results, nil
}

// TopUsersByTotalTransactionValue returns the top 10 users by total transaction value.
func (r *TransactionRepository) TopUsersByTotalTransactionValue(ctx context.Context) ([]model.TopUsersByTotalTransactionValueResponse, error) {
	query := `
        SELECT 
            user_id, 
            CAST(SUM(amount) AS Float64) AS total_amount
        FROM transactions
        GROUP BY user_id
        ORDER BY total_amount DESC
        LIMIT 10;
    `
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []model.TopUsersByTotalTransactionValueResponse
	for rows.Next() {
		var response model.TopUsersByTotalTransactionValueResponse
		if err := rows.Scan(&response.UserID, &response.TotalAmount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, response)
	}
	return results, nil
}

// GetTransactionSummary fetches the transaction summary from ClickHouse and fills missing dates in the past 30 days.
func (r *TransactionRepository) GetTransactionSummary(ctx context.Context) ([]model.TransactionDailySummaryResponse, error) {
	query := `
        SELECT toStartOfDay(timestamp)                AS date,
		   CAST(SUM(amount) AS Float64)                            as amount,
		   SUM(transaction_type = 'purchase')     AS purchase,
		   SUM(transaction_type = 'refund')       AS refund
        FROM transactions
        WHERE timestamp >= now() - INTERVAL 30 DAY
        GROUP BY date
        ORDER BY date;
    `
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var summaries []model.TransactionDailySummaryResponse
	for rows.Next() {
		var summary model.TransactionDailySummaryResponse
		var date time.Time
		if err := rows.Scan(&date, &summary.Amount, &summary.Purchase, &summary.Refund); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Ensure the date is in UTC and truncated to midnight (00:00:00 UTC)
		summary.Date = date.UTC().Truncate(24 * time.Hour)
		summaries = append(summaries, summary)
	}

	// Generate a full list of dates in time.Time for the last 30 days
	dates := generateDateRange(30)

	// Merge the fetched summaries with the full list of dates
	mergedSummaries := mergeSummaries(dates, summaries)

	return mergedSummaries, nil
}

// generateDateRange generates a list of dates (as time.Time) for the past 'days' days.
func generateDateRange(days int) []time.Time {
	now := time.Now().UTC().Truncate(24 * time.Hour) // Truncate to midnight UTC
	dates := make([]time.Time, days)
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		dates[i] = date
	}
	return dates
}

// mergeSummaries fills in missing dates in the transaction summary (dates are in time.Time).
func mergeSummaries(dates []time.Time, summaries []model.TransactionDailySummaryResponse) []model.TransactionDailySummaryResponse {
	summaryMap := make(map[time.Time]model.TransactionDailySummaryResponse)
	for _, s := range summaries {
		summaryMap[s.Date] = s
	}

	result := make([]model.TransactionDailySummaryResponse, len(dates))
	for i, date := range dates {
		if summary, exists := summaryMap[date]; exists {
			result[i] = summary
		} else {
			result[i] = model.TransactionDailySummaryResponse{Date: date}
		}
	}
	return result
}

// GetUsersWithMoreThanXTransactions retrieves users with more than X transactions in the last 30 days.
func (r *TransactionRepository) GetUsersWithMoreThanXTransactions(ctx context.Context, threshold int) ([]model.GetUsersWithMoreThanXTransactionsResponse, error) {
	query := `
        SELECT user_id AS UserID, COUNT(*) AS transaction_count
        FROM transactions
        WHERE timestamp >= now() - INTERVAL 30 DAY
        GROUP BY user_id
        HAVING transaction_count > ?
    `
	rows, err := r.conn.Query(ctx, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []model.GetUsersWithMoreThanXTransactionsResponse
	for rows.Next() {
		var user model.GetUsersWithMoreThanXTransactionsResponse
		if err := rows.Scan(&user.UserID, &user.TransactionCount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// GetSumOfTransactionAmountsPerDay retrieves the sum of transaction amounts per day for a specific type.
func (r *TransactionRepository) GetSumOfTransactionAmountsPerDay(ctx context.Context, transactionType string) ([]model.TransactionDailySummaryResponse, error) {
	query := `
        SELECT toStartOfDay(timestamp) AS date, CAST(SUM(amount) AS Float64) AS sum
        FROM transactions
        WHERE transaction_type = ? AND timestamp >= now() - INTERVAL 30 DAY
        GROUP BY toStartOfDay(timestamp)
        ORDER BY date;
    `
	rows, err := r.conn.Query(ctx, query, transactionType)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var summaries []model.TransactionDailySummaryResponse
	for rows.Next() {
		var summary model.TransactionDailySummaryResponse
		if err := rows.Scan(&summary.Date, &summary.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		summary.Date = summary.Date.UTC()
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (r *TransactionRepository) GetTotalRefundsProcessed(ctx context.Context) (float64, error) {
	query := `
        SELECT CAST(SUM(amount) AS Float64)
        FROM transactions
        WHERE transaction_type = 'refund' AND timestamp >= now() - INTERVAL 30 DAY;
    `
	var total float64
	err := r.conn.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return total, nil
}

func (r *TransactionRepository) GetPeakHourForTransactions(ctx context.Context) (uint8, error) {
	query := `
        SELECT toHour(timestamp) AS hour
        FROM transactions
        WHERE timestamp >= now() - INTERVAL 30 DAY
        GROUP BY hour
        ORDER BY COUNT(*) DESC
        LIMIT 1;
    `
	var peakHour uint8
	err := r.conn.QueryRow(ctx, query).Scan(&peakHour)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return peakHour, nil
}

// GetUsersWithBothPurchasesAndRefunds retrieves users who made both purchases and refunds in the last 30 days.
func (r *TransactionRepository) GetUsersWithBothPurchasesAndRefunds(ctx context.Context) ([]model.GetUsersWithBothPurchasesAndRefundsResponse, error) {
	query := `
        SELECT user_id AS UserID
        FROM transactions
        WHERE timestamp >= now() - INTERVAL 30 DAY
        GROUP BY user_id
        HAVING SUM(transaction_type = 'purchase') > 0 AND SUM(transaction_type = 'refund') > 0;
    `
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []model.GetUsersWithBothPurchasesAndRefundsResponse
	for rows.Next() {
		var user model.GetUsersWithBothPurchasesAndRefundsResponse
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
