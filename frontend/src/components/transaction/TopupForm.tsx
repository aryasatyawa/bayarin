import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { v4 as uuidv4 } from 'uuid';
import { transactionApi } from '@/api/transaction.api';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { toMinorUnit } from '@/utils/currency';

interface TopupFormData {
    amount: number;
    channel_code: string;
}

interface TopupFormProps {
    onSuccess?: () => void;
}

export const TopupForm: React.FC<TopupFormProps> = ({ onSuccess }) => {
    const [isLoading, setIsLoading] = useState(false);
    const queryClient = useQueryClient();

    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<TopupFormData>();

    const channels = [
        { code: 'BCA_VA', name: 'BCA Virtual Account' },
        { code: 'MANDIRI_VA', name: 'Mandiri Virtual Account' },
        { code: 'BNI_VA', name: 'BNI Virtual Account' },
        { code: 'BRI_VA', name: 'BRI Virtual Account' },
        { code: 'OVO', name: 'OVO' },
        { code: 'GOPAY', name: 'GoPay' },
        { code: 'DANA', name: 'DANA' },
    ];

    const onSubmit = async (data: TopupFormData) => {
        setIsLoading(true);
        try {
            // Convert to minor unit (multiply by 100)
            const amountInMinorUnit = toMinorUnit(data.amount);

            await transactionApi.topup({
                amount: amountInMinorUnit,
                channel_code: data.channel_code,
                idempotency_key: uuidv4(),
            });

            toast.success('Topup berhasil!');
            reset();

            // Invalidate profile and transactions queries to refresh data
            await queryClient.invalidateQueries({ queryKey: ['profile'] });
            await queryClient.invalidateQueries({ queryKey: ['transactions'] });

            onSuccess?.();
        } catch (error: any) {
            const errorMessage = error.response?.data?.error?.message || 'Topup gagal';
            toast.error(errorMessage);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
                label="Jumlah Topup"
                type="number"
                placeholder="100000"
                {...register('amount', {
                    required: 'Jumlah wajib diisi',
                    min: {
                        value: 10000,
                        message: 'Minimal topup Rp 10.000',
                    },
                    max: {
                        value: 10000000,
                        message: 'Maksimal topup Rp 10.000.000',
                    },
                })}
                error={errors.amount?.message}
                helperText="Minimal Rp 10.000, Maksimal Rp 10.000.000"
            />

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    Metode Pembayaran <span className="text-red-500">*</span>
                </label>
                <select
                    className="input-field"
                    {...register('channel_code', {
                        required: 'Metode pembayaran wajib dipilih',
                    })}
                >
                    <option value="">Pilih metode pembayaran</option>
                    {channels.map((channel) => (
                        <option key={channel.code} value={channel.code}>
                            {channel.name}
                        </option>
                    ))}
                </select>
                {errors.channel_code && (
                    <p className="mt-1 text-sm text-red-600">{errors.channel_code.message}</p>
                )}
            </div>

            <Button type="submit" variant="primary" className="w-full" isLoading={isLoading}>
                Topup Sekarang
            </Button>
        </form>
    );
};