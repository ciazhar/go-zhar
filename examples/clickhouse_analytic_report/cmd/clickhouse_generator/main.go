package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/clickhouse_analytic_report/db/migrations"
	clickhouse2 "github.com/ciazhar/go-start-small/pkg/clickhouse"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	EventID         string
	UserID          uint64
	TransactionType string
	Amount          decimal.Decimal // Use Decimal type from shopspring/decimal
	Timestamp       time.Time
	ProductID       uint32
	Category        string
	PaymentMethod   string
	IsFraudulent    uint8
	Metadata        string
	CreatedAt       time.Time
}

const (
	BatchSize    = 1000
	TotalRecords = 10_000
	DBAddr       = "localhost:9000"
	DBName       = "default"
	DBUser       = "default"
	DBPassword   = ""
)

var (
	transactionTypes = []string{"purchase", "refund"}
	categories       = []string{"electronics", "clothing", "home", "sports", "books"}
	paymentMethods   = []string{"credit_card", "paypal", "bank_transfer", "cash"}
)

func main() {
	ctx := context.Background()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{DBAddr},
		Auth: clickhouse.Auth{
			Database: DBName,
			Username: DBUser,
			Password: DBPassword,
		},
		Debug: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer conn.Close()

	// Initialize database
	clickhouse2.InitDBMigration(DBAddr, DBName, DBUser, DBPassword, migrations.MigrationsFS)

	totalBatches := TotalRecords / BatchSize
	log.Printf("Starting data generation: Total Records=%d, Batch Size=%d, Total Batches=%d", TotalRecords, BatchSize, totalBatches)

	for batchNum := 0; batchNum < totalBatches; batchNum++ {
		transactions := generateTransactions(BatchSize)

		batch, err := conn.PrepareBatch(ctx, `INSERT INTO transactions (
			event_id, user_id, transaction_type, amount, timestamp, 
			product_id, category, payment_method, is_fraudulent, metadata, created_at)`)
		if err != nil {
			log.Fatalf("Failed to prepare batch: %v", err)
		}

		for _, tx := range transactions {
			if err := batch.Append(
				tx.EventID, tx.UserID, tx.TransactionType, tx.Amount, tx.Timestamp,
				tx.ProductID, tx.Category, tx.PaymentMethod, tx.IsFraudulent, tx.Metadata, tx.CreatedAt,
			); err != nil {
				log.Fatalf("Failed to append transaction to batch: %v", err)
			}
		}

		if err := batch.Send(); err != nil {
			log.Fatalf("Failed to send batch: %v", err)
		}

		if batchNum%100 == 0 {
			log.Printf("Batch %d/%d inserted.", batchNum+1, totalBatches)
		}
	}

	log.Println("Data generation and insertion completed successfully.")
}

func generateTransactions(batchSize int) []Transaction {
	transactions := make([]Transaction, 0, batchSize)

	for i := 0; i < batchSize; i++ {
		userID := uint64(rand.Int63n(1_000_000))
		amount := decimal.NewFromFloat(generateRandomAmount()).Round(2) // Format as Decimal(18, 2)
		timestamp := time.Now().Add(time.Duration(rand.Intn(100000)) * time.Second * -1)
		productID := uint32(rand.Intn(10000))
		category := categories[rand.Intn(len(categories))]
		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		isFraudulent := uint8(rand.Intn(2))

		metadata := map[string]interface{}{
			"ip_address":   fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256)),
			"device_type":  "mobile",
			"geo_location": fmt.Sprintf("%.6f,%.6f", rand.Float64()*180-90, rand.Float64()*360-180),
		}
		metadataJSON, _ := json.Marshal(metadata)

		transactions = append(transactions, Transaction{
			EventID:         uuid.New().String(),
			UserID:          userID,
			TransactionType: transactionTypes[rand.Intn(len(transactionTypes))],
			Amount:          amount,
			Timestamp:       timestamp,
			ProductID:       productID,
			Category:        category,
			PaymentMethod:   paymentMethod,
			IsFraudulent:    isFraudulent,
			Metadata:        string(metadataJSON),
			CreatedAt:       time.Now(),
		})
	}
	return transactions
}

func generateRandomAmount() float64 {
	mean := 100.0
	stddev := 50.0
	amount := rand.NormFloat64()*stddev + mean
	if amount < 0 {
		amount = 0.01
	}
	return amount
}

func formatAmount(amount float64) string {
	// Format the amount as a string with two decimal places
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
