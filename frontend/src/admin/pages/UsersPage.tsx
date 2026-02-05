import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Search, Lock, Unlock, Wallet as WalletIcon } from 'lucide-react';
import { AdminLayout } from '@/admin/layouts/AdminLayout';
import { Input } from '@/components/common/Input';
import { Button } from '@/components/common/Button';
import { Modal } from '@/components/common/Modal';
import { Card } from '@/components/common/Card';
import { adminUsersApi } from '@/admin/api/users.api';
import { formatCurrency } from '@/utils/currency';
import { formatDateTime } from '@/utils/date';
import { toast } from 'react-hot-toast';
import { useAdminAuth } from '@/admin/hooks/useAdminAuth';

export const UsersPage: React.FC = () => {
    const { canFreezeWallet } = useAdminAuth();
    const [searchQuery, setSearchQuery] = useState('');
    const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
    const [freezeModal, setFreezeModal] = useState<{
        isOpen: boolean;
        walletId: string | null;
        action: 'freeze' | 'unfreeze';
    }>({ isOpen: false, walletId: null, action: 'freeze' });
    const [reason, setReason] = useState('');

    // Search users
    const { data: searchResults, isLoading: isSearching } = useQuery({
        queryKey: ['admin-search-users', searchQuery],
        queryFn: () => adminUsersApi.searchUsers(searchQuery, 20, 0),
        enabled: searchQuery.length >= 3,
    });

    // Get user details
    const { data: userDetails, refetch: refetchUserDetails } = useQuery({
        queryKey: ['admin-user-details', selectedUserId],
        queryFn: () => adminUsersApi.getUserDetails(selectedUserId!),
        enabled: !!selectedUserId,
    });

    const handleFreezeWallet = async () => {
        if (!freezeModal.walletId || !reason.trim()) {
            toast.error('Alasan wajib diisi');
            return;
        }

        try {
            if (freezeModal.action === 'freeze') {
                await adminUsersApi.freezeWallet(freezeModal.walletId, reason);
                toast.success('Wallet berhasil dibekukan');
            } else {
                await adminUsersApi.unfreezeWallet(freezeModal.walletId, reason);
                toast.success('Wallet berhasil diaktifkan kembali');
            }

            setFreezeModal({ isOpen: false, walletId: null, action: 'freeze' });
            setReason('');
            refetchUserDetails();
        } catch (error: unknown) {
            if (error instanceof Error) {
                toast.error(error.message || 'Gagal memproses');
            } else {
                toast.error('Gagal memproses');
            }
        }
    };

    return (
        <AdminLayout>
            <div className="space-y-6">
                {/* Header */}
                <div>
                    <h1 className="text-3xl font-bold text-gray-800 mb-2">User Management</h1>
                    <p className="text-gray-500">Search dan kelola user wallets</p>
                </div>

                {/* Search */}
                <Card>
                    <div className="flex gap-4">
                        <div className="flex-1">
                            <Input
                                placeholder="Search by email, phone, atau nama (min 3 karakter)..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                            />
                        </div>
                        <Button variant="primary" className="flex items-center gap-2">
                            <Search className="w-5 h-5" />
                            Search
                        </Button>
                    </div>

                    {isSearching && (
                        <p className="text-gray-500 mt-4">Searching...</p>
                    )}
                </Card>

                {/* Search Results */}
                {searchResults && searchResults.length > 0 && (
                    <div className="grid grid-cols-1 gap-4">
                        {searchResults.map((user: any) => (
                            <Card
                                key={user.id}
                                className="cursor-pointer hover:border-gray-300"
                                onClick={() => setSelectedUserId(user.id)}
                            >
                                <div className="flex justify-between items-start">
                                    <div className="flex-1">
                                        <h3 className="text-gray-800 font-semibold text-lg mb-2">
                                            {user.full_name}
                                        </h3>
                                        <div className="space-y-1 text-sm">
                                            <p className="text-gray-600">Email: {user.email}</p>
                                            <p className="text-gray-600">Phone: {user.phone}</p>
                                        </div>
                                    </div>
                                    <div>
                                        <span
                                            className={`px-3 py-1 rounded-full text-xs font-medium ${user.status === 'active'
                                                ? 'bg-green-100 text-green-700'
                                                : 'bg-red-100 text-red-700'
                                                }`}
                                        >
                                            {user.status}
                                        </span>
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                )}

                {/* User Details */}
                {userDetails && (
                    <div className="space-y-6">
                        {/* User Info */}
                        <Card>
                            <h2 className="text-gray-800 font-semibold text-xl mb-4">User Details</h2>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div>
                                    <p className="text-gray-500 text-sm">Full Name</p>
                                    <p className="text-gray-800 font-medium">{userDetails.user.full_name}</p>
                                </div>
                                <div>
                                    <p className="text-gray-500 text-sm">Email</p>
                                    <p className="text-gray-800 font-medium">{userDetails.user.email}</p>
                                </div>
                                <div>
                                    <p className="text-gray-500 text-sm">Phone</p>
                                    <p className="text-gray-800 font-medium">{userDetails.user.phone}</p>
                                </div>
                                <div>
                                    <p className="text-gray-500 text-sm">Status</p>
                                    <p className="text-gray-800 font-medium capitalize">{userDetails.user.status}</p>
                                </div>
                            </div>
                        </Card>

                        {/* Statistics */}
                        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                            <Card>
                                <p className="text-gray-500 text-sm mb-2">Total Balance</p>
                                <p className="text-gray-800 text-2xl font-bold">
                                    {formatCurrency(userDetails.total_balance)}
                                </p>
                            </Card>
                            <Card>
                                <p className="text-gray-500 text-sm mb-2">Total Transactions</p>
                                <p className="text-gray-800 text-2xl font-bold">
                                    {userDetails.total_transactions}
                                </p>
                            </Card>
                            <Card>
                                <p className="text-gray-500 text-sm mb-2">Success</p>
                                <p className="text-green-600 text-2xl font-bold">
                                    {userDetails.success_transactions}
                                </p>
                            </Card>
                            <Card>
                                <p className="text-gray-500 text-sm mb-2">Failed</p>
                                <p className="text-red-600 text-2xl font-bold">
                                    {userDetails.failed_transactions}
                                </p>
                            </Card>
                        </div>

                        {/* Wallets */}
                        <Card>
                            <h2 className="text-gray-800 font-semibold text-xl mb-4">Wallets</h2>
                            <div className="space-y-4">
                                {userDetails.wallets.map((wallet) => (
                                    <div
                                        key={wallet.id}
                                        className="bg-gray-50 rounded-lg p-4 border border-gray-200"
                                    >
                                        <div className="flex justify-between items-start">
                                            <div className="flex-1">
                                                <div className="flex items-center gap-3 mb-2">
                                                    <WalletIcon className="w-5 h-5 text-blue-600" />
                                                    <h3 className="text-gray-800 font-semibold capitalize">
                                                        {wallet.wallet_type}
                                                    </h3>
                                                    <span
                                                        className={`px-2 py-1 rounded text-xs font-medium ${wallet.status === 'active'
                                                            ? 'bg-green-100 text-green-700'
                                                            : wallet.status === 'frozen'
                                                                ? 'bg-yellow-100 text-yellow-700'
                                                                : 'bg-red-100 text-red-700'
                                                            }`}
                                                    >
                                                        {wallet.status}
                                                    </span>
                                                </div>

                                                <div className="grid grid-cols-2 gap-4 mt-3">
                                                    <div>
                                                        <p className="text-gray-500 text-sm">Balance</p>
                                                        <p className="text-gray-800 font-bold text-lg">
                                                            {formatCurrency(wallet.balance)}
                                                        </p>
                                                    </div>
                                                    <div>
                                                        .                                                        <p className="text-gray-500 text-sm">Transactions</p>
                                                        <p className="text-gray-800 font-medium">
                                                            {wallet.transaction_count}
                                                        </p>
                                                    </div>
                                                </div>

                                                {wallet.last_activity_at && (
                                                    <p className="text-gray-400 text-xs mt-2">
                                                        Last activity: {formatDateTime(wallet.last_activity_at)}
                                                    </p>
                                                )}
                                            </div>

                                            {/* Freeze/Unfreeze Button */}
                                            {canFreezeWallet() && (
                                                <div>
                                                    {wallet.status === 'active' ? (
                                                        <Button
                                                            variant="secondary"
                                                            size="sm"
                                                            onClick={() =>
                                                                setFreezeModal({
                                                                    isOpen: true,
                                                                    walletId: wallet.id,
                                                                    action: 'freeze',
                                                                })
                                                            }
                                                            className="text-yellow-600 border-yellow-300 hover:bg-yellow-50"
                                                        >
                                                            <Lock className="w-4 h-4 mr-2" />
                                                            Freeze
                                                        </Button>
                                                    ) : wallet.status === 'frozen' ? (
                                                        <Button
                                                            variant="secondary"
                                                            size="sm"
                                                            onClick={() =>
                                                                setFreezeModal({
                                                                    isOpen: true,
                                                                    walletId: wallet.id,
                                                                    action: 'unfreeze',
                                                                })
                                                            }
                                                            className="text-green-600 border-green-300 hover:bg-green-50"
                                                        >
                                                            <Unlock className="w-4 h-4 mr-2" />
                                                            Unfreeze
                                                        </Button>
                                                    ) : null}
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </Card>
                    </div>
                )}

                {/* Freeze/Unfreeze Modal */}
                <Modal
                    isOpen={freezeModal.isOpen}
                    onClose={() => setFreezeModal({ isOpen: false, walletId: null, action: 'freeze' })}
                    title={freezeModal.action === 'freeze' ? 'Freeze Wallet' : 'Unfreeze Wallet'}
                >
                    <div className="space-y-4">
                        <p className="text-gray-600">
                            {freezeModal.action === 'freeze'
                                ? 'Wallet yang dibekukan tidak dapat melakukan transaksi.'
                                : 'Wallet akan diaktifkan kembali dan dapat melakukan transaksi.'}
                        </p>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Alasan <span className="text-red-500">*</span>
                            </label>
                            <textarea
                                value={reason}
                                onChange={(e) => setReason(e.target.value)}
                                placeholder="Masukkan alasan (min 10 karakter)..."
                                className="w-full px-4 py-2 bg-white text-gray-800 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                                rows={4}
                            />
                        </div>

                        <div className="flex gap-3">
                            <Button
                                variant="secondary"
                                onClick={() => setFreezeModal({ isOpen: false, walletId: null, action: 'freeze' })}
                                className="flex-1"
                            >
                                Batal
                            </Button>
                            <Button
                                variant={freezeModal.action === 'freeze' ? 'danger' : 'primary'}
                                onClick={handleFreezeWallet}
                                className="flex-1"
                            >
                                {freezeModal.action === 'freeze' ? 'Freeze Wallet' : 'Unfreeze Wallet'}
                            </Button>
                        </div>
                    </div>
                </Modal>
            </div>
        </AdminLayout>
    );
};