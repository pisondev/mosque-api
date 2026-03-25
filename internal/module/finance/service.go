package finance

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
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

	CreateDonation(ctx context.Context, hostname string, req DonatePayload) (*TransactionResponse, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

// Helper untuk generate slug
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
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
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}
	return s.repo.CreateCampaign(ctx, tenantID, req)
}

func (s *service) GetCampaign(ctx context.Context, tenantID string, id int64) (*CampaignResponse, error) {
	return s.repo.GetCampaign(ctx, tenantID, id)
}

func (s *service) UpdateCampaign(ctx context.Context, tenantID string, id int64, req CampaignPayload) error {
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}
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

func (s *service) ListPublicDonors(ctx context.Context, hostname string, campaignID int64, q ListQuery) ([]TransactionResponse, int64, error) {
	// 1. Dapatkan tenant_id dari hostname menggunakan antarmuka Repository resmi (Tidak ada lagi leaky abstraction!)
	tenantID, err := s.repo.GetTenantIDByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err // Hostname tidak valid atau tidak ditemukan
	}

	// 2. Panggil repo yg sudah ada
	return s.repo.ListPublicDonors(ctx, tenantID, campaignID, q)
}

// ==========================================
// PUBLIC DONATION / CHECKOUT LOGIC
// ==========================================

func (s *service) CreateDonation(ctx context.Context, hostname string, req DonatePayload) (*TransactionResponse, error) {
	// 1. Dapatkan tenant_id dari hostname
	tenantID, err := s.repo.GetTenantIDByHostname(ctx, hostname)
	if err != nil {
		return nil, errors.New("hostname tidak valid")
	}

	// 2. Validasi Kampanye (Pastikan ada dan aktif)
	campaign, err := s.repo.GetCampaign(ctx, tenantID, req.CampaignID)
	if err != nil || !campaign.IsActive {
		return nil, errors.New("kampanye tidak ditemukan atau sudah tidak aktif")
	}

	// 3. Simpan Transaksi ke DB dengan status "pending"
	txRes, err := s.repo.CreateTransaction(ctx, tenantID, req, "pending")
	if err != nil {
		return nil, err
	}

	// 4. Siapkan Kunci Midtrans berdasarkan PG Config Masjid
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY") // Default pakai pusat
	envType := midtrans.Sandbox

	pgConfig, err := s.repo.GetPGConfig(ctx, tenantID)
	if err == nil && pgConfig != nil && !pgConfig.UseCentralPG && pgConfig.IsActive {
		// Jika masjid punya config mandiri dan aktif
		serverKey = pgConfig.ServerKey // Di tahap advanced, ini harus di-decrypt dulu
		if pgConfig.IsProduction {
			envType = midtrans.Production
		}
	}

	if serverKey == "" {
		return nil, errors.New("konfigurasi payment gateway (Server Key) belum diatur")
	}

	// 5. Inisialisasi Snap Client KHUSUS untuk request ini (Multi-tenant safe)
	var snapClient snap.Client
	snapClient.New(serverKey, envType)

	// 6. Buat Request ke Midtrans
	midtransReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  txRes.TransactionID,
			GrossAmt: int64(req.Amount),
		},
		Items: &[]midtrans.ItemDetails{ // <-- Ubah nama field di sini dari ItemDetails menjadi Items
			{
				ID:    "CMP-" + fmt.Sprintf("%d", campaign.ID),
				Price: int64(req.Amount),
				Qty:   1,
				Name:  campaign.Title,
			},
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: txRes.DonorName,
		},
	}

	// 7. Eksekusi Request ke Midtrans
	snapResp, midtransErr := snapClient.CreateTransaction(midtransReq)
	if midtransErr != nil {
		s.log.Error("Midtrans Error: ", midtransErr.Message)
		return nil, errors.New("gagal menghubungkan ke payment gateway")
	}

	// 8. Update DB kita dengan Snap Token dan Redirect URL
	err = s.repo.UpdateTransactionPGInfo(ctx, txRes.TransactionID, snapResp.Token, snapResp.RedirectURL)
	if err != nil {
		return nil, errors.New("donasi berhasil dibuat, tetapi gagal menyimpan token")
	}

	// 9. Kembalikan response lengkap ke Frontend
	txRes.PaymentURL = &snapResp.RedirectURL
	// Kita bisa sisipkan snap token sementara di PaymentMethod jika frontend butuh popup JS
	txRes.PaymentMethod = &snapResp.Token

	return txRes, nil
}
