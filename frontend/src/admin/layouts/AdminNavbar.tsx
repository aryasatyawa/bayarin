import React from 'react';
import { LogOut, Bell } from 'lucide-react';
import { useAdminAuth } from '@/admin/hooks/useAdminAuth';
import { Button } from '@/components/common/Button';
import { toast } from 'react-hot-toast';

export const AdminNavbar: React.FC = () => {
    const { admin, logout } = useAdminAuth();

    const handleLogout = () => {
        logout();
        toast.success('Berhasil logout');
    };

    return (
        <nav className="bg-white border-b border-gray-200 sticky top-0 z-40">
            <div className="px-6 py-4">
                <div className="flex justify-between items-center">
                    {/* Left */}
                    <div>
                        <h2 className="text-gray-900 font-semibold text-lg">Admin Dashboard</h2>
                        <p className="text-gray-500 text-sm">Kelola sistem Bayarin</p>
                    </div>

                    {/* Right */}
                    <div className="flex items-center gap-4">
                        {/* Notifications */}
                        <button className="relative p-2 text-gray-500 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors">
                            <Bell className="w-5 h-5" />
                            <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
                        </button>

                        {/* User Info */}
                        <div className="hidden sm:flex items-center gap-3 px-4 py-2 bg-gray-100 rounded-lg">
                            <div className="text-right">
                                <p className="text-gray-900 text-sm font-medium">{admin?.full_name}</p>
                                <p className="text-gray-500 text-xs capitalize">
                                    {admin?.role?.replace('_', ' ')}
                                </p>
                            </div>
                        </div>

                        {/* Logout */}
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={handleLogout}
                            className="text-red-600 hover:text-red-700 hover:bg-red-50"
                        >
                            <LogOut className="w-5 h-5" />
                        </Button>
                    </div>
                </div>
            </div>
        </nav>
    );
};