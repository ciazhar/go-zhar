CREATE TABLE transactions
(
    event_id           UUID,                -- Unique identifier for each transaction
    user_id            UInt64,              -- User ID (partition key for optimized user-based queries)
    transaction_type   Enum8('purchase' = 1, 'refund' = 2), -- Compact type for transaction categorization
    amount             Decimal(18, 2),      -- Monetary value with precision
    timestamp          DateTime,            -- Transaction timestamp
    product_id         UInt32,              -- Product ID for category-level analysis
    category           String,              -- Product category for filtering (if needed)
    payment_method     Enum8('credit_card' = 1, 'paypal' = 2, 'bank_transfer' = 3, 'cash' = 4), -- Optimized for known methods
    is_fraudulent      UInt8 DEFAULT 0,     -- Fraud detection flag (1 = Yes, 0 = No)
    metadata           String,              -- JSON string for flexible additional data
    created_at         DateTime DEFAULT now()
)
    ENGINE = ReplacingMergeTree(created_at)
PARTITION BY (toYYYYMM(timestamp), intHash64(user_id) % 100)  -- Partition by month and hash of user_id
ORDER BY (user_id, timestamp)            -- Ordered by user and time for efficient range queries
PRIMARY KEY (user_id, timestamp)         -- Defines uniqueness within partitions
TTL created_at + INTERVAL 1 YEAR         -- Automatically clean up data older than 1 year
SETTINGS index_granularity = 8192;       -- Granularity for index efficiency