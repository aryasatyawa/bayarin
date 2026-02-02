-- Drop views
DROP VIEW IF EXISTS v_failed_transactions;

DROP VIEW IF EXISTS v_pending_transactions;

DROP VIEW IF EXISTS v_daily_transaction_summary;

DROP VIEW IF EXISTS v_system_liability;

-- Drop function
DROP FUNCTION IF EXISTS log_admin_action;

-- Drop tables (reverse order)
DROP TABLE IF EXISTS admin_sessions;

DROP TABLE IF EXISTS idempotency_monitor;

DROP TABLE IF EXISTS qr_static_codes;

DROP TABLE IF EXISTS audit_logs;

DROP TABLE IF EXISTS admins;

-- Drop types
DROP TYPE IF EXISTS audit_action;

DROP TYPE IF EXISTS admin_status;

DROP TYPE IF EXISTS admin_role;