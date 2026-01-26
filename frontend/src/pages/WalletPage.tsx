import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { MainLayout } from '@/components/layout/MainLayout';
import { WalletCard } from '@/components/wallet/WalletCard';
import { WalletHistory } from '@/components/wallet/WalletHistory';
import { authApi } from '@/api/auth.api';
import { walletApi } from '@/api/wallet.api';

export const WalletPage: React.FC = () => {
    const [showBalance, setShowBalance] = useState(true);
    const [selectedWalletId, setSelectedWalletId] = useState<string | null>(null);

    // Fetch user profile with wallets
    const { data: profile } = useQuery({
        queryKey: ['profile'],
        queryFn: authApi.getProfile,
    });

    // Fetch wallet history for selected wallet
    const { data: walletHistory } = useQuery({
        queryKey: ['walletHistory', selectedWalletId],
        queryFn: () => walletApi.getHistory(selectedWalletId!, 20, 0),
        enabled: !!selectedWalletId,
    });

    // Auto-select main wallet on load
    React.useEffect(() => {
        if (profile?.wallets && !selectedWalletId) {
            const mainWallet = profile.wallets.find((w) => w.wallet_type === 'main');
            if (mainWallet) {
                setSelectedWalletId(mainWallet.wallet_id);
            }
        }
    }, [profile, selectedWalletId]);

    const selectedWallet = profile?.wallets.find((w) => w.wallet_id === selectedWalletId);

    return (
        <MainLayout>
            <div className="space-y-6">
                {/* Header */}
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Wallet Saya</h1>
                    <p className="text-gray-600 mt-1">Kelola saldo dan riwayat transaksi</p>
                </div>

                {/* Wallet Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {profile?.wallets.map((wallet) => (
                        <div
                            key={wallet.wallet_id}
                            onClick={() => setSelectedWalletId(wallet.wallet_id)}
                            className={`cursor-pointer transition-all ${selectedWalletId === wallet.wallet_id ? 'ring-2 ring-blue-500 ring-offset-2' : ''
                                }`}
                        >
                            <WalletCard
                                wallet={wallet}
                                showBalance={showBalance}
                                onToggleBalance={() => setShowBalance(!showBalance)}
                            />
                        </div>
                    ))}
                </div>

                {/* Wallet History */}
                {selectedWallet && (
                    <div>
                        <h2 className="text-lg font-semibold text-gray-900 mb-4">
                            Riwayat Transaksi - {selectedWallet.wallet_type === 'main' ? 'Saldo Utama' : selectedWallet.wallet_type}
                        </h2>

                        {walletHistory && walletHistory.entries.length > 0 ? (
                            <WalletHistory entries={walletHistory.entries} />
                        ) : (
                            <div className="card text-center py-8">
                                <p className="text-gray-500">Belum ada riwayat transaksi</p>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </MainLayout>
    );
};