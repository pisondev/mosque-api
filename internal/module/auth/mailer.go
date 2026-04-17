package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type EmailSender interface {
	SendPasswordReset(ctx context.Context, toEmail, resetURL string) error
}

type resendSender struct {
	apiKey     string
	fromEmail  string
	httpClient *http.Client
	log        *logrus.Logger
}

func NewResendSender(log *logrus.Logger) EmailSender {
	return &resendSender{
		apiKey:     strings.TrimSpace(os.Getenv("RESEND_API_KEY")),
		fromEmail:  defaultString(strings.TrimSpace(os.Getenv("RESEND_FROM_EMAIL")), "onboarding@resend.dev"),
		httpClient: &http.Client{Timeout: 15 * time.Second},
		log:        log,
	}
}

func (s *resendSender) SendPasswordReset(ctx context.Context, toEmail, resetURL string) error {
	if s.apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not configured")
	}

	payload := map[string]interface{}{
		"from":    s.fromEmail,
		"to":      []string{toEmail},
		"subject": "Reset password akun eTAKMIR",
		"html": fmt.Sprintf(`
			<div style="font-family:Arial,sans-serif;line-height:1.6;color:#111827">
			  <h2 style="margin-bottom:12px">Reset password eTAKMIR</h2>
			  <p>Kami menerima permintaan untuk mengganti password akun Anda.</p>
			  <p>Tautan ini berlaku selama 30 menit.</p>
			  <p><a href="%s" style="display:inline-block;padding:12px 18px;background:#059669;color:#ffffff;text-decoration:none;border-radius:10px;font-weight:700">Atur password baru</a></p>
			  <p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
			  <p style="font-size:12px;color:#6b7280">%s</p>
			</div>`, resetURL, resetURL),
		"text": fmt.Sprintf("Reset password akun eTAKMIR. Buka tautan berikut dalam 30 menit: %s", resetURL),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("resend returned status %d", resp.StatusCode)
	}

	return nil
}

func defaultString(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
