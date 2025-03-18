-- name: CreateProduct :execresult
INSERT INTO products (branch_id, loan_amount, repay_amount, interest_amount, updated_by) 
VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetProduct :one
SELECT 
    p.*, 
    b.name AS branch_name 
FROM products p
JOIN branches b ON p.branch_id = b.id
WHERE p.id = ? 
LIMIT 1;
-- SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: GetProductRepayAmount :one
SELECT repay_amount FROM products WHERE id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products LIMIT ? OFFSET ?;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;

-- name: ListProductsByBranch :many
SELECT * FROM products WHERE branch_id = ? LIMIT ? OFFSET ?;

-- name: HelperProduct :many
SELECT p.id AS productId,
    p.loan_amount AS loanAmount,
    b.name AS branchNAme
FROM products p JOIN branches b ON p.branch_id = b.id;

-- name: ListProductsByCategory :many
SELECT 
    p.*, 
    b.name AS branch_name
FROM products p
JOIN branches b ON p.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(b.name) LIKE ?
    )
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountLoansByCategory :one
SELECT COUNT(*) AS total_products
FROM products p
JOIN branches b ON p.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(b.name) LIKE ?
    )