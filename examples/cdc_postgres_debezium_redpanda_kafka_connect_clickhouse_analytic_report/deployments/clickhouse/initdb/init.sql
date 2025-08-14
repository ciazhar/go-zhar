CREATE TABLE IF NOT EXISTS default."orders" (
  id UInt64,
  user_id UInt64,
  product_id UInt64,
  quantity UInt32,
  total_amount Float64,
  created_at DateTime
) ENGINE = MergeTree()
ORDER BY id;