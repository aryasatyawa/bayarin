import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
    LayoutDashboard,
    BookOpen,
    ArrowLeftRight,
    Users,
    DollarSign,
    Shield,
    FileText,
    Settings,
} from 'lucide-react';
import { useAdminAuth } from '@/admin/hooks/useAdminAuth';

export const AdminSidebar: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { admin, canRefund, canFreezeWallet } = useAdminAuth();

    const menuItems = [
        {
            label: 'Dashboard',
            icon: LayoutDashboard,
            path: '/admin/dashboard',
            roles: ['super_admin', 'ops_admin', 'finance_admin'],
        },
        {
            label: 'Ledger',
            icon: BookOpen,
            path: '/admin/ledger',
            roles: ['super_admin', 'ops_admin', 'finance_admin'],
        },
        {
            label: 'Transactions',
            icon: ArrowLeftRight,
            path: '/admin/transactions',
            roles: ['super_admin', 'ops_admin', 'finance_admin'],
        },
        {
            label: 'Users',
            icon: Users,
            path: '/admin/users',
            roles: ['super_admin', 'ops_admin', 'finance_admin'],
        },
        {
            label: 'Refund',
            icon: DollarSign,
            path: '/admin/refund',
            roles: ['super_admin', 'finance_admin'],
            show: canRefund(),
        },
        {
            label: 'Admins',
            icon: Shield,
            path: '/admin/admins',
            roles: ['super_admin'],
            show: admin?.role === 'super_admin',
        },
        {
            label: 'Audit Logs',
            icon: FileText,
            path: '/admin/audit-logs',
            roles: ['super_admin', 'ops_admin', 'finance_admin'],
        },
    ];

    const isActive = (path: string) => location.pathname === path;

    return (
        <aside className="flex flex-col h-screen sticky top-0 shrink-0 w-64 bg-white border-r border-gray-200">
            {/* Logo */}
            <div className="shrink-0 p-6 border-b border-gray-200">
                <div className="flex items-center gap-3">
                    <div className="bg-blue-600 p-2 rounded-lg">
                        <Shield className="w-6 h-6 text-white" />
                    </div>
                    <div>
                        <h1 className="text-gray-900 font-bold text-lg">Bayarin</h1>
                        <p className="text-gray-500 text-xs">Admin Panel</p>
                    </div>
                </div>
            </div>

            {/* Navigation */}
            <nav className="flex-1 overflow-y-auto min-h-0 p-4 space-y-2">
                {menuItems.map((item) => {
                    // Check if item should be shown
                    if (item.show === false) return null;

                    const Icon = item.icon;
                    const active = isActive(item.path);

                    return (
                        <button
                            key={item.path}
                            onClick={() => navigate(item.path)}
                            className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${active
                                    ? 'bg-blue-600 text-white'
                                    : 'text-gray-600 hover:bg-gray-100'
                                }`}
                        >
                            <Icon className="w-5 h-5" />
                            <span className="font-medium">{item.label}</span>
                        </button>
                    );
                })}
            </nav>

            {/* Admin Info */}
            <div className="mt-auto shrink-0 p-4 border-t border-gray-200 bg-white">
                <div className="flex items-center gap-3">
                    <div className="bg-gray-100 p-2 rounded-full">
                        <Settings className="w-5 h-5 text-gray-600" />
                    </div>
                    <div className="flex-1 min-w-0">
                        <p className="text-gray-900 text-sm font-medium truncate">
                            {admin?.full_name}
                        </p>
                        <p className="text-gray-500 text-xs capitalize">{admin?.role?.replace('_', ' ')}</p>
                    </div>
                </div>
            </div>
        </aside>
    );
};