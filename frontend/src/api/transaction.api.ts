import { apiClient } from './axios';
import type { ApiResponse } from './axios';
import {
    TopupRequest,
    TransferRequest,
    TransactionResponse,
    TransactionDetail,
} from '@/types/transaction.types';

export const transactionApi = {
    // Topup wallet
    topup: async (data: TopupRequest): Promise<TransactionResponse> => {
        const response = await apiClient.post<ApiResponse<TransactionResponse>>(
            '/transaction/topup',
            data
        );
        return response.data.data!;
    },

    // Transfer to another user
    transfer: async (data: TransferRequest): Promise<TransactionResponse> => {
        const response = await apiClient.post<ApiResponse<TransactionResponse>>(
            '/transaction/transfer',
            data
        );
        return response.data.data!;
    },

    // Get transaction detail
    getTransaction: async (id: string): Promise<TransactionDetail> => {
        const response = await apiClient.get<ApiResponse<TransactionDetail>>(
            `/transaction/${id}`
        );
        return response.data.data!;
    },

    // Get user transaction history
    getHistory: async (
        limit: number = 20,
        offset: number = 0
    ): Promise<TransactionDetail[]> => {
        const response = await apiClient.get<ApiResponse<TransactionDetail[]>>(
            `/transaction/history?limit=${limit}&offset=${offset}`
        );
        return response.data.data!;
    },
};