CREATE TABLE products
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255),
    price      NUMERIC(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
