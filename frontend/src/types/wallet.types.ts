export interface WalletBalance {
    wallet_id: string;
    wallet_type: 'main' | 'bonus' | 'cashback';
    balance: number;
    balance_idr: string;
    currency: string;
    status: string;
}

export interface LedgerEntry {
    id: string;
    transaction_id: string;
    wallet_id: string;
    entry_type: 'debit' | 'credit';
    amount: number;
    balance_before: number;
    balance_after: number;
    description: string;
    created_at: string;
}

export interface WalletHistory {
    wallet_id: string;
    entries: LedgerEntry[];
    total: number;
    limit: number;
    offset: number;
}