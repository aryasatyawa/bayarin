-- ============================================
-- ADMIN SYSTEM SCHEMA
-- Version: 2.0
-- ============================================

-- ============================================
-- ENUM TYPES
-- ============================================
CREATE TYPE admin_role AS ENUM ('super_admin', 'ops_admin', 'finance_admin');

CREATE TYPE admin_status AS ENUM ('active', 'suspended', 'inactive');

CREATE TYPE audit_action AS ENUM (
    'login', 'logout', 
    'view_user', 'freeze_wallet', 'unfreeze_wallet',
    'view_transaction', 'refund_transaction', 'reverse_transaction',
    'view_ledger', 'daily_settlement',
    'create_qr', 'update_qr', 'delete_qr'
);

-- ============================================
-- TABLE: admins
-- Deskripsi: Admin user accounts
-- ============================================
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role admin_role NOT NULL DEFAULT 'ops_admin',
    status admin_status DEFAULT 'active',
    last_login_at TIMESTAMP,
    created_by UUID REFERENCES admins (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admins_username ON admins (username);

CREATE INDEX idx_admins_email ON admins (email);

CREATE INDEX idx_admins_role ON admins (role);

CREATE INDEX idx_admins_status ON admins (status);

-- ============================================
-- TABLE: audit_logs
-- Deskripsi: IMMUTABLE audit trail (write-only)
-- CRITICAL: Tidak boleh UPDATE atau DELETE
-- ============================================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    admin_id UUID NOT NULL REFERENCES admins (id),
    action audit_action NOT NULL,
    resource_type VARCHAR(50), -- user, wallet, transaction, ledger, qr
    resource_id UUID,
    description TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    before_value JSONB, -- Snapshot sebelum action
    after_value JSONB, -- Snapshot setelah action
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_admin_id ON audit_logs (admin_id);

CREATE INDEX idx_audit_action ON audit_logs (action);

CREATE INDEX idx_audit_resource ON audit_logs (resource_type, resource_id);

CREATE INDEX idx_audit_created_at ON audit_logs (created_at DESC);

-- Prevent UPDATE & DELETE on audit_logs
CREATE RULE audit_logs_no_update AS ON UPDATE TO audit_logs DO INSTEAD NOTHING;

CREATE RULE audit_logs_no_delete AS ON DELETE TO audit_logs DO INSTEAD NOTHING;

-- ============================================
-- TABLE: qr_static_codes
-- Deskripsi: QR code static untuk merchant
-- ============================================
CREATE TABLE qr_static_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    qr_code VARCHAR(255) UNIQUE NOT NULL, -- QR string identifier
    merchant_name VARCHAR(255) NOT NULL,
    merchant_wallet_id UUID NOT NULL REFERENCES wallets (id),
    amount BIGINT, -- NULL = dynamic amount, value = fixed amount
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    scan_count BIGINT DEFAULT 0,
    created_by UUID REFERENCES admins (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP
);

CREATE INDEX idx_qr_code ON qr_static_codes (qr_code);

CREATE INDEX idx_qr_merchant_wallet ON qr_static_codes (merchant_wallet_id);

CREATE INDEX idx_qr_active ON qr_static_codes (is_active);

-- ============================================
-- TABLE: settlements
-- Deskripsi: Daily settlement records (sudah ada di migration 001, tapi update)
-- ============================================
-- Update existing settlements table
ALTER TABLE settlements
ADD COLUMN IF NOT EXISTS settled_by UUID REFERENCES admins (id);

ALTER TABLE settlements ADD COLUMN IF NOT EXISTS notes TEXT;

-- ============================================
-- TABLE: idempotency_monitor
-- Deskripsi: Monitor idempotency key untuk retry
-- ============================================
CREATE TABLE idempotency_monitor (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    transaction_id UUID REFERENCES transactions (id),
    status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, failed
    retry_count INT DEFAULT 0,
    last_retry_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_idem_key ON idempotency_monitor (idempotency_key);

CREATE INDEX idx_idem_status ON idempotency_monitor (status);

CREATE INDEX idx_idem_created_at ON idempotency_monitor (created_at DESC);

-- ============================================
-- TABLE: admin_sessions
-- Deskripsi: Admin session tracking (untuk Redis fallback)
-- ============================================
CREATE TABLE admin_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    admin_id UUID NOT NULL REFERENCES admins (id),
    session_token VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_session_token ON admin_sessions (session_token);

CREATE INDEX idx_admin_session_admin_id ON admin_sessions (admin_id);

CREATE INDEX idx_admin_session_expires ON admin_sessions (expires_at);

-- ============================================
-- SEED DATA: Super Admin
-- ============================================
-- Password: admin123 (hash via bcrypt)
-- IMPORTANT: Ganti password setelah first login!
INSERT INTO
    admins (
        id,
        username,
        email,
        password_hash,
        full_name,
        role,
        status
    )
VALUES (
        uuid_generate_v4 (),
        'superadmin',
        'admin@bayarin.com',
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIq7E0u4LW', -- admin123
        'Super Administrator',
        'super_admin',
        'active'
    );

-- ============================================
-- FUNCTIONS: Audit Log Helper
-- ============================================
CREATE OR REPLACE FUNCTION log_admin_action(
    p_admin_id UUID,
    p_action audit_action,
    p_resource_type VARCHAR(50),
    p_resource_id UUID,
    p_description TEXT,
    p_ip_address VARCHAR(45),
    p_before_value JSONB DEFAULT NULL,
    p_after_value JSONB DEFAULT NULL
) RETURNS UUID AS $$
DECLARE
    v_log_id UUID;
BEGIN
    INSERT INTO audit_logs (
        admin_id, action, resource_type, resource_id,
        description, ip_address, before_value, after_value
    ) VALUES (
        p_admin_id, p_action, p_resource_type, p_resource_id,
        p_description, p_ip_address, p_before_value, p_after_value
    ) RETURNING id INTO v_log_id;
    
    RETURN v_log_id;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- VIEWS: Dashboard Metrics
-- ============================================

-- Total sistem liability (total saldo user)
CREATE OR REPLACE VIEW v_system_liability AS
SELECT
    SUM(balance) as total_liability,
    COUNT(DISTINCT user_id) as total_users,
    COUNT(*) as total_wallets
FROM wallets
WHERE
    status = 'active';

-- Daily transaction summary
CREATE OR REPLACE VIEW v_daily_transaction_summary AS
SELECT
    DATE (created_at) as transaction_date,
    transaction_type,
    status,
    COUNT(*) as transaction_count,
    SUM(amount) as total_amount
FROM transactions
GROUP BY
    DATE (created_at),
    transaction_type,
    status
ORDER BY transaction_date DESC;

-- Pending transactions
CREATE OR REPLACE VIEW v_pending_transactions AS
SELECT
    id,
    user_id,
    transaction_type,
    amount,
    created_at,
    idempotency_key
FROM transactions
WHERE
    status = 'pending'
ORDER BY created_at DESC;

-- Failed transactions (last 7 days)
CREATE OR REPLACE VIEW v_failed_transactions AS
SELECT
    id,
    user_id,
    transaction_type,
    amount,
    created_at,
    idempotency_key,
    description
FROM transactions
WHERE
    status = 'failed'
    AND created_at >= NOW() - INTERVAL '7 days'
ORDER BY created_at DESC;