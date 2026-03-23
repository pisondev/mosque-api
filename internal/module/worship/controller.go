package worship

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	GetPrayerTimeSettings(c *fiber.Ctx) error
	UpsertPrayerTimeSettings(c *fiber.Ctx) error
	ListPrayerTimesDaily(c *fiber.Ctx) error
	CreatePrayerTimesDaily(c *fiber.Ctx) error
	GetPrayerTimesDaily(c *fiber.Ctx) error
	UpdatePrayerTimesDaily(c *fiber.Ctx) error
	DeletePrayerTimesDaily(c *fiber.Ctx) error
	ListPrayerDuties(c *fiber.Ctx) error
	CreatePrayerDuty(c *fiber.Ctx) error
	GetPrayerDuty(c *fiber.Ctx) error
	UpdatePrayerDuty(c *fiber.Ctx) error
	DeletePrayerDuty(c *fiber.Ctx) error
	ListSpecialDays(c *fiber.Ctx) error
	CreateSpecialDay(c *fiber.Ctx) error
	GetSpecialDay(c *fiber.Ctx) error
	UpdateSpecialDay(c *fiber.Ctx) error
	DeleteSpecialDay(c *fiber.Ctx) error
	GetPrayerCalendar(c *fiber.Ctx) error
}

type controller struct {
	service Service
	log     *logrus.Logger
}

func NewController(service Service, log *logrus.Logger) Controller {
	return &controller{service: service, log: log}
}

func (ctrl *controller) GetPrayerTimeSettings(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	item, err := ctrl.service.GetPrayerTimeSettings(c.Context(), tenantID)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer time settings loaded", item, nil)
}

func (ctrl *controller) UpsertPrayerTimeSettings(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req PrayerTimeSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Timezone == "" || req.LocationMode == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "timezone", Message: "timezone is required"},
			{Field: "location_mode", Message: "location_mode is required"},
		})
	}
	item, err := ctrl.service.UpsertPrayerTimeSettings(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer time settings saved", item, nil)
}

func (ctrl *controller) ListPrayerTimesDaily(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	items, total, err := ctrl.service.ListPrayerTimesDaily(c.Context(), tenantID, PrayerTimesDailyQuery{
		From:  c.Query("from"),
		To:    c.Query("to"),
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer times daily loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreatePrayerTimesDaily(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req PrayerTimesDailyPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.DayDate == "" || req.SubuhTime == "" || req.DzuhurTime == "" || req.AsharTime == "" || req.MaghribTime == "" || req.IsyaTime == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "day_date and mandatory prayer times are required"}})
	}
	item, err := ctrl.service.CreatePrayerTimesDaily(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "prayer times daily created", item, nil)
}

func (ctrl *controller) GetPrayerTimesDaily(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	item, err := ctrl.service.GetPrayerTimesDaily(c.Context(), tenantID, id)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer times daily loaded", item, nil)
}

func (ctrl *controller) UpdatePrayerTimesDaily(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req PrayerTimesDailyPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdatePrayerTimesDaily(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer times daily updated", nil, nil)
}

func (ctrl *controller) DeletePrayerTimesDaily(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	if err := ctrl.service.DeletePrayerTimesDaily(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer times daily deleted", nil, nil)
}

func (ctrl *controller) ListPrayerDuties(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	items, total, err := ctrl.service.ListPrayerDuties(c.Context(), tenantID, PrayerDutiesQuery{
		From:     c.Query("from"),
		To:       c.Query("to"),
		Category: c.Query("category"),
		Prayer:   c.Query("prayer"),
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer duties loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreatePrayerDuty(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req PrayerDutyPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Category == "" || req.DutyDate == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "category and duty_date are required"}})
	}
	item, err := ctrl.service.CreatePrayerDuty(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "prayer duty created", item, nil)
}

func (ctrl *controller) GetPrayerDuty(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	item, err := ctrl.service.GetPrayerDuty(c.Context(), tenantID, id)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer duty loaded", item, nil)
}

func (ctrl *controller) UpdatePrayerDuty(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req PrayerDutyPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdatePrayerDuty(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer duty updated", nil, nil)
}

func (ctrl *controller) DeletePrayerDuty(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	if err := ctrl.service.DeletePrayerDuty(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer duty deleted", nil, nil)
}

func (ctrl *controller) ListSpecialDays(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	items, total, err := ctrl.service.ListSpecialDays(c.Context(), tenantID, SpecialDaysQuery{
		Year:  c.Query("year"),
		Kind:  c.Query("kind"),
		From:  c.Query("from"),
		To:    c.Query("to"),
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "special days loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateSpecialDay(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req SpecialDayPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Kind == "" || req.Title == "" || req.DayDate == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "kind, title, day_date are required"}})
	}
	item, err := ctrl.service.CreateSpecialDay(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "special day created", item, nil)
}

func (ctrl *controller) GetSpecialDay(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	item, err := ctrl.service.GetSpecialDay(c.Context(), tenantID, id)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "special day loaded", item, nil)
}

func (ctrl *controller) UpdateSpecialDay(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req SpecialDayPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateSpecialDay(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "special day updated", nil, nil)
}

func (ctrl *controller) DeleteSpecialDay(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	if err := ctrl.service.DeleteSpecialDay(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "special day deleted", nil, nil)
}

func (ctrl *controller) GetPrayerCalendar(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	data, err := ctrl.service.GetPrayerCalendar(c.Context(), tenantID, c.Query("from"), c.Query("to"))
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "prayer calendar loaded", data, nil)
}

func handleError(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrNotFound) {
		return response.Error(c, fiber.StatusNotFound, "resource not found")
	}
	if errors.Is(err, ErrConflict) {
		return response.Error(c, fiber.StatusConflict, "resource conflict")
	}
	if errors.Is(err, ErrValidation) {
		return response.Error(c, fiber.StatusBadRequest, "validation failed")
	}
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}

func tenantIDFromLocals(c *fiber.Ctx) (string, error) {
	tenantID := fmt.Sprint(c.Locals("tenant_id"))
	if tenantID == "" || tenantID == "<nil>" {
		return "", errors.New("tenant_id not found")
	}
	return tenantID, nil
}

func parseIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return v
}
