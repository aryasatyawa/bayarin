import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { v4 as uuidv4 } from 'uuid';
import { transactionApi } from '@/api/transaction.api';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { toMinorUnit } from '@/utils/currency';

interface TransferFormData {
    to_user_id: string;
    amount: number;
    description: string;
    pin: string;
}

interface TransferFormProps {
    onSuccess?: () => void;
}

export const TransferForm: React.FC<TransferFormProps> = ({ onSuccess }) => {
    const [isLoading, setIsLoading] = useState(false);
    const queryClient = useQueryClient();

    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<TransferFormData>();

    const onSubmit = async (data: TransferFormData) => {
        setIsLoading(true);
        try {
            // Convert to minor unit
            const amountInMinorUnit = toMinorUnit(data.amount);

            await transactionApi.transfer({
                to_user_id: data.to_user_id,
                amount: amountInMinorUnit,
                description: data.description,
                pin: data.pin,
                idempotency_key: uuidv4(),
            });

            toast.success('Transfer berhasil!');
            reset();

            // Invalidate profile and transactions queries to refresh data
            await queryClient.invalidateQueries({ queryKey: ['profile'] });
            await queryClient.invalidateQueries({ queryKey: ['transactions'] });

            onSuccess?.();
        } catch (error: any) {
            const errorMessage = error.response?.data?.error?.message || 'Transfer gagal';
            toast.error(errorMessage);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
                label="User ID Tujuan"
                type="text"
                placeholder="UUID pengguna tujuan"
                {...register('to_user_id', {
                    required: 'User ID wajib diisi',
                })}
                error={errors.to_user_id?.message}
            />

            <Input
                label="Jumlah Transfer"
                type="number"
                placeholder="50000"
                {...register('amount', {
                    required: 'Jumlah wajib diisi',
                    min: {
                        value: 1000,
                        message: 'Minimal transfer Rp 1.000',
                    },
                })}
                error={errors.amount?.message}
            />

            <Input
                label="Catatan (Opsional)"
                type="text"
                placeholder="Catatan transfer"
                {...register('description')}
            />

            <Input
                label="PIN Transaksi"
                type="password"
                placeholder="6 digit PIN"
                maxLength={6}
                {...register('pin', {
                    required: 'PIN wajib diisi',
                    pattern: {
                        value: /^\d{6}$/,
                        message: 'PIN harus 6 digit angka',
                    },
                })}
                error={errors.pin?.message}
            />

            <Button type="submit" variant="primary" className="w-full" isLoading={isLoading}>
                Transfer Sekarang
            </Button>
        </form>
    );
};