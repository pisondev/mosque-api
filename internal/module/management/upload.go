package management

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pisondev/mosque-api/internal/constant"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/pisondev/mosque-api/internal/storage"
)

const maxUploadBytes = 2 * 1024 * 1024 // 2 MB

// resolveSubdomain fetches the subdomain for a tenantID directly from DB.
func resolveSubdomain(ctx context.Context, db *pgxpool.Pool, tenantID string) string {
	var subdomain string
	_ = db.QueryRow(ctx, `SELECT COALESCE(subdomain,'') FROM tenants WHERE id=$1`, tenantID).Scan(&subdomain)
	return subdomain
}

func resolveSubscriptionPlan(ctx context.Context, db *pgxpool.Pool, tenantID string) string {
	var plan string
	_ = db.QueryRow(ctx, `SELECT COALESCE(subscription_plan,'free') FROM tenants WHERE id=$1`, tenantID).Scan(&plan)
	if plan == "" {
		return constant.PlanFree
	}
	return plan
}

// UploadHandler handles multipart file uploads for tenant media.
// POST /tenant/upload?kind=header|management_photo
// Returns: { url: "https://..." }
func UploadHandler(store *storage.Client, db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID, err := tenantIDFromLocals(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}

		subdomain := resolveSubdomain(c.Context(), db, tenantID)
		if subdomain == "" {
			if len(tenantID) >= 8 {
				subdomain = tenantID[:8]
			} else {
				subdomain = tenantID
			}
		}

		kind := c.Query("kind", "header")
		if kind != "header" && kind != "management_photo" && kind != "event_poster" && kind != "qris" {
			return response.Error(c, fiber.StatusBadRequest, "jenis upload tidak dikenal (header|management_photo|event_poster|qris)")
		}

		planCode := resolveSubscriptionPlan(c.Context(), db, tenantID)
		planDetail, ok := constant.SubscriptionPlans[planCode]
		if !ok {
			planDetail = constant.SubscriptionPlans[constant.PlanFree]
		}
		if planDetail.StorageLimitMB <= 0 {
			return response.Error(c, fiber.StatusForbidden, "Paket saat ini belum memiliki kuota media. Upgrade paket untuk mengunggah gambar.")
		}

		file, err := c.FormFile("file")
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "file tidak ditemukan dalam request")
		}

		if file.Size > maxUploadBytes {
			return response.Error(c, fiber.StatusBadRequest, "ukuran file melebihi batas maksimal 2MB. Kompres gambar terlebih dahulu.")
		}

		currentUsedBytes, err := store.GetBucketSizeBytes(c.Context(), storage.TenantFolder(subdomain))
		if err != nil {
			return response.Error(c, fiber.StatusBadGateway, "gagal membaca kuota media saat ini")
		}

		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			return response.Error(c, fiber.StatusBadRequest, "hanya file gambar yang diizinkan")
		}

		src, err := file.Open()
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "gagal membaca file")
		}
		defer src.Close()

		ext := "jpg"
		if strings.Contains(contentType, "png") {
			ext = "png"
		} else if strings.Contains(contentType, "webp") {
			ext = "webp"
		} else if strings.Contains(contentType, "gif") {
			ext = "gif"
		}

		ts := time.Now().UnixMilli()
		var key string
		switch kind {
		case "management_photo":
			key = storage.ManagementPhotoKey(subdomain, fmt.Sprintf("%d.%s", ts, ext))
		case "event_poster":
			key = storage.EventPosterKey(subdomain, fmt.Sprintf("%d.%s", ts, ext))
		case "qris":
			key = storage.QrisImageKey(subdomain, fmt.Sprintf("%d.%s", ts, ext))
		default:
			key = storage.HeaderImageKey(subdomain, fmt.Sprintf("%d.%s", ts, ext))
		}

		// Evaluate quota first (replace-aware)
		oldURL := c.FormValue("old_url", "")
		var oldSizeBytes int64
		if oldURL != "" {
			if oldKey := storage.KeyFromURL(oldURL); oldKey != "" {
				oldSizeBytes, _ = store.GetObjectSizeBytes(c.Context(), oldKey)
			}
		}

		limitBytes := int64(planDetail.StorageLimitMB * 1024 * 1024)
		effectiveUsed := currentUsedBytes - oldSizeBytes + file.Size
		if effectiveUsed > limitBytes {
			return response.Error(c, fiber.StatusForbidden, fmt.Sprintf("Kuota media tidak cukup. Sisa kuota sekitar %.2f MB.", float64(limitBytes-currentUsedBytes)/(1024*1024)))
		}

		// Delete old file if provided (best-effort), then upload new one
		if oldURL != "" {
			if oldKey := storage.KeyFromURL(oldURL); oldKey != "" {
				_ = store.Delete(c.Context(), oldKey)
			}
		}

		publicURL, err := store.UploadStream(c.Context(), key, src, contentType)
		if err != nil {
			return response.Error(c, fiber.StatusBadGateway, "gagal mengunggah file ke storage: "+err.Error())
		}

		return response.Success(c, fiber.StatusOK, "file berhasil diunggah", fiber.Map{"url": publicURL}, nil)
	}
}

// StorageQuotaHandler returns real bucket usage (bytes + MB) for the tenant.
// GET /tenant/storage-quota
func StorageQuotaHandler(store *storage.Client, db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID, err := tenantIDFromLocals(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}

		subdomain := resolveSubdomain(c.Context(), db, tenantID)

		empty := fiber.Map{"used_bytes": int64(0), "used_mb": 0.0}

		if subdomain == "" {
			return response.Success(c, fiber.StatusOK, "storage quota", empty, nil)
		}

		used, err := store.GetBucketSizeBytes(c.Context(), storage.TenantFolder(subdomain))
		if err != nil {
			return response.Success(c, fiber.StatusOK, "storage quota (estimate)", empty, nil)
		}

		usedMB := float64(used) / (1024 * 1024)
		return response.Success(c, fiber.StatusOK, "storage quota", fiber.Map{
			"used_bytes": used,
			"used_mb":    usedMB,
		}, nil)
	}
}
