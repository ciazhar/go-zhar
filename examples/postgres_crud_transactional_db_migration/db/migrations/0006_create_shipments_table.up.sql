-- Shipment table to store shipment details
CREATE TABLE shipments
(
    id              SERIAL PRIMARY KEY,
    order_id        INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    tracking_number VARCHAR(255),
    carrier         VARCHAR(255),                  -- e.g., FedEx, DHL
    status          VARCHAR(50) DEFAULT 'Pending', -- e.g., Pending, Shipped, Delivered
    shipped_at      TIMESTAMP,
    delivered_at    TIMESTAMP
);