-- Drop tables in reverse order (respect foreign keys)
DROP TABLE IF EXISTS settlements;

DROP TABLE IF EXISTS payment_methods;

DROP TABLE IF EXISTS topup_channels;

DROP TABLE IF EXISTS ledger_entries;

DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS wallets;

DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";