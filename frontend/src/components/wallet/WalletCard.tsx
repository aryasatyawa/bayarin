import React from 'react';
import { Wallet, Eye, EyeOff } from 'lucide-react';
import { WalletBalance } from '@/types/wallet.types';
import { Card } from '@/components/common/Card';
import { formatCurrency } from '@/utils/currency';

interface WalletCardProps {
    wallet: WalletBalance;
    showBalance?: boolean;
    onToggleBalance?: () => void;
}

export const WalletCard: React.FC<WalletCardProps> = ({
    wallet,
    showBalance = true,
    onToggleBalance,
}) => {
    const walletTypeLabels = {
        main: 'Saldo Utama',
        bonus: 'Saldo Bonus',
        cashback: 'Saldo Cashback',
    };

    const walletTypeColors = {
        main: 'from-blue-600 to-blue-800',
        bonus: 'from-blue-400 to-blue-600',
        cashback: 'from-blue-500 to-blue-700',
    };

    return (
        <Card className={`bg-linear-to-br ${walletTypeColors[wallet.wallet_type]} text-white`}>
            <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-2">
                    <Wallet className="w-6 h-6" />
                    <span className="font-medium">{walletTypeLabels[wallet.wallet_type]}</span>
                </div>
                {onToggleBalance && (
                    <button
                        onClick={onToggleBalance}
                        className="text-white/80 hover:text-white transition-colors"
                    >
                        {showBalance ? <Eye className="w-5 h-5" /> : <EyeOff className="w-5 h-5" />}
                    </button>
                )}
            </div>

            <div className="space-y-1">
                <p className="text-sm text-white/80">Total Saldo</p>
                <p className="text-3xl font-bold">
                    {showBalance ? formatCurrency(wallet.balance) : 'Rp ••••••••'}
                </p>
            </div>

            <div className="mt-4 pt-4 border-t border-white/20">
                <p className="text-xs text-white/60">
                    Status: <span className="text-white font-medium capitalize">{wallet.status}</span>
                </p>
            </div>
        </Card>
    );
};