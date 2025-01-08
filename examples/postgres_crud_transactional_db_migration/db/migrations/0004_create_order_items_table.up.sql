-- Order items table to store items in each order
CREATE TABLE order_items
(
    id          SERIAL PRIMARY KEY,
    order_id    INT            NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id  INT            NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
    quantity    INT            NOT NULL,
    price       NUMERIC(10, 2) NOT NULL, -- Price at the time of the order
    total_price NUMERIC(10, 2) NOT NULL, -- Calculated as quantity * price
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);