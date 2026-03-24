package engagement

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	ListStaticPaymentMethods(c *fiber.Ctx) error
	CreateStaticPaymentMethod(c *fiber.Ctx) error
	GetStaticPaymentMethod(c *fiber.Ctx) error
	UpdateStaticPaymentMethod(c *fiber.Ctx) error
	DeleteStaticPaymentMethod(c *fiber.Ctx) error
	ListPublicStaticPaymentMethods(c *fiber.Ctx) error
	ListSocialLinks(c *fiber.Ctx) error
	CreateSocialLink(c *fiber.Ctx) error
	GetSocialLink(c *fiber.Ctx) error
	UpdateSocialLink(c *fiber.Ctx) error
	DeleteSocialLink(c *fiber.Ctx) error
	ListPublicSocialLinks(c *fiber.Ctx) error
	ListExternalLinks(c *fiber.Ctx) error
	CreateExternalLink(c *fiber.Ctx) error
	GetExternalLink(c *fiber.Ctx) error
	UpdateExternalLink(c *fiber.Ctx) error
	DeleteExternalLink(c *fiber.Ctx) error
	ListPublicExternalLinks(c *fiber.Ctx) error
	ListFeatureCatalog(c *fiber.Ctx) error
	ListWebsiteFeatures(c *fiber.Ctx) error
	UpsertWebsiteFeature(c *fiber.Ctx) error
	BulkUpsertWebsiteFeatures(c *fiber.Ctx) error
}

type controller struct {
	service Service
	log     *logrus.Logger
}

func NewController(service Service, log *logrus.Logger) Controller {
	return &controller{service: service, log: log}
}

func (ctrl *controller) ListStaticPaymentMethods(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	page, limit := pageLimit(c)
	items, total, err := ctrl.service.ListStaticPaymentMethods(c.Context(), tenantID, ListQuery{Page: page, Limit: limit})
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "donation channels loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateStaticPaymentMethod(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	var req StaticPaymentMethodPayload
	if err := c.BodyParser(&req); err != nil || req.ChannelType == "" || req.Label == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "channel_type and label are required"}})
	}
	item, err := ctrl.service.CreateStaticPaymentMethod(c.Context(), tenantID, req)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "donation channel created", item, nil)
}

func (ctrl *controller) GetStaticPaymentMethod(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	item, err := ctrl.service.GetStaticPaymentMethod(c.Context(), tenantID, id)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "donation channel loaded", item, nil)
}

func (ctrl *controller) UpdateStaticPaymentMethod(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	var req StaticPaymentMethodPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateStaticPaymentMethod(c.Context(), tenantID, id, req); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "donation channel updated", nil, nil)
}

func (ctrl *controller) DeleteStaticPaymentMethod(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteStaticPaymentMethod(c.Context(), tenantID, id); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "donation channel deleted", nil, nil)
}

func (ctrl *controller) ListPublicStaticPaymentMethods(c *fiber.Ctx) error {
	page, limit := pageLimit(c)
	items, total, err := ctrl.service.ListPublicStaticPaymentMethods(c.Context(), c.Params("hostname"), ListQuery{Page: page, Limit: limit})
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public donation channels loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) ListSocialLinks(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	page, limit := pageLimit(c)
	items, total, err := ctrl.service.ListSocialLinks(c.Context(), tenantID, ListQuery{Page: page, Limit: limit})
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "social links loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateSocialLink(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	var req SocialLinkPayload
	if err := c.BodyParser(&req); err != nil || req.Platform == "" || req.URL == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "platform and url are required"}})
	}
	item, err := ctrl.service.CreateSocialLink(c.Context(), tenantID, req)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "social link created", item, nil)
}

func (ctrl *controller) GetSocialLink(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	item, err := ctrl.service.GetSocialLink(c.Context(), tenantID, id)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "social link loaded", item, nil)
}

func (ctrl *controller) UpdateSocialLink(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	var req SocialLinkPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateSocialLink(c.Context(), tenantID, id, req); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "social link updated", nil, nil)
}

func (ctrl *controller) DeleteSocialLink(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteSocialLink(c.Context(), tenantID, id); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "social link deleted", nil, nil)
}

func (ctrl *controller) ListPublicSocialLinks(c *fiber.Ctx) error {
	page, limit := pageLimit(c)
	items, total, err := ctrl.service.ListPublicSocialLinks(c.Context(), c.Params("hostname"), ListQuery{Page: page, Limit: limit})
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public social links loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) ListExternalLinks(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	page, limit := pageLimit(c)
	items, total, err := ctrl.service.ListExternalLinks(c.Context(), tenantID, ListQuery{Page: page, Limit: limit})
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "external links loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) CreateExternalLink(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	var req ExternalLinkPayload
	if err := c.BodyParser(&req); err != nil || req.LinkType == "" || req.Label == "" || req.URL == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "payload", Message: "link_type, label, url are required"}})
	}
	item, err := ctrl.service.CreateExternalLink(c.Context(), tenantID, req)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "external link created", item, nil)
}

func (ctrl *controller) GetExternalLink(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	item, err := ctrl.service.GetExternalLink(c.Context(), tenantID, id)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "external link loaded", item, nil)
}

func (ctrl *controller) UpdateExternalLink(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	var req ExternalLinkPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpdateExternalLink(c.Context(), tenantID, id, req); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "external link updated", nil, nil)
}

func (ctrl *controller) DeleteExternalLink(c *fiber.Ctx) error {
	tenantID, id, ok := tenantAndID(c)
	if !ok {
		return nil
	}
	if err := ctrl.service.DeleteExternalLink(c.Context(), tenantID, id); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "external link deleted", nil, nil)
}

func (ctrl *controller) ListPublicExternalLinks(c *fiber.Ctx) error {
	page, limit := pageLimit(c)
	items, total, err := ctrl.service.ListPublicExternalLinks(c.Context(), c.Params("hostname"), ListQuery{Page: page, Limit: limit})
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "public external links loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

func (ctrl *controller) ListFeatureCatalog(c *fiber.Ctx) error {
	items, err := ctrl.service.ListFeatureCatalog(c.Context())
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "feature catalog loaded", items, nil)
}

func (ctrl *controller) ListWebsiteFeatures(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	items, err := ctrl.service.ListWebsiteFeatures(c.Context(), tenantID)
	if err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "website features loaded", items, nil)
}

func (ctrl *controller) UpsertWebsiteFeature(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	featureID, err := strconv.ParseInt(c.Params("feature_id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "feature_id", Message: "invalid feature_id"}})
	}
	var req WebsiteFeatureUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if err := ctrl.service.UpsertWebsiteFeature(c.Context(), tenantID, featureID, req); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "website feature updated", nil, nil)
}

func (ctrl *controller) BulkUpsertWebsiteFeatures(c *fiber.Ctx) error {
	tenantID, ok := getTenantID(c)
	if !ok {
		return nil
	}
	var req WebsiteFeatureBulkRequest
	if err := c.BodyParser(&req); err != nil || len(req.Items) == 0 {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "items", Message: "items is required"}})
	}
	if err := ctrl.service.BulkUpsertWebsiteFeatures(c.Context(), tenantID, req.Items); err != nil {
		return handleErr(c, err)
	}
	return response.Success(c, fiber.StatusOK, "website features updated", nil, nil)
}

func getTenantID(c *fiber.Ctx) (string, bool) {
	tenantID := fmt.Sprint(c.Locals("tenant_id"))
	if tenantID == "" || tenantID == "<nil>" {
		_ = response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		return "", false
	}
	return tenantID, true
}

func tenantAndID(c *fiber.Ctx) (string, int64, bool) {
	tenantID, ok := getTenantID(c)
	if !ok {
		return "", 0, false
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		_ = response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
		return "", 0, false
	}
	return tenantID, id, true
}

func pageLimit(c *fiber.Ctx) (int, int) {
	page := 1
	limit := 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}
	return page, limit
}

func handleErr(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrNotFound) {
		return response.Error(c, fiber.StatusNotFound, "resource not found")
	}
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
