import React from 'react';
import { useNavigate } from 'react-router-dom';
import { LoginForm } from '@/components/auth/LoginForm';

export const LoginPage: React.FC = () => {
    const navigate = useNavigate();

    return (
        <div className="min-h-screen bg-linear-to-br from-white to-blue-800 flex items-center justify-center p-4">
            <div className="w-full max-w-md">
                {/* Logo */}
                <div className="text-center mb-8">
                    <h1
                        className="text-4xl font-bold text-white mb-2 tracking-tighter cursor-pointer hover:opacity-80 transition-opacity"
                        onClick={() => navigate('/information')}
                    >
                        Bayarin<span className="text-blue-600">.</span>
                    </h1>
                    <p className="text-blue-100">Digital Wallet & Payment Gateway</p>
                </div>

                {/* Login Card */}
                <div className="bg-white rounded-2xl shadow-xl p-8">
                    <h2 className="text-2xl font-bold text-gray-900 mb-6">Login</h2>
                    <LoginForm />
                </div>

                {/* Footer */}
                <p className="text-center text-blue-100 text-sm mt-6">
                    &copy; 2024 Bayarin. All rights reserved.
                </p>
            </div>
        </div>
    );
};