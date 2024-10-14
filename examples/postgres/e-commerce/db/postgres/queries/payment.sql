
-- name: CreatePayment :exec
INSERT INTO payment (order_id, payment_method_id, amount)
VALUES (@order_id::int, @payment_method_id::int, @amount::float);

-- name: GetPayments :many
SELECT p.payment_id, pm.name AS payment_method, p.amount, p.payment_date
FROM payment p
         INNER JOIN payment_method pm ON p.payment_method_id = pm.payment_method_id
WHERE p.order_id = @order_id::int;