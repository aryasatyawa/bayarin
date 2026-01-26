import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { LayoutDashboard, Wallet, ArrowLeftRight, History, User } from 'lucide-react';

export const BottomNav: React.FC = () => {
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
            label: 'Profil',
            icon: User,
            path: '/profile',
        },
    ];

    const isActive = (path: string) => location.pathname === path;

    return (
        <div className="md:hidden fixed bottom-0 left-0 right-0 bg-linear-to-t from-white to-gray-50 border-t border-gray-300 z-40 shadow-xl">
            <nav className="flex justify-around items-center h-16">
                {menuItems.map((item) => {
                    const Icon = item.icon;
                    const active = isActive(item.path);

                    return (
                        <button
                            key={item.path}
                            onClick={() => navigate(item.path)}
                            className={`flex flex-col items-center justify-center flex-1 h-full transition-all duration-300 relative group ${active
                                    ? 'text-blue-600'
                                    : 'text-gray-500 hover:text-blue-500'
                                }`}
                        >
                            {active && (
                                <div className="absolute top-0 left-1/2 transform -translate-x-1/2 w-2 h-1 bg-blue-600 rounded-b-full"></div>
                            )}
                            <div className={`p-2 rounded-lg transition-all duration-300 ${active
                                    ? 'bg-blue-100 shadow-md'
                                    : 'group-hover:bg-gray-100'
                                }`}>
                                <Icon className="w-6 h-6 transition-transform duration-300 group-hover:scale-110" />
                            </div>
                            <span className={`text-xs mt-1 font-medium transition-colors duration-300 ${active ? 'text-blue-600' : 'text-gray-600'
                                }`}>
                                {item.label}
                            </span>
                        </button>
                    );
                })}
            </nav>
        </div>
    );
};