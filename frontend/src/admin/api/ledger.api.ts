import { ApiResponse } from '@/api/axios';
import { adminApiClient } from '@/admin/api/admin-api';
import { LedgerEntry, LedgerFilter } from '@/admin/types/ledger.types';

export const ledgerApi = {
    // Get ledger entries with filter
    getLedgerEntries: async (filter: Partial<LedgerFilter>): Promise<{
        entries: LedgerEntry[];
        total: number;
    }> => {
        const response = await adminApiClient.get<ApiResponse<{
            entries: LedgerEntry[];
            total: number;
        }>>('/admin/ledger', { params: filter });
        return response.data.data!;
    },

    // Get ledger by transaction
    getLedgerByTransaction: async (transactionId: string): Promise<LedgerEntry[]> => {
        const response = await adminApiClient.get<ApiResponse<LedgerEntry[]>>(
            `/admin/ledger/transaction/${transactionId}`
        );
        return response.data.data!;
    },

    // Get ledger by wallet
    getLedgerByWallet: async (
        walletId: string,
        limit: number = 20,
        offset: number = 0
    ): Promise<{
        wallet_id: string;
        current_balance: number;
        entries: LedgerEntry[];
        total: number;
    }> => {
        const response = await adminApiClient.get<ApiResponse<{
            wallet_id: string;
            current_balance: number;
            entries: LedgerEntry[];
            total: number;
        }>>(`/admin/ledger/wallet/${walletId}`, {
            params: { limit, offset },
        });
        return response.data.data!;
    },

    // Validate wallet balance
    validateBalance: async (walletId: string): Promise<{
        wallet_id: string;
        current_balance: number;
        calculated_balance: number;
        is_valid: boolean;
        difference: number;
        message: string;
    }> => {
        const response = await adminApiClient.get<ApiResponse<any>>(
            `/admin/ledger/wallet/${walletId}/validate`
        );
        return response.data.data!;
    },
};