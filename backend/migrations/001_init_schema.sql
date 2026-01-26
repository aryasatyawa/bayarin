-- ============================================
-- BAYARIN - Digital Wallet & Payment Gateway
-- Schema Version: 1.0
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- TABLE: users
-- Deskripsi: Menyimpan data pengguna
-- ============================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    pin_hash VARCHAR(255), -- PIN untuk transaksi
    status VARCHAR(20) DEFAULT 'active', -- active, suspended, blocked
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_phone ON users (phone);

CREATE INDEX idx_users_status ON users (status);

-- ============================================
-- TABLE: wallets
-- Deskripsi: Saldo wallet per user
-- PENTING: balance dalam INTEGER (minor unit)
-- Contoh: Rp 100.000,00 = 10000000 (satuan sen)
-- ============================================
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    wallet_type VARCHAR(50) DEFAULT 'main', -- main, bonus, cashback
    balance BIGINT DEFAULT 0 CHECK (balance >= 0), -- WAJIB INTEGER
    currency VARCHAR(3) DEFAULT 'IDR',
    status VARCHAR(20) DEFAULT 'active', -- active, frozen, closed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, wallet_type)
);

CREATE INDEX idx_wallets_user_id ON wallets (user_id);

CREATE INDEX idx_wallets_status ON wallets (status);

-- ============================================
-- TABLE: transactions
-- Deskripsi: Semua transaksi (topup, transfer, payment)
-- PENTING: amount dalam INTEGER (minor unit)
-- ============================================
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL, -- Untuk idempotency
    user_id UUID NOT NULL REFERENCES users (id),
    transaction_type VARCHAR(50) NOT NULL, -- topup, transfer, payment, withdrawal
    amount BIGINT NOT NULL CHECK (amount > 0), -- WAJIB INTEGER
    currency VARCHAR(3) DEFAULT 'IDR',
    status VARCHAR(20) DEFAULT 'pending', -- pending, success, failed, reversed
    from_wallet_id UUID REFERENCES wallets (id),
    to_wallet_id UUID REFERENCES wallets (id),
    reference_id VARCHAR(255), -- External reference (bank, VA, dll)
    description TEXT,
    metadata JSONB, -- Data tambahan (fee, promo, dll)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_transactions_user_id ON transactions (user_id);

CREATE INDEX idx_transactions_status ON transactions (status);

CREATE INDEX idx_transactions_type ON transactions (transaction_type);

CREATE INDEX idx_transactions_idempotency ON transactions (idempotency_key);

CREATE INDEX idx_transactions_reference ON transactions (reference_id);

CREATE INDEX idx_transactions_created_at ON transactions (created_at DESC);

-- ============================================
-- TABLE: ledger_entries
-- Deskripsi: Double-entry bookkeeping
-- Setiap transaksi menghasilkan minimal 2 entries (debit & credit)
-- PENTING: amount dalam INTEGER
-- ============================================
CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    transaction_id UUID NOT NULL REFERENCES transactions (id),
    wallet_id UUID NOT NULL REFERENCES wallets (id),
    entry_type VARCHAR(10) NOT NULL CHECK (
        entry_type IN ('debit', 'credit')
    ),
    amount BIGINT NOT NULL CHECK (amount > 0), -- WAJIB INTEGER
    balance_before BIGINT NOT NULL, -- Snapshot balance sebelum transaksi
    balance_after BIGINT NOT NULL, -- Snapshot balance setelah transaksi
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ledger_transaction_id ON ledger_entries (transaction_id);

CREATE INDEX idx_ledger_wallet_id ON ledger_entries (wallet_id);

CREATE INDEX idx_ledger_created_at ON ledger_entries (created_at DESC);

-- ============================================
-- TABLE: topup_channels
-- Deskripsi: Channel untuk topup (VA Bank, E-Wallet, dll)
-- ============================================
CREATE TABLE topup_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    channel_code VARCHAR(50) UNIQUE NOT NULL, -- BCA_VA, MANDIRI_VA, OVO, GOPAY
    channel_name VARCHAR(255) NOT NULL,
    channel_type VARCHAR(50) NOT NULL, -- bank_transfer, ewallet, retail
    fee_type VARCHAR(20) DEFAULT 'fixed', -- fixed, percentage
    fee_amount BIGINT DEFAULT 0, -- WAJIB INTEGER
    min_amount BIGINT DEFAULT 1000000, -- Min Rp 10.000 = 1000000 sen
    max_amount BIGINT DEFAULT 1000000000, -- Max Rp 10jt = 1000000000 sen
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TABLE: payment_methods
-- Deskripsi: Metode pembayaran untuk merchant
-- ============================================
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    method_code VARCHAR(50) UNIQUE NOT NULL, -- wallet, qris, bank_transfer
    method_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TABLE: settlements
-- Deskripsi: Settlement record untuk reconciliation
-- ============================================
CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    settlement_date DATE NOT NULL,
    total_transactions INT DEFAULT 0,
    total_amount BIGINT DEFAULT 0, -- WAJIB INTEGER
    total_fee BIGINT DEFAULT 0,
    net_amount BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending', -- pending, completed, failed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_settlements_date ON settlements (settlement_date);

CREATE INDEX idx_settlements_status ON settlements (status);

-- ============================================
-- SEED DATA
-- ============================================

-- Topup Channels
INSERT INTO
    topup_channels (
        channel_code,
        channel_name,
        channel_type,
        fee_amount,
        min_amount,
        max_amount
    )
VALUES (
        'BCA_VA',
        'BCA Virtual Account',
        'bank_transfer',
        0,
        1000000,
        5000000000
    ),
    (
        'MANDIRI_VA',
        'Mandiri Virtual Account',
        'bank_transfer',
        0,
        1000000,
        5000000000
    ),
    (
        'BNI_VA',
        'BNI Virtual Account',
        'bank_transfer',
        0,
        1000000,
        5000000000
    ),
    (
        'BRI_VA',
        'BRI Virtual Account',
        'bank_transfer',
        0,
        1000000,
        5000000000
    ),
    (
        'OVO',
        'OVO E-Wallet',
        'ewallet',
        50000,
        1000000,
        1000000000
    ),
    (
        'GOPAY',
        'GoPay E-Wallet',
        'ewallet',
        50000,
        1000000,
        1000000000
    ),
    (
        'DANA',
        'DANA E-Wallet',
        'ewallet',
        50000,
        1000000,
        1000000000
    );

-- Payment Methods
INSERT INTO
    payment_methods (method_code, method_name)
VALUES ('wallet', 'Bayarin Wallet'),
    ('qris', 'QRIS'),
    (
        'bank_transfer',
        'Transfer Bank'
    );