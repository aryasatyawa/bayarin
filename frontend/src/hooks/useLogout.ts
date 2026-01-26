import { useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { storage } from '@/utils/storage';

export const useLogout = () => {
    const navigate = useNavigate();
    const queryClient = useQueryClient();

    const logout = (showToast: boolean = true) => {
        try {
            // Clear local storage
            storage.clear();

            // Clear React Query cache
            queryClient.clear();

            if (showToast) {
                toast.success('Berhasil logout');
            }

            // Redirect to login
            navigate('/login');
        } catch (error) {
            console.error('Logout error:', error);
        }
    };

    return { logout };
};
