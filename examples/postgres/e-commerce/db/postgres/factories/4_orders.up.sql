CREATE TABLE orders
(
    order_id     SERIAL PRIMARY KEY,
    customer_id  INT            NOT NULL,
    order_date   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status       VARCHAR    NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customer (customer_id)
);