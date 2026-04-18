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

// ==========================================
// WEBHOOK / MIDTRANS NOTIFICATION
// ==========================================

type MidtransNotificationPayload struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	TransactionTime   string `json:"transaction_time"`
	FraudStatus       string `json:"fraud_status"`
}

type SubscriptionPlanItemResponse struct {
	PlanCode           string   `json:"plan_code"`
	Name               string   `json:"name"`
	Price              float64  `json:"price"`
	Currency           string   `json:"currency"`
	FeaturesUnlocked   []string `json:"features_unlocked"`
	AttributionEnabled bool     `json:"attribution_enabled"`
}

type CreateSubscriptionCheckoutPayload struct {
	PlanCode      string `json:"plan_code"`
	DurationMonth int    `json:"duration_month"`
}

type SubscriptionCheckoutResponse struct {
	TransactionID string     `json:"transaction_id"`
	OrderID       string     `json:"order_id"`
	PlanCode      string     `json:"plan_code"`
	Action        string     `json:"action,omitempty"`
	DurationMonth int        `json:"duration_month,omitempty"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	PaymentURL    *string    `json:"payment_url"`
	SnapToken     *string    `json:"snap_token"`
	CreatedAt     *time.Time `json:"created_at"`
	ExpiredAt     *time.Time `json:"expired_at"`

	ActivePlan         string     `json:"active_plan,omitempty"`
	RemainingDays      *int       `json:"remaining_days,omitempty"`
	NextBillingDueDate *time.Time `json:"next_billing_due_date,omitempty"`
	CurrentPeriodStart *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
	WarningMessage     *string    `json:"warning_message,omitempty"`
}

type SubscriptionTransactionResponse struct {
	TransactionID string     `json:"transaction_id"`
	OrderID       string     `json:"order_id"`
	PlanCode      string     `json:"plan_code"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	PaymentMethod *string    `json:"payment_method"`
	PaidAt        *time.Time `json:"paid_at"`
	ExpiredAt     *time.Time `json:"expired_at"`
	CreatedAt     time.Time  `json:"created_at"`
	PaymentURL    *string    `json:"payment_url"`
}

type TenantSubscriptionStatus struct {
	ActivePlan         string     `json:"active_plan"`
	RemainingDays      *int       `json:"remaining_days,omitempty"`
	NextBillingDueDate *time.Time `json:"next_billing_due_date,omitempty"`
	CurrentPeriodStart *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
}

type SubscriptionHistoryResponse struct {
	Status       TenantSubscriptionStatus          `json:"status"`
	Transactions []SubscriptionTransactionResponse `json:"transactions"`
}

type SubscriptionQuotePayload struct {
	PlanCode      string `json:"plan_code"`
	DurationMonth int    `json:"duration_month"`
}

type SubscriptionQuoteResponse struct {
	TargetPlanCode string                   `json:"target_plan_code"`
	Action         string                   `json:"action"`
	DurationMonth  int                      `json:"duration_month"`
	Amount         float64                  `json:"amount"`
	ProrateAmount  float64                  `json:"prorate_amount,omitempty"`
	WarningMessage *string                  `json:"warning_message,omitempty"`
	Status         TenantSubscriptionStatus `json:"status"`
}
