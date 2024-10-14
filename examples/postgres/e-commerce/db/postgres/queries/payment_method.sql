-- name: GetPaymentMethods :many
SELECT payment_method_id, name, description
FROM payment_method;
