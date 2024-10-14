-- name: CreateProduct :exec
INSERT INTO products (name, price)
VALUES (@name::text, @price::float);

-- name: GetProducts :many
SELECT id, coalesce(name, ''), coalesce(price, 0)::float as price, (EXTRACT(epoch FROM created_at) * 1000)::bigint as created_at
FROM products
WHERE case when @name::varchar = '' then true else name ILIKE @name::varchar end
  AND case when @price::float = 0.0 then true else price = @price::float end
ORDER BY @sort_by
LIMIT @si OFFSET @offs;

-- name: CountProducts :one
SELECT COUNT(*)
FROM products
WHERE case when @name::varchar = '' then true else name ILIKE @name::varchar end
  AND case when @price::float = 0 then true else price = @price::float end;

-- name: GetProductsCursor :many
SELECT id, coalesce(name, ''), coalesce(price, 0)::float as price, (EXTRACT(epoch FROM created_at) * 1000)::bigint as created_at
FROM products
WHERE case when @name::varchar = '' then true else name ILIKE @name::varchar end
  AND case when @price::float = 0 then true else price = @price::float end
ORDER BY id
LIMIT @si::int;

-- name: GetProductsNextCursor :many
SELECT id, coalesce(name, ''), coalesce(price, 0)::float as price, (EXTRACT(epoch FROM created_at) * 1000)::bigint as created_at
FROM products
WHERE id > @cursor
  AND case when @name::varchar = '' then true else name ILIKE @name::varchar end
  AND case when @price::float = 0 then true else price = @price::float end
ORDER BY id
LIMIT @si::int;

-- name: GetProductsPrevCursor :many
SELECT id, coalesce(name, ''), coalesce(price, 0)::float as price, (EXTRACT(epoch FROM created_at) * 1000)::bigint as created_at
FROM products
WHERE id < @cursor
  AND case when @name::varchar = '' then true else name ILIKE @name::varchar end
  AND case when @price::float = 0 then true else price = @price::float end
ORDER BY id desc
LIMIT @si::int;

-- name: UpdateProduct :exec
UPDATE products
SET price = @price::float,
    name  = @name::text
WHERE id = @id::int;

-- name: DeleteProduct :exec
DELETE
FROM products
WHERE id = $1;