package finance

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
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

// ==========================================
// PG CONFIGURATIONS
// ==========================================

func (r *repository) GetPGConfig(ctx context.Context, tenantID string) (*PGConfigResponse, error) {
	query := `
		SELECT id, use_central_pg, provider, client_key, is_production, is_active
		FROM pg_configs
		WHERE tenant_id = $1
	`
	var res PGConfigResponse
	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&res.ID, &res.UseCentralPG, &res.Provider, &res.ClientKey, &res.IsProduction, &res.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil tanpa error jika masjid belum punya config
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
