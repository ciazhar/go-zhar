CREATE TABLE payments
(
    id         SERIAL PRIMARY KEY,
    order_id   INT            NOT NULL REFERENCES orders (id),
    amount     NUMERIC(10, 2) NOT NULL,
    status     VARCHAR(50)    NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);