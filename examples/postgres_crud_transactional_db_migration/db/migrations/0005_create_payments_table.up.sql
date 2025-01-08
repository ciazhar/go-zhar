-- Payments table to store payment-related information
CREATE TABLE payments
(
    id             SERIAL PRIMARY KEY,
    order_id       INT            NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    payment_method VARCHAR(50),             -- e.g., Credit Card, PayPal, Cash
    amount         NUMERIC(10, 2) NOT NULL,
    status         VARCHAR(50)    NOT NULL, -- e.g., Success, Failed, Pending
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);