-- name: CreateProduct :execresult
INSERT INTO products (branch_id, loan_amount, repay_amount, interest_amount, updated_by) 
VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetProduct :one
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products LIMIT ? OFFSET ?;

-- name: UpdateProduct :execresult
UPDATE products 
    SET loan_amount = coalesce(sqlc.narg("loan_amount"), loan_amount),
    repay_amount = coalesce(sqlc.narg("repay_amount"), repay_amount),
    interest_amount = coalesce(sqlc.narg("interest_amount"), interest_amount),
    updated_by = sqlc.arg("updated_by")
WHERE id = sqlc.arg("id");

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;

-- name: ListProductsByBranch :many
SELECT * FROM products WHERE branch_id = ? LIMIT ? OFFSET ?;