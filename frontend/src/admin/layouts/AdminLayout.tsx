import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { AdminSidebar } from './AdminSidebar';
import { AdminNavbar } from './AdminNavbar';
import { useAdminAuth } from '@/admin/hooks/useAdminAuth';

interface AdminLayoutProps {
    children: React.ReactNode;
}

export const AdminLayout: React.FC<AdminLayoutProps> = ({ children }) => {
    const navigate = useNavigate();
    const { isAuthenticated, isLoading } = useAdminAuth();

    useEffect(() => {
        if (!isLoading && !isAuthenticated) {
            navigate('/admin/login');
        }
    }, [isAuthenticated, isLoading, navigate]);

    if (isLoading) {
        return (
            <div className="min-h-screen bg-white flex items-center justify-center">
                <div className="text-gray-900">Loading...</div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return null;
    }

    return (
        <div className="h-screen overflow-hidden bg-gray-50 flex">
            <AdminSidebar />

            <div className="flex-1 flex flex-col min-h-0">
                <div className="shrink-0">
                    <AdminNavbar />
                </div>
                <main className="flex-1 min-h-0 p-6 overflow-auto">
                    <div className="max-w-7xl mx-auto">
                        {children}
                    </div>
                </main>
            </div>
        </div>
    );
};