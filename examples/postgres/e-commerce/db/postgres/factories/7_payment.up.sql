CREATE TABLE payment
(
    payment_id        SERIAL PRIMARY KEY,
    order_id          INT            NOT NULL,
    payment_method_id INT            NOT NULL,
    amount            DECIMAL(10, 2) NOT NULL,
    payment_date      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders (order_id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_method (payment_method_id)
);