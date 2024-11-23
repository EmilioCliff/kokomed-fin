CREATE TABLE `branches` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL
);

CREATE TABLE `users` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `full_name` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `refresh_token` VARCHAR(255) NOT NULL,
  `role` ENUM('ADMIN', 'AGENT') NOT NULL,
  `branch_id` INT NOT NULL,
  `updated_by` INT NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  `created_by` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_users_branch_id FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`),
  CONSTRAINT fk_users_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`),
  CONSTRAINT fk_users_created_by FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
);

CREATE TABLE `clients` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `full_name` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(255) NOT NULL,
  `id_number` VARCHAR(255) NULL,
  `dob` DATE NULL,
  `gender` ENUM('MALE', 'FEMALE') NOT NULL,
  `active` BOOLEAN NOT NULL DEFAULT FALSE,
  `branch_id` INT NOT NULL,
  `assigned_staff` INT NOT NULL,
  `overpayment` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `updated_by` INT NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  `created_by` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_clients_branch_id FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`),
  CONSTRAINT fk_clients_assigned_staff FOREIGN KEY (`assigned_staff`) REFERENCES `users` (`id`),
  CONSTRAINT fk_clients_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`),
  CONSTRAINT fk_clients_created_by FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
);

CREATE TABLE `products` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `branch_id` INT NOT NULL,
  `loan_amount` DECIMAL(10, 2) NOT NULL,
  `repay_amount` DECIMAL(10, 2) NOT NULL,
  `interest_amount` DECIMAL(10, 2) NOT NULL,
  `updated_by` INT NOT NULL,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_products_branch_id FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`),
  CONSTRAINT fk_products_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`),
  CONSTRAINT fk_products_created_by FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
);

CREATE TABLE `loans` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `product_id` INT NOT NULL,
  `client_id` INT NOT NULL,
  `loan_officer` INT NOT NULL,
  `loan_purpose` TEXT NULL,
  `due_date` DATE NULL,
  `approved_by` INT NOT NULL,
  `disbursed_on` DATE NULL,
  `disbursed_by` INT NULL,
  `total_installments` INT NOT NULL DEFAULT 4,
  `installments_period` INT NOT NULL DEFAULT 7,
  `status` ENUM('INACTIVE', 'ACTIVE', 'COMPLETED', 'DEFAULTED') NOT NULL,
  `processing_fee` DECIMAL(10,2) NOT NULL DEFAULT 400.00,
  `paid_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `updated_by` INT NULL,
  `created_by` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_loans_product_id FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT fk_loans_client_id FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`),
  CONSTRAINT fk_loans_loan_officer FOREIGN KEY (`loan_officer`) REFERENCES `users` (`id`),
  CONSTRAINT fk_loans_approved_by FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`),
  CONSTRAINT fk_loans_disbursed_by FOREIGN KEY (`disbursed_by`) REFERENCES `users` (`id`),
  CONSTRAINT fk_loans_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`),
  CONSTRAINT fk_loans_created_by FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
);

CREATE TABLE `installments` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `loan_id` INT NOT NULL,
  `installment_number` INT NOT NULL,
  `amount_due` DECIMAL(10,2) NOT NULL,
  `remaining_amount` DECIMAL(10,2) NOT NULL,
  `paid` BOOLEAN NOT NULL DEFAULT FALSE,
  `paid_at` TIMESTAMP NULL,
  `due_date` DATE NOT NULL,

  CONSTRAINT fk_installments_loan_id FOREIGN KEY (`loan_id`) REFERENCES `loans` (`id`)
);

CREATE TABLE `non_posted` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `transaction_number` VARCHAR(255) NOT NULL,
  `account_number` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(255) NOT NULL,
  `paying_name` VARCHAR(255) NOT NULL,
  `amount` DECIMAL(10,2) NOT NULL,
  `paid_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `assign_to` INT NULL,

  CONSTRAINT fk_non_posted_assign_to FOREIGN KEY (`assign_to`) REFERENCES `clients` (`id`)
);

-- Indexes
CREATE INDEX idx_clients_phone_number ON `clients` (`phone_number`);
CREATE INDEX idx_clients_active ON `clients` (`active`);
CREATE INDEX idx_users_email ON `users` (`email`);
CREATE INDEX idx_users_branch_id ON `users` (`branch_id`);
CREATE INDEX idx_loans_client_id ON `loans` (`client_id`);
CREATE INDEX idx_loans_loan_officer ON `loans` (`loan_officer`);
CREATE INDEX idx_loans_status ON `loans` (`status`);
CREATE INDEX idx_installments_loan_id ON `installments` (`loan_id`);
CREATE INDEX idx_installments_paid ON `installments` (`paid`);
