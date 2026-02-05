export interface LedgerEntry {
    id: string;
    transaction_id: string;
    wallet_id: string;
    user_id: string;
    entry_type: 'debit' | 'credit';
    amount: number;
    balance_before: number;
    balance_after: number;
    description: string;
    transaction_type: string;
    created_at: string;
}

export interface LedgerFilter {
    user_id?: string;
    wallet_id?: string;
    transaction_id?: string;
    entry_type?: 'debit' | 'credit';
    start_date?: string;
    end_date?: string;
    limit: number;
    offset: number;
}

export interface TransactionDetail {
    id: string;
    idempotency_key: string;
    user_id: string;
    user_email: string;
    transaction_type: string;
    amount: number;
    currency: string;
    status: string;
    from_wallet_id?: string;
    to_wallet_id?: string;
    reference_id?: string;
    description: string;
    ledger_entries?: LedgerEntry[];
    created_at: string;
    updated_at: string;
    completed_at?: string;
}

export interface UserInspectorDetail {
    user: {
        id: string;
        email: string;
        phone: string;
        full_name: string;
        status: string;
    };
    wallets: WalletDetail[];
    total_balance: number;
    total_transactions: number;
    success_transactions: number;
    failed_transactions: number;
    last_transaction_at?: string;
}

export interface WalletDetail {
    id: string;
    wallet_type: string;
    balance: number;
    status: string;
    transaction_count: number;
    last_activity_at?: string;
    created_at: string;
}