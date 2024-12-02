CREATE TABLE inventory
(
    product_id   SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    stock        INT          NOT NULL
);