export interface User {
    id: string;
    email: string;
    phone: string;
    full_name: string;
    status: string;
    has_pin: boolean;
    created_at: string;
}

export interface Wallet {
    wallet_id: string;
    wallet_type: 'main' | 'bonus' | 'cashback';
    balance: number;
    balance_idr: string;
    currency: string;
    status: string;
}

export interface UserProfile {
    id: string;
    email: string;
    phone: string;
    full_name: string;
    status: string;
    has_pin: boolean;
    wallets: Wallet[];
    created_at: string;
}

export interface RegisterRequest {
    email: string;
    phone: string;
    full_name: string;
    password: string;
}

export interface RegisterResponse {
    user_id: string;
    email: string;
    phone: string;
    token: string;
}

export interface LoginRequest {
    identifier: string;
    password: string;
}

export interface LoginResponse {
    user_id: string;
    email: string;
    token: string;
}

export interface SetPINRequest {
    pin: string;
}