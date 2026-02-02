package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/jmoiron/sqlx"
)

type DashboardUsecase interface {
	GetOverview(ctx context.Context) (*DashboardOverview, error)
	GetDailyStats(ctx context.Context, date time.Time) (*DailyStats, error)
	GetTransactionSummary(ctx context.Context, startDate, endDate time.Time) (*TransactionSummary, error)
}

type dashboardUsecase struct {
	db              *sqlx.DB
	userRepo        repository.UserRepository
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
}

func NewDashboardUsecase(
	db *sqlx.DB,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		db:              db,
		userRepo:        userRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

// DTOs
type DashboardOverview struct {
	TotalUsers           int64 `json:"total_users"`
	TotalActiveWallets   int64 `json:"total_active_wallets"`
	TotalSystemLiability int64 `json:"total_system_liability"` // Total saldo semua user
	TodayTransactions    int64 `json:"today_transactions"`
	TodayVolume          int64 `json:"today_volume"`
	PendingTransactions  int64 `json:"pending_transactions"`
	FailedTransactions   int64 `json:"failed_transactions"`
	TodayTopups          int64 `json:"today_topups"`
	TodayTransfers       int64 `json:"today_transfers"`
}

type DailyStats struct {
	Date              time.Time              `json:"date"`
	TotalTransactions int64                  `json:"total_transactions"`
	TotalVolume       int64                  `json:"total_volume"`
	ByType            map[string]TypeStats   `json:"by_type"`
	ByStatus          map[string]StatusStats `json:"by_status"`
}

type TypeStats struct {
	Count  int64 `json:"count"`
	Volume int64 `json:"volume"`
}

type StatusStats struct {
	Count int64 `json:"count"`
}

type TransactionSummary struct {
	StartDate         time.Time              `json:"start_date"`
	EndDate           time.Time              `json:"end_date"`
	TotalTransactions int64                  `json:"total_transactions"`
	TotalVolume       int64                  `json:"total_volume"`
	ByType            map[string]TypeStats   `json:"by_type"`
	ByStatus          map[string]StatusStats `json:"by_status"`
	DailyBreakdown    []DailyBreakdown       `json:"daily_breakdown"`
}

type DailyBreakdown struct {
	Date   time.Time `json:"date"`
	Count  int64     `json:"count"`
	Volume int64     `json:"volume"`
}

// GetOverview returns dashboard overview metrics
func (uc *dashboardUsecase) GetOverview(ctx context.Context) (*DashboardOverview, error) {
	overview := &DashboardOverview{}

	// Get total users
	var totalUsers int64
	queryUsers := `SELECT COUNT(*) FROM users WHERE status = 'active'`
	if err := uc.db.GetContext(ctx, &totalUsers, queryUsers); err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}
	overview.TotalUsers = totalUsers

	// Get total active wallets & system liability
	var walletStats struct {
		TotalWallets int64 `db:"total_wallets"`
		TotalBalance int64 `db:"total_balance"`
	}
	queryWallets := `
		SELECT 
			COUNT(*) as total_wallets,
			COALESCE(SUM(balance), 0) as total_balance
		FROM wallets
		WHERE status = 'active'
	`
	if err := uc.db.GetContext(ctx, &walletStats, queryWallets); err != nil {
		return nil, fmt.Errorf("failed to get wallet stats: %w", err)
	}
	overview.TotalActiveWallets = walletStats.TotalWallets
	overview.TotalSystemLiability = walletStats.TotalBalance

	// Get today's transactions
	today := time.Now().Format("2006-01-02")

	var todayStats struct {
		TotalCount    int64 `db:"total_count"`
		TotalVolume   int64 `db:"total_volume"`
		TopupCount    int64 `db:"topup_count"`
		TransferCount int64 `db:"transfer_count"`
	}
	queryToday := `
		SELECT 
			COUNT(*) as total_count,
			COALESCE(SUM(amount), 0) as total_volume,
			COUNT(CASE WHEN transaction_type = 'topup' THEN 1 END) as topup_count,
			COUNT(CASE WHEN transaction_type = 'transfer' THEN 1 END) as transfer_count
		FROM transactions
		WHERE DATE(created_at) = $1 AND status = 'success'
	`
	if err := uc.db.GetContext(ctx, &todayStats, queryToday, today); err != nil {
		return nil, fmt.Errorf("failed to get today stats: %w", err)
	}
	overview.TodayTransactions = todayStats.TotalCount
	overview.TodayVolume = todayStats.TotalVolume
	overview.TodayTopups = todayStats.TopupCount
	overview.TodayTransfers = todayStats.TransferCount

	// Get pending transactions
	var pendingCount int64
	queryPending := `SELECT COUNT(*) FROM transactions WHERE status = 'pending'`
	if err := uc.db.GetContext(ctx, &pendingCount, queryPending); err != nil {
		return nil, fmt.Errorf("failed to get pending count: %w", err)
	}
	overview.PendingTransactions = pendingCount

	// Get failed transactions (last 7 days)
	var failedCount int64
	queryFailed := `
		SELECT COUNT(*) 
		FROM transactions 
		WHERE status = 'failed' 
		AND created_at >= NOW() - INTERVAL '7 days'
	`
	if err := uc.db.GetContext(ctx, &failedCount, queryFailed); err != nil {
		return nil, fmt.Errorf("failed to get failed count: %w", err)
	}
	overview.FailedTransactions = failedCount

	return overview, nil
}

// GetDailyStats returns daily transaction statistics
func (uc *dashboardUsecase) GetDailyStats(ctx context.Context, date time.Time) (*DailyStats, error) {
	dateStr := date.Format("2006-01-02")

	stats := &DailyStats{
		Date:     date,
		ByType:   make(map[string]TypeStats),
		ByStatus: make(map[string]StatusStats),
	}

	// Get total transactions and volume
	var totalStats struct {
		Count  int64 `db:"count"`
		Volume int64 `db:"volume"`
	}
	queryTotal := `
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as volume
		FROM transactions
		WHERE DATE(created_at) = $1
	`
	if err := uc.db.GetContext(ctx, &totalStats, queryTotal, dateStr); err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}
	stats.TotalTransactions = totalStats.Count
	stats.TotalVolume = totalStats.Volume

	// Get stats by transaction type
	type TypeResult struct {
		Type   string `db:"transaction_type"`
		Count  int64  `db:"count"`
		Volume int64  `db:"volume"`
	}
	var typeResults []TypeResult
	queryByType := `
		SELECT 
			transaction_type,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as volume
		FROM transactions
		WHERE DATE(created_at) = $1
		GROUP BY transaction_type
	`
	if err := uc.db.SelectContext(ctx, &typeResults, queryByType, dateStr); err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}

	for _, result := range typeResults {
		stats.ByType[result.Type] = TypeStats{
			Count:  result.Count,
			Volume: result.Volume,
		}
	}

	// Get stats by status
	type StatusResult struct {
		Status string `db:"status"`
		Count  int64  `db:"count"`
	}
	var statusResults []StatusResult
	queryByStatus := `
		SELECT 
			status,
			COUNT(*) as count
		FROM transactions
		WHERE DATE(created_at) = $1
		GROUP BY status
	`
	if err := uc.db.SelectContext(ctx, &statusResults, queryByStatus, dateStr); err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}

	for _, result := range statusResults {
		stats.ByStatus[result.Status] = StatusStats{
			Count: result.Count,
		}
	}

	return stats, nil
}

// GetTransactionSummary returns transaction summary for date range
func (uc *dashboardUsecase) GetTransactionSummary(ctx context.Context, startDate, endDate time.Time) (*TransactionSummary, error) {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	summary := &TransactionSummary{
		StartDate: startDate,
		EndDate:   endDate,
		ByType:    make(map[string]TypeStats),
		ByStatus:  make(map[string]StatusStats),
	}

	// Get total transactions and volume
	var totalStats struct {
		Count  int64 `db:"count"`
		Volume int64 `db:"volume"`
	}
	queryTotal := `
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as volume
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
	`
	if err := uc.db.GetContext(ctx, &totalStats, queryTotal, startStr, endStr); err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}
	summary.TotalTransactions = totalStats.Count
	summary.TotalVolume = totalStats.Volume

	// Get stats by type
	type TypeResult struct {
		Type   string `db:"transaction_type"`
		Count  int64  `db:"count"`
		Volume int64  `db:"volume"`
	}
	var typeResults []TypeResult
	queryByType := `
		SELECT 
			transaction_type,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as volume
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
		GROUP BY transaction_type
	`
	if err := uc.db.SelectContext(ctx, &typeResults, queryByType, startStr, endStr); err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}

	for _, result := range typeResults {
		summary.ByType[result.Type] = TypeStats{
			Count:  result.Count,
			Volume: result.Volume,
		}
	}

	// Get stats by status
	type StatusResult struct {
		Status string `db:"status"`
		Count  int64  `db:"count"`
	}
	var statusResults []StatusResult
	queryByStatus := `
		SELECT 
			status,
			COUNT(*) as count
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
		GROUP BY status
	`
	if err := uc.db.SelectContext(ctx, &statusResults, queryByStatus, startStr, endStr); err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}

	for _, result := range statusResults {
		summary.ByStatus[result.Status] = StatusStats{
			Count: result.Count,
		}
	}

	// Get daily breakdown
	type DailyResult struct {
		Date   time.Time `db:"date"`
		Count  int64     `db:"count"`
		Volume int64     `db:"volume"`
	}
	var dailyResults []DailyResult
	queryDaily := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as volume
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`
	if err := uc.db.SelectContext(ctx, &dailyResults, queryDaily, startStr, endStr); err != nil {
		return nil, fmt.Errorf("failed to get daily breakdown: %w", err)
	}

	summary.DailyBreakdown = make([]DailyBreakdown, 0, len(dailyResults))
	for _, result := range dailyResults {
		summary.DailyBreakdown = append(summary.DailyBreakdown, DailyBreakdown{
			Date:   result.Date,
			Count:  result.Count,
			Volume: result.Volume,
		})
	}

	return summary, nil
}
