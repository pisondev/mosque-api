package finance

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service interface {
	GetPGConfig(ctx context.Context, tenantID string) (*PGConfigResponse, error)
	UpsertPGConfig(ctx context.Context, tenantID string, req PGConfigPayload) error

	ListCampaigns(ctx context.Context, tenantID string, q ListQuery) ([]CampaignResponse, int64, error)
	ListPublicCampaigns(ctx context.Context, hostname string, q ListQuery) ([]CampaignResponse, int64, error) // Beda parameter dengan admin
	CreateCampaign(ctx context.Context, tenantID string, req CampaignPayload) (*CampaignResponse, error)
	GetCampaign(ctx context.Context, tenantID string, id int64) (*CampaignResponse, error)
	GetPublicCampaignBySlug(ctx context.Context, hostname string, slug string) (*CampaignResponse, error)
	UpdateCampaign(ctx context.Context, tenantID string, id int64, req CampaignPayload) error

	ListTransactions(ctx context.Context, tenantID string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error)
	ListPublicDonors(ctx context.Context, hostname string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error)

	// CreateDonation akan memanggil API Midtrans (Nanti di Tahap 4)
	// CreateDonation(ctx context.Context, hostname string, req DonatePayload) (*TransactionResponse, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

// ==========================================
// PG CONFIGURATIONS
// ==========================================

func (s *service) GetPGConfig(ctx context.Context, tenantID string) (*PGConfigResponse, error) {
	return s.repo.GetPGConfig(ctx, tenantID)
}

func (s *service) UpsertPGConfig(ctx context.Context, tenantID string, req PGConfigPayload) error {
	// TODO: Di masa depan, tambahkan logika enkripsi untuk req.ServerKey menggunakan AES-GCM di sini
	// sebelum diteruskan ke repository.
	return s.repo.UpsertPGConfig(ctx, tenantID, req)
}

// ==========================================
// DONATION CAMPAIGNS
// ==========================================

func (s *service) CreateCampaign(ctx context.Context, tenantID string, req CampaignPayload) (*CampaignResponse, error) {
	return s.repo.CreateCampaign(ctx, tenantID, req)
}

func (s *service) GetCampaign(ctx context.Context, tenantID string, id int64) (*CampaignResponse, error) {
	return s.repo.GetCampaign(ctx, tenantID, id)
}

func (s *service) UpdateCampaign(ctx context.Context, tenantID string, id int64, req CampaignPayload) error {
	return s.repo.UpdateCampaign(ctx, tenantID, id, req)
}

func (s *service) GetPublicCampaignBySlug(ctx context.Context, hostname string, slug string) (*CampaignResponse, error) {
	return s.repo.GetPublicCampaignBySlug(ctx, hostname, slug)
}

// ==========================================
// LIST METHODS (Placeholder)
// ==========================================

func (s *service) ListCampaigns(ctx context.Context, tenantID string, q ListQuery) ([]CampaignResponse, int64, error) {
	return s.repo.ListCampaigns(ctx, tenantID, q)
}

func (s *service) ListPublicCampaigns(ctx context.Context, hostname string, q ListQuery) ([]CampaignResponse, int64, error) {
	return s.repo.ListPublicCampaigns(ctx, hostname, q)
}

func (s *service) ListTransactions(ctx context.Context, tenantID string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error) {
	return s.repo.ListTransactions(ctx, tenantID, campaignID, q)
}

// Untuk ListPublicDonors, kita butuh trick dikit untuk dapetin tenantID dari hostname dulu,
// baru panggil fungsi ListPublicDonors yang udah dibikin sebelumnya.
func (s *service) ListPublicDonors(ctx context.Context, hostname string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error) {
	// 1. Dapatkan tenant_id dari hostname
	var tenantID string
	err := s.repo.(*repository).db.QueryRow(ctx, `SELECT tenant_id FROM website_domains WHERE domain_name = $1 LIMIT 1`, hostname).Scan(&tenantID)
	if err != nil {
		return nil, 0, err // Hostname tidak valid
	}

	// 2. Panggil repo yg sudah ada
	return s.repo.ListPublicDonors(ctx, tenantID, campaignID, q)
}
