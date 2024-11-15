package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Transaction struct {
	Timestamp time.Time
	Type      string
	UserID    string
	Amount    float64
}

// Helper function to generate random transactions
func generateTransactions(batchSize int) []Transaction {
	transactionTypes := []string{"Purchase", "Refund", "Subscription", "Cancellation"}
	transactions := make([]Transaction, 0, batchSize)

	for i := 0; i < batchSize; i++ {
		transactions = append(transactions, Transaction{
			Timestamp: time.Now().Add(time.Duration(rand.Intn(100000)) * time.Second * -1), // Random timestamp in the past
			Type:      transactionTypes[rand.Intn(len(transactionTypes))],
			UserID:    generateUserID(),
			Amount:    float64(rand.Intn(10000)) + rand.Float64(),
		})
	}
	return transactions
}

// Helper function to generate random UserID
func generateUserID() string {
	return fmt.Sprintf("user_%c%c%d", 'A'+rand.Intn(26), 'A'+rand.Intn(26), rand.Intn(10))
}

func main() {
	ctx := context.Background()

	// Step 1: Connect to ClickHouse
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"}, // Adjust the address accordingly
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Debug: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	// Set batch size and total number of records to insert
	batchSize := 1000             // You can adjust the batch size based on your system's memory
	totalRecords := 1_000_000_000 // 1 billion records
	totalBatches := totalRecords / batchSize

	// Start inserting batches
	for batchNum := 0; batchNum < totalBatches; batchNum++ {
		transactions := generateTransactions(batchSize)

		// Step 3: Insert data in batches
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO transactions (Timestamp, Type, UserID, Amount)")
		if err != nil {
			log.Fatalf("Failed to prepare batch: %v", err)
		}

		for _, tx := range transactions {
			if err := batch.Append(tx.Timestamp, tx.Type, tx.UserID, tx.Amount); err != nil {
				log.Fatalf("Failed to append transaction to batch: %v", err)
			}
		}

		if err := batch.Send(); err != nil {
			log.Fatalf("Failed to send batch: %v", err)
		}

		// Log the progress every 10,000 batches
		if batchNum%10000 == 0 {
			log.Printf("Inserted %d records so far...", batchNum*batchSize)
		}
	}

	log.Println("Successfully inserted 1 billion records!")
}
