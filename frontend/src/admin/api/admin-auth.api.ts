import { ApiResponse } from '@/api/axios';
import { adminApiClient } from '@/admin/api/admin-api';
import { AdminLoginRequest, AdminLoginResponse } from '@/admin/types/admin.types';

export const adminAuthApi = {
    // Admin login
    login: async (data: AdminLoginRequest): Promise<AdminLoginResponse> => {
        const response = await adminApiClient.post<ApiResponse<AdminLoginResponse>>(
            '/admin/auth/login',
            data
        );
        return response.data.data!;
    },
};