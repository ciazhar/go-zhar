CREATE TABLE payment_method
(
    payment_method_id SERIAL PRIMARY KEY,
    name              VARCHAR NOT NULL,
    description       TEXT
);