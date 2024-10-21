CREATE TABLE transactions
(
    Timestamp DateTime CODEC (Delta, ZSTD),
    Type      Enum8('Purchase' = 1, 'Refund' = 2, 'Subscription' = 3, 'Cancellation' = 4) CODEC (LZ4),
    Value     Int32 CODEC (Delta, ZSTD),
    UserID    String CODEC (ZSTD), -- Keep UserID as String
    Amount    Float64 CODEC (ZSTD)
)
    ENGINE = MergeTree()
        PARTITION BY toYYYYMM(Timestamp)
        ORDER BY (Timestamp, cityHash64(UserID)) -- Order by hashed UserID for sampling
        SAMPLE BY cityHash64(UserID) -- Sampling by hashed UserID
        TTL Timestamp + INTERVAL 1 YEAR DELETE
        SETTINGS index_granularity = 8192;