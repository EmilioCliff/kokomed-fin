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

-- name: GetClientsNonPosted :many
SELECT 
    np.id,
    np.transaction_source,
    np.transaction_number,
    np.account_number,
    np.phone_number,
    np.paying_name,
    np.amount,
    np.paid_date,
    np.assign_to,
    np.assigned_by
FROM non_posted np
WHERE 
    (np.assign_to = sqlc.narg('assign_to'))
    OR (np.account_number = sqlc.narg('account_number'))
ORDER BY np.paid_date DESC
LIMIT ? OFFSET ?;

-- name: CountClientsNonPosted :one
SELECT COUNT(*) AS total_non_posted 
FROM non_posted
WHERE 
    (
        (assign_to = sqlc.narg('assign_to'))
        OR (account_number = sqlc.narg('account_number'))
    );

-- name: GetTotalPaidByIDorAccountNo :one
SELECT SUM(amount) 
    FROM non_posted
    WHERE 
        (assign_to = sqlc.narg('assign_to'))
        OR (account_number = sqlc.narg('account_number'));

-- name: ListNonPostedByCategory :many
SELECT 
    np.*, 
    -- Client Details (if assigned)
    c.id AS client_id,
    c.full_name AS client_name,
    c.phone_number AS client_phone,
    c.overpayment AS client_overpayment,
    b.name AS client_branch_name

FROM non_posted np
LEFT JOIN clients c ON np.assign_to = c.id
LEFT JOIN branches b ON c.branch_id = b.id

WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(np.paying_name) LIKE ?
        OR LOWER(np.account_number) LIKE ?
        OR LOWER(np.transaction_number) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(np.transaction_source, ?) > 0
    )
ORDER BY np.paid_date DESC
LIMIT ? OFFSET ?;
-- SELECT *
-- FROM non_posted
-- WHERE 
--     (
--         COALESCE(?, '') = '' 
--         OR LOWER(paying_name) LIKE ?
--         OR LOWER(account_number) LIKE ?
--         OR LOWER(transaction_number) LIKE ?
--     )
--     AND (
--         COALESCE(?, '') = '' 
--         OR FIND_IN_SET(transaction_source, ?) > 0
--     )
--  ORDER BY paid_date DESC
-- LIMIT ? OFFSET ?;


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