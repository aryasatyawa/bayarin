import { apiClient } from './axios';
import type { ApiResponse } from './axios';
import {
    RegisterRequest,
    RegisterResponse,
    LoginRequest,
    LoginResponse,
    UserProfile,
    SetPINRequest,
} from '@/types/auth.types';

export const authApi = {
    // Register new user
    register: async (data: RegisterRequest): Promise<RegisterResponse> => {
        const response = await apiClient.post<ApiResponse<RegisterResponse>>(
            '/auth/register',
            data
        );
        return response.data.data!;
    },

    // Login user
    login: async (data: LoginRequest): Promise<LoginResponse> => {
        const response = await apiClient.post<ApiResponse<LoginResponse>>(
            '/auth/login',
            data
        );
        return response.data.data!;
    },

    // Get user profile
    getProfile: async (): Promise<UserProfile> => {
        const response = await apiClient.get<ApiResponse<UserProfile>>(
            '/user/profile'
        );
        return response.data.data!;
    },

    // Set transaction PIN
    setPIN: async (data: SetPINRequest): Promise<void> => {
        await apiClient.post('/user/pin', data);
    },

    // Verify PIN
    verifyPIN: async (data: SetPINRequest): Promise<void> => {
        await apiClient.post('/user/pin/verify', data);
    },
};