-- name: CreateNonPosted :execresult
INSERT INTO non_posted (transaction_number, account_number, phone_number, paying_name, amount, paid_date) 
VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: ListAllNonPosted :many
SELECT * FROM non_posted LIMIT ? OFFSET ?;

-- name: ListUnassignedNonPosted :many
SELECT * FROM non_posted WHERE assign_to IS NULL LIMIT ? OFFSET ?;

-- name: GetNonPosted :one
SELECT * FROM non_posted WHERE id = ? LIMIT 1;

-- name: AssignNonPosted :execresult
UPDATE non_posted 
    SET assign_to = sqlc.arg("assign_to")
WHERE id = sqlc.arg("id");

-- name: DeleteNonPosted :exec
DELETE FROM non_posted WHERE id = ?;