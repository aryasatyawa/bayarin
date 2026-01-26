import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';
import { authApi } from '@/api/auth.api';
import { RegisterRequest } from '@/types/auth.types';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { storage } from '@/utils/storage';

export const RegisterForm: React.FC = () => {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<RegisterRequest>();

    const onSubmit = async (data: RegisterRequest) => {
        setIsLoading(true);
        try {
            const response = await authApi.register(data);

            // Save token and redirect
            storage.setToken(response.token);
            storage.setUser({
                id: response.user_id,
                email: response.email,
                phone: response.phone,
            });

            toast.success('Registrasi berhasil!');
            navigate('/dashboard');
        } catch (error: any) {
            const errorMessage = error.response?.data?.error?.message || 'Registrasi gagal';
            toast.error(errorMessage);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
                label="Nama Lengkap"
                type="text"
                placeholder="John Doe"
                {...register('full_name', {
                    required: 'Nama lengkap wajib diisi',
                    minLength: {
                        value: 3,
                        message: 'Nama minimal 3 karakter',
                    },
                })}
                error={errors.full_name?.message}
            />

            <Input
                label="Email"
                type="email"
                placeholder="user@example.com"
                {...register('email', {
                    required: 'Email wajib diisi',
                    pattern: {
                        value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                        message: 'Format email tidak valid',
                    },
                })}
                error={errors.email?.message}
            />

            <Input
                label="Nomor HP"
                type="tel"
                placeholder="081234567890"
                {...register('phone', {
                    required: 'Nomor HP wajib diisi',
                    pattern: {
                        value: /^(\+62|62|0)[0-9]{9,12}$/,
                        message: 'Format nomor HP tidak valid',
                    },
                })}
                error={errors.phone?.message}
                helperText="Contoh: 081234567890"
            />

            <Input
                label="Password"
                type="password"
                placeholder="Minimal 8 karakter"
                {...register('password', {
                    required: 'Password wajib diisi',
                    minLength: {
                        value: 8,
                        message: 'Password minimal 8 karakter',
                    },
                })}
                error={errors.password?.message}
            />

            <Button
                type="submit"
                variant="primary"
                className="w-full"
                isLoading={isLoading}
            >
                Daftar
            </Button>

            <p className="text-center text-sm text-gray-600">
                Sudah punya akun?{' '}
                <a href="/login" className="text-blue-600 hover:text-blue-700 font-medium">
                    Login di sini
                </a>
            </p>
        </form>
    );
};