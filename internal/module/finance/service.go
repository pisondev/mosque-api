package finance

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/pisondev/mosque-api/internal/constant"
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
	ListSubscriptionPlans(ctx context.Context) ([]SubscriptionPlanItemResponse, error)
	CreateSubscriptionCheckout(ctx context.Context, tenantID string, req CreateSubscriptionCheckoutPayload) (*SubscriptionCheckoutResponse, error)
	GetSubscriptionQuote(ctx context.Context, tenantID string, req SubscriptionQuotePayload) (*SubscriptionQuoteResponse, error)
	CreateSubscriptionCheckoutFromQuote(ctx context.Context, tenantID string, req SubscriptionQuotePayload) (*SubscriptionCheckoutResponse, error)
	ListSubscriptionTransactions(ctx context.Context, tenantID string, q ListQuery) (*SubscriptionHistoryResponse, int64, error)
	GetSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error)
	GetActiveSubscriptionTransaction(ctx context.Context, tenantID string) (*SubscriptionTransactionResponse, error)
	CancelSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error)
	ActivateFreePlan(ctx context.Context, tenantID string) error
	HandleMidtransWebhook(ctx context.Context, payload MidtransNotificationPayload) error
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

type midtransStatusResponse struct {
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	TransactionID     string `json:"transaction_id"`
}

type subscriptionQuoteInternal struct {
	TargetPlanCode string
	Action         string
	DurationMonth  int
	Amount         float64
	ProrateAmount  float64
	WarningMessage *string
	Status         TenantSubscriptionStatus
}

func (s *service) isCentralMidtransProduction() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("MIDTRANS_IS_PRODUCTION")))
	return v == "true" || v == "1" || v == "yes"
}

func (s *service) midtransCoreAPIBaseURL() string {
	if s.isCentralMidtransProduction() {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
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
	if s.isCentralMidtransProduction() {
		envType = midtrans.Production
	}

	pgConfig, err := s.repo.GetPGConfig(ctx, tenantID)
	if err == nil && pgConfig != nil && !pgConfig.UseCentralPG && pgConfig.IsActive {
		// Jika masjid punya config mandiri dan aktif
		serverKey = pgConfig.ServerKey // Di tahap advanced, ini harus di-decrypt dulu
		if pgConfig.IsProduction {
			envType = midtrans.Production
		} else {
			envType = midtrans.Sandbox
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

func (s *service) ListSubscriptionPlans(ctx context.Context) ([]SubscriptionPlanItemResponse, error) {
	_ = ctx
	plans := []string{
		constant.PlanFree,
		constant.PlanPremiumPlus,
		constant.PlanProPlus,
		constant.PlanMaxPlus,
	}

	res := make([]SubscriptionPlanItemResponse, 0, len(plans))
	for _, planCode := range plans {
		plan := constant.SubscriptionPlans[planCode]
		res = append(res, SubscriptionPlanItemResponse{
			PlanCode:           planCode,
			Name:               plan.Name,
			Price:              plan.Price,
			Currency:           "IDR",
			FeaturesUnlocked:   plan.FeaturesUnlocked,
			AttributionEnabled: plan.AttributionEnabled,
		})
	}

	return res, nil
}

func (s *service) ActivateFreePlan(ctx context.Context, tenantID string) error {
	return s.repo.ActivateFreePlan(ctx, tenantID)
}

func (s *service) CreateSubscriptionCheckout(ctx context.Context, tenantID string, req CreateSubscriptionCheckoutPayload) (*SubscriptionCheckoutResponse, error) {
	quoteReq := SubscriptionQuotePayload{
		PlanCode:      req.PlanCode,
		DurationMonth: req.DurationMonth,
	}
	return s.CreateSubscriptionCheckoutFromQuote(ctx, tenantID, quoteReq)
}

func (s *service) ListSubscriptionTransactions(ctx context.Context, tenantID string, q ListQuery) (*SubscriptionHistoryResponse, int64, error) {
	items, total, err := s.repo.ListSubscriptionTransactions(ctx, tenantID, q)
	if err != nil {
		return nil, 0, err
	}

	status, err := s.buildTenantSubscriptionStatus(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}

	return &SubscriptionHistoryResponse{
		Status:       status,
		Transactions: items,
	}, total, nil
}

func (s *service) GetSubscriptionQuote(ctx context.Context, tenantID string, req SubscriptionQuotePayload) (*SubscriptionQuoteResponse, error) {
	q, err := s.prepareSubscriptionQuote(ctx, tenantID, req)
	if err != nil {
		return nil, err
	}
	return &SubscriptionQuoteResponse{
		TargetPlanCode: q.TargetPlanCode,
		Action:         q.Action,
		DurationMonth:  q.DurationMonth,
		Amount:         q.Amount,
		ProrateAmount:  q.ProrateAmount,
		WarningMessage: q.WarningMessage,
		Status:         q.Status,
	}, nil
}

func (s *service) CreateSubscriptionCheckoutFromQuote(ctx context.Context, tenantID string, req SubscriptionQuotePayload) (*SubscriptionCheckoutResponse, error) {
	q, err := s.prepareSubscriptionQuote(ctx, tenantID, req)
	if err != nil {
		return nil, err
	}

	if q.TargetPlanCode == constant.PlanFree {
		if err := s.repo.ActivateFreePlan(ctx, tenantID); err != nil {
			return nil, err
		}
		return &SubscriptionCheckoutResponse{
			TransactionID:      "",
			OrderID:            "",
			PlanCode:           q.TargetPlanCode,
			Action:             q.Action,
			DurationMonth:      q.DurationMonth,
			Amount:             0,
			Status:             "free",
			ActivePlan:         q.Status.ActivePlan,
			RemainingDays:      q.Status.RemainingDays,
			NextBillingDueDate: q.Status.NextBillingDueDate,
			CurrentPeriodStart: q.Status.CurrentPeriodStart,
			CurrentPeriodEnd:   q.Status.CurrentPeriodEnd,
		}, nil
	}

	if q.Action == "downgrade" {
		txRes, err := s.repo.ApplyImmediateSubscriptionDowngrade(ctx, tenantID, q.TargetPlanCode)
		if err != nil {
			return nil, err
		}
		updatedStatus, statusErr := s.buildTenantSubscriptionStatus(ctx, tenantID)
		if statusErr == nil {
			q.Status = updatedStatus
		}
		return &SubscriptionCheckoutResponse{
			TransactionID:      txRes.TransactionID,
			OrderID:            txRes.OrderID,
			PlanCode:           q.TargetPlanCode,
			Action:             q.Action,
			DurationMonth:      q.DurationMonth,
			Amount:             0,
			Status:             "paid",
			CreatedAt:          &txRes.CreatedAt,
			ActivePlan:         q.Status.ActivePlan,
			RemainingDays:      q.Status.RemainingDays,
			NextBillingDueDate: q.Status.NextBillingDueDate,
			CurrentPeriodStart: q.Status.CurrentPeriodStart,
			CurrentPeriodEnd:   q.Status.CurrentPeriodEnd,
			WarningMessage:     q.WarningMessage,
		}, nil
	}

	if q.Amount <= 0 {
		return nil, errors.New("nominal checkout tidak valid")
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return nil, errors.New("konfigurasi midtrans server key belum diatur")
	}

	activeTx, err := s.repo.GetLatestPendingSubscriptionTransaction(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if activeTx != nil {
		expiresAt := activeTx.CreatedAt.Add(15 * time.Minute)
		if activeTx.ExpiredAt != nil {
			expiresAt = *activeTx.ExpiredAt
		}
		return &SubscriptionCheckoutResponse{
			TransactionID:      activeTx.TransactionID,
			OrderID:            activeTx.OrderID,
			PlanCode:           activeTx.PlanCode,
			Amount:             activeTx.Amount,
			Status:             activeTx.Status,
			PaymentURL:         activeTx.PaymentURL,
			CreatedAt:          &activeTx.CreatedAt,
			ExpiredAt:          &expiresAt,
			Action:             q.Action,
			DurationMonth:      q.DurationMonth,
			ActivePlan:         q.Status.ActivePlan,
			RemainingDays:      q.Status.RemainingDays,
			NextBillingDueDate: q.Status.NextBillingDueDate,
			CurrentPeriodStart: q.Status.CurrentPeriodStart,
			CurrentPeriodEnd:   q.Status.CurrentPeriodEnd,
			WarningMessage:     q.WarningMessage,
		}, nil
	}

	orderID := "SUB-" + strings.ToUpper(uuid.NewString())
	txRes, err := s.repo.CreateSubscriptionTransaction(ctx, tenantID, orderID, q.TargetPlanCode, q.Amount, "pending")
	if err != nil {
		return nil, err
	}

	envType := midtrans.Sandbox
	if s.isCentralMidtransProduction() {
		envType = midtrans.Production
	}

	var snapClient snap.Client
	snapClient.New(serverKey, envType)

	plan := constant.SubscriptionPlans[q.TargetPlanCode]
	midtransReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(q.Amount),
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "SUB-" + strings.ToUpper(q.TargetPlanCode),
				Price: int64(q.Amount),
				Qty:   1,
				Name:  "Langganan " + plan.Name,
			},
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "minute",
			Duration: 15,
		},
	}

	snapResp, midtransErr := snapClient.CreateTransaction(midtransReq)
	if midtransErr != nil {
		s.log.Error("Midtrans subscription error: ", midtransErr.Message)
		return nil, errors.New("gagal membuat transaksi subscription")
	}

	if err := s.repo.UpdateSubscriptionPGInfo(ctx, txRes.TransactionID, snapResp.Token, snapResp.RedirectURL); err != nil {
		return nil, err
	}

	expiresAt := txRes.CreatedAt.Add(15 * time.Minute)

	return &SubscriptionCheckoutResponse{
		TransactionID:      txRes.TransactionID,
		OrderID:            txRes.OrderID,
		PlanCode:           txRes.PlanCode,
		Action:             q.Action,
		DurationMonth:      q.DurationMonth,
		Amount:             txRes.Amount,
		Status:             txRes.Status,
		PaymentURL:         &snapResp.RedirectURL,
		SnapToken:          &snapResp.Token,
		CreatedAt:          &txRes.CreatedAt,
		ExpiredAt:          &expiresAt,
		ActivePlan:         q.Status.ActivePlan,
		RemainingDays:      q.Status.RemainingDays,
		NextBillingDueDate: q.Status.NextBillingDueDate,
		CurrentPeriodStart: q.Status.CurrentPeriodStart,
		CurrentPeriodEnd:   q.Status.CurrentPeriodEnd,
		WarningMessage:     q.WarningMessage,
	}, nil
}

func (s *service) GetSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error) {
	txData, err := s.repo.GetSubscriptionTransaction(ctx, tenantID, transactionID)
	if err != nil {
		return nil, err
	}
	if txData.Status != "pending" {
		return txData, nil
	}

	midStatus, err := s.fetchMidtransStatus(txData.OrderID)
	if err == nil && midStatus.TransactionStatus != "" {
		_ = s.repo.ProcessSubscriptionWebhookTransaction(
			ctx,
			txData.OrderID,
			midStatus.TransactionStatus,
			midStatus.PaymentType,
			midStatus.TransactionID,
			midStatus,
		)
	}

	return s.repo.GetSubscriptionTransaction(ctx, tenantID, transactionID)
}

func (s *service) GetActiveSubscriptionTransaction(ctx context.Context, tenantID string) (*SubscriptionTransactionResponse, error) {
	txData, err := s.repo.GetLatestPendingSubscriptionTransaction(ctx, tenantID)
	if err != nil || txData == nil {
		return txData, err
	}

	midStatus, statusErr := s.fetchMidtransStatus(txData.OrderID)
	if statusErr == nil && midStatus.TransactionStatus != "" {
		_ = s.repo.ProcessSubscriptionWebhookTransaction(
			ctx,
			txData.OrderID,
			midStatus.TransactionStatus,
			midStatus.PaymentType,
			midStatus.TransactionID,
			midStatus,
		)
		return s.repo.GetLatestPendingSubscriptionTransaction(ctx, tenantID)
	}

	return txData, nil
}

func (s *service) CancelSubscriptionTransaction(ctx context.Context, tenantID, transactionID string) (*SubscriptionTransactionResponse, error) {
	txData, err := s.repo.GetSubscriptionTransaction(ctx, tenantID, transactionID)
	if err != nil {
		return nil, err
	}
	if txData.Status != "pending" {
		return nil, errors.New("hanya transaksi pending yang dapat dibatalkan")
	}

	_ = s.cancelMidtransTransaction(txData.OrderID)
	return s.repo.CancelSubscriptionTransaction(ctx, tenantID, transactionID)
}

func (s *service) prepareSubscriptionQuote(ctx context.Context, tenantID string, req SubscriptionQuotePayload) (*subscriptionQuoteInternal, error) {
	planCode := strings.TrimSpace(strings.ToLower(req.PlanCode))
	if planCode == "" {
		return nil, errors.New("paket langganan tidak valid")
	}
	plan, exists := constant.SubscriptionPlans[planCode]
	if !exists {
		return nil, errors.New("paket langganan tidak valid")
	}

	duration := req.DurationMonth
	if duration == 0 {
		duration = 1
	}
	if duration != 1 && duration != 3 && duration != 6 && duration != 12 {
		return nil, errors.New("durasi langganan hanya boleh 1, 3, 6, atau 12 bulan")
	}

	status, err := s.buildTenantSubscriptionStatus(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	amount := plan.Price
	action := "new"
	var warning *string
	var prorate float64

	if planCode == constant.PlanFree {
		action = "free"
		duration = 1
		amount = 0
	} else if planCode == status.ActivePlan {
		action = "renew"
		amount = plan.Price * float64(duration)
	} else {
		currentPlan, existsCurrent := constant.SubscriptionPlans[status.ActivePlan]
		if !existsCurrent {
			currentPlan = constant.SubscriptionPlans[constant.PlanFree]
		}

		if plan.Price > currentPlan.Price {
			action = "upgrade"
			duration = 1
			if status.RemainingDays != nil && currentPlan.Price > 0 {
				prorate = (plan.Price / 30.0) * float64(*status.RemainingDays)
				amount = math.Round(prorate)
			} else {
				amount = plan.Price
			}
		} else {
			action = "downgrade"
			duration = 1
			amount = 0
			msg := "Downgrade tidak mendapat kompensasi. Paket saat ini tetap berlaku sampai akhir periode aktif."
			warning = &msg
		}
	}

	return &subscriptionQuoteInternal{
		TargetPlanCode: planCode,
		Action:         action,
		DurationMonth:  duration,
		Amount:         amount,
		ProrateAmount:  prorate,
		WarningMessage: warning,
		Status:         status,
	}, nil
}

func (s *service) buildTenantSubscriptionStatus(ctx context.Context, tenantID string) (TenantSubscriptionStatus, error) {
	activePlan, err := s.repo.GetTenantSubscriptionPlan(ctx, tenantID)
	if err != nil {
		return TenantSubscriptionStatus{}, err
	}
	if _, ok := constant.SubscriptionPlans[activePlan]; !ok {
		activePlan = constant.PlanFree
	}

	status := TenantSubscriptionStatus{ActivePlan: activePlan}
	paidTransactions, err := s.repo.ListPaidSubscriptionTransactions(ctx, tenantID)
	if err != nil {
		return status, err
	}
	if len(paidTransactions) == 0 {
		return status, nil
	}

	var periodStart *time.Time
	var periodEnd *time.Time
	for _, tx := range paidTransactions {
		if tx.PaidAt == nil {
			continue
		}

		planDetail, ok := constant.SubscriptionPlans[tx.PlanCode]
		if !ok || planDetail.Price <= 0 || tx.Amount <= 0 {
			continue
		}

		monthsFloat := tx.Amount / planDetail.Price
		monthsRounded := int(math.Round(monthsFloat))
		if monthsRounded <= 0 || math.Abs(monthsFloat-float64(monthsRounded)) > 0.001 {
			continue
		}

		paidAt := tx.PaidAt.UTC()
		if periodStart == nil || periodEnd == nil {
			start := paidAt
			end := start.AddDate(0, monthsRounded, 0)
			periodStart = &start
			periodEnd = &end
			continue
		}

		if paidAt.Before(*periodEnd) || paidAt.Equal(*periodEnd) {
			end := periodEnd.AddDate(0, monthsRounded, 0)
			periodEnd = &end
			continue
		}

		start := paidAt
		end := start.AddDate(0, monthsRounded, 0)
		periodStart = &start
		periodEnd = &end
	}

	if periodStart == nil || periodEnd == nil {
		// Fallback: derive from latest successful paid transaction with amount > 0.
		for i := len(paidTransactions) - 1; i >= 0; i-- {
			tx := paidTransactions[i]
			if tx.PaidAt == nil || tx.Amount <= 0 {
				continue
			}
			planDetail, ok := constant.SubscriptionPlans[tx.PlanCode]
			if !ok || planDetail.Price <= 0 {
				continue
			}
			monthsFloat := tx.Amount / planDetail.Price
			monthsRounded := int(math.Round(monthsFloat))
			if monthsRounded < 1 {
				monthsRounded = 1
			}
			start := tx.PaidAt.UTC()
			end := start.AddDate(0, monthsRounded, 0)
			periodStart = &start
			periodEnd = &end
			break
		}
	}
	if periodStart == nil || periodEnd == nil {
		return status, nil
	}

	now := time.Now().UTC()
	remaining := int(math.Ceil(periodEnd.Sub(now).Hours() / 24.0))
	if remaining < 0 {
		remaining = 0
	}
	status.RemainingDays = &remaining
	status.CurrentPeriodStart = periodStart
	status.CurrentPeriodEnd = periodEnd
	status.NextBillingDueDate = periodEnd

	return status, nil
}

// ==========================================
// WEBHOOK SIGNATURE VALIDATION
// ==========================================

func (s *service) HandleMidtransWebhook(ctx context.Context, payload MidtransNotificationPayload) error {
	if strings.HasPrefix(payload.OrderID, "SUB-") {
		tenantID, err := s.repo.GetTenantIDBySubscriptionOrderID(ctx, payload.OrderID)
		if err != nil {
			s.log.Warn("Webhook subscription ditolak: transaksi tidak ditemukan untuk order ", payload.OrderID)
			return errors.New("transaksi tidak ditemukan")
		}

		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		if serverKey == "" {
			s.log.Error("Webhook subscription gagal: server key tidak ditemukan untuk tenant ", tenantID)
			return errors.New("server key tidak dikonfigurasi")
		}

		hashInput := payload.OrderID + payload.StatusCode + payload.GrossAmount + serverKey
		hasher := sha512.New()
		hasher.Write([]byte(hashInput))
		expectedSignature := hex.EncodeToString(hasher.Sum(nil))
		if payload.SignatureKey != expectedSignature {
			s.log.Warnf("Webhook subscription signature tidak cocok. OrderID: %s", payload.OrderID)
			return errors.New("invalid signature key")
		}

		if err := s.repo.ProcessSubscriptionWebhookTransaction(ctx, payload.OrderID, payload.TransactionStatus, payload.PaymentType, payload.TransactionID, payload); err != nil {
			s.log.Error("Gagal memproses webhook subscription: ", err)
			return err
		}
		s.log.Infof("Webhook subscription sukses diproses. OrderID: %s, Status: %s", payload.OrderID, payload.TransactionStatus)
		return nil
	}

	// 1. Cari Tenant ID dari tabel transaksi
	tenantID, err := s.repo.GetTenantIDByTransactionID(ctx, payload.OrderID)
	if err != nil {
		s.log.Warn("Webhook ditolak: Transaksi tidak ditemukan untuk OrderID ", payload.OrderID)
		return errors.New("transaksi tidak ditemukan")
	}

	// 2. Tentukan Server Key yang dipakai (Pusat atau Mandiri)
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	pgConfig, err := s.repo.GetPGConfig(ctx, tenantID)
	if err == nil && pgConfig != nil && !pgConfig.UseCentralPG && pgConfig.IsActive {
		serverKey = pgConfig.ServerKey
	}

	if serverKey == "" {
		s.log.Error("Webhook gagal: Server Key tidak ditemukan untuk Tenant ", tenantID)
		return errors.New("server key tidak dikonfigurasi")
	}

	// 3. Validasi Keamanan (Signature Key SHA512)
	// Rumus Midtrans: SHA512(order_id + status_code + gross_amount + server_key)
	hashInput := payload.OrderID + payload.StatusCode + payload.GrossAmount + serverKey
	hasher := sha512.New()
	hasher.Write([]byte(hashInput))
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	if payload.SignatureKey != expectedSignature {
		s.log.Warnf("Webhook HACK ATTEMPT! Signature tidak cocok. OrderID: %s", payload.OrderID)
		return errors.New("invalid signature key")
	}

	// 4. Jika aman, lanjutkan eksekusi DB Transaction (Row Locking)
	err = s.repo.ProcessWebhookTransaction(ctx, payload.OrderID, payload.TransactionStatus, payload.PaymentType)
	if err != nil {
		s.log.Error("Gagal mengeksekusi DB Transaction Webhook: ", err)
		return err
	}

	s.log.Infof("Webhook sukses diproses. OrderID: %s, Status: %s", payload.OrderID, payload.TransactionStatus)
	return nil
}

func (s *service) fetchMidtransStatus(orderID string) (*midtransStatusResponse, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return nil, errors.New("konfigurasi midtrans server key belum diatur")
	}

	url := s.midtransCoreAPIBaseURL() + "/v2/" + orderID + "/status"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, errors.New("gagal mengambil status dari midtrans")
	}

	var out midtransStatusResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *service) cancelMidtransTransaction(orderID string) error {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return errors.New("konfigurasi midtrans server key belum diatur")
	}

	url := s.midtransCoreAPIBaseURL() + "/v2/" + orderID + "/cancel"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return errors.New("gagal membatalkan transaksi di midtrans")
	}
	return nil
}
