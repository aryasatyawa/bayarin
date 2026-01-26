import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';
import { Lock, User, Shield, Copy } from 'lucide-react';
import { MainLayout } from '@/components/layout/MainLayout';
import { Card } from '@/components/common/Card';
import { Input } from '@/components/common/Input';
import { Button } from '@/components/common/Button';
import { authApi } from '@/api/auth.api';

interface PINFormData {
    pin: string;
    confirmPin: string;
}

export const SettingsPage: React.FC = () => {
    const [isSettingPIN, setIsSettingPIN] = useState(false);

    const { data: profile, refetch } = useQuery({
        queryKey: ['profile'],
        queryFn: authApi.getProfile,
    });

    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
        watch,
    } = useForm<PINFormData>();

    const pin = watch('pin');

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        toast.success('Copied to clipboard!');
    };

    const onSubmitPIN = async (data: PINFormData) => {
        if (data.pin !== data.confirmPin) {
            toast.error('PIN tidak cocok');
            return;
        }

        setIsSettingPIN(true);
        try {
            await authApi.setPIN({ pin: data.pin });
            toast.success('PIN berhasil diatur!');
            reset();
            refetch();
        } catch (error: any) {
            const errorMessage = error.response?.data?.error?.message || 'Gagal mengatur PIN';
            toast.error(errorMessage);
        } finally {
            setIsSettingPIN(false);
        }
    };

    return (
        <MainLayout>
            <div className="space-y-6 max-w-2xl mx-auto">
                {/* Header */}
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Pengaturan</h1>
                    <p className="text-gray-600 mt-1">Kelola akun dan keamanan</p>
                </div>

                {/* Profile Info */}
                <Card>
                    <div className="flex items-start gap-4">
                        <div className="bg-blue-100 p-3 rounded-full">
                            <User className="w-6 h-6 text-blue-600" />
                        </div>
                        <div className="flex-1">
                            <h3 className="font-semibold text-gray-900">Informasi Profil</h3>
                            <div className="mt-3 space-y-2 text-sm">
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600">User ID</span>
                                    <div className="flex items-center gap-2">
                                        <span className="font-mono text-xs break-all max-w-xs">{profile?.id}</span>
                                        <button
                                            onClick={() => copyToClipboard(profile?.id || '')}
                                            className="p-1 hover:bg-gray-100 rounded transition-colors"
                                            title="Copy User ID"
                                        >
                                            <Copy className="w-4 h-4 text-gray-600" />
                                        </button>
                                    </div>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-600">Nama Lengkap</span>
                                    <span className="font-medium">{profile?.full_name}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-600">Email</span>
                                    <span className="font-medium">{profile?.email}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-600">Nomor HP</span>
                                    <span className="font-medium">{profile?.phone}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-600">Status</span>
                                    <span className="font-medium capitalize text-blue-600">{profile?.status}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </Card>

                {/* Wallet IDs */}
                <Card>
                    <h3 className="font-semibold text-gray-900 mb-4">Wallet IDs</h3>
                    <div className="space-y-3">
                        {profile?.wallets.map((wallet) => (
                            <div key={wallet.wallet_id} className="p-3 bg-gray-50 rounded-lg">
                                <div className="flex items-start justify-between gap-3">
                                    <div className="flex-1">
                                        <p className="text-sm font-medium text-gray-900 capitalize">{wallet.wallet_type} Wallet</p>
                                        <p className="text-xs font-mono text-gray-600 mt-1 break-all">{wallet.wallet_id}</p>
                                    </div>
                                    <button
                                        onClick={() => copyToClipboard(wallet.wallet_id)}
                                        className="p-1 hover:bg-gray-200 rounded transition-colors shrink-0"
                                        title="Copy Wallet ID"
                                    >
                                        <Copy className="w-4 h-4 text-gray-600" />
                                    </button>
                                </div>
                                <p className="text-xs text-gray-500 mt-2">
                                    Status: <span className="capitalize font-medium">{wallet.status}</span>
                                </p>
                            </div>
                        ))}
                    </div>
                </Card>

                {/* PIN Settings */}
                <Card>
                    <div className="flex items-start gap-4 mb-4">
                        <div className="bg-blue-100 p-3 rounded-full">
                            <Lock className="w-6 h-6 text-blue-600" />
                        </div>
                        <div className="flex-1">
                            <h3 className="font-semibold text-gray-900">PIN Transaksi</h3>
                            <p className="text-sm text-gray-600 mt-1">
                                {profile?.has_pin
                                    ? 'PIN sudah diatur. Anda dapat mengubah PIN kapan saja.'
                                    : 'PIN digunakan untuk mengamankan transaksi Anda.'}
                            </p>
                        </div>
                    </div>

                    <form onSubmit={handleSubmit(onSubmitPIN)} className="space-y-4 mt-4">
                        <Input
                            label="PIN Baru (6 digit)"
                            type="password"
                            maxLength={6}
                            placeholder="123456"
                            {...register('pin', {
                                required: 'PIN wajib diisi',
                                pattern: {
                                    value: /^\d{6}$/,
                                    message: 'PIN harus 6 digit angka',
                                },
                            })}
                            error={errors.pin?.message}
                        />

                        <Input
                            label="Konfirmasi PIN"
                            type="password"
                            maxLength={6}
                            placeholder="123456"
                            {...register('confirmPin', {
                                required: 'Konfirmasi PIN wajib diisi',
                                validate: (value) => value === pin || 'PIN tidak cocok',
                            })}
                            error={errors.confirmPin?.message}
                        />
                        <Button
                            type="submit"
                            variant="primary"
                            className="w-full"
                            isLoading={isSettingPIN}
                        >
                            {profile?.has_pin ? 'Ubah PIN' : 'Atur PIN'}
                        </Button>
                    </form>
                </Card>

                {/* Security Info */}
                <Card className="bg-blue-50 border border-blue-200">
                    <div className="flex items-start gap-3">
                        <Shield className="w-5 h-5 text-blue-600 mt-0.5" />
                        <div>
                            <h3 className="font-medium text-blue-900">Tips Keamanan</h3>
                            <ul className="text-sm text-blue-700 mt-2 space-y-1 list-disc list-inside">
                                <li>Gunakan PIN yang unik dan mudah diingat</li>
                                <li>Jangan bagikan PIN kepada siapapun</li>
                                <li>Ubah PIN secara berkala untuk keamanan</li>
                                <li>Logout setelah selesai menggunakan aplikasi</li>
                            </ul>
                        </div>
                    </div>
                </Card>
            </div>
        </MainLayout>
    );
};