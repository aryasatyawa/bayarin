import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { AdminRole } from '@/admin/types/admin.types';

interface AdminUser {
    admin_id: string;
    username: string;
    full_name: string;
    role: AdminRole;
}

export const useAdminAuth = () => {
    const navigate = useNavigate();
    const [admin, setAdmin] = useState<AdminUser | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Load admin from localStorage
        const token = localStorage.getItem('admin_token');
        const adminData = localStorage.getItem('admin_user');

        if (token && adminData) {
            setAdmin(JSON.parse(adminData));
        }

        setIsLoading(false);
    }, []);

    const login = (token: string, adminData: AdminUser) => {
        localStorage.setItem('admin_token', token);
        localStorage.setItem('admin_user', JSON.stringify(adminData));
        setAdmin(adminData);
    };

    const logout = () => {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user');
        setAdmin(null);
        navigate('/admin/login');
    };

    const hasPermission = (requiredRole: AdminRole): boolean => {
        if (!admin) return false;
        if (admin.role === 'super_admin') return true;
        return admin.role === requiredRole;
    };

    const canRefund = (): boolean => {
        return hasPermission('super_admin') || hasPermission('finance_admin');
    };

    const canFreezeWallet = (): boolean => {
        return hasPermission('super_admin') || hasPermission('ops_admin');
    };

    return {
        admin,
        isLoading,
        isAuthenticated: !!admin,
        login,
        logout,
        hasPermission,
        canRefund,
        canFreezeWallet,
    };
};