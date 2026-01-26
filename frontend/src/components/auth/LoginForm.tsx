import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { authApi } from '@/api/auth.api';
import { LoginRequest } from '@/types/auth.types';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { storage } from '@/utils/storage';

export const LoginForm: React.FC = () => {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const [isLoading, setIsLoading] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginRequest>();

    const onSubmit = async (data: LoginRequest) => {
        setIsLoading(true);
        try {
            // Clear any existing cache and storage before login
            storage.clear();
            queryClient.clear();

            const response = await authApi.login(data);

            // Save token and user data
            storage.setToken(response.token);
            storage.setUser({
                id: response.user_id,
                email: response.email,
            });

            toast.success('Login berhasil!');
            navigate('/dashboard');
        } catch (error: any) {
            const errorMessage = error.response?.data?.error?.message || 'Login gagal';
            toast.error(errorMessage);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
                label="Email atau Nomor HP"
                type="text"
                placeholder="user@example.com atau 081234567890"
                {...register('identifier', {
                    required: 'Email atau nomor HP wajib diisi',
                })}
                error={errors.identifier?.message}
            />

            <Input
                label="Password"
                type="password"
                placeholder="Masukkan password"
                {...register('password', {
                    required: 'Password wajib diisi',
                })}
                error={errors.password?.message}
            />

            <Button
                type="submit"
                variant="primary"
                className="w-full"
                isLoading={isLoading}
            >
                Login
            </Button>

            <p className="text-center text-sm text-gray-600">
                Belum punya akun?{' '}
                <a href="/register" className="text-blue-600 hover:text-blue-700 font-medium">
                    Daftar di sini
                </a>
            </p>
        </form>
    );
};