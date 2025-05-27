ALTER TABLE non_posted DROP COLUMN deleted_at;
ALTER TABLE client_overpayment_transactions DROP FOREIGN KEY fk_cot_client_id;
ALTER TABLE payment_allocations DROP FOREIGN KEY fk_payment_allocations_non_posted_id;

DROP TABLE IF EXISTS client_overpayment_transactions;
DROP TABLE IF EXISTS payment_allocations;