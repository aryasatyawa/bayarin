import { apiClient } from './axios';
import type { ApiResponse } from './axios';
import { WalletBalance, WalletHistory } from '@/types/wallet.types';

export const walletApi = {
    // Get wallet balance
    getBalance: async (type: string = 'main'): Promise<WalletBalance> => {
        const response = await apiClient.get<ApiResponse<WalletBalance>>(
            `/wallet/balance?type=${type}`
        );
        return response.data.data!;
    },

    // Get all wallets
    getAllWallets: async (): Promise<WalletBalance[]> => {
        const response = await apiClient.get<ApiResponse<WalletBalance[]>>(
            '/wallet/all'
        );
        return response.data.data!;
    },

    // Get wallet history
    getHistory: async (
        walletId: string,
        limit: number = 20,
        offset: number = 0
    ): Promise<WalletHistory> => {
        const response = await apiClient.get<ApiResponse<WalletHistory>>(
            `/wallet/${walletId}/history?limit=${limit}&offset=${offset}`
        );
        return response.data.data!;
    },
};