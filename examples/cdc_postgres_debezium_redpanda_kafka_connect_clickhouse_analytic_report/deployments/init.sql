CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    product TEXT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tambahkan dummy data
INSERT INTO orders (user_id, product, quantity, price, status)
VALUES
    (1, 'Keyboard', 2, 150.00, 'confirmed'),
    (2, 'Mouse', 1, 50.00, 'shipped'),
    (3, 'Monitor', 1, 300.00, 'pending');
