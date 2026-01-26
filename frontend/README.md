# Bayarin Frontend

React + TypeScript frontend untuk Bayarin Digital Wallet & Payment Gateway.

## 🚀 Quick Start

### 1. Install Dependencies
```bash
npm install
```

### 2. Setup Environment
```bash
# Copy .env example
cp .env.example .env

# Edit .env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 3. Run Development Server
```bash
npm run dev
```

Frontend akan berjalan di: http://localhost:3000

## 🛠️ Tech Stack

- **React 18** - UI Library
- **TypeScript** - Type Safety
- **Vite** - Build Tool
- **TailwindCSS** - Styling
- **React Router** - Routing
- **React Query** - Data Fetching
- **React Hook Form** - Form Management
- **Axios** - HTTP Client
- **Lucide React** - Icons
- **React Hot Toast** - Notifications
- **Date-fns** - Date Formatting

## 📁 Project Structure
```
src/
├── api/              # API clients
├── components/       # Reusable components
│   ├── auth/        # Authentication components
│   ├── wallet/      # Wallet components
│   ├── transaction/ # Transaction components
│   ├── layout/      # Layout components
│   └── common/      # Common UI components
├── pages/           # Page components
├── hooks/           # Custom hooks
├── types/           # TypeScript types
├── utils/           # Utility functions
├── App.tsx          # Main app component
└── main.tsx         # Entry point
```

## 🎨 Features

- ✅ User Authentication (Register/Login)
- ✅ JWT Token Management
- ✅ Wallet Balance Display
- ✅ Topup Wallet
- ✅ Transfer to Other Users
- ✅ Transaction History
- ✅ PIN Management
- ✅ Responsive Design (Mobile & Desktop)
- ✅ Protected Routes
- ✅ Toast Notifications
- ✅ Form Validation
- ✅ Error Handling

## 📱 Pages

- `/login` - Login page
- `/register` - Register page
- `/dashboard` - Dashboard with overview
- `/wallet` - Wallet management
- `/transfer` - Transfer form
- `/history` - Transaction history
- `/settings` - Settings & PIN management
- `/profile` - User profile

## 🔧 Available Scripts
```bash
npm run dev      # Start development server
npm run build    # Build for production
npm run preview  # Preview production build
npm run lint     # Run ESLint
```

## 🌐 API Integration

Backend API: http://localhost:8080/api/v1

Endpoints:
- `POST /auth/register` - Register user
- `POST /auth/login` - Login user
- `GET /user/profile` - Get profile
- `POST /user/pin` - Set PIN
- `GET /wallet/balance` - Get balance
- `GET /wallet/all` - Get all wallets
- `POST /transaction/topup` - Topup
- `POST /transaction/transfer` - Transfer
- `GET /transaction/history` - Get history

## 💡 Tips

1. **Currency Format**: Semua amount di frontend dalam **major unit** (Rupiah), akan dikonversi ke **minor unit** (sen) saat kirim ke backend.

2. **Idempotency Key**: Setiap transaksi menggunakan UUID sebagai idempotency key untuk mencegah duplikasi.

3. **PIN Validation**: PIN harus 6 digit angka dan divalidasi di frontend sebelum dikirim ke backend.

## 🔐 Security

- JWT token disimpan di localStorage
- Token dikirim via Authorization header
- Auto redirect ke login jika token expired
- PIN validation untuk transaksi sensitif

## 📄 License

MIT License
```

---

## ✅ Final Project Structure
```
bayarin-fe/
├── public/
│   └── vite.svg
├── src/
│   ├── api/
│   │   ├── axios.ts
│   │   ├── auth.api.ts
│   │   ├── wallet.api.ts
│   │   └── transaction.api.ts
│   ├── components/
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   └── RegisterForm.tsx
│   │   ├── wallet/
│   │   │   ├── WalletCard.tsx
│   │   │   └── WalletHistory.tsx
│   │   ├── transaction/
│   │   │   ├── TopupForm.tsx
│   │   │   ├── TransferForm.tsx
│   │   │   └── TransactionList.tsx
│   │   ├── layout/
│   │   │   ├── Navbar.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── BottomNav.tsx
│   │   │   └── MainLayout.tsx
│   │   ├── common/
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Modal.tsx
│   │   │   └── Card.tsx
│   │   ├── ProtectedRoute.tsx
│   │   └── PublicRoute.tsx
│   ├── pages/
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── DashboardPage.tsx
│   │   ├── WalletPage.tsx
│   │   ├── TransferPage.tsx
│   │   ├── HistoryPage.tsx
│   │   ├── SettingsPage.tsx
│   │   ├── ProfilePage.tsx
│   │   └── NotFoundPage.tsx
│   ├── types/
│   │   ├── auth.types.ts
│   │   ├── wallet.types.ts
│   │   └── transaction.types.ts
│   ├── utils/
│   │   ├── currency.ts
│   │   ├── storage.ts
│   │   └── date.ts
│   ├── App.tsx
│   ├── main.tsx
│   ├── index.css
│   └── vite-env.d.ts
├── .env
├── .gitignore
├── index.html
├── package.json
├── tailwind.config.js
├── tsconfig.json
├── tsconfig.node.json
├── vite.config.ts
└── README.md