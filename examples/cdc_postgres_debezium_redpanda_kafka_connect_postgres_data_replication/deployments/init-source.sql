-- Buat tabel sample untuk CDC
CREATE TABLE public.orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Tambahkan sample data awal
INSERT INTO public.orders (user_id, product, quantity, price)
VALUES
(1, 'Keyboard', 2, 350000),
(2, 'Mouse', 1, 150000),
(3, 'Monitor', 1, 2750000);
