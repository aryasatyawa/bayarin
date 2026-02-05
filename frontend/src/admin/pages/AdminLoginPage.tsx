import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';
import { Shield, Lock, User } from 'lucide-react';
import { adminAuthApi } from '@/admin/api/admin-auth.api';
import { AdminLoginRequest } from '@/admin/types/admin.types';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { useAdminAuth } from '@/admin/hooks/useAdminAuth';

export const AdminLoginPage: React.FC = () => {
    const navigate = useNavigate();
    const { login } = useAdminAuth();
    const [isLoading, setIsLoading] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<AdminLoginRequest>();

    const onSubmit = async (data: AdminLoginRequest) => {
        setIsLoading(true);
        try {
            const response = await adminAuthApi.login(data);

            // Save to localStorage
            login(response.token, {
                admin_id: response.admin_id,
                username: response.username,
                full_name: response.full_name,
                role: response.role,
            });

            toast.success('Login berhasil!');
            navigate('/admin/dashboard');
        } catch (error: any) {
            const errorMessage = error.response?.data?.error?.message || 'Login gagal';
            toast.error(errorMessage);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-linear-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center p-4">
            <div className="w-full max-w-md">
                {/* Logo */}
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-20 h-20 bg-linear-to-br from-blue-500 to-blue-700 rounded-2xl shadow-2xl mb-4">
                        <Shield className="w-10 h-10 text-white" />
                    </div>
                    <h1 className="text-3xl font-bold text-white mb-2">Bayarin Admin</h1>
                    <p className="text-gray-400">Admin Dashboard & Management</p>
                </div>

                {/* Login Card */}
                <div className="bg-gray-800 rounded-2xl shadow-2xl p-8 border border-gray-700">
                    <div className="mb-6">
                        <h2 className="text-2xl font-bold text-white">Admin Login</h2>
                        <p className="text-gray-400 text-sm mt-1">
                            Masuk dengan akun admin Anda
                        </p>
                    </div>

                    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                <User className="w-4 h-4 inline mr-2" />
                                Username
                            </label>
                            <Input
                                type="text"
                                placeholder="superadmin"
                                {...register('username', {
                                    required: 'Username wajib diisi',
                                })}
                                error={errors.username?.message}
                                className="bg-gray-700 text-white border-gray-600 focus:ring-blue-500"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                <Lock className="w-4 h-4 inline mr-2" />
                                Password
                            </label>
                            <Input
                                type="password"
                                placeholder="••••••••"
                                {...register('password', {
                                    required: 'Password wajib diisi',
                                })}
                                error={errors.password?.message}
                                className="bg-gray-700 text-white border-gray-600 focus:ring-blue-500"
                            />
                        </div>

                        <Button
                            type="submit"
                            variant="primary"
                            className="w-full bg-blue-600 hover:bg-blue-700"
                            isLoading={isLoading}
                        >
                            Login
                        </Button>
                    </form>

                    {/* Info */}
                    <div className="mt-6 p-4 bg-gray-700 rounded-lg border border-gray-600">
                        <p className="text-xs text-gray-400">
                            🔒 Akses terbatas hanya untuk admin. Login menggunakan kredensial admin yang telah diberikan.
                        </p>
                    </div>
                </div>

                {/* Footer */}
                <p className="text-center text-gray-500 text-sm mt-6">
                    &copy; 2024 Bayarin. Admin Panel v1.0
                </p>
            </div>
        </div>
    );
};