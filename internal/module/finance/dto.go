package finance

import (
	"time"
)

// ==========================================
// PG CONFIGURATIONS
// ==========================================

type PGConfigPayload struct {
	UseCentralPG bool   `json:"use_central_pg"`
	Provider     string `json:"provider"`
	ClientKey    string `json:"client_key"`
	ServerKey    string `json:"server_key"`
	IsProduction bool   `json:"is_production"`
}

type PGConfigResponse struct {
	ID           int64  `json:"id"`
	UseCentralPG bool   `json:"use_central_pg"`
	Provider     string `json:"provider"`
	ClientKey    string `json:"client_key"`
	ServerKey    string `json:"-"` // <-- Tambahkan ini. json:"-" = diabaikan saat jadi JSON
	IsProduction bool   `json:"is_production"`
	IsActive     bool   `json:"is_active"`
}

// ==========================================
// DONATION CAMPAIGNS
// ==========================================

type CampaignPayload struct {
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	Description  *string    `json:"description"`
	ImageURL     *string    `json:"image_url"`
	TargetAmount float64    `json:"target_amount"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsActive     bool       `json:"is_active"`
}

type CampaignResponse struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Description     *string    `json:"description"`
	ImageURL        *string    `json:"image_url"`
	TargetAmount    float64    `json:"target_amount"`
	CollectedAmount float64    `json:"collected_amount"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	IsActive        bool       `json:"is_active"`
}

// ==========================================
// DONATION TRANSACTIONS & DONORS
// ==========================================

type DonatePayload struct {
	CampaignID  int64   `json:"campaign_id"`
	DonorName   *string `json:"donor_name"`
	IsAnonymous bool    `json:"is_anonymous"`
	Amount      float64 `json:"amount"`
}

type TransactionResponse struct {
	TransactionID string     `json:"transaction_id"` // UUID as string
	DonorName     string     `json:"donor_name"`
	IsAnonymous   bool       `json:"is_anonymous"`
	Amount        float64    `json:"amount"`
	PaymentMethod *string    `json:"payment_method"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
	PaymentURL    *string    `json:"payment_url"` // Diisi saat create donasi baru
}

// Common pagination query
type ListQuery struct {
	Page  int
	Limit int
}
