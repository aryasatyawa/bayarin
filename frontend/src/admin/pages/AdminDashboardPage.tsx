import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
    Users,
    WalletCards,
    TrendingUp,
    AlertCircle,
    XCircle,
    ArrowUpRight,
    ArrowDownLeft,
} from 'lucide-react';
import { AdminLayout } from '@/admin/layouts/AdminLayout';
import { StatsCard } from '@/admin/components/dashboard/StatsCard';
import { SimpleChart } from '@/admin/components/dashboard/SimpleChart';
import { dashboardApi } from '@/admin/api/dashboard.api';
import { formatCurrency } from '@/utils/currency';

export const AdminDashboardPage: React.FC = () => {
    // Fetch dashboard overview
    const { data: overview, isLoading } = useQuery({
        queryKey: ['admin-dashboard-overview'],
        queryFn: dashboardApi.getOverview,
        refetchInterval: 30000, // Refresh every 30s
    });

    // Fetch transaction summary (last 7 days)
    const { data: summary } = useQuery({
        queryKey: ['admin-transaction-summary'],
        queryFn: () => {
            const endDate = new Date();
            const startDate = new Date();
            startDate.setDate(startDate.getDate() - 7);

            return dashboardApi.getTransactionSummary(
                startDate.toISOString().split('T')[0],
                endDate.toISOString().split('T')[0]
            );
        },
        refetchInterval: 60000, // Refresh every minute
    });

    if (isLoading) {
        return (
            <AdminLayout>
                <div className="flex items-center justify-center h-64">
                    <p className="text-gray-500">Loading dashboard...</p>
                </div>
            </AdminLayout>
        );
    }

    return (
        <AdminLayout>
            <div className="space-y-6">
                {/* Header */}
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 mb-2">Dashboard Overview</h1>
                    <p className="text-gray-500">Monitoring sistem real-time</p>
                </div>

                {/* Stats Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <StatsCard
                        title="Total Users"
                        value={overview?.total_users.toLocaleString() || '0'}
                        icon={Users}
                        color="blue"
                    />
                    <StatsCard
                        title="System Liability"
                        value={overview ? formatCurrency(overview.total_system_liability) : 'Rp 0'}
                        icon={WalletCards}
                        color="green"
                    />
                    <StatsCard
                        title="Today Volume"
                        value={overview ? formatCurrency(overview.today_volume) : 'Rp 0'}
                        icon={TrendingUp}
                        color="purple"
                    />
                    <StatsCard
                        title="Pending Transactions"
                        value={overview?.pending_transactions || 0}
                        icon={AlertCircle}
                        color="yellow"
                    />
                </div>


                {/* Secondary Stats */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                        <div className="flex items-center gap-3 mb-4">
                            <div className="bg-green-500/10 p-3 rounded-lg border border-green-500/20">
                                <ArrowDownLeft className="w-5 h-5 text-green-500" />
                            </div>
                            <div>
                                <p className="text-gray-500 text-sm">Today Topups</p>
                                <p className="text-gray-900 text-2xl font-bold">
                                    {overview?.today_topups || 0}
                                </p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                        <div className="flex items-center gap-3 mb-4">
                            <div className="bg-blue-500/10 p-3 rounded-lg border border-blue-500/20">
                                <ArrowUpRight className="w-5 h-5 text-blue-500" />
                            </div>
                            <div>
                                <p className="text-gray-500 text-sm">Today Transfers</p>
                                <p className="text-gray-900 text-2xl font-bold">
                                    {overview?.today_transfers || 0}
                                </p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                        <div className="flex items-center gap-3 mb-4">
                            <div className="bg-red-500/10 p-3 rounded-lg border border-red-500/20">
                                <XCircle className="w-5 h-5 text-red-500" />
                            </div>
                            <div>
                                <p className="text-gray-500 text-sm">Failed (7d)</p>
                                <p className="text-gray-900 text-2xl font-bold">
                                    {overview?.failed_transactions || 0}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Charts */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Transaction Volume Chart */}
                    {summary?.daily_breakdown && (
                        <SimpleChart
                            data={summary.daily_breakdown}
                            title="Transaction Volume (7 Days)"
                        />
                    )}

                    {/* Transaction by Type */}
                    <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                        <h3 className="text-gray-900 font-semibold text-lg mb-6">
                            Transaction by Type (Today)
                        </h3>
                        <div className="space-y-4">
                            <div className="flex justify-between items-center p-4 bg-gray-50 rounded-lg">
                                <div className="flex items-center gap-3">
                                    <div className="bg-green-500/20 p-2 rounded">
                                        <ArrowDownLeft className="w-5 h-5 text-green-500" />
                                    </div>
                                    <span className="text-gray-900 font-medium">Topup</span>
                                </div>
                                <span className="text-gray-900 font-bold">
                                    {overview?.today_topups || 0}
                                </span>
                            </div>

                            <div className="flex justify-between items-center p-4 bg-gray-50 rounded-lg">
                                <div className="flex items-center gap-3">
                                    <div className="bg-blue-500/20 p-2 rounded">
                                        <ArrowUpRight className="w-5 h-5 text-blue-500" />
                                    </div>
                                    <span className="text-gray-900 font-medium">Transfer</span>
                                </div>
                                <span className="text-gray-900 font-bold">
                                    {overview?.today_transfers || 0}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Quick Actions */}
                <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                    <h3 className="text-gray-900 font-semibold text-lg mb-4">Quick Actions</h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <button className="p-4 bg-gray-50 hover:bg-gray-100 rounded-lg text-left transition-colors">
                            <p className="text-gray-900 font-medium">View Pending</p>
                            <p className="text-gray-500 text-sm mt-1">
                                {overview?.pending_transactions || 0} pending transactions
                            </p>
                        </button>

                        <button className="p-4 bg-gray-50 hover:bg-gray-100 rounded-lg text-left transition-colors">
                            <p className="text-gray-900 font-medium">Failed Transactions</p>
                            <p className="text-gray-500 text-sm mt-1">
                                {overview?.failed_transactions || 0} failed (7 days)
                            </p>
                        </button>

                        <button className="p-4 bg-gray-50 hover:bg-gray-100 rounded-lg text-left transition-colors">
                            <p className="text-gray-900 font-medium">Active Wallets</p>
                            <p className="text-gray-500 text-sm mt-1">
                                {overview?.total_active_wallets || 0} wallets
                            </p>
                        </button>
                    </div>
                </div>
            </div>
        </AdminLayout>
    );
};