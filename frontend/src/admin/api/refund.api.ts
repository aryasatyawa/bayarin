import { ApiResponse } from '@/api/axios';
import { adminApiClient } from '@/admin/api/admin-api';

export interface RefundRequest {
    original_transaction_id: string;
    reason: string;
    amount?: number;
    idempotency_key: string;
}

export interface RefundResponse {
    refund_transaction_id: string;
    original_transaction_id: string;
    amount: number;
    status: string;
    reason: string;
    created_at: string;
}

export const refundApi = {
    // Refund transaction
    refundTransaction: async (data: RefundRequest): Promise<RefundResponse> => {
        const response = await adminApiClient.post<ApiResponse<RefundResponse>>(
            '/admin/refund',
            data
        );
        return response.data.data!;
    },

    // Reverse transaction
    reverseTransaction: async (data: {
        original_transaction_id: string;
        reason: string;
        idempotency_key: string;
    }): Promise<RefundResponse> => {
        const response = await adminApiClient.post<ApiResponse<RefundResponse>>(
            '/admin/refund/reverse',
            data
        );
        return response.data.data!;
    },

    // Get refund history
    getRefundHistory: async (transactionId: string): Promise<any[]> => {
        const response = await adminApiClient.get<ApiResponse<any[]>>(
            `/admin/refund/history/${transactionId}`
        );
        return response.data.data!;
    },
};