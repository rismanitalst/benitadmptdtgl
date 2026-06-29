-- ============================================
-- E-Money 2FA — Database Migration
-- Untuk Laragon MySQL (phpMyAdmin)
-- ============================================

CREATE DATABASE IF NOT EXISTS emoney_2fa;
USE emoney_2fa;

-- ============================================
-- Users
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at    DATETIME(3)     NULL,
    updated_at    DATETIME(3)     NULL,
    deleted_at    DATETIME(3)     NULL,
    firebase_uid  VARCHAR(191)    NOT NULL,
    email         VARCHAR(191)    NOT NULL,
    name          VARCHAR(255)    NULL,
    role          VARCHAR(50)     NULL DEFAULT 'user',
    email_verified BOOLEAN        NULL DEFAULT FALSE,
    fcm_token     TEXT            NULL,
    totp_secret   TEXT            NULL,
    totp_enabled  BOOLEAN         NULL DEFAULT FALSE,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_users_firebase_uid (firebase_uid),
    UNIQUE INDEX idx_users_email (email),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- Accounts (Saldo)
-- ============================================
CREATE TABLE IF NOT EXISTS accounts (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at DATETIME(3)     NULL,
    updated_at DATETIME(3)     NULL,
    deleted_at DATETIME(3)     NULL,
    user_id    BIGINT UNSIGNED NOT NULL,
    balance    DECIMAL(15,2)   NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_accounts_user_id (user_id),
    INDEX idx_accounts_deleted_at (deleted_at),
    CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- Transactions
-- ============================================
CREATE TABLE IF NOT EXISTS transactions (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at     DATETIME(3)     NULL,
    updated_at     DATETIME(3)     NULL,
    deleted_at     DATETIME(3)     NULL,
    account_id     BIGINT UNSIGNED NOT NULL,
    amount         DECIMAL(15,2)   NULL,
    type           VARCHAR(10)     NULL,       -- "debit" | "credit"
    description    VARCHAR(255)    NULL,
    balance_before DECIMAL(15,2)   NULL,
    balance_after  DECIMAL(15,2)   NULL,
    PRIMARY KEY (id),
    INDEX idx_transactions_account_id (account_id),
    INDEX idx_transactions_deleted_at (deleted_at),
    CONSTRAINT fk_transactions_account FOREIGN KEY (account_id) REFERENCES accounts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- OTPs
-- ============================================
CREATE TABLE IF NOT EXISTS otps (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at DATETIME(3)     NULL,
    updated_at DATETIME(3)     NULL,
    deleted_at DATETIME(3)     NULL,
    user_id    BIGINT UNSIGNED NOT NULL,
    code       VARCHAR(10)     NOT NULL,
    type       VARCHAR(20)     NOT NULL,       -- "firebase" | "email"
    expires_at DATETIME(3)     NULL,
    used       BOOLEAN         NULL DEFAULT FALSE,
    PRIMARY KEY (id),
    INDEX idx_otps_user_id (user_id),
    INDEX idx_otps_deleted_at (deleted_at),
    CONSTRAINT fk_otps_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
