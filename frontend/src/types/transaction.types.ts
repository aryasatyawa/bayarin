export interface TopupRequest {
    amount: number;
    channel_code: string;
    idempotency_key: string;
}

export interface TransferRequest {
    to_user_id: string;
    amount: number;
    description: string;
    pin: string;
    idempotency_key: string;
}

export interface TransactionResponse {
    transaction_id: string;
    type: 'topup' | 'transfer' | 'payment' | 'withdrawal';
    amount: number;
    amount_idr: string;
    status: 'pending' | 'success' | 'failed' | 'reversed';
    description: string;
    created_at: string;
}

export interface TransactionDetail {
    id: string;
    type: 'topup' | 'transfer' | 'payment' | 'withdrawal';
    amount: number;
    amount_idr: string;
    status: 'pending' | 'success' | 'failed' | 'reversed';
    from_wallet_id?: string;
    to_wallet_id?: string;
    description: string;
    reference_id?: string;
    created_at: string;
    completed_at?: string;
}

export interface TopupChannel {
    id: string;
    channel_code: string;
    channel_name: string;
    channel_type: string;
    fee_amount: number;
    min_amount: number;
    max_amount: number;
    is_active: boolean;
}