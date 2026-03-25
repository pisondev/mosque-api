package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pisondev/mosque-api/internal/constant"
	"github.com/pisondev/mosque-api/internal/response"
)

// RequireFeature adalah guard/satpam yang memastikan tenant memiliki paket yang sesuai
func RequireFeature(db *pgxpool.Pool, featureKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil tenant_id dari JWT yang sudah divalidasi oleh middleware Protected()
		val := c.Locals("tenant_id")
		if val == nil {
			return response.Error(c, fiber.StatusUnauthorized, "Sesi tidak valid atau tenant ID hilang")
		}

		// Safe type assertion
		tenantID, ok := val.(string)
		if !ok {
			tenantID = fmt.Sprintf("%v", val)
		}

		// 2. Cek paket langganan saat ini secara real-time ke database
		var currentPlan string
		err := db.QueryRow(c.Context(), "SELECT subscription_plan FROM tenants WHERE id = $1", tenantID).Scan(&currentPlan)
		if err != nil {
			// Jika error database atau tenant tidak ditemukan
			return response.Error(c, fiber.StatusInternalServerError, "Gagal memverifikasi status langganan")
		}

		// 3. Validasi dengan Kamus Single Source of Truth kita
		if !constant.HasFeature(currentPlan, featureKey) {
			// Jika paketnya tidak punya fitur ini, tolak dengan 403 Forbidden!
			return response.Error(c, fiber.StatusForbidden, "Paket langganan Anda tidak mencakup fitur ini. Silakan upgrade paket Anda.")
		}

		// 4. Jika lolos, izinkan lanjut ke Controller
		return c.Next()
	}
}
