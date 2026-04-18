package community

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	ListEvents(c *fiber.Ctx) error
	CreateEvent(c *fiber.Ctx) error
	GetEvent(c *fiber.Ctx) error
	UpdateEvent(c *fiber.Ctx) error
	UpdateEventStatus(c *fiber.Ctx) error
	DeleteEvent(c *fiber.Ctx) error
	ListPublicEvents(c *fiber.Ctx) error

	ListGalleryAlbums(c *fiber.Ctx) error
	CreateGalleryAlbum(c *fiber.Ctx) error
	GetGalleryAlbum(c *fiber.Ctx) error
	UpdateGalleryAlbum(c *fiber.Ctx) error
	DeleteGalleryAlbum(c *fiber.Ctx) error
	ListPublicGalleryAlbums(c *fiber.Ctx) error

	ListGalleryItems(c *fiber.Ctx) error
	CreateGalleryItem(c *fiber.Ctx) error
	GetGalleryItem(c *fiber.Ctx) error
	UpdateGalleryItem(c *fiber.Ctx) error
	DeleteGalleryItem(c *fiber.Ctx) error
	ListPublicGalleryItems(c *fiber.Ctx) error

	ListManagementMembers(c *fiber.Ctx) error
	CreateManagementMember(c *fiber.Ctx) error
	GetManagementMember(c *fiber.Ctx) error
	UpdateManagementMember(c *fiber.Ctx) error
	DeleteManagementMember(c *fiber.Ctx) error
	ListPublicManagementMembers(c *fiber.Ctx) error
}

type controller struct {
	service Service
	log     *logrus.Logger
}

func NewController(service Service, log *logrus.Logger) Controller {
	return &controller{service: service, log: log}
}

func (ctrl *controller) ListEvents(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListEvents(c.Context(), tenantID, EventListQuery{Status: c.Query("status"), Category: c.Query("category"), Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "events loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateEvent(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req EventPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Title == "" || req.StartDate == "" || req.Category == "" || req.TimeMode == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "title, category, start_date, time_mode are required"}})
	}
	item, err := ctrl.service.CreateEvent(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "event created", item, nil)
}

func (ctrl *controller) GetEvent(c *fiber.Ctx) error {
	return ctrl.getEventByID(c, func(tenantID string, id int64) (interface{}, error) {
		return ctrl.service.GetEvent(c.Context(), tenantID, id)
	}, "event loaded")
}

func (ctrl *controller) UpdateEvent(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	var req EventPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateEvent(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "event updated", nil, nil)
}

func (ctrl *controller) UpdateEventStatus(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	var req UpdateEventStatusRequest
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "status", Message: "status is required"}})
	}
	if err := ctrl.service.UpdateEventStatus(c.Context(), tenantID, id, req.Status); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "event status updated", nil, nil)
}

func (ctrl *controller) DeleteEvent(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteEvent(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "event deleted", nil, nil)
}

func (ctrl *controller) ListPublicEvents(c *fiber.Ctx) error {
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListPublicEvents(c.Context(), c.Params("hostname"), page, limit)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public events loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) ListGalleryAlbums(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListGalleryAlbums(c.Context(), tenantID, BaseListQuery{Page: page, Limit: limit, Search: c.Query("search")})
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "gallery albums loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateGalleryAlbum(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req GalleryAlbumPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Title == "" || req.MediaKind == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "title and media_kind are required"}})
	}
	item, err := ctrl.service.CreateGalleryAlbum(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "gallery album created", item, nil)
}

func (ctrl *controller) GetGalleryAlbum(c *fiber.Ctx) error {
	return ctrl.getGenericByID(c, func(tenantID string, id int64) (interface{}, error) {
		return ctrl.service.GetGalleryAlbum(c.Context(), tenantID, id)
	}, "gallery album loaded")
}

func (ctrl *controller) UpdateGalleryAlbum(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	var req GalleryAlbumPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateGalleryAlbum(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "gallery album updated", nil, nil)
}

func (ctrl *controller) DeleteGalleryAlbum(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteGalleryAlbum(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "gallery album deleted", nil, nil)
}

func (ctrl *controller) ListPublicGalleryAlbums(c *fiber.Ctx) error {
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListPublicGalleryAlbums(c.Context(), c.Params("hostname"), page, limit)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public gallery albums loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) ListGalleryItems(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListGalleryItems(c.Context(), tenantID, BaseListQuery{Page: page, Limit: limit, Search: c.Query("search")}, c.Query("album_id"))
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "gallery items loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateGalleryItem(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req GalleryItemPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.MediaType == "" || req.MediaURL == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "media_type and media_url are required"}})
	}
	item, err := ctrl.service.CreateGalleryItem(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "gallery item created", item, nil)
}

func (ctrl *controller) GetGalleryItem(c *fiber.Ctx) error {
	return ctrl.getGenericByID(c, func(tenantID string, id int64) (interface{}, error) {
		return ctrl.service.GetGalleryItem(c.Context(), tenantID, id)
	}, "gallery item loaded")
}

func (ctrl *controller) UpdateGalleryItem(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	var req GalleryItemPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateGalleryItem(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "gallery item updated", nil, nil)
}

func (ctrl *controller) DeleteGalleryItem(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteGalleryItem(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "gallery item deleted", nil, nil)
}

func (ctrl *controller) ListPublicGalleryItems(c *fiber.Ctx) error {
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListPublicGalleryItems(c.Context(), c.Params("hostname"), page, limit)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public gallery items loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) ListManagementMembers(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListManagementMembers(c.Context(), tenantID, BaseListQuery{Page: page, Limit: limit, Search: c.Query("search")})
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "management members loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateManagementMember(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req ManagementMemberPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.FullName == "" || req.RoleTitle == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "full_name and role_title are required"}})
	}
	item, err := ctrl.service.CreateManagementMember(c.Context(), tenantID, req)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "management member created", item, nil)
}

func (ctrl *controller) GetManagementMember(c *fiber.Ctx) error {
	return ctrl.getGenericByID(c, func(tenantID string, id int64) (interface{}, error) {
		return ctrl.service.GetManagementMember(c.Context(), tenantID, id)
	}, "management member loaded")
}

func (ctrl *controller) UpdateManagementMember(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	var req ManagementMemberPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateManagementMember(c.Context(), tenantID, id, req); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "management member updated", nil, nil)
}

func (ctrl *controller) DeleteManagementMember(c *fiber.Ctx) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteManagementMember(c.Context(), tenantID, id); err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "management member deleted", nil, nil)
}

func (ctrl *controller) ListPublicManagementMembers(c *fiber.Ctx) error {
	page, limit := parsePageLimit(c)
	items, total, err := ctrl.service.ListPublicManagementMembers(c.Context(), c.Params("hostname"), page, limit)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public management members loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) getEventByID(c *fiber.Ctx, f func(string, int64) (interface{}, error), msg string) error {
	return ctrl.getGenericByID(c, f, msg)
}

func (ctrl *controller) getGenericByID(c *fiber.Ctx, f func(string, int64) (interface{}, error), msg string) error {
	tenantID, id, ok := parseTenantIDAndID(c)
	if !ok {
		return nil
	}
	data, err := f(tenantID, id)
	if err != nil {
		return handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, msg, data, nil)
}

func parseTenantIDAndID(c *fiber.Ctx) (string, int64, bool) {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		_ = response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		return "", 0, false
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		_ = response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
		return "", 0, false
	}
	return tenantID, id, true
}

func parsePageLimit(c *fiber.Ctx) (int, int) {
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	return page, limit
}

func parseIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	v := c.Query(key)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}

func tenantIDFromLocals(c *fiber.Ctx) (string, error) {
	tenantID := fmt.Sprint(c.Locals("tenant_id"))
	if tenantID == "" || tenantID == "<nil>" {
		return "", errors.New("tenant_id not found")
	}
	return tenantID, nil
}

func handleError(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrNotFound) {
		return response.Error(c, fiber.StatusNotFound, "resource not found")
	}
	if errors.Is(err, ErrValidation) {
		return response.Error(c, fiber.StatusBadRequest, "validation failed")
	}
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
