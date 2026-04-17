package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pisondev/mosque-api/internal/testutil"
	"github.com/sirupsen/logrus"
)

type mockEmailSender struct {
	lastEmail string
	lastURL   string
}

func (m *mockEmailSender) SendPasswordReset(_ context.Context, toEmail, resetURL string) error {
	m.lastEmail = toEmail
	m.lastURL = resetURL
	return nil
}

func TestAuthRegisterLoginAndResetFlow(t *testing.T) {
	db := testutil.OpenTestDB(t)
	defer db.Close()

	t.Setenv("JWT_SECRET", "auth-secret")
	t.Setenv("APP_BASE_URL", "http://localhost:3000")

	repo := NewRepository(db)
	sender := &mockEmailSender{}
	service := NewService(repo, logrus.New(), sender)
	controller := NewController(service, logrus.New())

	app := fiber.New()
	group := app.Group("/api/v1/auth")
	group.Post("/register", controller.Register)
	group.Post("/login", controller.Login)
	group.Post("/forgot-password", controller.ForgotPassword)
	group.Post("/reset-password", controller.ResetPassword)

	email := "auth.integration@test.local"

	status, body := doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/register", `{"email":"`+email+`","password":"Password123"}`)
	if status != http.StatusCreated {
		t.Fatalf("register expected 201 got %d body=%v", status, body)
	}
	data := body["data"].(map[string]interface{})
	if data["access_token"] == "" {
		t.Fatalf("register should return access token body=%v", body)
	}

	user, err := repo.FindByEmail(context.Background(), email)
	if err != nil {
		t.Fatalf("find user after register failed: %v", err)
	}
	defer cleanupAuthTenant(t, db, user.TenantID)
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		t.Fatalf("password hash should be stored after register")
	}

	status, _ = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{"email":"`+email+`","password":"Password123"}`)
	if status != http.StatusOK {
		t.Fatalf("login expected 200 got %d", status)
	}

	status, _ = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{"email":"`+email+`","password":"Wrong123"}`)
	if status != http.StatusUnauthorized {
		t.Fatalf("invalid login expected 401 got %d", status)
	}

	status, body = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/forgot-password", `{"email":"`+email+`"}`)
	if status != http.StatusOK {
		t.Fatalf("forgot password expected 200 got %d body=%v", status, body)
	}
	if sender.lastEmail != email || sender.lastURL == "" {
		t.Fatalf("forgot password should send email, got email=%s url=%s", sender.lastEmail, sender.lastURL)
	}
	token := sender.lastURL[strings.LastIndex(sender.lastURL, "token=")+6:]

	status, body = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/reset-password", `{"token":"`+token+`","password":"NewPassword123"}`)
	if status != http.StatusOK {
		t.Fatalf("reset password expected 200 got %d body=%v", status, body)
	}

	status, _ = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{"email":"`+email+`","password":"Password123"}`)
	if status != http.StatusUnauthorized {
		t.Fatalf("old password should fail after reset, got %d", status)
	}

	status, _ = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{"email":"`+email+`","password":"NewPassword123"}`)
	if status != http.StatusOK {
		t.Fatalf("new password should work after reset, got %d", status)
	}

	status, _ = doAuthRequest(t, app, http.MethodPost, "/api/v1/auth/reset-password", `{"token":"`+token+`","password":"AnotherPassword123"}`)
	if status != http.StatusBadRequest {
		t.Fatalf("used token should be rejected, got %d", status)
	}
}

func cleanupAuthTenant(t *testing.T, db *pgxpool.Pool, tenantID string) {
	t.Helper()
	_, _ = db.Exec(context.Background(), `DELETE FROM tenants WHERE id = $1::uuid`, tenantID)
}

func doAuthRequest(t *testing.T, app *fiber.App, method, path, rawBody string) (int, map[string]interface{}) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(rawBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body failed: %v", err)
	}
	return resp.StatusCode, body
}
