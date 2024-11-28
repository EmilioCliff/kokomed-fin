-- name: CreateNonPosted :execresult
INSERT INTO non_posted (transaction_source, transaction_number, account_number, phone_number, paying_name, amount, paid_date, assign_to) 
VALUES (
    sqlc.arg("transaction_source"),
    sqlc.arg("transaction_number"),
    sqlc.arg("account_number"),
    sqlc.arg("phone_number"),
    sqlc.arg("paying_name"),
    sqlc.arg("amount"),
    sqlc.arg("paid_date"),
    sqlc.narg("assign_to")
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
    transaction_source = sqlc.arg("transaction_source")
WHERE id = sqlc.arg("id");

-- name: DeleteNonPosted :exec
DELETE FROM non_posted WHERE id = ?;