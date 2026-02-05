import { ApiResponse } from '@/api/axios';
import { adminApiClient } from '@/admin/api/admin-api';
import { TransactionDetail } from '@/admin/types/ledger.types';

export const adminTransactionsApi = {
    // Get all transactions with filter
    getAllTransactions: async (filter: any): Promise<{
        transactions: TransactionDetail[];
        total: number;
    }> => {
        const response = await adminApiClient.get<ApiResponse<{
            transactions: TransactionDetail[];
            total: number;
        }>>('/admin/transactions', { params: filter });
        return response.data.data!;
    },

    // Get transaction detail
    getTransactionDetail: async (transactionId: string): Promise<TransactionDetail> => {
        const response = await adminApiClient.get<ApiResponse<TransactionDetail>>(
            `/admin/transactions/${transactionId}`
        );
        return response.data.data!;
    },

    // Get pending transactions
    getPendingTransactions: async (
        limit: number = 20,
        offset: number = 0
    ): Promise<TransactionDetail[]> => {
        const response = await adminApiClient.get<ApiResponse<TransactionDetail[]>>(
            '/admin/transactions/pending',
            { params: { limit, offset } }
        );
        return response.data.data!;
    },

    // Get failed transactions
    getFailedTransactions: async (
        days: number = 7,
        limit: number = 20,
        offset: number = 0
    ): Promise<TransactionDetail[]> => {
        const response = await adminApiClient.get<ApiResponse<TransactionDetail[]>>(
            '/admin/transactions/failed',
            { params: { days, limit, offset } }
        );
        return response.data.data!;
    },
};