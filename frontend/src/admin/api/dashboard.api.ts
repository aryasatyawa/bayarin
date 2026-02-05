import { ApiResponse } from '@/api/axios';
import { adminApiClient } from '@/admin/api/admin-api';
import {
    DashboardOverview,
    DailyStats,
    TransactionSummary,
} from '@/admin/types/dashboard.types';

export const dashboardApi = {
    // Get dashboard overview
    getOverview: async (): Promise<DashboardOverview> => {
        const response = await adminApiClient.get<ApiResponse<DashboardOverview>>(
            '/admin/dashboard/overview'
        );
        return response.data.data!;
    },

    // Get daily stats
    getDailyStats: async (date?: string): Promise<DailyStats> => {
        const params = date ? { date } : {};
        const response = await adminApiClient.get<ApiResponse<DailyStats>>(
            '/admin/dashboard/daily-stats',
            { params }
        );
        return response.data.data!;
    },

    // Get transaction summary
    getTransactionSummary: async (
        startDate?: string,
        endDate?: string
    ): Promise<TransactionSummary> => {
        const params: any = {};
        if (startDate) params.start_date = startDate;
        if (endDate) params.end_date = endDate;

        const response = await adminApiClient.get<ApiResponse<TransactionSummary>>(
            '/admin/dashboard/transaction-summary',
            { params }
        );
        return response.data.data!;
    },
};