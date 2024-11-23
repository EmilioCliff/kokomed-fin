-- name: CreateInstallment :execresult
INSERT INTO installments (loan_id, installment_number, amount_due, remaining_amount, due_date) 
VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetInstallment :one
SELECT * FROM installments WHERE id = ? LIMIT 1;

-- name: ListInstallmentsByLoan :many
SELECT * FROM installments WHERE loan_id = ? ORDER BY due_date ASC LIMIT ? OFFSET ?;

-- name: UpdateInstallment :execresult
UPDATE installments 
    SET remaining_amount =  sqlc.arg("remaining_amount"),
    paid =  coalesce(sqlc.narg("paid"), paid),
    paid_at =  coalesce(sqlc.narg("paid_at"), paid_at)
WHERE id = sqlc.arg("id");
