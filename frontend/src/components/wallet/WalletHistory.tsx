import React from 'react';
import { ArrowUpRight, ArrowDownLeft } from 'lucide-react';
import { LedgerEntry } from '@/types/wallet.types';
import { formatCurrency } from '@/utils/currency';
import { formatDateTime } from '@/utils/date';
import { Card } from '@/components/common/Card';

interface WalletHistoryProps {
    entries: LedgerEntry[];
}

export const WalletHistory: React.FC<WalletHistoryProps> = ({ entries }) => {
    if (entries.length === 0) {
        return (
            <Card>
                <p className="text-center text-gray-500 py-8">Belum ada riwayat transaksi</p>
            </Card>
        );
    }

    return (
        <div className="space-y-3">
            {entries.map((entry) => (
                <Card key={entry.id} className="hover:shadow-md transition-shadow">
                    <div className="flex items-start justify-between">
                        <div className="flex items-start gap-3">
                            <div
                                className={`p-2 rounded-lg ${entry.entry_type === 'credit'
                                    ? 'bg-blue-100 text-blue-600'
                                    : 'bg-blue-100 text-blue-600'
                                    }`}
                            >
                                {entry.entry_type === 'credit' ? (
                                    <ArrowDownLeft className="w-5 h-5" />
                                ) : (
                                    <ArrowUpRight className="w-5 h-5" />
                                )}
                            </div>

                            <div>
                                <p className="font-medium text-gray-900">{entry.description}</p>
                                <p className="text-sm text-gray-500 mt-1">
                                    {formatDateTime(entry.created_at)}
                                </p>
                                <div className="flex gap-4 mt-2 text-xs text-gray-600">
                                    <span>Sebelum: {formatCurrency(entry.balance_before)}</span>
                                    <span>Sesudah: {formatCurrency(entry.balance_after)}</span>
                                </div>
                            </div>
                        </div>

                        <div className="text-right">
                            <p
                                className={`text-lg font-semibold ${entry.entry_type === 'credit' ? 'text-blue-600' : 'text-blue-600'
                                    }`}
                            >
                                {entry.entry_type === 'credit' ? '+' : '-'}
                                {formatCurrency(entry.amount)}
                            </p>
                        </div>
                    </div>
                </Card>
            ))}
        </div>
    );
};