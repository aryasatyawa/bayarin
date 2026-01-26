import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Filter } from 'lucide-react';
import { MainLayout } from '@/components/layout/MainLayout';
import { TransactionList } from '@/components/transaction/TransactionList';
import { Button } from '@/components/common/Button';
import { transactionApi } from '@/api/transaction.api';

export const HistoryPage: React.FC = () => {
    const [limit] = useState(20);
    const [offset, setOffset] = useState(0);

    const { data: transactions, isLoading } = useQuery({
        queryKey: ['transactions', limit, offset],
        queryFn: () => transactionApi.getHistory(limit, offset),
    });

    const handleLoadMore = () => {
        setOffset((prev) => prev + limit);
    };

    return (
        <MainLayout>
            <div className="space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900">Riwayat Transaksi</h1>
                        <p className="text-gray-600 mt-1">Semua aktivitas transaksi Anda</p>
                    </div>

                    <Button variant="ghost" className="flex items-center gap-2">
                        <Filter className="w-5 h-5" />
                        Filter
                    </Button>
                </div>

                {/* Transaction List */}
                {isLoading ? (
                    <div className="card text-center py-8">
                        <p className="text-gray-500">Memuat transaksi...</p>
                    </div>
                ) : transactions && transactions.length > 0 ? (
                    <>
                        <TransactionList transactions={transactions} />

                        {transactions.length >= limit && (
                            <div className="text-center">
                                <Button variant="secondary" onClick={handleLoadMore}>
                                    Muat Lebih Banyak
                                </Button>
                            </div>
                        )}
                    </>
                ) : (
                    <div className="card text-center py-8">
                        <p className="text-gray-500">Belum ada transaksi</p>
                    </div>
                )}
            </div>
        </MainLayout>
    );
};