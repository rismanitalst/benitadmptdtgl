-- ============================================
-- Database: benitadmptdgtl
-- Import ini ke phpMyAdmin atau jalankan di MySQL
-- ============================================

USE `benitadmptdgtl`;

-- Tabel Users
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  `firebase_uid` VARCHAR(191) NOT NULL,
  `email` VARCHAR(191) NOT NULL,
  `name` VARCHAR(255) DEFAULT NULL,
  `role` VARCHAR(50) DEFAULT 'user',
  `email_verified` TINYINT(1) DEFAULT 0,
  `fcm_token` TEXT,
  `totp_secret` TEXT,
  `totp_enabled` TINYINT(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_firebase_uid` (`firebase_uid`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel OTPs
CREATE TABLE IF NOT EXISTS `otps` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `code` VARCHAR(10) NOT NULL,
  `type` VARCHAR(20) NOT NULL,
  `expires_at` DATETIME(3) DEFAULT NULL,
  `used` TINYINT(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_otps_user_id` (`user_id`),
  KEY `idx_otps_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_otps_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel Accounts
CREATE TABLE IF NOT EXISTS `accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `balance` DECIMAL(15,2) DEFAULT 0.00,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_accounts_user_id` (`user_id`),
  KEY `idx_accounts_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_accounts_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel Transactions
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `amount` DECIMAL(15,2) DEFAULT NULL,
  `type` VARCHAR(10) DEFAULT NULL,
  `description` VARCHAR(255) DEFAULT NULL,
  `balance_before` DECIMAL(15,2) DEFAULT NULL,
  `balance_after` DECIMAL(15,2) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transactions_account_id` (`account_id`),
  KEY `idx_transactions_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_transactions_account` FOREIGN KEY (`account_id`) REFERENCES `accounts`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
