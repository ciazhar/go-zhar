package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	delayExchange   = "order.delay.exchange"
	delayQueue      = "order.delay.queue"
	dlxExchange     = "order.dlx"
	dlqQueue        = "order.dlq.queue"
	delayRoutingKey = "order.delay"
	dlqRoutingKey   = "order.dlq"
	ttlMilliseconds = 20000 // 10 seconds for testing
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	db   *sql.DB
)

func setupRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	ch.ExchangeDeclare(delayExchange, "direct", true, false, false, false, nil)
	ch.ExchangeDeclare(dlxExchange, "direct", true, false, false, false, nil)

	ch.QueueDeclare(delayQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    dlxExchange,
		"x-dead-letter-routing-key": dlqRoutingKey,
		"x-message-ttl":             int32(ttlMilliseconds),
	})
	ch.QueueBind(delayQueue, delayRoutingKey, delayExchange, false, nil)
	ch.QueueDeclare(dlqQueue, true, false, false, false, nil)
	ch.QueueBind(dlqQueue, dlqRoutingKey, dlxExchange, false, nil)
}

func setupPostgres() {
	var err error
	db, err = sql.Open("postgres", "postgres://user:password@localhost:5432/demo?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS stock (
			id SERIAL PRIMARY KEY,
			product_id INT NOT NULL REFERENCES product(id),
			quantity INT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			product_id INT NOT NULL REFERENCES product(id),
			quantity INT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	seedData()

	log.Println("PostgreSQL tables created or verified")
}

func seedData() {
	// Check if there's already a product
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM product`).Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count products: %v", err)
	}

	if count == 0 {
		// Insert dummy products
		_, err := db.Exec(`
			INSERT INTO product (name) VALUES 
			('Apple'), 
			('Banana'), 
			('Orange')
		`)
		if err != nil {
			log.Fatalf("Failed to seed products: %v", err)
		}

		// Insert corresponding stock
		_, err = db.Exec(`
			INSERT INTO stock (product_id, quantity) VALUES 
			(1, 100), 
			(2, 150), 
			(3, 200)
		`)
		if err != nil {
			log.Fatalf("Failed to seed stock: %v", err)
		}
		log.Println("Seeded initial product and stock data")
	} else {
		log.Println("Product data already exists, skipping seed")
	}
}

func reserveStock(productID, qty int) error {
	res, err := db.Exec(`UPDATE stock SET quantity = quantity - $1 WHERE product_id = $2 AND quantity >= $1`, qty, productID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("not enough stock")
	}
	return nil
}

func revertStock(orderID string) {
	var productID, qty int
	err := db.QueryRow(`SELECT product_id, quantity FROM orders WHERE id = $1`, orderID).Scan(&productID, &qty)
	if err != nil {
		log.Printf("Failed to get order data: %v", err)
		return
	}
	_, err = db.Exec(`UPDATE stock SET quantity = quantity + $1 WHERE product_id = $2`, qty, productID)
	if err != nil {
		log.Printf("Failed to revert stock: %v", err)
		return
	}
	log.Printf("StockService: Reverted %d stock for product %d (order %s)", qty, productID, orderID)
}

func updateOrderStatus(orderID, status string) {
	_, err := db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, orderID)
	if err != nil {
		log.Printf("OrderService: Failed to update order status: %v", err)
		return
	}
	log.Printf("OrderService: Updated order %s status to %s", orderID, status)
}

func publishOrder(orderID string) {
	err := ch.PublishWithContext(context.Background(),
		delayExchange,
		delayRoutingKey,
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(orderID),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return
	}
	log.Printf("OrderService: Published order %s to delay queue", orderID)
}

func consumeDLQ() {
	msgs, err := ch.Consume(dlqQueue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to consume from DLQ: %v", err)
	}
	go func() {
		for msg := range msgs {
			orderID := string(msg.Body)

			var currentStatus string
			err := db.QueryRow(`SELECT status FROM orders WHERE id = $1`, orderID).Scan(&currentStatus)
			if err != nil {
				log.Printf("DLQConsumer: Failed to fetch order %s status: %v", orderID, err)
				continue
			}

			if currentStatus != "PENDING" {
				log.Printf("DLQConsumer: Order %s already processed (status: %s), skipping", orderID, currentStatus)
				continue
			}

			revertStock(orderID)
			updateOrderStatus(orderID, "EXPIRED")
		}
	}()
}

func triggerPayment(orderID string) {
	log.Printf("PaymentService: Triggering payment for order %s (simulated)", orderID)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	productID := 1
	qty := 1
	orderID := fmt.Sprintf("ORDER-%d", time.Now().UnixNano())

	err := reserveStock(productID, qty)
	if err != nil {
		http.Error(w, "Stock not available", http.StatusConflict)
		return
	}

	_, err = db.Exec(`INSERT INTO orders (id, product_id, quantity, status) VALUES ($1, $2, $3, 'PENDING')`, orderID, productID, qty)
	if err != nil {
		log.Printf("Failed to insert order: %v", err)
		return
	}

	publishOrder(orderID)
	triggerPayment(orderID)
	fmt.Fprintf(w, "Order %s placed\n", orderID)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")
	if orderID == "" {
		http.Error(w, "Missing orderID", http.StatusBadRequest)
		return
	}
	var status string
	err := db.QueryRow(`SELECT status FROM orders WHERE id = $1`, orderID).Scan(&status)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Order %s status: %s\n", orderID, status)
}

func main() {
	setupPostgres()
	setupRabbitMQ()
	consumeDLQ()

	// Only close when server stops
	defer conn.Close()
	defer ch.Close()
	defer db.Close()

	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/status", statusHandler)

	log.Println("API is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
