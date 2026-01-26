# Bayarin - Digital Wallet & Payment Gateway Simulator

Production-ready digital wallet and payment gateway simulator built with Golang, PostgreSQL, and Redis.

## 🚀 Features

- ✅ User Authentication (JWT)
- ✅ Digital Wallet Management
- ✅ Topup via Multiple Channels
- ✅ Transfer Between Users
- ✅ Transaction History
- ✅ Double-Entry Bookkeeping
- ✅ ACID Compliance
- ✅ Idempotency Support
- ✅ PIN Protection

## 🛠️ Tech Stack

**Backend:**
- Golang 1.22+
- PostgreSQL 14+
- Redis 6+
- Gin Framework
- SQLX
- JWT Authentication

**Frontend:**
- React + TypeScript
- Vite
- TailwindCSS
- Axios
- React Query

## 📋 Prerequisites

- Go 1.22 or higher
- PostgreSQL 14 or higher
- Redis 6 or higher
- golang-migrate (for migrations)

## 🔧 Installation

### 1. Clone Repository
```bash
git clone https://github.com/yourusername/bayarin.git
cd bayarin
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 4. Setup Database
```bash
# Create database and user
make db-setup

# Run migrations
make migrate-up
```

### 5. Run Application
```bash
make run
```

Server will start at `http://localhost:8080`

## 📖 API Documentation

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "phone": "081234567890",
  "full_name": "John Doe",
  "password": "password123"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "identifier": "user@example.com",
  "password": "password123"
}
```

### Wallet

#### Get Balance
```http
GET /api/v1/wallet/balance?type=main
Authorization: Bearer <token>
```

#### Get All Wallets
```http
GET /api/v1/wallet/all
Authorization: Bearer <token>
```

### Transaction

#### Topup
```http
POST /api/v1/transaction/topup
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": 10000000,
  "channel_code": "BCA_VA",
  "idempotency_key": "unique-key-123"
}
```

#### Transfer
```http
POST /api/v1/transaction/transfer
Authorization: Bearer <token>
Content-Type: application/json

{
  "to_user_id": "uuid-here",
  "amount": 5000000,
  "description": "Payment",
  "pin": "123456",
  "idempotency_key": "unique-key-456"
}
```

#### Get Transaction History
```http
GET /api/v1/transaction/history?limit=20&offset=0
Authorization: Bearer <token>
```

## 🧪 Testing
```bash
# Run all tests
make test

# Run with coverage
go test -v -cover ./...
```

## 📊 Database Schema

### Money Handling
- All amounts stored as **INTEGER** (minor unit)
- Example: Rp 100.000 = 10000000 (in cents)
- NO floating point for money

### Key Tables
- `users` - User accounts
- `wallets` - User wallets (balance)
- `transactions` - All transactions
- `ledger_entries` - Double-entry bookkeeping

## 🔐 Security Features

- JWT authentication
- Password hashing (bcrypt)
- PIN protection for transactions
- Row-level locking (prevent race conditions)
- Idempotency keys (prevent duplicates)

## 🎯 Architecture
```
cmd/
  api/
    main.go           # Entry point
internal/
  config/             # Configuration
  domain/             # Domain entities
  repository/         # Data access layer
  usecase/            # Business logic
  handler/            # HTTP handlers
  middleware/         # HTTP middleware
  pkg/                # Shared packages
migrations/           # Database migrations
```

## 📝 Make Commands
```bash
make db-setup       # Setup database
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations
make db-reset       # Reset database
make run            # Run application
make build          # Build binary
make test           # Run tests
```

## 🚧 Roadmap

- [ ] QR Code Payment
- [ ] Payment Gateway Integration
- [ ] Merchant Dashboard
- [ ] Settlement System
- [ ] Admin Panel
- [ ] Webhook Support

## 📄 License

MIT License

## 👨‍💻 Author

Your Name - [GitHub](https://github.com/yourusername)

---

**Note:** This is a simulator for portfolio purposes. Not for production use with real money.