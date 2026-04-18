package finance

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	GetPGConfig(c *fiber.Ctx) error
	UpsertPGConfig(c *fiber.Ctx) error

	ListCampaigns(c *fiber.Ctx) error
	ListPublicCampaigns(c *fiber.Ctx) error
	CreateCampaign(c *fiber.Ctx) error
	GetCampaign(c *fiber.Ctx) error
	GetPublicCampaignBySlug(c *fiber.Ctx) error
	UpdateCampaign(c *fiber.Ctx) error

	ListTransactions(c *fiber.Ctx) error
	ListPublicDonors(c *fiber.Ctx) error
	CreateDonation(c *fiber.Ctx) error

	ListSubscriptionPlans(c *fiber.Ctx) error
	GetSubscriptionQuote(c *fiber.Ctx) error
	CreateSubscriptionCheckout(c *fiber.Ctx) error
	CreateSubscriptionCheckoutFromQuote(c *fiber.Ctx) error
	ListSubscriptionTransactions(c *fiber.Ctx) error
	GetSubscriptionTransaction(c *fiber.Ctx) error
	GetActiveSubscriptionTransaction(c *fiber.Ctx) error
	CancelSubscriptionTransaction(c *fiber.Ctx) error
	ActivateFreePlan(c *fiber.Ctx) error

	MidtransWebhook(c *fiber.Ctx) error
}

type controller struct {
	svc Service
	log *logrus.Logger
}

func NewController(svc Service, log *logrus.Logger) Controller {
	return &controller{svc: svc, log: log}
}

// Helper untuk mengambil tenant_id dari middleware auth
func getTenantID(c *fiber.Ctx) string {
	val := c.Locals("tenant_id")
	if val == nil {
		return ""
	}

	// Safe type assertion
	tenantID, ok := val.(string)
	if !ok {
		// Fallback jika karena suatu alasan JWT mem-parsingnya bukan sebagai string murni
		return fmt.Sprintf("%v", val)
	}
	return tenantID
}

// ==========================================
// PG CONFIGURATIONS
// ==========================================

func (ctrl *controller) GetPGConfig(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	res, err := ctrl.svc.GetPGConfig(c.Context(), tenantID)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil konfigurasi payment gateway")
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil konfigurasi payment gateway", res, nil)
}

func (ctrl *controller) UpsertPGConfig(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	var req PGConfigPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}

	err := ctrl.svc.UpsertPGConfig(c.Context(), tenantID, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal menyimpan konfigurasi payment gateway")
	}
	return response.Success(c, fiber.StatusOK, "Konfigurasi payment gateway berhasil disimpan", nil, nil)
}

// ==========================================
// DONATION CAMPAIGNS
// ==========================================

func (ctrl *controller) CreateCampaign(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	var req CampaignPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}

	res, err := ctrl.svc.CreateCampaign(c.Context(), tenantID, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal membuat kampanye donasi")
	}
	return response.Success(c, fiber.StatusCreated, "Kampanye donasi berhasil dibuat", res, nil)
}

func (ctrl *controller) GetCampaign(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID kampanye tidak valid")
	}

	res, err := ctrl.svc.GetCampaign(c.Context(), tenantID, id)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusNotFound, "Kampanye donasi tidak ditemukan")
	}
	return response.Success(c, fiber.StatusOK, "Detail kampanye donasi", res, nil)
}

func (ctrl *controller) UpdateCampaign(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID kampanye tidak valid")
	}

	var req CampaignPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}

	err = ctrl.svc.UpdateCampaign(c.Context(), tenantID, id, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal memperbarui kampanye donasi")
	}
	return response.Success(c, fiber.StatusOK, "Kampanye donasi berhasil diperbarui", nil, nil)
}

func (ctrl *controller) GetPublicCampaignBySlug(c *fiber.Ctx) error {
	hostname := c.Params("hostname")
	slug := c.Params("slug")

	res, err := ctrl.svc.GetPublicCampaignBySlug(c.Context(), hostname, slug)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusNotFound, "Kampanye donasi tidak ditemukan")
	}
	return response.Success(c, fiber.StatusOK, "Detail kampanye donasi", res, nil)
}

// ==========================================
// LIST METHODS (Placeholder)
// ==========================================

// Helper untuk extract pagination
func getPagination(c *fiber.Ctx) ListQuery {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	return ListQuery{Page: page, Limit: limit}
}

func (ctrl *controller) ListCampaigns(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	q := getPagination(c)

	data, total, err := ctrl.svc.ListCampaigns(c.Context(), tenantID, q)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil daftar kampanye")
	}

	meta := fiber.Map{"page": q.Page, "limit": q.Limit, "total": total}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar kampanye", data, meta)
}

func (ctrl *controller) ListPublicCampaigns(c *fiber.Ctx) error {
	hostname := c.Params("hostname")
	q := getPagination(c)

	data, total, err := ctrl.svc.ListPublicCampaigns(c.Context(), hostname, q)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil daftar kampanye")
	}

	meta := fiber.Map{"page": q.Page, "limit": q.Limit, "total": total}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar kampanye", data, meta)
}

func (ctrl *controller) ListTransactions(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	campaignID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID kampanye tidak valid")
	}
	q := getPagination(c)

	data, total, err := ctrl.svc.ListTransactions(c.Context(), tenantID, campaignID, q)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil daftar transaksi")
	}

	meta := fiber.Map{"page": q.Page, "limit": q.Limit, "total": total}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar transaksi", data, meta)
}

func (ctrl *controller) ListPublicDonors(c *fiber.Ctx) error {
	hostname := c.Params("hostname")
	campaignID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID kampanye tidak valid")
	}
	q := getPagination(c)

	data, total, err := ctrl.svc.ListPublicDonors(c.Context(), hostname, campaignID, q)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil daftar donatur")
	}

	meta := fiber.Map{"page": q.Page, "limit": q.Limit, "total": total}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar donatur", data, meta)
}

func (ctrl *controller) CreateDonation(c *fiber.Ctx) error {
	hostname := c.Params("hostname")
	campaignID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID kampanye tidak valid")
	}

	var req DonatePayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}
	req.CampaignID = campaignID // Pastikan ID dari URL dimasukkan ke payload

	// Validasi dasar
	if req.Amount < 10000 {
		return response.Error(c, fiber.StatusBadRequest, "Minimal donasi adalah Rp 10.000")
	}

	res, err := ctrl.svc.CreateDonation(c.Context(), hostname, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Checkout berhasil, silakan lanjutkan pembayaran", res, nil)
}

func (ctrl *controller) ListSubscriptionPlans(c *fiber.Ctx) error {
	data, err := ctrl.svc.ListSubscriptionPlans(c.Context())
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil paket langganan")
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil paket langganan", data, nil)
}

func (ctrl *controller) CreateSubscriptionCheckout(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	var req CreateSubscriptionCheckoutPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}
	res, err := ctrl.svc.CreateSubscriptionCheckout(c.Context(), tenantID, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Checkout paket berhasil dibuat", res, nil)
}

func (ctrl *controller) GetSubscriptionQuote(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	var req SubscriptionQuotePayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}
	res, err := ctrl.svc.GetSubscriptionQuote(c.Context(), tenantID, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Quote paket berhasil dihitung", res, nil)
}

func (ctrl *controller) CreateSubscriptionCheckoutFromQuote(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	var req SubscriptionQuotePayload
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}
	res, err := ctrl.svc.CreateSubscriptionCheckoutFromQuote(c.Context(), tenantID, req)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Checkout paket berhasil dibuat", res, nil)
}

func (ctrl *controller) ListSubscriptionTransactions(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	q := getPagination(c)
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 50 {
		q.Limit = 50
	}
	res, total, err := ctrl.svc.ListSubscriptionTransactions(c.Context(), tenantID, q)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil riwayat transaksi subscription")
	}
	meta := fiber.Map{"page": q.Page, "limit": q.Limit, "total": total}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil riwayat transaksi subscription", res, meta)
}

func (ctrl *controller) GetSubscriptionTransaction(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id := c.Params("id")
	res, err := ctrl.svc.GetSubscriptionTransaction(c.Context(), tenantID, id)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusNotFound, "Transaksi subscription tidak ditemukan")
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil status transaksi subscription", res, nil)
}

func (ctrl *controller) GetActiveSubscriptionTransaction(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	res, err := ctrl.svc.GetActiveSubscriptionTransaction(c.Context(), tenantID)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengambil transaksi subscription aktif")
	}
	if res == nil {
		return response.Success(c, fiber.StatusOK, "Tidak ada transaksi subscription aktif", nil, nil)
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil transaksi subscription aktif", res, nil)
}

func (ctrl *controller) CancelSubscriptionTransaction(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id := c.Params("id")
	res, err := ctrl.svc.CancelSubscriptionTransaction(c.Context(), tenantID, id)
	if err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Transaksi subscription berhasil dibatalkan", res, nil)
}

func (ctrl *controller) ActivateFreePlan(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	if err := ctrl.svc.ActivateFreePlan(c.Context(), tenantID); err != nil {
		ctrl.log.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengaktifkan paket free")
	}
	return response.Success(c, fiber.StatusOK, "Paket free berhasil diaktifkan", nil, nil)
}

// ==========================================
// WEBHOOK ENDPOINT
// ==========================================

func (ctrl *controller) MidtransWebhook(c *fiber.Ctx) error {
	var payload MidtransNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		ctrl.log.Error("Gagal parsing webhook payload: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid payload format"})
	}

	err := ctrl.svc.HandleMidtransWebhook(c.Context(), payload)
	if err != nil {
		// Jika signature salah, tolak dengan 403 Forbidden
		if err.Error() == "invalid signature key" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "invalid signature"})
		}
		if err.Error() == "transaksi tidak ditemukan" {
			ctrl.log.Warn("Webhook diabaikan karena order_id tidak terdaftar di sistem: ", payload.OrderID)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "webhook ignored: unknown order"})
		}
		// Midtrans menyarankan mengembalikan 200 OK meskipun ada error internal agar mereka tidak spam retry berulang kali
		// Tapi untuk strict development, kita bisa kembalikan 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "internal server error"})
	}

	// Balas dengan HTTP 200 OK agar Midtrans tahu kita sudah menerimanya dengan baik
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "webhook processed successfully"})
}
