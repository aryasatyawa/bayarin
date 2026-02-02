import React from 'react';
import { useNavigate } from 'react-router-dom';
import { LogOut, User, Wallet } from 'lucide-react';
import { storage } from '@/utils/storage';
import { useLogout } from '@/hooks/useLogout';

export const Navbar: React.FC = () => {
    const navigate = useNavigate();
    const { logout } = useLogout();
    const user = storage.getUser();

    const handleLogout = () => {
        logout(true);
    };

    return (
        <nav className="bg-linear-to-r from-blue-600 via-blue-700 to-blue-800 sticky top-0 z-40 shadow-lg border-b border-blue-900">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                    {/* Logo */}
                    <div className="flex items-center gap-3 cursor-pointer hover:opacity-80 transition-opacity duration-200 group" onClick={() => navigate('/dashboard')}>
                        <div className="bg-white bg-opacity-20 p-2 rounded-lg group-hover:bg-opacity-30 transition-all duration-200 shadow-lg">
                            <Wallet className="w-6 h-6 text-white" />
                        </div>
                        <div>
                            <span className="text-2xl font-bold text-white block leading-none">Bayarin</span>
                            <span className="text-xs text-blue-100">Smart Wallet</span>
                        </div>
                    </div>

                    {/* User Menu */}
                    <div className="flex items-center gap-6">
                        <div className="hidden sm:block text-right">
                            <p className="text-sm font-semibold text-white">{user?.full_name || user?.email}</p>
                            <p className="text-xs text-blue-100">{user?.email}</p>
                        </div>

                        <div className="flex items-center gap-2 pl-6 border-l border-blue-400">
                            <button
                                onClick={() => navigate('/profile')}
                                className="p-2 rounded-lg hover:bg-blue-500 hover:bg-opacity-20 transition-all duration-200 hover:shadow-lg transform hover:scale-110 active:scale-95"
                                title="Profil"
                            >
                                <User className="w-5 h-5 text-white" />
                            </button>

                            <button
                                onClick={handleLogout}
                                className="p-2 rounded-lg hover:bg-red-500 hover:bg-opacity-30 transition-all duration-200 hover:shadow-lg transform hover:scale-110 active:scale-95"
                                title="Logout"
                            >
                                <LogOut className="w-5 h-5 text-white" />
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
    );
};