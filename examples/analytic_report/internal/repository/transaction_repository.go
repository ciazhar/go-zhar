package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/analytic_report/internal/model"
)

type TransactionRepository struct {
	conn clickhouse.Conn
}

// NewTransactionRepository initializes a new TransactionRepository
func NewTransactionRepository(conn clickhouse.Conn) *TransactionRepository {
	return &TransactionRepository{conn: conn}
}

// AggregateTotalAmountLastMonth calculates the total Amount spent on each Type in the last month.
func (r *TransactionRepository) AggregateTotalAmountLastMonth(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            Type, 
            SUM(Amount) AS TotalAmount
        FROM 
            transactions
        WHERE 
            Timestamp >= NOW() - INTERVAL '1 month'
        GROUP BY 
            Type
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var typ string
		var totalAmount float64
		if err := rows.Scan(&typ, &totalAmount); err != nil {
			return nil, err
		}
		result := map[string]interface{}{
			"Type":        typ,
			"TotalAmount": totalAmount,
		}
		results = append(results, result)
	}
	return results, nil
}

// CountTransactionsPerTypeLastMonth returns the count of transactions per type in the last month.
func (r *TransactionRepository) CountTransactionsPerTypeLastMonth(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            Type, 
            COUNT(*) AS TransactionCount
        FROM 
            transactions
        WHERE 
            Timestamp >= NOW() - INTERVAL '1 month'
        GROUP BY 
            Type
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var typ string
		var count int64
		if err := rows.Scan(&typ, &count); err != nil {
			return nil, err
		}
		result := map[string]interface{}{
			"Type":             typ,
			"TransactionCount": count,
		}
		results = append(results, result)
	}
	return results, nil
}

// AverageTransactionValuePerUser calculates the average transaction value for each user.
func (r *TransactionRepository) AverageTransactionValuePerUser(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            UserID, 
            AVG(Amount) AS AvgTransactionValue
        FROM 
            transactions
        GROUP BY 
            UserID
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var userID string
		var avgValue float64
		if err := rows.Scan(&userID, &avgValue); err != nil {
			return nil, err
		}
		result := map[string]interface{}{
			"UserID":              userID,
			"AvgTransactionValue": avgValue,
		}
		results = append(results, result)
	}
	return results, nil
}

// TopUsersByTotalTransactionValue returns the top 10 users by total transaction value.
func (r *TransactionRepository) TopUsersByTotalTransactionValue(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            UserID, 
            SUM(Amount) AS TotalAmount
        FROM 
            transactions
        GROUP BY 
            UserID
        ORDER BY 
            TotalAmount DESC
        LIMIT 10
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var userID string
		var totalAmount float64
		if err := rows.Scan(&userID, &totalAmount); err != nil {
			return nil, err
		}
		result := map[string]interface{}{
			"UserID":      userID,
			"TotalAmount": totalAmount,
		}
		results = append(results, result)
	}
	return results, nil
}

// GetTransactionSummary fetches the transaction summary from ClickHouse and fills missing dates in the past 30 days.
func (r *TransactionRepository) GetTransactionSummary(ctx context.Context) ([]model.TransactionSummary, error) {
    query := `
        SELECT toStartOfDay(Timestamp) AS date,
               SUM(Type = 'Purchase')     AS Purchase,
               SUM(Type = 'Refund')       AS Refund,
               SUM(Type = 'Subscription') AS Subscription,
               SUM(Type = 'Cancellation') AS Cancellation
        FROM transactions
        WHERE Timestamp >= now() - INTERVAL 30 DAY
          AND Timestamp < now() + INTERVAL 1 DAY
        GROUP BY toStartOfDay(Timestamp)
        ORDER BY date;
    `

    rows, err := r.conn.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }
    defer rows.Close()

    var summaries []model.TransactionSummary
    for rows.Next() {
        var summary model.TransactionSummary
        var date time.Time

        if err := rows.Scan(&date, &summary.Purchase, &summary.Refund, &summary.Subscription, &summary.Cancellation); err != nil {
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
func mergeSummaries(dates []time.Time, summaries []model.TransactionSummary) []model.TransactionSummary {
    summaryMap := make(map[time.Time]model.TransactionSummary)
    for _, s := range summaries {
        summaryMap[s.Date] = s
    }

    result := make([]model.TransactionSummary, len(dates))
    for i, date := range dates {
        if summary, exists := summaryMap[date]; exists {
            result[i] = summary
        } else {
            result[i] = model.TransactionSummary{Date: date}
        }
    }
    return result
}

func (r *TransactionRepository) GetUsersWithMoreThanXTransactions(ctx context.Context, threshold int) ([]model.User, error) {
	query := `
        SELECT UserID, COUNT(*) AS transaction_count
        FROM transactions
        WHERE Timestamp >= now() - INTERVAL 30 DAY
        GROUP BY UserID
        HAVING transaction_count > ?
    `
	rows, err := r.conn.Query(ctx, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.TransactionCount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *TransactionRepository) GetSumOfTransactionAmountsPerDay(ctx context.Context, transactionType string) ([]model.TransactionSummary, error) {
	query := `
        SELECT toStartOfDay(Timestamp) AS date, SUM(Amount) AS sum
        FROM transactions
        WHERE Type = ? AND Timestamp >= now() - INTERVAL 30 DAY
        GROUP BY toStartOfDay(Timestamp)
        ORDER BY date;
    `
	rows, err := r.conn.Query(ctx, query, transactionType)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var summaries []model.TransactionSummary
	for rows.Next() {
		var summary model.TransactionSummary
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
        SELECT SUM(Amount)
        FROM transactions
        WHERE Type = 'Refund' AND Timestamp >= now() - INTERVAL 30 DAY;
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
        SELECT toHour(Timestamp) AS hour
        FROM transactions
        WHERE Timestamp >= now() - INTERVAL 30 DAY
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

func (r *TransactionRepository) GetUsersWithBothPurchasesAndRefunds(ctx context.Context) ([]model.User, error) {
	query := `
        SELECT UserID
        FROM transactions
        WHERE Timestamp >= now() - INTERVAL 30 DAY
        GROUP BY UserID
        HAVING SUM(Type = 'Purchase') > 0 AND SUM(Type = 'Refund') > 0;
    `
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
