-- name: CreateOrder :exec
INSERT INTO orders (customer_id, total_amount, status)
VALUES (@customer_id::int, @total_amount::float, @status::varchar);

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = 'Shipped'
WHERE order_id = @order_id::int;

-- name: GetOrdersByCustomer :many
SELECT o.order_id, o.order_date, o.total_amount, o.status
FROM orders o
WHERE o.customer_id = @customer_id::int;
