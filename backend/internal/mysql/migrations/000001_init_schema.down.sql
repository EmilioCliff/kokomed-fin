-- Drop foreign keys for `clients`
ALTER TABLE clients DROP FOREIGN KEY fk_branch_id;
ALTER TABLE clients DROP FOREIGN KEY fk_assigned_staff;
ALTER TABLE clients DROP FOREIGN KEY fk_updated_by;
ALTER TABLE clients DROP FOREIGN KEY fk_created_by;

-- Drop foreign keys for `users`
ALTER TABLE users DROP FOREIGN KEY fk_branch_id;
ALTER TABLE users DROP FOREIGN KEY fk_updated_by;
ALTER TABLE users DROP FOREIGN KEY fk_created_by;

-- Drop foreign keys for `products`
ALTER TABLE products DROP FOREIGN KEY fk_branch_id;
ALTER TABLE products DROP FOREIGN KEY fk_updated_by;
ALTER TABLE products DROP FOREIGN KEY fk_created_by;

-- Drop foreign keys for `loans`
ALTER TABLE loans DROP FOREIGN KEY fk_product_id;
ALTER TABLE loans DROP FOREIGN KEY fk_client_id;
ALTER TABLE loans DROP FOREIGN KEY fk_loan_officer;
ALTER TABLE loans DROP FOREIGN KEY fk_approved_by;
ALTER TABLE loans DROP FOREIGN KEY fk_disbursed_by;
ALTER TABLE loans DROP FOREIGN KEY fk_updated_by;

-- Drop foreign keys for `installments`
ALTER TABLE installments DROP FOREIGN KEY fk_loan_id;

-- Drop foreign keys for `non_posted`
ALTER TABLE non_posted DROP FOREIGN KEY fk_assigned_to;

-- Drop tables in the correct order to handle dependencies
DROP TABLE IF EXISTS installments;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS non_posted;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS branches;
