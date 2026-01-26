import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Home, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/common/Button';

export const NotFoundPage: React.FC = () => {
    const navigate = useNavigate();

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
            <div className="text-center">
                <h1 className="text-9xl font-bold text-blue-600">404</h1>
                <h2 className="text-3xl font-bold text-gray-900 mt-4">Halaman Tidak Ditemukan</h2>
                <p className="text-gray-600 mt-2 mb-8">
                    Halaman yang Anda cari tidak ada atau telah dipindahkan.
                </p>

                <div className="flex gap-4 justify-center">
                    <Button
                        variant="secondary"
                        onClick={() => navigate(-1)}
                        className="flex items-center gap-2"
                    >
                        <ArrowLeft className="w-5 h-5" />
                        Kembali
                    </Button>

                    <Button
                        variant="primary"
                        onClick={() => navigate('/dashboard')}
                        className="flex items-center gap-2"
                    >
                        <Home className="w-5 h-5" />
                        Ke Dashboard
                    </Button>
                </div>
            </div>
        </div>
    );
};