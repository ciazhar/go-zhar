-- Orders table to store order-specific details
CREATE TABLE orders
(
    id           SERIAL PRIMARY KEY,
    customer_id  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    status       VARCHAR(50) DEFAULT 'Pending', -- e.g., Pending, Completed, Canceled
    total_amount NUMERIC(10, 2),
    created_at   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);