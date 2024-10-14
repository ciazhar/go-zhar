-- name: GetProductsByCategory :many
SELECT p.product_id, p.name, p.description, p.price, c.name AS category_name
FROM product p
         INNER JOIN category c ON p.category_id = c.category_id
WHERE c.name = @category_name::varchar;
