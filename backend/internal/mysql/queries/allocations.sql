-- name: CreatePaymentAllocation :execresult
INSERT INTO payment_allocations (non_posted_id, loan_id, installment_id, amount)
VALUES (sqlc.arg("non_posted_id"), sqlc.narg("loan_id"), sqlc.narg("installment_id"), sqlc.arg("amount"));

-- name: ListPaymentAllocationsByNonPostedId :many
SELECT * FROM payment_allocations WHERE non_posted_id = sqlc.arg("non_posted_id") AND deleted_at IS NULL;

-- name: ListPaymentAllocationsByLoanId :many
SELECT * FROM payment_allocations WHERE loan_id = sqlc.arg("loan_id") AND deleted_at IS NULL;

-- name: DeletePaymentAllocation :execresult
UPDATE payment_allocations SET deleted_at = NOW() WHERE id = sqlc.arg("id");

-- name: DeletePaymentAllocationsByNonPostedId :execresult
UPDATE payment_allocations SET deleted_at = NOW() WHERE non_posted_id = sqlc.arg("non_posted_id") AND deleted_at IS NULL;
