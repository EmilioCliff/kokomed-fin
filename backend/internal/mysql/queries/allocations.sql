-- name: CreatePaymentAllocation :execresult
INSERT INTO payment_allocations (non_posted_id, loan_id, installment_id, amount, description)
VALUES (sqlc.arg("non_posted_id"), sqlc.narg("loan_id"), sqlc.narg("installment_id"), sqlc.arg("amount"), sqlc.arg("description"));

-- name: ListPaymentAllocationsByNonPostedId :many
SELECT * FROM payment_allocations WHERE non_posted_id = sqlc.arg("non_posted_id") AND deleted_at IS NULL;

-- name: ListPaymentAllocationsByLoanId :many
SELECT 
  pa.*,
  np.transaction_source,
  np.transaction_number,
  np.account_number,
  np.paying_name,
  np.amount,
  np.paid_date
FROM payment_allocations pa
JOIN non_posted np ON pa.non_posted_id = np.id
WHERE pa.loan_id = sqlc.arg("loan_id")
AND pa.deleted_at IS NULL;

-- name: ListPaymentAllocationsByNonPostedID :many
SELECT 
  pa.*,
  np.transaction_source,
  np.transaction_number,
  np.account_number,
  np.paying_name,
  np.amount,
  np.paid_date
FROM payment_allocations pa
JOIN non_posted np ON pa.non_posted_id = np.id
WHERE loan_id IS NULL
  AND non_posted_id = sqlc.arg("non_posted_id") AND pa.deleted_at IS NULL;



-- name: DeletePaymentAllocation :execresult
UPDATE payment_allocations SET deleted_at = NOW() WHERE id = sqlc.arg("id");

-- name: DeletePaymentAllocationsByNonPostedId :execresult
UPDATE payment_allocations SET deleted_at = NOW(), deleted_description = sqlc.arg("deleted_description") WHERE non_posted_id = sqlc.arg("non_posted_id") AND deleted_at IS NULL;
