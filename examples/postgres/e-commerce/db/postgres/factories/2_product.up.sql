CREATE TABLE product
(
    product_id  SERIAL PRIMARY KEY,
    name        VARCHAR   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    category_id INT            NOT NULL,
    FOREIGN KEY (category_id) REFERENCES category (category_id)
);