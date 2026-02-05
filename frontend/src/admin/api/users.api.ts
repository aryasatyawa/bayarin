import { apiClient, ApiResponse } from '@/api/axios';
import { UserInspectorDetail } from '@/admin/types/ledger.types';

export const adminUsersApi = {
    // Search users
    searchUsers: async (
        query: string,
        limit: number = 20,
        offset: number = 0
    ): Promise<any[]> => {
        const response = await apiClient.get<ApiResponse<any[]>>(
            '/admin/users/search',
            { params: { q: query, limit, offset } }
        );
        return response.data.data!;
    },

    // Get user details
    getUserDetails: async (userId: string): Promise<UserInspectorDetail> => {
        const response = await apiClient.get<ApiResponse<UserInspectorDetail>>(
            `/admin/users/${userId}`
        );
        return response.data.data!;
    },

    // Freeze wallet
    freezeWallet: async (walletId: string, reason: string): Promise<void> => {
        await apiClient.post(`/admin/wallets/${walletId}/freeze`, { reason });
    },

    // Unfreeze wallet
    unfreezeWallet: async (walletId: string, reason: string): Promise<void> => {
        await apiClient.post(`/admin/wallets/${walletId}/unfreeze`, { reason });
    },
};