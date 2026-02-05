import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Eye } from 'lucide-react';
import { AdminLayout } from '@/admin/layouts/AdminLayout';
import { adminTransactionsApi } from '@/admin/api/transactions.api';

import { Button } from '@/components/common/Button';
import { formatCurrency } from '@/utils/currency';
import { formatDateTime } from '@/utils/date';
import { Modal } from '@/components/common/Modal';
import { TransactionDetail } from '@/admin/types/ledger.types';

export const TransactionsPage: React.FC = () => {
    const [filter] = useState({
        limit: 20,
        offset: 0,
    });

    const [selectedTx, setSelectedTx] = useState<TransactionDetail | null>(null);
    const [showDetail, setShowDetail] = useState(false);

    const { data, isLoading } = useQuery({
        queryKey: ['admin-transactions', filter],
        queryFn: () => adminTransactionsApi.getAllTransactions(filter),
    });

    const handleViewDetail = async (txId: string) => {
        const detail = await adminTransactionsApi.getTransactionDetail(txId);
        setSelectedTx(detail);
        setShowDetail(true);
    };

    const getStatusColor = (status: string) => {
        const colors: Record<string, string> = {
            success: 'bg-green-100 text-green-800',
            pending: 'bg-yellow-100 text-yellow-800',
            failed: 'bg-red-100 text-red-800',
            reversed: 'bg-gray-100 text-gray-800',
        };
        return colors[status] || 'bg-gray-100 text-gray-800';
    };

    return (
        <AdminLayout>
            <div className="space-y-6">
                {/* Header */}
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 mb-2">Transactions</h1>
                    <p className="text-gray-500">Monitor all transactions</p>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="bg-white rounded-xl p-6 border border-gray-200">
                        <p className="text-gray-500 text-sm mb-2">Total Transactions</p>
                        <p className="text-gray-900 text-3xl font-bold">{data?.total || 0}</p>
                    </div>

                    <div className="bg-white rounded-xl p-6 border border-gray-200">
                        <p className="text-gray-500 text-sm mb-2">Success Rate</p>
                        <p className="text-green-600 text-3xl font-bold">
                            {data?.total ?
                                Math.round((data.transactions.filter(t => t.status === 'success').length / data.total) * 100)
                                : 0}%
                        </p>
                    </div>

                    <div className="bg-white rounded-xl p-6 border border-gray-200">
                        <p className="text-gray-500 text-sm mb-2">Failed</p>
                        <p className="text-red-600 text-3xl font-bold">
                            {data?.transactions.filter(t => t.status === 'failed').length || 0}
                        </p>
                    </div>
                </div>

                {/* Table */}
                <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Type
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Amount
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Status
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        User
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Date
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Actions
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {isLoading ? (
                                    <tr>
                                        <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                                            Loading transactions...
                                        </td>
                                    </tr>
                                ) : data?.transactions.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                                            No transactions found
                                        </td>
                                    </tr>
                                ) : (
                                    data?.transactions.map((tx) => (
                                        <tr key={tx.id} className="hover:bg-gray-50">
                                            <td className="px-6 py-4">
                                                <span className="text-gray-900 capitalize">{tx.transaction_type}</span>
                                            </td>
                                            <td className="px-6 py-4 text-gray-900 font-medium">
                                                {formatCurrency(tx.amount)}
                                            </td>
                                            <td className="px-6 py-4">
                                                <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(tx.status)}`}>
                                                    {tx.status.toUpperCase()}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 text-gray-500">
                                                {tx.user_email}
                                            </td>
                                            <td className="px-6 py-4 text-gray-500 text-sm">
                                                {formatDateTime(tx.created_at)}
                                            </td>
                                            <td className="px-6 py-4">
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleViewDetail(tx.id)}
                                                >
                                                    <Eye className="w-4 h-4" />
                                                </Button>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Transaction Detail Modal */}
                <Modal
                    isOpen={showDetail}
                    onClose={() => setShowDetail(false)}
                    title="Transaction Detail"
                    size="lg"
                >
                    {selectedTx && (
                        <div className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <p className="text-gray-500 text-sm">Transaction ID</p>
                                    <p className="text-gray-900 font-mono text-sm">{selectedTx.id}</p>
                                </div>
                                <div>
                                    <p className="text-gray-500 text-sm">Type</p>
                                    <p className="text-gray-900 capitalize">{selectedTx.transaction_type}</p>
                                </div>
                                <div>
                                    <p className="text-gray-500 text-sm">Amount</p>
                                    <p className="text-gray-900 font-bold">{formatCurrency(selectedTx.amount)}</p>
                                </div>
                                <div>
                                    <p className="text-gray-500 text-sm">Status</p>
                                    <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(selectedTx.status)}`}>
                                        {selectedTx.status.toUpperCase()}
                                    </span>
                                </div>
                            </div>

                            {selectedTx.ledger_entries && selectedTx.ledger_entries.length > 0 && (
                                <div>
                                    <p className="text-gray-900 font-medium mb-2">Ledger Entries</p>
                                    <div className="space-y-2">
                                        {selectedTx.ledger_entries.map((entry) => (
                                            <div key={entry.id} className="bg-gray-100 p-3 rounded">
                                                <div className="flex justify-between">
                                                    <span className={`text-sm font-medium ${entry.entry_type === 'debit' ? 'text-red-600' : 'text-green-600'
                                                        }`}>
                                                        {entry.entry_type.toUpperCase()}
                                                    </span>
                                                    <span className="text-gray-900">{formatCurrency(entry.amount)}</span>
                                                </div>
                                                <p className="text-gray-500 text-xs mt-1">{entry.description}</p>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </Modal>
            </div>
        </AdminLayout>
    );
};