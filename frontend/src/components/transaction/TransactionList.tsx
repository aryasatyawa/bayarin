import React from 'react';
import {
    ArrowUpRight,
    ArrowDownLeft,
    Wallet,
    CheckCircle,
    XCircle,
    Clock,
} from 'lucide-react';
import { TransactionDetail } from '@/types/transaction.types';
import { formatCurrency } from '@/utils/currency';
import { formatDateTime } from '@/utils/date';
import { Card } from '@/components/common/Card';

interface TransactionListProps {
    transactions: TransactionDetail[];
}

export const TransactionList: React.FC<TransactionListProps> = ({ transactions }) => {
    if (transactions.length === 0) {
        return (
            <Card>
                <p className="text-center text-gray-500 py-8">Belum ada transaksi</p>
            </Card>
        );
    }

    const getTransactionIcon = (type: string) => {
        switch (type) {
            case 'topup':
                return <ArrowDownLeft className="w-5 h-5" />;
            case 'transfer':
                return <ArrowUpRight className="w-5 h-5" />;
            case 'payment':
                return <Wallet className="w-5 h-5" />;
            default:
                return <Wallet className="w-5 h-5" />;
        }
    };

    const getTransactionColor = (type: string) => {
        switch (type) {
            case 'topup':
                return 'bg-blue-100 text-blue-600';
            case 'transfer':
                return 'bg-blue-100 text-blue-600';
            case 'payment':
                return 'bg-blue-100 text-blue-600';
            default:
                return 'bg-gray-100 text-gray-600';
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'success':
                return <CheckCircle className="w-5 h-5 text-blue-600" />;
            case 'failed':
                return <XCircle className="w-5 h-5 text-blue-600" />;
            case 'pending':
                return <Clock className="w-5 h-5 text-blue-600" />;
            default:
                return <Clock className="w-5 h-5 text-gray-600" />;
        }
    };

    const getTransactionLabel = (type: string) => {
        const labels: Record<string, string> = {
            topup: 'Topup',
            transfer: 'Transfer',
            payment: 'Pembayaran',
            withdrawal: 'Penarikan',
        };
        return labels[type] || type;
    };

    return (
        <div className="space-y-3">
            {transactions.map((transaction) => (
                <Card key={transaction.id} className="hover:shadow-md transition-shadow">
                    <div className="flex items-start justify-between">
                        <div className="flex items-start gap-3 flex-1">
                            <div className={`p-2 rounded-lg ${getTransactionColor(transaction.type)}`}>
                                {getTransactionIcon(transaction.type)}
                            </div>

                            <div className="flex-1">
                                <div className="flex items-center gap-2">
                                    <p className="font-medium text-gray-900">
                                        {getTransactionLabel(transaction.type)}
                                    </p>
                                    {getStatusIcon(transaction.status)}
                                </div>
                                <p className="text-sm text-gray-600 mt-1">{transaction.description}</p>
                                <p className="text-xs text-gray-500 mt-1">
                                    {formatDateTime(transaction.created_at)}
                                </p>
                                {transaction.reference_id && (
                                    <p className="text-xs text-gray-400 mt-1">
                                        Ref: {transaction.reference_id}
                                    </p>
                                )}
                            </div>
                        </div>

                        <div className="text-right">
                            <p
                                className={`text-lg font-semibold ${transaction.type === 'topup' ? 'text-blue-600' : 'text-blue-600'
                                    }`}
                            >
                                {transaction.type === 'topup' ? '+' : '-'}
                                {formatCurrency(transaction.amount)}
                            </p>
                            <p className="text-xs text-gray-500 capitalize mt-1">{transaction.status}</p>
                        </div>
                    </div>
                </Card>
            ))}
        </div>
    );
};