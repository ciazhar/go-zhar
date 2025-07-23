-- Tidak perlu buat tabel manual, karena JDBC Sink Connector bisa auto-create jika auto.create=true
-- Tapi bisa ditambahkan jika ingin kontrol skema

CREATE TABLE IF NOT EXISTS public.orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    product VARCHAR(255),
    quantity INTEGER,
    price NUMERIC(10, 2),
    status VARCHAR(50),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
