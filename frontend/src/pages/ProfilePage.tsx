import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { User, Mail, Phone, Calendar, Shield, Copy } from 'lucide-react';
import { toast } from 'react-hot-toast';
import { MainLayout } from '@/components/layout/MainLayout';
import { Card } from '@/components/common/Card';
import { authApi } from '@/api/auth.api';
import { formatDate } from '@/utils/date';

export const ProfilePage: React.FC = () => {
    const { data: profile } = useQuery({
        queryKey: ['profile'],
        queryFn: authApi.getProfile,
    });

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        toast.success('Copied to clipboard!');
    };

    return (
        <MainLayout>
            <div className="space-y-6 max-w-2xl mx-auto">
                {/* Header */}
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Profil Saya</h1>
                    <p className="text-gray-600 mt-1">Informasi akun Anda</p>
                </div>

                {/* Profile Card */}
                <Card>
                    <div className="flex items-center gap-4 mb-6">
                        <div className="bg-blue-100 p-4 rounded-full">
                            <User className="w-12 h-12 text-blue-600" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-gray-900">{profile?.full_name}</h2>
                            <p className="text-sm text-gray-600 capitalize">Status: {profile?.status}</p>
                        </div>
                    </div>

                    <div className="space-y-4">
                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                            <User className="w-5 h-5 text-gray-600" />
                            <div className="flex-1">
                                <p className="text-xs text-gray-600">User ID</p>
                                <p className="font-medium text-sm break-all">{profile?.id}</p>
                            </div>
                            <button
                                onClick={() => copyToClipboard(profile?.id || '')}
                                className="p-2 hover:bg-gray-200 rounded transition-colors"
                                title="Copy User ID"
                            >
                                <Copy className="w-4 h-4 text-gray-600" />
                            </button>
                        </div>

                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                            <Mail className="w-5 h-5 text-gray-600" />
                            <div>
                                <p className="text-xs text-gray-600">Email</p>
                                <p className="font-medium">{profile?.email}</p>
                            </div>
                        </div>

                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                            <Phone className="w-5 h-5 text-gray-600" />
                            <div>
                                <p className="text-xs text-gray-600">Nomor HP</p>
                                <p className="font-medium">{profile?.phone}</p>
                            </div>
                        </div>

                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                            <Calendar className="w-5 h-5 text-gray-600" />
                            <div>
                                <p className="text-xs text-gray-600">Bergabung Sejak</p>
                                <p className="font-medium">
                                    {profile?.created_at && formatDate(profile.created_at, 'dd MMMM yyyy')}
                                </p>
                            </div>
                        </div>

                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                            <Shield className="w-5 h-5 text-gray-600" />
                            <div>
                                <p className="text-xs text-gray-600">PIN Transaksi</p>
                                <p className="font-medium">
                                    {profile?.has_pin ? (
                                        <span className="text-blue-600">✓ Sudah diatur</span>
                                    ) : (
                                        <span className="text-blue-600">✗ Belum diatur</span>
                                    )}
                                </p>
                            </div>
                        </div>
                    </div>
                </Card>

                {/* Wallets Info */}
                <Card>
                    <h3 className="font-semibold text-gray-900 mb-4">Wallet Aktif</h3>
                    <div className="space-y-3">
                        {profile?.wallets.map((wallet) => (
                            <div key={wallet.wallet_id} className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                                <div>
                                    <p className="font-medium capitalize">{wallet.wallet_type}</p>
                                    <p className="text-sm text-gray-600">{wallet.wallet_id}</p>
                                </div>
                                <span className={`text-xs font-medium px-2 py-1 rounded ${wallet.status === 'active' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'
                                    }`}>
                                    {wallet.status}
                                </span>
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </MainLayout>
    );
};