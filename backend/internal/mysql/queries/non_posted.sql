-- name: CreateNonPosted :execresult
INSERT INTO non_posted (transaction_source, transaction_number, account_number, phone_number, paying_name, amount, paid_date, assign_to, assigned_by) 
VALUES (
    sqlc.arg("transaction_source"),
    sqlc.arg("transaction_number"),
    sqlc.arg("account_number"),
    sqlc.arg("phone_number"),
    sqlc.arg("paying_name"),
    sqlc.arg("amount"),
    sqlc.arg("paid_date"),
    sqlc.narg("assign_to"),
    sqlc.arg("assigned_by")
);

-- name: ListAllNonPosted :many
SELECT * FROM non_posted LIMIT ? OFFSET ?;

-- name: ListAllNonPostedByTransactionSource :many
SELECT * FROM non_posted WHERE transaction_source = ?;

-- name: ListNonPostedByTransactionSource :many
SELECT * FROM non_posted WHERE transaction_source = ? LIMIT ? OFFSET ?;

-- name: ListUnassignedNonPosted :many
SELECT * FROM non_posted WHERE assign_to IS NULL LIMIT ? OFFSET ?;

-- name: GetNonPosted :one
SELECT * FROM non_posted WHERE id = ? LIMIT 1;

-- name: AssignNonPosted :execresult
UPDATE non_posted 
    SET assign_to = sqlc.arg("assign_to"),
    transaction_source = sqlc.arg("transaction_source"),
    assigned_by = COALESCE(sqlc.narg("assigned_by"), assigned_by)
WHERE id = sqlc.arg("id");

-- name: DeleteNonPosted :exec
DELETE FROM non_posted WHERE id = ?;

-- name: ListNonPostedByCategory :many
SELECT *
FROM non_posted
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(paying_name) LIKE ?
        OR LOWER(account_number) LIKE ?
        OR LOWER(transaction_number) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(transaction_source, ?) > 0
    )
 ORDER BY paid_date DESC
LIMIT ? OFFSET ?;


-- name: CountNonPostedByCategory :one
SELECT COUNT(*) AS total_non_posted 
FROM non_posted
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(paying_name) LIKE ?
        OR LOWER(account_number) LIKE ?
        OR LOWER(transaction_number) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(transaction_source, ?) > 0
    );