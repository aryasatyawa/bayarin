import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';

// User Components
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { PublicRoute } from '@/components/PublicRoute';

// User Pages
import { LoginPage } from '@/pages/LoginPage';
import { RegisterPage } from '@/pages/RegisterPage';
import { DashboardPage } from '@/pages/DashboardPage';
import { WalletPage } from '@/pages/WalletPage';
import { TransferPage } from '@/pages/TransferPage';
import { HistoryPage } from '@/pages/HistoryPage';
import { SettingsPage } from '@/pages/SettingsPage';
import { ProfilePage } from '@/pages/ProfilePage';
import { NotFoundPage } from '@/pages/NotFoundPage';
import { InformationPage } from '@/pages/InformationPage';

// Admin Pages
import { AdminLoginPage } from '@/admin/pages/AdminLoginPage';
import { AdminDashboardPage } from '@/admin/pages/AdminDashboardPage';
import { LedgerPage } from '@/admin/pages/LedgerPage';
import { TransactionsPage } from '@/admin/pages/TransactionsPage';
import { UsersPage } from '@/admin/pages/UsersPage';
import { RefundPage } from '@/admin/pages/RefundPage';

// Create QueryClient
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* ============================================ */}
          {/* USER ROUTES */}
          {/* ============================================ */}
          {/* Default redirect */}
          <Route path="/" element={<Navigate to="/information" replace />} />

          {/* Information/Landing Page */}
          <Route path="/information" element={<InformationPage />} />

          {/* Public routes */}
          <Route
            path="/login"
            element={
              <PublicRoute>
                <LoginPage />
              </PublicRoute>
            }
          />
          <Route
            path="/register"
            element={
              <PublicRoute>
                <RegisterPage />
              </PublicRoute>
            }
          />

          {/* Protected user routes */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/wallet"
            element={
              <ProtectedRoute>
                <WalletPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/transfer"
            element={
              <ProtectedRoute>
                <TransferPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/history"
            element={
              <ProtectedRoute>
                <HistoryPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/settings"
            element={
              <ProtectedRoute>
                <SettingsPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            }
          />

          {/* ============================================ */}
          {/* ADMIN ROUTES */}
          {/* ============================================ */}
          <Route path="/admin/login" element={<AdminLoginPage />} />
          <Route path="/admin/dashboard" element={<AdminDashboardPage />} />
          <Route path="/admin/ledger" element={<LedgerPage />} />
          <Route path="/admin/transactions" element={<TransactionsPage />} />
          <Route path="/admin/users" element={<UsersPage />} />
          <Route path="/admin/refund" element={<RefundPage />} />

          {/* 404 */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>

      {/* Toast notifications */}
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 3000,
          style: {
            background: '#363636',
            color: '#fff',
          },
          success: {
            duration: 3000,
            iconTheme: {
              primary: '#10b981',
              secondary: '#fff',
            },
          },
          error: {
            duration: 4000,
            iconTheme: {
              primary: '#ef4444',
              secondary: '#fff',
            },
          },
        }}
      />
    </QueryClientProvider>
  );
};

export default App;