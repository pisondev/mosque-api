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
	// CreateDonation(c *fiber.Ctx) error // Nanti di Tahap 4
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
