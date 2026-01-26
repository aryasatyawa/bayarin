import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
    ArrowUpRight,
    ArrowDownLeft,
    Wallet as WalletIcon,
    Plus,
    Copy,
    CheckCircle,
    AlertCircle,
    QrCode,
    TrendingDown,
    TrendingUp,
} from 'lucide-react';
import { toast } from 'react-hot-toast';
import { MainLayout } from '@/components/layout/MainLayout';
import { WalletCard } from '@/components/wallet/WalletCard';
import { Button } from '@/components/common/Button';
import { Card } from '@/components/common/Card';
import { Modal } from '@/components/common/Modal';
import { TopupForm } from '@/components/transaction/TopupForm';
import { authApi } from '@/api/auth.api';
import { transactionApi } from '@/api/transaction.api';
import { storage } from '@/utils/storage';
import { formatCurrency } from '@/utils/currency';
import type { UserProfile } from '@/types/auth.types';
import type { TransactionDetail } from '@/types/transaction.types';

export const DashboardPage: React.FC = () => {
    const navigate = useNavigate();
    const [showBalance, setShowBalance] = useState(true);
    const [showTopupModal, setShowTopupModal] = useState(false);
    const [showRequestModal, setShowRequestModal] = useState(false);

    const { data: profile, refetch: refetchProfile } = useQuery<UserProfile>({
        queryKey: ['profile'],
        queryFn: authApi.getProfile,
        staleTime: 1000 * 30,
        refetchOnWindowFocus: true,
    });

    const { data: transactions, refetch: refetchTransactions } = useQuery<TransactionDetail[]>({
        queryKey: ['transactions'],
        queryFn: () => transactionApi.getHistory(5, 0),
        staleTime: 1000 * 30,
        refetchOnWindowFocus: true,
    });

    useEffect(() => {
        if (profile) {
            storage.setUser({
                id: profile.id,
                email: profile.email,
                phone: profile.phone,
                full_name: profile.full_name,
            });
        }
    }, [profile]);

    const handleTopupSuccess = () => {
        setShowTopupModal(false);
        refetchProfile();
        refetchTransactions();
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        toast.success('Copied to clipboard!');
    };

    const mainWallet = profile?.wallets.find((w: any) => w.wallet_type === 'main');

    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const monthStart = new Date(today.getFullYear(), today.getMonth(), 1);

    const todayTransactions = transactions?.filter(
        (t) => new Date(t.created_at).getTime() >= today.getTime()
    ) || [];

    const monthTransactions = transactions?.filter(
        (t) => new Date(t.created_at).getTime() >= monthStart.getTime()
    ) || [];

    const todayExpense = todayTransactions
        .filter((t) => t.type === 'transfer' && t.status === 'success')
        .reduce((sum, t) => sum + t.amount, 0);

    const monthExpense = monthTransactions
        .filter((t) => t.type === 'transfer' && t.status === 'success')
        .reduce((sum, t) => sum + t.amount, 0);

    const lastTransaction = transactions?.[0];

    const securityStatus = {
        emailVerified: !!profile?.email,
        phoneVerified: !!profile?.phone,
        pinActive: profile?.has_pin,
        twoFaEnabled: false,
    };

    const securityScore = Object.values(securityStatus).filter(Boolean).length;

    const getTransactionIcon = (type: string) => {
        switch (type) {
            case 'topup':
                return <Plus className="w-5 h-5" />;
            case 'transfer':
                return <ArrowUpRight className="w-5 h-5" />;
            default:
                return <ArrowDownLeft className="w-5 h-5" />;
        }
    };

    const getTransactionColor = (type: string) => {
        switch (type) {
            case 'topup':
                return 'bg-green-100 text-green-600';
            case 'transfer':
                return 'bg-red-100 text-red-600';
            default:
                return 'bg-blue-100 text-blue-600';
        }
    };

    const getAmountColor = (type: string) => {
        return type === 'topup' ? 'text-green-600' : 'text-red-600';
    };

    const getRelativeTime = (date: string) => {
        const now = new Date();
        const transDate = new Date(date);
        const diffMs = now.getTime() - transDate.getTime();
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'Baru saja';
        if (diffMins < 60) return `${diffMins}m lalu`;
        if (diffHours < 24) return `${diffHours}h lalu`;
        if (diffDays < 7) return `${diffDays}d lalu`;
        return transDate.toLocaleDateString('id-ID');
    };

    return (
        <MainLayout>
            <div className="space-y-6">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
                    <p className="text-gray-600 mt-1">Selamat datang, {profile?.full_name}!</p>
                </div>

                {profile?.id && (
                    <div className="bg-linear-to-r from-blue-50 to-blue-100 rounded-lg p-4 border border-blue-200">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-xs text-blue-700 font-medium">User ID Anda</p>
                                <p className="text-sm font-mono text-blue-900 mt-1 break-all">{profile.id}</p>
                                <p className="text-xs text-blue-600 mt-2">Gunakan ID ini untuk menerima transfer</p>
                            </div>
                            <button
                                onClick={() => copyToClipboard(profile.id)}
                                className="p-2 hover:bg-blue-200 rounded-lg transition-colors shrink-0 ml-2"
                                title="Copy User ID"
                            >
                                <Copy className="w-5 h-5 text-blue-600" />
                            </button>
                        </div>
                    </div>
                )}

                {mainWallet && (
                    <WalletCard
                        wallet={mainWallet}
                        showBalance={showBalance}
                        onToggleBalance={() => setShowBalance(!showBalance)}
                    />
                )}

                <div className="grid grid-cols-4 gap-3">
                    <button
                        onClick={() => setShowTopupModal(true)}
                        className="flex flex-col items-center gap-2 p-4 rounded-lg bg-linear-to-br from-green-50 to-green-100 hover:from-green-100 hover:to-green-200 transition-all duration-300 group border border-green-200 hover:shadow-md hover:scale-105 active:scale-95"
                    >
                        <div className="bg-linear-to-br from-green-400 to-green-600 p-3 rounded-full group-hover:shadow-lg transition-all duration-300 shadow-md">
                            <Plus className="w-5 h-5 text-white" />
                        </div>
                        <span className="text-xs font-bold text-green-900">Top Up</span>
                    </button>

                    <button
                        onClick={() => navigate('/transfer')}
                        className="flex flex-col items-center gap-2 p-4 rounded-lg bg-linear-to-br from-purple-50 to-purple-100 hover:from-purple-100 hover:to-purple-200 transition-all duration-300 group border border-purple-200 hover:shadow-md hover:scale-105 active:scale-95"
                    >
                        <div className="bg-linear-to-br from-purple-400 to-purple-600 p-3 rounded-full group-hover:shadow-lg transition-all duration-300 shadow-md">
                            <ArrowUpRight className="w-5 h-5 text-white" />
                        </div>
                        <span className="text-xs font-bold text-purple-900">Transfer</span>
                    </button>

                    <button
                        onClick={() => navigate('/history')}
                        className="flex flex-col items-center gap-2 p-4 rounded-lg bg-linear-to-br from-cyan-50 to-cyan-100 hover:from-cyan-100 hover:to-cyan-200 transition-all duration-300 group border border-cyan-200 hover:shadow-md hover:scale-105 active:scale-95"
                    >
                        <div className="bg-linear-to-br from-cyan-400 to-cyan-600 p-3 rounded-full group-hover:shadow-lg transition-all duration-300 shadow-md">
                            <TrendingDown className="w-5 h-5 text-white" />
                        </div>
                        <span className="text-xs font-bold text-cyan-900">Riwayat</span>
                    </button>

                    <button
                        onClick={() => setShowRequestModal(true)}
                        className="flex flex-col items-center gap-2 p-4 rounded-lg bg-linear-to-br from-amber-50 to-amber-100 hover:from-amber-100 hover:to-amber-200 transition-all duration-300 group border border-amber-200 hover:shadow-md hover:scale-105 active:scale-95"
                    >
                        <div className="bg-linear-to-br from-amber-400 to-amber-600 p-3 rounded-full group-hover:shadow-lg transition-all duration-300 shadow-md">
                            <QrCode className="w-5 h-5 text-white" />
                        </div>
                        <span className="text-xs font-bold text-amber-900">Request</span>
                    </button>
                </div>

                <div className="grid grid-cols-3 gap-3">
                    <Card variant="interactive" className="bg-linear-to-br from-orange-50 to-orange-100 border border-orange-300 hover:border-orange-400">
                        <div className="flex items-center gap-3">
                            <div className="bg-linear-to-br from-orange-300 to-orange-500 p-3 rounded-lg shadow-md">
                                <TrendingDown className="w-5 h-5 text-white" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-xs text-orange-700 font-semibold">Hari Ini</p>
                                <p className="text-sm font-bold text-orange-900 truncate">
                                    {formatCurrency(todayExpense)}
                                </p>
                            </div>
                        </div>
                    </Card>

                    <Card variant="interactive" className="bg-linear-to-br from-red-50 to-red-100 border border-red-300 hover:border-red-400">
                        <div className="flex items-center gap-3">
                            <div className="bg-linear-to-br from-red-300 to-red-500 p-3 rounded-lg shadow-md">
                                <TrendingUp className="w-5 h-5 text-white" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-xs text-red-700 font-semibold">Bulan Ini</p>
                                <p className="text-sm font-bold text-red-900 truncate">
                                    {formatCurrency(monthExpense)}
                                </p>
                            </div>
                        </div>
                    </Card>

                    <Card variant="interactive" className="bg-linear-to-br from-indigo-50 to-indigo-100 border border-indigo-300 hover:border-indigo-400">
                        <div className="flex items-center gap-3">
                            <div className="bg-linear-to-br from-indigo-300 to-indigo-500 p-3 rounded-lg shadow-md">
                                <ArrowDownLeft className="w-5 h-5 text-white" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-xs text-indigo-700 font-semibold">Terakhir</p>
                                <p className="text-xs font-semibold text-indigo-900 truncate">
                                    {lastTransaction?.description || 'Belum ada'}
                                </p>
                            </div>
                        </div>
                    </Card>
                </div>

                <Card variant="elevated" className="bg-linear-to-r from-blue-50 via-blue-100 to-blue-50 border border-blue-300">
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <h3 className="font-bold text-lg text-gray-900 flex items-center gap-2">
                                <div className="w-2 h-2 bg-blue-600 rounded-full"></div>
                                Status Keamanan
                            </h3>
                            <div className="text-sm font-bold bg-blue-600 text-white px-3 py-1 rounded-full shadow-md">
                                {securityScore}/4 ✅
                            </div>
                        </div>

                        <div className="space-y-3">
                            <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-blue-200 transition-colors duration-200">
                                {securityStatus.emailVerified ? (
                                    <div className="w-6 h-6 rounded-full bg-green-100 flex items-center justify-center">
                                        <CheckCircle className="w-5 h-5 text-green-600 shadow-sm" />
                                    </div>
                                ) : (
                                    <div className="w-6 h-6 rounded-full bg-amber-100 flex items-center justify-center">
                                        <AlertCircle className="w-5 h-5 text-amber-600 shadow-sm" />
                                    </div>
                                )}
                                <span className="text-gray-700 font-medium flex-1">Email Terverifikasi</span>
                            </div>

                            <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-blue-200 transition-colors duration-200">
                                {securityStatus.phoneVerified ? (
                                    <div className="w-6 h-6 rounded-full bg-green-100 flex items-center justify-center">
                                        <CheckCircle className="w-5 h-5 text-green-600 shadow-sm" />
                                    </div>
                                ) : (
                                    <div className="w-6 h-6 rounded-full bg-amber-100 flex items-center justify-center">
                                        <AlertCircle className="w-5 h-5 text-amber-600 shadow-sm" />
                                    </div>
                                )}
                                <span className="text-gray-700 font-medium flex-1">Nomor Telepon Terverifikasi</span>
                            </div>

                            <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-blue-200 transition-colors duration-200">
                                {securityStatus.pinActive ? (
                                    <div className="w-6 h-6 rounded-full bg-green-100 flex items-center justify-center">
                                        <CheckCircle className="w-5 h-5 text-green-600 shadow-sm" />
                                    </div>
                                ) : (
                                    <div className="w-6 h-6 rounded-full bg-amber-100 flex items-center justify-center">
                                        <AlertCircle className="w-5 h-5 text-amber-600 shadow-sm" />
                                    </div>
                                )}
                                <span className="text-gray-700 font-medium flex-1">PIN Transaksi Aktif</span>
                            </div>

                            <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-blue-200 transition-colors duration-200">
                                {securityStatus.twoFaEnabled ? (
                                    <div className="w-6 h-6 rounded-full bg-green-100 flex items-center justify-center">
                                        <CheckCircle className="w-5 h-5 text-green-600 shadow-sm" />
                                    </div>
                                ) : (
                                    <div className="w-6 h-6 rounded-full bg-amber-100 flex items-center justify-center">
                                        <AlertCircle className="w-5 h-5 text-amber-600 shadow-sm" />
                                    </div>
                                )}
                                <span className="text-gray-700 font-medium flex-1">2FA: Nonaktif</span>
                            </div>
                        </div>

                        {securityScore < 4 && (
                            <Button
                                variant="primary"
                                size="sm"
                                className="w-full mt-2"
                                onClick={() => navigate('/settings')}
                            >
                                Tingkatkan Keamanan
                            </Button>
                        )}
                    </div>
                </Card>

                <Card variant="elevated" className="bg-linear-to-r from-emerald-50 via-emerald-100 to-emerald-50 border border-emerald-300">
                    <div className="space-y-4">
                        <div className="flex items-center gap-2">
                            <div className="w-2 h-2 bg-emerald-600 rounded-full"></div>
                            <h3 className="font-bold text-lg text-gray-900">Wallet Health</h3>
                        </div>

                        <div className="space-y-3 text-sm">
                            <div>
                                <div className="flex justify-between mb-2">
                                    <span className="text-emerald-700 font-semibold">Daily Limit</span>
                                    <span className="font-bold text-emerald-900 bg-emerald-200 px-2 py-1 rounded">Rp 5.000.000</span>
                                </div>

                                <div className="flex justify-between mb-2">
                                    <span className="text-emerald-700 font-semibold">Sisa Limit</span>
                                    <span className="font-bold text-emerald-900">
                                        {formatCurrency(5000000 - todayExpense)}
                                    </span>
                                </div>
                                <div className="w-full bg-emerald-200 rounded-full h-3 shadow-sm overflow-hidden">
                                    <div
                                        className="bg-linear-to-r from-emerald-400 to-emerald-600 h-3 rounded-full transition-all duration-500 shadow-md"
                                        style={{
                                            width: `${Math.min((todayExpense / 5000000) * 100, 100)}%`,
                                        }}
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </Card>

                <div>
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-lg font-bold text-gray-900 flex items-center gap-2">
                            <div className="w-2 h-2 bg-gray-900 rounded-full"></div>
                            Transaksi Terakhir
                        </h2>
                        <Button variant="ghost" size="sm" onClick={() => navigate('/history')}>
                            Lihat Semua →
                        </Button>
                    </div>

                    {transactions && transactions.length > 0 ? (
                        <div className="space-y-2">
                            {transactions.slice(0, 5).map((transaction) => (
                                <Card key={transaction.id} variant="interactive" className="border border-gray-200 hover:border-gray-300">
                                    <div className="flex items-center gap-4 p-1">
                                        <div className={`p-3 rounded-full shrink-0 shadow-md ${getTransactionColor(transaction.type)}`}>
                                            {getTransactionIcon(transaction.type)}
                                        </div>

                                        <div className="flex-1 min-w-0">
                                            <p className="font-bold text-gray-900 capitalize">
                                                {transaction.type === 'topup'
                                                    ? 'Top Up'
                                                    : transaction.type === 'transfer'
                                                        ? 'Transfer'
                                                        : transaction.description}
                                            </p>
                                            <div className="flex items-center gap-2 mt-1">
                                                <p className="text-xs text-gray-500">
                                                    {getRelativeTime(transaction.created_at)}
                                                </p>
                                                <span className={`text-xs font-bold px-2 py-1 rounded-full ${transaction.status === 'success'
                                                        ? 'bg-green-100 text-green-700'
                                                        : 'bg-amber-100 text-amber-700'
                                                    }`}>
                                                    {transaction.status === 'success' ? '✅ Berhasil' : '⏳ Proses'}
                                                </span>
                                            </div>
                                        </div>

                                        <div className="text-right shrink-0">
                                            <p className={`font-bold text-lg ${getAmountColor(transaction.type)}`}>
                                                {transaction.type === 'topup' ? '+' : '-'}
                                                {formatCurrency(transaction.amount)}
                                            </p>
                                        </div>
                                    </div>
                                </Card>
                            ))}
                        </div>
                    ) : (
                        <Card variant="elevated" className="text-center py-12 border-2 border-dashed border-gray-300 bg-gray-50">
                            <div className="flex flex-col items-center gap-4">
                                <div className="bg-linear-to-br from-blue-100 to-blue-200 p-6 rounded-full shadow-md">
                                    <WalletIcon className="w-10 h-10 text-blue-600" />
                                </div>
                                <div>
                                    <p className="text-gray-900 font-bold text-lg">Belum ada transaksi</p>
                                    <p className="text-gray-600 mt-2">Mulai dengan Top Up atau Transfer untuk melihat riwayat</p>
                                </div>
                                <Button size="sm" variant="primary" onClick={() => setShowTopupModal(true)}>
                                    🚀 Top Up Sekarang
                                </Button>
                            </div>
                        </Card>
                    )}
                </div>
            </div>

            <Modal isOpen={showTopupModal} onClose={() => setShowTopupModal(false)} title="Topup Saldo">
                <TopupForm onSuccess={handleTopupSuccess} />
            </Modal>

            <Modal isOpen={showRequestModal} onClose={() => setShowRequestModal(false)} title="Request Dana">
                <div className="space-y-4">
                    <div className="bg-linear-to-br from-purple-50 to-purple-100 border border-purple-300 rounded-lg p-4 text-center shadow-md">
                        <div className="bg-linear-to-br from-purple-400 to-purple-600 w-16 h-16 rounded-full mx-auto mb-3 flex items-center justify-center shadow-md">
                            <QrCode className="w-8 h-8 text-white" />
                        </div>
                        <p className="text-sm text-gray-700 font-semibold">User ID Anda:</p>
                        <p className="font-mono font-bold text-purple-900 mt-3 break-all bg-white px-3 py-2 rounded-lg border border-purple-200 text-sm">{profile?.id}</p>
                        <button
                            onClick={() => {
                                copyToClipboard(profile?.id || '');
                            }}
                            className="w-full mt-4 px-4 py-3 bg-linear-to-r from-purple-600 to-purple-700 text-white rounded-lg hover:from-purple-700 hover:to-purple-800 transition-all duration-200 hover:shadow-lg active:scale-95 flex items-center justify-center gap-2 font-semibold shadow-md"
                        >
                            <Copy className="w-4 h-4" />
                            Salin ID
                        </button>
                    </div>

                    <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 text-center">
                        <p className="text-xs text-gray-600 mb-3">Bagikan ID ini kepada orang yang ingin mengirim dana:</p>
                        <div className="space-y-2">
                            <p className="text-sm font-medium text-gray-900">📱 Share via Pesan / WhatsApp</p>
                            <p className="text-xs text-gray-600">atau biarkan mereka scan QR Code Anda</p>
                        </div>
                    </div>
                </div>
            </Modal>
        </MainLayout>
    );
};
