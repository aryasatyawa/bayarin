import React from 'react';
import { Wallet } from 'lucide-react';
import { RegisterForm } from '@/components/auth/RegisterForm';

export const RegisterPage: React.FC = () => {
    return (
        <div className="min-h-screen bg-linear-to-br from-blue-600 to-blue-800 flex items-center justify-center p-4">
            <div className="w-full max-w-md">
                {/* Logo */}
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-16 h-16 bg-white rounded-2xl shadow-lg mb-4">
                        <Wallet className="w-8 h-8 text-blue-600" />
                    </div>
                    <h1 className="text-3xl font-bold text-white mb-2">Bayarin</h1>
                    <p className="text-blue-100">Digital Wallet & Payment Gateway</p>
                </div>

                {/* Register Card */}
                <div className="bg-white rounded-2xl shadow-xl p-8">
                    <h2 className="text-2xl font-bold text-gray-900 mb-6">Daftar Akun</h2>
                    <RegisterForm />
                </div>

                {/* Footer */}
                <p className="text-center text-blue-100 text-sm mt-6">
                    &copy; 2024 Bayarin. All rights reserved.
                </p>
            </div>
        </div>
    );
};