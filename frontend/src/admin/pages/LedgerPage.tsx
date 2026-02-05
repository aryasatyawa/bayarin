import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Search, Download } from 'lucide-react';
import { AdminLayout } from '@/admin/layouts/AdminLayout';
import { ledgerApi } from '@/admin/api/ledger.api';
import { LedgerFilter } from '@/admin/types/ledger.types';
import { Input } from '@/components/common/Input';
import { Button } from '@/components/common/Button';
import { formatCurrency } from '@/utils/currency';
import { formatDateTime } from '@/utils/date';

export const LedgerPage: React.FC = () => {
    const [filter, setFilter] = useState<Partial<LedgerFilter>>({
        limit: 50,
        offset: 0,
    });

    const [searchTerm, setSearchTerm] = useState('');

    const { data, isLoading, refetch } = useQuery({
        queryKey: ['admin-ledger', filter],
        queryFn: () => ledgerApi.getLedgerEntries(filter),
    });

    const handleSearch = () => {
        // Update filter based on search term
        if (searchTerm) {
            setFilter({
                ...filter,
                transaction_id: searchTerm,
                user_id: searchTerm,
                wallet_id: searchTerm,
                offset: 0,
            });
        }
    };

    const handleClearFilter = () => {
        setFilter({
            limit: 50,
            offset: 0,
            transaction_id: undefined,
            user_id: undefined,
            wallet_id: undefined,
            entry_type: undefined
        });
        setSearchTerm('');
        refetch();
    };

    return (
        <AdminLayout>
            <div className="space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold text-gray-800 mb-2">Ledger Viewer</h1>
                        <p className="text-gray-500">Read-only ledger entries (double-entry bookkeeping)</p>
                    </div>
                    <Button variant="ghost">
                        <Download className="w-5 h-5 mr-2" />
                        Export
                    </Button>
                </div>

                {/* Filters */}
                <div className="bg-white rounded-xl p-6 border border-gray-200">
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div className="md:col-span-2">
                            <Input
                                placeholder="Search by Transaction ID, Wallet ID, or User ID"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>

                        <Button variant="primary" onClick={handleSearch}>
                            <Search className="w-5 h-5 mr-2" />
                            Search
                        </Button>

                        <Button variant="secondary" onClick={handleClearFilter}>
                            Clear Filter
                        </Button>
                    </div>

                    <div className="mt-4 flex gap-2">
                        <button
                            onClick={() => setFilter({ ...filter, entry_type: 'debit' })}
                            className={`px-4 py-2 rounded-lg text-sm font-medium ${filter.entry_type === 'debit'
                                ? 'bg-red-500 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                }`}
                        >
                            Debit Only
                        </button>
                        <button
                            onClick={() => setFilter({ ...filter, entry_type: 'credit' })}
                            className={`px-4 py-2 rounded-lg text-sm font-medium ${filter.entry_type === 'credit'
                                ? 'bg-green-500 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                }`}
                        >
                            Credit Only
                        </button>
                        <button
                            onClick={() => setFilter({ ...filter, entry_type: undefined })}
                            className={`px-4 py-2 rounded-lg text-sm font-medium ${!filter.entry_type
                                ? 'bg-blue-500 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                }`}
                        >
                            All
                        </button>
                    </div>
                </div>

                {/* Ledger Table */}
                <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Entry Type
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Amount
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Balance Before
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Balance After
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Description
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Date
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {isLoading ? (
                                    <tr>
                                        <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                                            Loading ledger entries...
                                        </td>
                                    </tr>
                                ) : data?.entries.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                                            No ledger entries found
                                        </td>
                                    </tr>
                                ) : (
                                    data?.entries.map((entry) => (
                                        <tr key={entry.id} className="hover:bg-gray-50 transition-colors">
                                            <td className="px-6 py-4">
                                                <span
                                                    className={`px-2 py-1 rounded text-xs font-medium ${entry.entry_type === 'debit'
                                                        ? 'bg-red-100 text-red-700'
                                                        : 'bg-green-100 text-green-700'
                                                        }`}
                                                >
                                                    {entry.entry_type.toUpperCase()}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 text-gray-800 font-medium">
                                                {formatCurrency(entry.amount)}
                                            </td>
                                            <td className="px-6 py-4 text-gray-600">
                                                {formatCurrency(entry.balance_before)}
                                            </td>
                                            <td className="px-6 py-4 text-gray-600">
                                                {formatCurrency(entry.balance_after)}
                                            </td>
                                            <td className="px-6 py-4 text-gray-700">
                                                <p className="max-w-xs truncate">{entry.description}</p>
                                                <p className="text-xs text-gray-400 mt-1">
                                                    TX: {entry.transaction_id.substring(0, 8)}...
                                                </p>
                                            </td>
                                            <td className="px-6 py-4 text-gray-600 text-sm">
                                                {formatDateTime(entry.created_at)}
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>

                    {/* Pagination */}
                    {data && data.total > 0 && (
                        <div className="px-6 py-4 bg-gray-50 border-t border-gray-200 flex items-center justify-between">
                            <p className="text-gray-600 text-sm">
                                Showing {data.entries.length} of {data.total} entries
                            </p>
                            <div className="flex gap-2">
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    disabled={filter.offset === 0}
                                    onClick={() =>
                                        setFilter({
                                            ...filter,
                                            offset: Math.max(0, (filter.offset || 0) - (filter.limit || 50)),
                                        })
                                    }
                                >
                                    Previous
                                </Button>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    disabled={(filter.offset || 0) + (filter.limit || 50) >= data.total}
                                    onClick={() =>
                                        setFilter({
                                            ...filter,
                                            offset: (filter.offset || 0) + (filter.limit || 50),
                                        })
                                    }
                                >
                                    Next
                                </Button>
                            </div>
                        </div>
                    )}
                </div>

                {/* Info */}
                <div className="bg-blue-50 border border-blue-200 rounded-xl p-4">
                    <p className="text-blue-700 text-sm">
                        ℹ️ <strong>Read-Only Mode:</strong> Ledger entries cannot be modified or deleted.
                        This ensures data integrity and audit trail compliance.
                    </p>
                </div>
            </div>
        </AdminLayout>
    );
};