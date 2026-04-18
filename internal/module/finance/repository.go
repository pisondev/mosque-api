package finance

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	// Helper untuk rute publik
	GetTenantIDByHostname(ctx context.Context, hostname string) (string, error)

	// PG Config
	GetPGConfig(ctx context.Context, tenantID string) (*PGConfigResponse, error)
	UpsertPGConfig(ctx context.Context, tenantID string, req PGConfigPayload) error

	// Campaigns
	ListCampaigns(ctx context.Context, tenantID string, q ListQuery) ([]CampaignResponse, int64, error)
	CreateCampaign(ctx context.Context, tenantID string, req CampaignPayload) (*CampaignResponse, error)
	GetCampaign(ctx context.Context, tenantID string, id int64) (*CampaignResponse, error)
	GetCampaignBySlug(ctx context.Context, tenantID string, slug string) (*CampaignResponse, error)
	UpdateCampaign(ctx context.Context, tenantID string, id int64, req CampaignPayload) error

	// Transactions (Read only for now, Create & Update will be in Tahap 4)
	ListTransactions(ctx context.Context, tenantID string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error)
	ListPublicDonors(ctx context.Context, tenantID string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error)

	CreateTransaction(ctx context.Context, tenantID string, req DonatePayload, status string) (*TransactionResponse, error)
	UpdateTransactionPGInfo(ctx context.Context, transactionID string, snapToken string, paymentURL string) error
	CreateSubscriptionTransaction(ctx context.Context, tenantID, orderID, planCode string, amount float64, status string) (*SubscriptionTransactionResponse, error)
	UpdateSubscriptionPGInfo(ctx context.Context, transactionID string, snapToken string, paymentURL string) error
	GetSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error)
	GetLatestPendingSubscriptionTransaction(ctx context.Context, tenantID string) (*SubscriptionTransactionResponse, error)
	GetLatestPaidSubscriptionTransaction(ctx context.Context, tenantID string) (*SubscriptionTransactionResponse, error)
	ListPaidSubscriptionTransactions(ctx context.Context, tenantID string) ([]SubscriptionTransactionResponse, error)
	ListSubscriptionTransactions(ctx context.Context, tenantID string, q ListQuery) ([]SubscriptionTransactionResponse, int64, error)
	GetTenantSubscriptionPlan(ctx context.Context, tenantID string) (string, error)
	ApplyImmediateSubscriptionDowngrade(ctx context.Context, tenantID, planCode string) (*SubscriptionTransactionResponse, error)
	CancelSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error)
	GetTenantIDBySubscriptionOrderID(ctx context.Context, orderID string) (string, error)
	ProcessSubscriptionWebhookTransaction(ctx context.Context, orderID string, status string, paymentMethod string, transactionID string, rawPayload interface{}) error
	ActivateFreePlan(ctx context.Context, tenantID string) error

	ListPublicCampaigns(ctx context.Context, hostname string, q ListQuery) ([]CampaignResponse, int64, error)
	GetPublicCampaignBySlug(ctx context.Context, hostname string, slug string) (*CampaignResponse, error)

	// Webhook Handler
	ProcessWebhookTransaction(ctx context.Context, orderID string, status string, paymentMethod string) error
	GetTenantIDByTransactionID(ctx context.Context, transactionID string) (string, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// Helper untuk pagination
func getOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return (page - 1) * limit
}

func (r *repository) GetTenantIDByHostname(ctx context.Context, hostname string) (string, error) {
	var tenantID string
	err := r.db.QueryRow(ctx, `SELECT tenant_id FROM website_domains WHERE hostname = $1 LIMIT 1`, hostname).Scan(&tenantID)
	return tenantID, err
}

// ==========================================
// PG CONFIGURATIONS
// ==========================================

func (r *repository) GetPGConfig(ctx context.Context, tenantID string) (*PGConfigResponse, error) {
	// Tambahkan server_key pada query
	query := `
		SELECT id, use_central_pg, provider, client_key, server_key, is_production, is_active
		FROM pg_configs
		WHERE tenant_id = $1
	`
	var res PGConfigResponse
	// Tambahkan &res.ServerKey pada bagian Scan
	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&res.ID, &res.UseCentralPG, &res.Provider, &res.ClientKey, &res.ServerKey, &res.IsProduction, &res.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *repository) UpsertPGConfig(ctx context.Context, tenantID string, req PGConfigPayload) error {
	query := `
		INSERT INTO pg_configs (tenant_id, use_central_pg, provider, client_key, server_key, is_production)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tenant_id) DO UPDATE SET
			use_central_pg = EXCLUDED.use_central_pg,
			provider = EXCLUDED.provider,
			client_key = EXCLUDED.client_key,
			server_key = EXCLUDED.server_key,
			is_production = EXCLUDED.is_production,
			updated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query, tenantID, req.UseCentralPG, req.Provider, req.ClientKey, req.ServerKey, req.IsProduction)
	return err
}

// ==========================================
// DONATION CAMPAIGNS
// ==========================================

func (r *repository) CreateCampaign(ctx context.Context, tenantID string, req CampaignPayload) (*CampaignResponse, error) {
	query := `
		INSERT INTO donation_campaigns (tenant_id, title, slug, description, image_url, target_amount, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, slug, description, image_url, target_amount, collected_amount, start_date, end_date, is_active
	`
	var res CampaignResponse
	err := r.db.QueryRow(ctx, query, tenantID, req.Title, req.Slug, req.Description, req.ImageURL, req.TargetAmount, req.StartDate, req.EndDate, req.IsActive).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.ImageURL, &res.TargetAmount, &res.CollectedAmount, &res.StartDate, &res.EndDate, &res.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) GetCampaign(ctx context.Context, tenantID string, id int64) (*CampaignResponse, error) {
	query := `
		SELECT id, title, slug, description, image_url, target_amount, collected_amount, start_date, end_date, is_active
		FROM donation_campaigns
		WHERE tenant_id = $1 AND id = $2
	`
	var res CampaignResponse
	err := r.db.QueryRow(ctx, query, tenantID, id).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.ImageURL, &res.TargetAmount, &res.CollectedAmount, &res.StartDate, &res.EndDate, &res.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) GetCampaignBySlug(ctx context.Context, tenantID string, slug string) (*CampaignResponse, error) {
	query := `
		SELECT id, title, slug, description, image_url, target_amount, collected_amount, start_date, end_date, is_active
		FROM donation_campaigns
		WHERE tenant_id = $1 AND slug = $2 AND is_active = true
	`
	var res CampaignResponse
	err := r.db.QueryRow(ctx, query, tenantID, slug).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.ImageURL, &res.TargetAmount, &res.CollectedAmount, &res.StartDate, &res.EndDate, &res.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) UpdateCampaign(ctx context.Context, tenantID string, id int64, req CampaignPayload) error {
	query := `
		UPDATE donation_campaigns
		SET title = $1, slug = $2, description = $3, image_url = $4, target_amount = $5, start_date = $6, end_date = $7, is_active = $8, updated_at = NOW()
		WHERE tenant_id = $9 AND id = $10
	`
	cmdTag, err := r.db.Exec(ctx, query, req.Title, req.Slug, req.Description, req.ImageURL, req.TargetAmount, req.StartDate, req.EndDate, req.IsActive, tenantID, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ==========================================
// LIST METHODS (Pagination)
// ==========================================

func (r *repository) ListCampaigns(ctx context.Context, tenantID string, q ListQuery) ([]CampaignResponse, int64, error) {
	offset := getOffset(q.Page, q.Limit)

	// 1. Hitung total data
	var total int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(id) FROM donation_campaigns WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Jika kosong, langsung return agar tidak perlu eksekusi query kedua
	if total == 0 {
		return []CampaignResponse{}, 0, nil
	}

	// 2. Ambil data dengan Limit & Offset
	query := `
		SELECT id, title, slug, description, image_url, target_amount, collected_amount, start_date, end_date, is_active
		FROM donation_campaigns
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var campaigns []CampaignResponse
	for rows.Next() {
		var c CampaignResponse
		err := rows.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.ImageURL, &c.TargetAmount, &c.CollectedAmount, &c.StartDate, &c.EndDate, &c.IsActive)
		if err != nil {
			return nil, 0, err
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, total, rows.Err()
}

func (r *repository) ListTransactions(ctx context.Context, tenantID string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error) {
	offset := getOffset(q.Page, q.Limit)

	var total int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(id) FROM donation_transactions WHERE tenant_id = $1 AND campaign_id = $2`, tenantID, campaignID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []TransactionResponse{}, 0, nil
	}

	query := `
		SELECT id, donor_name, is_anonymous, amount, payment_method, status, paid_at
		FROM donation_transactions
		WHERE tenant_id = $1 AND campaign_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(ctx, query, tenantID, campaignID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var txs []TransactionResponse
	for rows.Next() {
		var tx TransactionResponse
		err := rows.Scan(&tx.TransactionID, &tx.DonorName, &tx.IsAnonymous, &tx.Amount, &tx.PaymentMethod, &tx.Status, &tx.PaidAt)
		if err != nil {
			return nil, 0, err
		}
		txs = append(txs, tx)
	}

	return txs, total, rows.Err()
}

func (r *repository) ListPublicDonors(ctx context.Context, tenantID string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error) {
	// Hampir sama dengan ListTransactions, tapi kita HANYA ambil yang statusnya 'paid' untuk transparansi publik
	offset := getOffset(q.Page, q.Limit)

	var total int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(id) FROM donation_transactions WHERE tenant_id = $1 AND campaign_id = $2 AND status = 'paid'`, tenantID, campaignID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []TransactionResponse{}, 0, nil
	}

	query := `
		SELECT id, donor_name, is_anonymous, amount, payment_method, status, paid_at
		FROM donation_transactions
		WHERE tenant_id = $1 AND campaign_id = $2 AND status = 'paid'
		ORDER BY paid_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(ctx, query, tenantID, campaignID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var txs []TransactionResponse
	for rows.Next() {
		var tx TransactionResponse
		err := rows.Scan(&tx.TransactionID, &tx.DonorName, &tx.IsAnonymous, &tx.Amount, &tx.PaymentMethod, &tx.Status, &tx.PaidAt)
		if err != nil {
			return nil, 0, err
		}

		// Logika Anonim: Ganti nama menjadi "Hamba Allah" jika is_anonymous = true
		if tx.IsAnonymous {
			tx.DonorName = "Hamba Allah"
		}

		txs = append(txs, tx)
	}

	return txs, total, rows.Err()
}

// ==========================================
// TRANSACTIONS (CREATE & UPDATE)
// ==========================================

func (r *repository) CreateTransaction(ctx context.Context, tenantID string, req DonatePayload, status string) (*TransactionResponse, error) {
	query := `
		INSERT INTO donation_transactions (tenant_id, campaign_id, donor_name, is_anonymous, amount, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, NULL, $6)
		RETURNING id, donor_name, is_anonymous, amount, payment_method, status, paid_at
	`
	// Jika donor_name kosong dan anonymous, kita set default
	if req.IsAnonymous && (req.DonorName == nil || *req.DonorName == "") {
		hambaAllah := "Hamba Allah"
		req.DonorName = &hambaAllah
	}

	var tx TransactionResponse
	err := r.db.QueryRow(ctx, query, tenantID, req.CampaignID, req.DonorName, req.IsAnonymous, req.Amount, status).Scan(
		&tx.TransactionID, &tx.DonorName, &tx.IsAnonymous, &tx.Amount, &tx.PaymentMethod, &tx.Status, &tx.PaidAt,
	)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *repository) UpdateTransactionPGInfo(ctx context.Context, transactionID string, snapToken string, paymentURL string) error {
	query := `
		UPDATE donation_transactions 
		SET snap_token = $1, payment_url = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, snapToken, paymentURL, transactionID)
	return err
}

func (r *repository) CreateSubscriptionTransaction(ctx context.Context, tenantID, orderID, planCode string, amount float64, status string) (*SubscriptionTransactionResponse, error) {
	query := `
		INSERT INTO subscription_transactions (tenant_id, order_id, plan_code, amount, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
	`
	var res SubscriptionTransactionResponse
	err := r.db.QueryRow(ctx, query, tenantID, orderID, planCode, amount, status).Scan(
		&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) UpdateSubscriptionPGInfo(ctx context.Context, transactionID string, snapToken string, paymentURL string) error {
	rawJSON, _ := json.Marshal(map[string]interface{}{
		"source":      "checkout_create",
		"snap_token":  snapToken,
		"payment_url": paymentURL,
		"created_at":  time.Now().UTC(),
	})
	query := `
		UPDATE subscription_transactions
		SET snap_token = $1, payment_url = $2, raw_notification = COALESCE($4::jsonb, raw_notification), updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, snapToken, paymentURL, transactionID, string(rawJSON))
	return err
}

func (r *repository) GetSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error) {
	query := `
		SELECT id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
		FROM subscription_transactions
		WHERE tenant_id = $1 AND id = $2
	`
	var res SubscriptionTransactionResponse
	err := r.db.QueryRow(ctx, query, tenantID, transactionID).Scan(
		&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) GetLatestPendingSubscriptionTransaction(ctx context.Context, tenantID string) (*SubscriptionTransactionResponse, error) {
	query := `
		SELECT id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
		FROM subscription_transactions
		WHERE tenant_id = $1 AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1
	`
	var res SubscriptionTransactionResponse
	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *repository) GetLatestPaidSubscriptionTransaction(ctx context.Context, tenantID string) (*SubscriptionTransactionResponse, error) {
	query := `
		SELECT id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
		FROM subscription_transactions
		WHERE tenant_id = $1 AND status = 'paid'
		ORDER BY paid_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`
	var res SubscriptionTransactionResponse
	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *repository) ListPaidSubscriptionTransactions(ctx context.Context, tenantID string) ([]SubscriptionTransactionResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
		FROM subscription_transactions
		WHERE tenant_id = $1 AND status = 'paid'
		ORDER BY paid_at ASC NULLS LAST, created_at ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]SubscriptionTransactionResponse, 0)
	for rows.Next() {
		var res SubscriptionTransactionResponse
		if err := rows.Scan(&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL); err != nil {
			return nil, err
		}
		items = append(items, res)
	}
	return items, rows.Err()
}

func (r *repository) ListSubscriptionTransactions(ctx context.Context, tenantID string, q ListQuery) ([]SubscriptionTransactionResponse, int64, error) {
	offset := getOffset(q.Page, q.Limit)

	var total int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(id) FROM subscription_transactions WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []SubscriptionTransactionResponse{}, 0, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
		FROM subscription_transactions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]SubscriptionTransactionResponse, 0, q.Limit)
	for rows.Next() {
		var res SubscriptionTransactionResponse
		if err := rows.Scan(
			&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, res)
	}

	return items, total, rows.Err()
}

func (r *repository) GetTenantSubscriptionPlan(ctx context.Context, tenantID string) (string, error) {
	var planCode string
	err := r.db.QueryRow(ctx, `SELECT subscription_plan FROM tenants WHERE id = $1`, tenantID).Scan(&planCode)
	if err != nil {
		return "", err
	}
	return planCode, nil
}

func (r *repository) ApplyImmediateSubscriptionDowngrade(ctx context.Context, tenantID, planCode string) (*SubscriptionTransactionResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	orderID := "SUB-DOWN-" + strings.ToUpper(uuid.NewString())
	var res SubscriptionTransactionResponse
	err = tx.QueryRow(ctx, `
		INSERT INTO subscription_transactions (tenant_id, order_id, plan_code, amount, status, payment_method, paid_at, raw_notification)
		VALUES ($1, $2, $3, 0, 'paid', 'internal_downgrade', NOW(), '{"source":"internal_downgrade"}')
		RETURNING id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
	`, tenantID, orderID, planCode).Scan(&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL)
	if err != nil {
		return nil, err
	}

	if _, err = tx.Exec(ctx, `
		UPDATE tenants
		SET subscription_plan = $1,
			onboarding_payment_status = 'paid',
			updated_at = NOW()
		WHERE id = $2
	`, planCode, tenantID); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) CancelSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error) {
	query := `
		UPDATE subscription_transactions
		SET status = 'failed',
			expired_at = COALESCE(expired_at, NOW()),
			updated_at = NOW()
		WHERE tenant_id = $1 AND id = $2 AND status = 'pending'
		RETURNING id, order_id, plan_code, amount, status, payment_method, paid_at, expired_at, created_at, payment_url
	`
	var res SubscriptionTransactionResponse
	err := r.db.QueryRow(ctx, query, tenantID, transactionID).Scan(
		&res.TransactionID, &res.OrderID, &res.PlanCode, &res.Amount, &res.Status, &res.PaymentMethod, &res.PaidAt, &res.ExpiredAt, &res.CreatedAt, &res.PaymentURL,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ==========================================
// PUBLIC METHODS (Menggunakan Hostname & JOIN)
// ==========================================

func (r *repository) ListPublicCampaigns(ctx context.Context, hostname string, q ListQuery) ([]CampaignResponse, int64, error) {
	offset := getOffset(q.Page, q.Limit)

	var total int64
	// Ingat: Cuma ambil yang is_active = true
	countQuery := `
		SELECT COUNT(dc.id) 
		FROM donation_campaigns dc
		JOIN website_domains wd ON dc.tenant_id = wd.tenant_id
		WHERE wd.hostname = $1 AND dc.is_active = true
	`
	err := r.db.QueryRow(ctx, countQuery, hostname).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []CampaignResponse{}, 0, nil
	}

	query := `
		SELECT dc.id, dc.title, dc.slug, dc.description, dc.image_url, dc.target_amount, dc.collected_amount, dc.start_date, dc.end_date, dc.is_active
		FROM donation_campaigns dc
		JOIN website_domains wd ON dc.tenant_id = wd.tenant_id
		WHERE wd.hostname = $1 AND dc.is_active = true
		ORDER BY dc.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, hostname, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var campaigns []CampaignResponse
	for rows.Next() {
		var c CampaignResponse
		err := rows.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.ImageURL, &c.TargetAmount, &c.CollectedAmount, &c.StartDate, &c.EndDate, &c.IsActive)
		if err != nil {
			return nil, 0, err
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, total, rows.Err()
}

func (r *repository) GetPublicCampaignBySlug(ctx context.Context, hostname string, slug string) (*CampaignResponse, error) {
	query := `
		SELECT dc.id, dc.title, dc.slug, dc.description, dc.image_url, dc.target_amount, dc.collected_amount, dc.start_date, dc.end_date, dc.is_active
		FROM donation_campaigns dc
		JOIN website_domains wd ON dc.tenant_id = wd.tenant_id
		WHERE wd.hostname = $1 AND dc.slug = $2 AND dc.is_active = true
	`
	var res CampaignResponse
	err := r.db.QueryRow(ctx, query, hostname, slug).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.ImageURL, &res.TargetAmount, &res.CollectedAmount, &res.StartDate, &res.EndDate, &res.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ==========================================
// WEBHOOK HANDLER (DB TRANSACTION & ROW LOCKING)
// ==========================================

func (r *repository) ProcessWebhookTransaction(ctx context.Context, orderID string, status string, paymentMethod string) error {
	// 1. Mulai DB Transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Pastikan rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// 2. Kunci baris transaksi ini agar tidak diubah proses lain (Row Locking)
	var currentStatus string
	var campaignID int64
	var amount float64

	lockQuery := `
		SELECT status, campaign_id, amount 
		FROM donation_transactions 
		WHERE id = $1 
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, lockQuery, orderID).Scan(&currentStatus, &campaignID, &amount)
	if err != nil {
		return err // Transaksi tidak ditemukan
	}

	// 3. Cegah update berulang (Idempotency)
	// Jika status di DB sudah paid/settlement, abaikan webhook ini (Midtrans sering kirim webhook multiple kali)
	if currentStatus == "paid" {
		return nil
	}

	// 4. Update status di tabel transaksi
	// Mapping status Midtrans ke status lokal kita
	localStatus := "pending"
	switch status {
	case "settlement", "capture":
		localStatus = "paid"
	case "expire", "cancel", "deny":
		localStatus = "expired"
	}

	updateTxQuery := `
		UPDATE donation_transactions
		SET status = $1, payment_method = $2, paid_at = CASE WHEN $1 = 'paid' THEN NOW() ELSE paid_at END, updated_at = NOW()
		WHERE id = $3
	`
	_, err = tx.Exec(ctx, updateTxQuery, localStatus, paymentMethod, orderID)
	if err != nil {
		return err
	}

	// 5. Jika sukses dibayar (paid), tambah saldo di tabel kampanye
	if localStatus == "paid" {
		updateCampQuery := `
			UPDATE donation_campaigns 
			SET collected_amount = collected_amount + $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = tx.Exec(ctx, updateCampQuery, amount, campaignID)
		if err != nil {
			return err
		}
	}

	// 6. Selesai dan Commit!
	return tx.Commit(ctx)
}

func (r *repository) GetTenantIDBySubscriptionOrderID(ctx context.Context, orderID string) (string, error) {
	var tenantID string
	err := r.db.QueryRow(ctx, `SELECT tenant_id FROM subscription_transactions WHERE order_id = $1`, orderID).Scan(&tenantID)
	return tenantID, err
}

func (r *repository) ProcessSubscriptionWebhookTransaction(ctx context.Context, orderID string, status string, paymentMethod string, transactionID string, rawPayload interface{}) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var currentStatus string
	var tenantID string
	var planCode string
	lockQuery := `
		SELECT status, tenant_id, plan_code
		FROM subscription_transactions
		WHERE order_id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, lockQuery, orderID).Scan(&currentStatus, &tenantID, &planCode)
	if err != nil {
		return err
	}

	if currentStatus == "paid" {
		return nil
	}

	localStatus := "pending"
	switch status {
	case "settlement", "capture":
		localStatus = "paid"
	case "expire":
		localStatus = "expired"
	case "cancel", "deny":
		localStatus = "failed"
	}

	var paidAt interface{}
	var expiredAt interface{}
	if localStatus == "paid" {
		paidAt = time.Now()
	}
	if localStatus == "expired" {
		expiredAt = time.Now()
	}

	rawJSON, _ := json.Marshal(rawPayload)
	_, err = tx.Exec(ctx, `
		UPDATE subscription_transactions
		SET status = $1,
			payment_method = $2,
			midtrans_transaction_id = NULLIF($3, ''),
			paid_at = COALESCE($4, paid_at),
			expired_at = COALESCE($5, expired_at),
			raw_notification = COALESCE($6::jsonb, raw_notification),
			updated_at = NOW()
		WHERE order_id = $7
	`, localStatus, paymentMethod, transactionID, paidAt, expiredAt, string(rawJSON), orderID)
	if err != nil {
		return err
	}

	if localStatus == "paid" {
		_, err = tx.Exec(ctx, `
			UPDATE tenants
			SET subscription_plan = $1,
				onboarding_payment_status = 'paid',
				onboarding_completed = false,
				status = 'pending',
				updated_at = NOW()
			WHERE id = $2
		`, planCode, tenantID)
		if err != nil {
			return err
		}
	} else if localStatus == "expired" || localStatus == "failed" {
		_, err = tx.Exec(ctx, `
			UPDATE tenants
			SET onboarding_payment_status = $1,
				updated_at = NOW()
			WHERE id = $2
		`, localStatus, tenantID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) GetTenantIDByTransactionID(ctx context.Context, transactionID string) (string, error) {
	var tenantID string
	err := r.db.QueryRow(ctx, `SELECT tenant_id FROM donation_transactions WHERE id = $1`, transactionID).Scan(&tenantID)
	return tenantID, err
}

func (r *repository) ActivateFreePlan(ctx context.Context, tenantID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tenants
		SET subscription_plan = 'free',
			onboarding_payment_status = 'free',
			onboarding_completed = false,
			status = 'pending',
			updated_at = NOW()
		WHERE id = $1
	`, tenantID)
	return err
}
