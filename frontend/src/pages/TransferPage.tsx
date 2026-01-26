import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { MainLayout } from '@/components/layout/MainLayout';
import { TransferForm } from '@/components/transaction/TransferForm';
import { Card } from '@/components/common/Card';
import { authApi } from '@/api/auth.api';
import { formatCurrency } from '@/utils/currency';
import type { UserProfile } from '@/types/auth.types';

export const TransferPage: React.FC = () => {
    const navigate = useNavigate();

    const { data: profile, refetch: refetchProfile } = useQuery<UserProfile>({
        queryKey: ['profile'],
        queryFn: authApi.getProfile,
        staleTime: 1000 * 30, // 30 seconds
        refetchOnWindowFocus: true,
    });

    const mainWallet = profile?.wallets.find((w: any) => w.wallet_type === 'main');

    const handleTransferSuccess = async () => {
        // Refetch profile to update wallet balance
        await refetchProfile();
        navigate('/history');
    };

    return (
        <MainLayout>
            <div className="space-y-6 max-w-2xl mx-auto">
                {/* Header */}
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Transfer</h1>
                    <p className="text-gray-600 mt-1">Transfer ke pengguna lain</p>
                </div>

                {/* Balance Info */}
                {mainWallet && (
                    <Card className="bg-blue-50 border border-blue-200">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-blue-700">Saldo Tersedia</p>
                                <p className="text-2xl font-bold text-blue-900 mt-1">
                                    {formatCurrency(mainWallet.balance)}
                                </p>
                            </div>
                        </div>
                    </Card>
                )}

                {/* PIN Warning */}
                {!profile?.has_pin && (
                    <Card className="bg-blue-50 border border-blue-200">
                        <p className="text-sm text-blue-800">
                            ⚠️ Anda belum mengatur PIN. Silakan atur PIN terlebih dahulu di menu Pengaturan.
                        </p>
                    </Card>
                )}

                {/* Transfer Form */}
                <Card>
                    <TransferForm onSuccess={handleTransferSuccess} />
                </Card>

                {/* Info */}
                <Card className="bg-gray-50">
                    <h3 className="font-medium text-gray-900 mb-2">Informasi Transfer</h3>
                    <ul className="text-sm text-gray-600 space-y-1 list-disc list-inside">
                        <li>Minimal transfer Rp 1.000</li>
                        <li>Tidak ada biaya transfer antar pengguna Bayarin</li>
                        <li>Transfer akan diproses secara real-time</li>
                        <li>Pastikan User ID tujuan sudah benar</li>
                    </ul>
                </Card>
            </div>
        </MainLayout>
    );
};