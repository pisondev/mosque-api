package management

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	TenantMe(c *fiber.Ctx) error
	ListDomains(c *fiber.Ctx) error
	CreateDomain(c *fiber.Ctx) error
	UpdateDomain(c *fiber.Ctx) error
	DeleteDomain(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	UpsertProfile(c *fiber.Ctx) error
	ListTags(c *fiber.Ctx) error
	CreateTag(c *fiber.Ctx) error
	UpdateTag(c *fiber.Ctx) error
	DeleteTag(c *fiber.Ctx) error
	ListPosts(c *fiber.Ctx) error
	CreatePost(c *fiber.Ctx) error
	GetPost(c *fiber.Ctx) error
	UpdatePost(c *fiber.Ctx) error
	UpdatePostStatus(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
	ListStaticPages(c *fiber.Ctx) error
	UpsertStaticPage(c *fiber.Ctx) error
	SetupTenant(c *fiber.Ctx) error
}

type controller struct {
	service Service
	log     *logrus.Logger
}

func NewController(service Service, log *logrus.Logger) Controller {
	return &controller{service: service, log: log}
}

// TenantMe godoc
// @Summary Get tenant context
// @Tags Tenant
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/me [get]
func (ctrl *controller) TenantMe(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	email := strings.TrimSpace(fmt.Sprint(c.Locals("email")))
	role := strings.TrimSpace(fmt.Sprint(c.Locals("role")))

	data, err := ctrl.service.GetTenantMe(c.Context(), tenantID, email, role)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "tenant context loaded", data, nil)
}

// ListDomains godoc
// @Summary List tenant domains
// @Tags Domains
// @Security BearerAuth
// @Param status query string false "status"
// @Param domain_type query string false "domain_type"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/domains [get]
func (ctrl *controller) ListDomains(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	items, total, err := ctrl.service.ListDomains(c.Context(), tenantID, DomainListQuery{
		Status:     c.Query("status"),
		DomainType: c.Query("domain_type"),
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "domains loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

// CreateDomain godoc
// @Summary Create tenant domain
// @Tags Domains
// @Security BearerAuth
// @Param payload body CreateDomainRequest true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/domains [post]
func (ctrl *controller) CreateDomain(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req CreateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.DomainType == "" || req.Hostname == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "domain_type", Message: "domain_type is required"},
			{Field: "hostname", Message: "hostname is required"},
		})
	}
	item, err := ctrl.service.CreateDomain(c.Context(), tenantID, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "domain created", item, nil)
}

// UpdateDomain godoc
// @Summary Update tenant domain
// @Tags Domains
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body UpdateDomainRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/domains/{id} [patch]
func (ctrl *controller) UpdateDomain(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req UpdateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Status == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "status", Message: "status is required"}})
	}
	if err := ctrl.service.UpdateDomain(c.Context(), tenantID, id, req); err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "domain updated", fiber.Map{"id": id}, nil)
}

// DeleteDomain godoc
// @Summary Delete tenant domain
// @Tags Domains
// @Security BearerAuth
// @Param id path int true "id"
// @Success 204
// @Router /tenant/domains/{id} [delete]
func (ctrl *controller) DeleteDomain(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	if err := ctrl.service.DeleteDomain(c.Context(), tenantID, id); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetProfile godoc
// @Summary Get masjid profile
// @Tags Profile
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/profile [get]
func (ctrl *controller) GetProfile(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	item, err := ctrl.service.GetProfile(c.Context(), tenantID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "profile loaded", item, nil)
}

// UpsertProfile godoc
// @Summary Create or update masjid profile
// @Tags Profile
// @Security BearerAuth
// @Param payload body ProfileRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/profile [put]
func (ctrl *controller) UpsertProfile(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req ProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.OfficialName == "" || req.Kind == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "official_name", Message: "official_name is required"},
			{Field: "kind", Message: "kind is required"},
		})
	}
	item, err := ctrl.service.UpsertProfile(c.Context(), tenantID, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "profile updated", item, nil)
}

// ListTags godoc
// @Summary List tags
// @Tags Tags
// @Security BearerAuth
// @Param scope query string false "scope"
// @Param search query string false "search"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/tags [get]
func (ctrl *controller) ListTags(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	items, total, err := ctrl.service.ListTags(c.Context(), tenantID, c.Query("scope"), c.Query("search"), page, limit)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "tags loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

// CreateTag godoc
// @Summary Create tag
// @Tags Tags
// @Security BearerAuth
// @Param payload body CreateTagRequest true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/tags [post]
func (ctrl *controller) CreateTag(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Scope == "" || req.Name == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "scope", Message: "scope is required"},
			{Field: "name", Message: "name is required"},
		})
	}
	item, err := ctrl.service.CreateTag(c.Context(), tenantID, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "tag created", item, nil)
}

// UpdateTag godoc
// @Summary Update tag
// @Tags Tags
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body UpdateTagRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/tags/{id} [patch]
func (ctrl *controller) UpdateTag(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req UpdateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Name == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "name", Message: "name is required"}})
	}
	if err := ctrl.service.UpdateTag(c.Context(), tenantID, id, req); err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "tag updated", fiber.Map{"id": id}, nil)
}

// DeleteTag godoc
// @Summary Delete tag
// @Tags Tags
// @Security BearerAuth
// @Param id path int true "id"
// @Success 204
// @Router /tenant/tags/{id} [delete]
func (ctrl *controller) DeleteTag(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	if err := ctrl.service.DeleteTag(c.Context(), tenantID, id); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListPosts godoc
// @Summary List posts
// @Tags Posts
// @Security BearerAuth
// @Param category query string false "category"
// @Param status query string false "status"
// @Param search query string false "search"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/posts [get]
func (ctrl *controller) ListPosts(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	items, total, err := ctrl.service.ListPosts(c.Context(), tenantID, PostListQuery{
		Category:  c.Query("category"),
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		Page:      page,
		Limit:     limit,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	})
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "posts loaded", items, fiber.Map{"page": page, "limit": limit, "total": total})
}

// CreatePost godoc
// @Summary Create post
// @Tags Posts
// @Security BearerAuth
// @Param payload body PostPayload true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/posts [post]
func (ctrl *controller) CreatePost(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	var req PostPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Title == "" || req.Category == "" || req.ContentMarkdown == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "title", Message: "title is required"},
			{Field: "category", Message: "category is required"},
			{Field: "content_markdown", Message: "content_markdown is required"},
		})
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	item, err := ctrl.service.CreatePost(c.Context(), tenantID, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "post created", item, nil)
}

// GetPost godoc
// @Summary Get post detail
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/posts/{id} [get]
func (ctrl *controller) GetPost(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	item, err := ctrl.service.GetPost(c.Context(), tenantID, id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "post loaded", item, nil)
}

// UpdatePost godoc
// @Summary Update post
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body PostPayload true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/posts/{id} [put]
func (ctrl *controller) UpdatePost(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req PostPayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Title == "" || req.Category == "" || req.ContentMarkdown == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "title", Message: "title is required"},
			{Field: "category", Message: "category is required"},
			{Field: "content_markdown", Message: "content_markdown is required"},
		})
	}
	if err := ctrl.service.UpdatePost(c.Context(), tenantID, id, req); err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "post updated", fiber.Map{"id": id}, nil)
}

// UpdatePostStatus godoc
// @Summary Update post status
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body UpdatePostStatusRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/posts/{id}/status [patch]
func (ctrl *controller) UpdatePostStatus(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	var req UpdatePostStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Status == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "status", Message: "status is required"}})
	}
	if err := ctrl.service.UpdatePostStatus(c.Context(), tenantID, id, req); err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "post status updated", fiber.Map{"id": id}, nil)
}

// DeletePost godoc
// @Summary Delete post
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "id"
// @Success 204
// @Router /tenant/posts/{id} [delete]
func (ctrl *controller) DeletePost(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "id", Message: "invalid id"}})
	}
	if err := ctrl.service.DeletePost(c.Context(), tenantID, id); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListStaticPages godoc
// @Summary List static pages
// @Tags StaticPages
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/static-pages [get]
func (ctrl *controller) ListStaticPages(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	items, err := ctrl.service.ListStaticPages(c.Context(), tenantID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "static pages loaded", items, nil)
}

// UpsertStaticPage godoc
// @Summary Upsert static page by slug
// @Tags StaticPages
// @Security BearerAuth
// @Param slug path string true "slug"
// @Param payload body StaticPagePayload true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/static-pages/{slug} [put]
func (ctrl *controller) UpsertStaticPage(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}
	slug := c.Params("slug")
	if slug == "" {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "slug", Message: "slug is required"}})
	}
	var req StaticPagePayload
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Title == "" || req.ContentMarkdown == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "title", Message: "title is required"},
			{Field: "content_markdown", Message: "content_markdown is required"},
		})
	}
	item, err := ctrl.service.UpsertStaticPageBySlug(c.Context(), tenantID, slug, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "static page upserted", item, nil)
}

func handleServiceError(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrNotFound) {
		return response.Error(c, fiber.StatusNotFound, "resource not found")
	}
	if errors.Is(err, ErrConflict) {
		return response.Error(c, fiber.StatusConflict, "resource conflict")
	}
	if errors.Is(err, ErrTagInUse) {
		return response.Error(c, fiber.StatusConflict, "tag is still used by post")
	}
	if errors.Is(err, ErrValidation) {
		return response.Error(c, fiber.StatusBadRequest, "validation failed")
	}
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}

// ==========================================
// UTILITY FUNCTIONS
// ==========================================

func tenantIDFromLocals(c *fiber.Ctx) (string, error) {
	tenantID := fmt.Sprint(c.Locals("tenant_id"))
	if tenantID == "" || tenantID == "<nil>" {
		return "", errors.New("tenant_id not found in context")
	}
	return tenantID, nil
}

func parseIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

// SetupTenant godoc
// @Summary Setup tenant name and subdomain (Onboarding)
// @Tags Tenant
// @Security BearerAuth
// @Param payload body SetupTenantRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/setup [patch]
func (ctrl *controller) SetupTenant(c *fiber.Ctx) error {
	tenantID, err := tenantIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req SetupTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}
	if req.Name == "" || req.Subdomain == "" {
		return response.Validation(c, "validation failed", []response.FieldError{
			{Field: "name", Message: "name is required"},
			{Field: "subdomain", Message: "subdomain is required"},
		})
	}

	if err := ctrl.service.SetupTenant(c.Context(), tenantID, req.Name, req.Subdomain); err != nil {
		return handleServiceError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "tenant setup successful", nil, nil)
}
