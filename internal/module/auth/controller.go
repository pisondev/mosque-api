package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	GoogleLogin(c *fiber.Ctx) error
	GetAccountProfile(c *fiber.Ctx) error
	UpdateAccountProfile(c *fiber.Ctx) error
}

type controller struct {
	service Service
	log     *logrus.Logger
}

func NewController(service Service, log *logrus.Logger) Controller {
	return &controller{service: service, log: log}
}

func (ctrl *controller) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}

	resp, err := ctrl.service.Register(c.Context(), req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "registration successful", resp, nil)
}

func (ctrl *controller) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}

	resp, err := ctrl.service.Login(c.Context(), req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "login successful", resp, nil)
}

func (ctrl *controller) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}

	if err := ctrl.service.ForgotPassword(c.Context(), req); err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "password reset email queued", fiber.Map{"accepted": true}, nil)
}

func (ctrl *controller) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}

	resp, err := ctrl.service.ResetPassword(c.Context(), req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "password reset successful", resp, nil)
}

func (ctrl *controller) GoogleLogin(c *fiber.Ctx) error {
	ctrl.log.Info("received google login request")

	var req LoginGoogleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}

	resp, err := ctrl.service.HandleGoogleLogin(c.Context(), req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "google login successful", resp, nil)
}

func (ctrl *controller) GetAccountProfile(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	resp, err := ctrl.service.GetAccountProfile(c.Context(), userID)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "account profile loaded", resp, nil)
}

func (ctrl *controller) UpdateAccountProfile(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req UpdateAccountProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Validation(c, "validation failed", []response.FieldError{{Field: "body", Message: "invalid request format"}})
	}

	resp, err := ctrl.service.UpdateAccountProfile(c.Context(), userID, req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "account profile updated", resp, nil)
}

func getUserID(c *fiber.Ctx) string {
	if value := c.Locals("user_id"); value != nil {
		return strings.TrimSpace(fmt.Sprintf("%v", value))
	}
	return ""
}

func handleAuthError(c *fiber.Ctx, err error) error {
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		return response.Validation(c, validationErr.Error(), validationErr.Fields)
	}

	switch {
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInvalidGoogleToken):
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrEmailAlreadyExists):
		return response.Error(c, fiber.StatusConflict, err.Error())
	case errors.Is(err, ErrInvalidResetToken):
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, ErrEmailDeliveryFailure):
		return response.Error(c, fiber.StatusBadGateway, err.Error())
	default:
		return response.Error(c, fiber.StatusInternalServerError, "internal server error")
	}
}
