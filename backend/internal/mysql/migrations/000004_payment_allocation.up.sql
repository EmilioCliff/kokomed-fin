CREATE TABLE `payment_allocations` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `non_posted_id` INT NOT NULL,
  `loan_id` INT NULL,
  `installment_id` INT NULL,
  `amount` DECIMAL(10, 2) NOT NULL,
  `description` TEXT NOT NULL,
  `deleted_at` TIMESTAMP NULL,
  `deleted_description` TEXT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_payment_allocations_non_posted_id FOREIGN KEY (`non_posted_id`) REFERENCES `non_posted` (`id`)
);

CREATE TABLE client_overpayment_transactions (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `client_id` INT NOT NULL,
  `payment_id` INT NULL, 
  `amount` DECIMAL(10,2) NOT NULL,
  `description` TEXT NOT NULL,
  `created_by` VARCHAR(255) NOT NULL DEFAULT 'SYSTEM',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_cot_client_id FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`)
);

ALTER TABLE `non_posted` ADD COLUMN `deleted_at` TIMESTAMP DEFAULT NULL;
ALTER TABLE `non_posted` ADD COLUMN `deleted_description` TEXT DEFAULT NULL;

CREATE INDEX idx_non_posted_deleted_at ON `non_posted` (`deleted_at`);

CREATE INDEX idx_payment_allocations_loan_id ON `payment_allocations` (`loan_id`);
CREATE INDEX idx_payment_allocations_non_posted_id ON `payment_allocations` (`non_posted_id`);
CREATE INDEX idx_payment_allocations_deleted_at ON `payment_allocations` (`deleted_at`);

CREATE INDEX idx_client_overpayment_transactions_client_id ON `client_overpayment_transactions` (`client_id`);
CREATE INDEX idx_client_overpayment_transactions_payment_id ON `client_overpayment_transactions` (`payment_id`);