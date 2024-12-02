CREATE TABLE orders
(
    id            SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);