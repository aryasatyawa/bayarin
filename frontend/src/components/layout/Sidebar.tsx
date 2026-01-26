import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { LayoutDashboard, Wallet, ArrowLeftRight, History, Settings } from 'lucide-react';

export const Sidebar: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const menuItems = [
        {
            label: 'Dashboard',
            icon: LayoutDashboard,
            path: '/dashboard',
        },
        {
            label: 'Wallet',
            icon: Wallet,
            path: '/wallet',
        },
        {
            label: 'Transfer',
            icon: ArrowLeftRight,
            path: '/transfer',
        },
        {
            label: 'Riwayat',
            icon: History,
            path: '/history',
        },
        {
            label: 'Pengaturan',
            icon: Settings,
            path: '/settings',
        },
    ];

    const isActive = (path: string) => location.pathname === path;

    return (
        <aside className="hidden md:block w-64 bg-white border-r border-gray-200 min-h-screen">
            <nav className="p-4 space-y-2">
                {menuItems.map((item) => {
                    const Icon = item.icon;
                    const active = isActive(item.path);

                    return (
                        <button
                            key={item.path}
                            onClick={() => navigate(item.path)}
                            className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${active
                                ? 'bg-blue-50 text-blue-700 font-medium'
                                : 'text-gray-600 hover:bg-gray-50'
                                }`}
                        >
                            <Icon className="w-5 h-5" />
                            <span>{item.label}</span>
                        </button>
                    );
                })}
            </nav>
        </aside>
    );
};