import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { DollarSign, AlertCircle, Search } from 'lucide-react';
import { AdminLayout } from '@/admin/layouts/AdminLayout';
import { Card } from '@/components/common/Card';
import { Input } from '@/components/common/Input';
import { Button } from '@/components/common/Button';
import { Modal } from '@/components/common/Modal';
import { adminTransactionsApi } from '@/admin/api/transactions.api';
import { refundApi, RefundRequest } from '@/admin/api/refund.api';
import { formatCurrency, toMinorUnit } from '@/utils/currency';
import { formatDateTime } from '@/utils/date';
import { toast } from 'react-hot-toast';
import { v4 as uuidv4 } from 'uuid';

export const RefundPage: React.FC = () => {
    const [transactionId, setTransactionId] = useState('');
    const [selectedTransaction, setSelectedTransaction] = useState<any>(null);
    const [refundModal, setRefundModal] = useState(false);
    const [refundAmount, setRefundAmount] = useState('');
    const [refundReason, setRefundReason] = useState('');
    const [refundType, setRefundType] = useState<'full' | 'partial'>('full');
    const [isProcessing, setIsProcessing] = useState(false);

    // Get transaction detail
    const { data: transaction, refetch } = useQuery({
        queryKey: ['admin-transaction-detail', transactionId],
        queryFn: () => adminTransactionsApi.getTransactionDetail(transactionId),
        enabled: false,
    });

    // Get refund history
    const { data: refundHistory } = useQuery({
        queryKey: ['refund-history', selectedTransaction?.id],
        queryFn: () => refundApi.getRefundHistory(selectedTransaction.id),
        enabled: !!selectedTransaction,
    });

    const handleSearch = async () => {
        if (!transactionId.trim()) {
            toast.error('Masukkan Transaction ID');
            return;
        }

        try {
            await refetch();
            if (transaction) {
                setSelectedTransaction(transaction);
            }
        } catch (error) {
            toast.error('Transaction tidak ditemukan');
        }
    };

    const handleRefund = async () => {
        if (!refundReason.trim() || refundReason.length < 10) {
            toast.error('Alasan minimal 10 karakter');
            return;
        }

        if (refundType === 'partial') {
            const amount = parseFloat(refundAmount);
            if (!amount || amount <= 0) {
                toast.error('Jumlah refund tidak valid');
                return;
            }
            if (toMinorUnit(amount) > selectedTransaction.amount) {
                toast.error('Jumlah refund melebihi amount transaksi');
                return;
            }
        }

        setIsProcessing(true);
        try {
            const refundReq: RefundRequest = {
                original_transaction_id: selectedTransaction.id,
                reason: refundReason,
                idempotency_key: uuidv4(),
            };

            if (refundType === 'partial') {
                refundReq.amount = toMinorUnit(parseFloat(refundAmount));
            }

            await refundApi.refundTransaction(refundReq);
            toast.success('Refund berhasil diproses');

            setRefundModal(false);
            setRefundReason('');
            setRefundAmount('');
            setRefundType('full');

            // Refresh data
            await refetch();
        } catch (error: any) {
            toast.error(error.response?.data?.error?.message || 'Refund gagal');
        } finally {
            setIsProcessing(false);
        }
    };

    return (
        <AdminLayout>
            <div className="space-y-6">
                {/* Header */}
                <div>
                    <h1 className="text-3xl font-bold text-black mb-2">Refund Management</h1>
                    <p className="text-gray-400">Process refund & reversal transactions</p>
                </div>

                {/* Warning */}
                <Card className="bg-yellow-500/10 border-yellow-500/20">
                    <div className="flex items-start gap-3">
                        <AlertCircle className="w-5 h-5 text-yellow-500 mt-0.5" />
                        <div>
                            <h3 className="text-yellow-500 font-medium mb-1">Finance Admin Only</h3>
                            <p className="text-yellow-400 text-sm">
                                Refund akan membuat transaksi BARU. Pastikan alasan jelas dan valid.
                                Refund tidak bisa dibatalkan setelah diproses.
                            </p>
                        </div>
                    </div>
                </Card>

                {/* Search Transaction */}
                <Card className="bg-gray-800 border-gray-200">
                    <h2 className="text-black font-semibold text-lg mb-4">Search Transaction</h2>
                    <div className="flex gap-4">
                        <Input
                            placeholder="Transaction ID (UUID)"
                            value={transactionId}
                            onChange={(e) => setTransactionId(e.target.value)}
                            className="flex-1 bg-white text-gray-400 border-gray-600"
                        />
                        <Button onClick={handleSearch} className="flex items-center gap-2">
                            <Search className="w-5 h-5" />
                            Search
                        </Button>
                    </div>
                </Card>

                {/* Transaction Details */}
                {selectedTransaction && (
                    <Card className="bg-gray-800 border-gray-700">
                        <h2 className="text-white font-semibold text-lg mb-4">Transaction Details</h2>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                            <div>
                                <p className="text-gray-400 text-sm">Transaction ID</p>
                                <p className="text-white font-mono text-sm">{selectedTransaction.id}</p>
                            </div>
                            <div>
                                <p className="text-gray-400 text-sm">Type</p>
                                <p className="text-white font-medium capitalize">
                                    {selectedTransaction.transaction_type}
                                </p>
                            </div>
                            <div>
                                <p className="text-gray-400 text-sm">Amount</p>
                                <p className="text-white font-bold text-lg">
                                    {formatCurrency(selectedTransaction.amount)}
                                </p>
                            </div>
                            <div>
                                <p className="text-gray-400 text-sm">Status</p>
                                <span
                                    className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${selectedTransaction.status === 'success'
                                        ? 'bg-green-500/20 text-green-500'
                                        : selectedTransaction.status === 'failed'
                                            ? 'bg-red-500/20 text-red-500'
                                            : 'bg-yellow-500/20 text-yellow-500'
                                        }`}
                                >
                                    {selectedTransaction.status}
                                </span>
                            </div>
                            <div>
                                <p className="text-gray-400 text-sm">User Email</p>
                                <p className="text-white">{selectedTransaction.user_email}</p>
                            </div>
                            <div>
                                <p className="text-gray-400 text-sm">Created At</p>
                                <p className="text-white">{formatDateTime(selectedTransaction.created_at)}</p>
                            </div>
                        </div>

                        <div className="mb-6">
                            <p className="text-gray-400 text-sm mb-2">Description</p>
                            <p className="text-white">{selectedTransaction.description}</p>
                        </div>

                        {selectedTransaction.status === 'success' && (
                            <Button
                                variant="danger"
                                onClick={() => setRefundModal(true)}
                                className="w-full flex items-center justify-center gap-2"
                            >
                                <DollarSign className="w-5 h-5" />
                                Process Refund
                            </Button>
                        )}

                        {selectedTransaction.status !== 'success' && (
                            <div className="p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
                                <p className="text-red-500 text-sm">
                                    Hanya transaksi dengan status SUCCESS yang bisa di-refund.
                                </p>
                            </div>
                        )}
                    </Card>
                )}

                {/* Refund History */}
                {refundHistory && refundHistory.length > 0 && (
                    <Card className="bg-gray-800 border-gray-700">
                        <h2 className="text-white font-semibold text-lg mb-4">Refund History</h2>
                        <div className="space-y-3">
                            {refundHistory.map((item: any) => (
                                <div key={item.refund_transaction_id} className="bg-gray-700 rounded-lg p-4">
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <p className="text-white font-medium">
                                                {formatCurrency(item.amount)}
                                            </p>
                                            <p className="text-gray-400 text-sm mt-1">{item.reason}</p>
                                            <p className="text-gray-500 text-xs mt-2">
                                                {formatDateTime(item.created_at)}
                                            </p>
                                        </div>
                                        <span
                                            className={`px-2 py-1 rounded text-xs font-medium ${item.status === 'success'
                                                ? 'bg-green-500/20 text-green-500'
                                                : 'bg-yellow-500/20 text-yellow-500'
                                                }`}
                                        >
                                            {item.status}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </Card>
                )}

                {/* Refund Modal */}
                <Modal
                    isOpen={refundModal}
                    onClose={() => setRefundModal(false)}
                    title="Process Refund"
                    size="lg"
                >
                    <div className="space-y-4">
                        {/* Refund Type */}
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                Refund Type
                            </label>
                            <div className="grid grid-cols-2 gap-3">
                                <button
                                    onClick={() => setRefundType('full')}
                                    className={`p-3 rounded-lg border transition-all ${refundType === 'full'
                                        ? 'bg-blue-500/20 border-blue-500 text-blue-500'
                                        : 'bg-gray-700 border-gray-600 text-gray-400'
                                        }`}
                                >
                                    Full Refund
                                </button>
                                <button
                                    onClick={() => setRefundType('partial')}
                                    className={`p-3 rounded-lg border transition-all ${refundType === 'partial'
                                        ? 'bg-blue-500/20 border-blue-500 text-blue-500'
                                        : 'bg-gray-700 border-gray-600 text-gray-400'
                                        }`}
                                >
                                    Partial Refund
                                </button>
                            </div>
                        </div>

                        {/* Amount (if partial) */}
                        {refundType === 'partial' && (
                            <div>
                                <label className="block text-sm font-medium text-gray-300 mb-2">
                                    Refund Amount (Rupiah)
                                </label>
                                <Input
                                    type="number"
                                    placeholder="100000"
                                    value={refundAmount}
                                    onChange={(e) => setRefundAmount(e.target.value)}
                                    className="bg-gray-700 text-white border-gray-600"
                                />
                                <p className="text-xs text-gray-500 mt-1">
                                    Max: {selectedTransaction && formatCurrency(selectedTransaction.amount)}
                                </p>
                            </div>
                        )}

                        {/* Reason */}
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                Reason <span className="text-red-500">*</span>
                            </label>
                            <textarea
                                value={refundReason}
                                onChange={(e) => setRefundReason(e.target.value)}
                                placeholder="Alasan refund (min 10 karakter)..."
                                className="w-full px-4 py-2 bg-gray-700 text-white border border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                                rows={4}
                            />
                        </div>

                        {/* Actions */}
                        <div className="flex gap-3 pt-4">
                            <Button
                                variant="secondary"
                                onClick={() => setRefundModal(false)}
                                className="flex-1"
                            >
                                Cancel
                            </Button>
                            <Button
                                variant="danger"
                                onClick={handleRefund}
                                isLoading={isProcessing}
                                className="flex-1"
                            >
                                Process Refund
                            </Button>
                        </div>
                    </div>
                </Modal>
            </div>
        </AdminLayout>
    );
};