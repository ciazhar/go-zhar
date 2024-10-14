CREATE TABLE customer
(
    customer_id SERIAL PRIMARY KEY,
    name        VARCHAR NOT NULL,
    email       VARCHAR NOT NULL UNIQUE,
    password    VARCHAR NOT NULL,
    address     TEXT
);