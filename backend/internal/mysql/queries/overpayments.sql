-- name: CreateClientOverpaymentTransaction :execresult
INSERT INTO client_overpayment_transactions (client_id, payment_id, amount, created_by)
VALUES (sqlc.arg("client_id"), sqlc.narg("payment_id"), sqlc.arg("amount"), sqlc.arg("created_by"));

-- name: GetClientOverpaymentTransaction :one
SELECT * FROM client_overpayment_transactions WHERE id = sqlc.arg("id");

-- name: GetClientOverpaymentTransactions :many
SELECT * FROM client_overpayment_transactions WHERE client_id = sqlc.arg("client_id");

-- name: GetClientOverpaymentTransactionByPaymentId :one
SELECT * FROM client_overpayment_transactions WHERE payment_id = sqlc.arg("payment_id");