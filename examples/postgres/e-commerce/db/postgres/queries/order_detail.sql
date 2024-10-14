-- name: CreateOrderDetail :exec
INSERT INTO order_detail (order_id, product_id, quantity, price)
VALUES (@order_id::int, @product_id::int, @quantity::int, @price::float);

-- name: GetOrderDetails :many
SELECT od.order_detail_id, p.name AS product_name, od.quantity, od.price
FROM order_detail od
         INNER JOIN product p ON od.product_id = p.product_id
WHERE od.order_id = @order_id::int;

-- name: GetOrderTotalAmount :one
SELECT SUM(quantity * price) AS total_amount
FROM order_detail
WHERE order_id = @order_id::int;
