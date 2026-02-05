export interface DashboardOverview {
    total_users: number;
    total_active_wallets: number;
    total_system_liability: number;
    today_transactions: number;
    today_volume: number;
    pending_transactions: number;
    failed_transactions: number;
    today_topups: number;
    today_transfers: number;
}

export interface DailyStats {
    date: string;
    total_transactions: number;
    total_volume: number;
    by_type: Record<string, TypeStats>;
    by_status: Record<string, StatusStats>;
}

export interface TypeStats {
    count: number;
    volume: number;
}

export interface StatusStats {
    count: number;
}

export interface TransactionSummary {
    start_date: string;
    end_date: string;
    total_transactions: number;
    total_volume: number;
    by_type: Record<string, TypeStats>;
    by_status: Record<string, StatusStats>;
    daily_breakdown: DailyBreakdown[];
}

export interface DailyBreakdown {
    date: string;
    count: number;
    volume: number;
}