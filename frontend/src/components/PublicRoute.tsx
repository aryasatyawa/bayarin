import React from 'react';
import { Navigate } from 'react-router-dom';
import { storage } from '@/utils/storage';

interface PublicRouteProps {
    children: React.ReactNode;
}

export const PublicRoute: React.FC<PublicRouteProps> = ({ children }) => {
    const token = storage.getToken();

    if (token) {
        return <Navigate to="/dashboard" replace />;
    }

    return <>{children}</>;
};