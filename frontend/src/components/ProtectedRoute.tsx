import React from 'react';
import { Navigate } from 'react-router-dom';
import { storage } from '@/utils/storage';

interface ProtectedRouteProps {
    children: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const token = storage.getToken();

    if (!token) {
        return <Navigate to="/login" replace />;
    }

    return <>{children}</>;
};