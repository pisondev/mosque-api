package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/jackc/pgx/v5"
	"github.com/pisondev/mosque-api/internal/response"
	"github.com/pisondev/mosque-api/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) (*AuthResponse, error)
	HandleGoogleLogin(ctx context.Context, req LoginGoogleRequest) (*AuthResponse, error)
	GetAccountProfile(ctx context.Context, userID string) (*AccountProfileResponse, error)
	UpdateAccountProfile(ctx context.Context, userID string, req UpdateAccountProfileRequest) (*AccountProfileResponse, error)
}

type service struct {
	repo              Repository
	log               *logrus.Logger
	emailSender       EmailSender
	verifyGoogleToken func(context.Context, string, string) (*idtoken.Payload, error)
}

type ValidationError struct {
	Fields []response.FieldError
}

func (e ValidationError) Error() string {
	return "validation failed"
}

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailAlreadyExists   = errors.New("email already registered")
	ErrInvalidResetToken    = errors.New("invalid or expired reset token")
	ErrInvalidGoogleToken   = errors.New("invalid google token")
	ErrEmailDeliveryFailure = errors.New("failed to send password reset email")
)

func NewService(repo Repository, log *logrus.Logger, emailSender EmailSender) Service {
	return &service{
		repo:              repo,
		log:               log,
		emailSender:       emailSender,
		verifyGoogleToken: idtoken.Validate,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	email := normalizeEmail(req.Email)
	password := strings.TrimSpace(req.Password)
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		s.log.Errorf("database error during register lookup: %v", err)
		return nil, errors.New("internal server error")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Errorf("failed to hash password during register: %v", err)
		return nil, errors.New("internal server error")
	}

	user, err := s.repo.CreateTenantAndUser(ctx, CreateUserParams{
		Email:        email,
		DisplayName:  stringPtr(deriveDisplayName(email)),
		PasswordHash: stringPtr(string(hash)),
		TenantName:   deriveTenantName(email),
	})
	if err != nil {
		s.log.Errorf("failed to create tenant and user during register: %v", err)
		return nil, errors.New("internal server error")
	}

	return buildAuthResponse(user)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	email := normalizeEmail(req.Email)
	password := strings.TrimSpace(req.Password)
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if password == "" {
		return nil, ValidationError{Fields: []response.FieldError{{Field: "password", Message: "password is required"}}}
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		s.log.Errorf("database error during login lookup: %v", err)
		return nil, errors.New("internal server error")
	}
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return buildAuthResponse(user)
}

func (s *service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	email := normalizeEmail(req.Email)
	if err := validateEmail(email); err != nil {
		return err
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		s.log.Errorf("database error during forgot password lookup: %v", err)
		return errors.New("internal server error")
	}

	token, tokenHash, err := generateResetToken()
	if err != nil {
		s.log.Errorf("failed to generate password reset token: %v", err)
		return errors.New("internal server error")
	}

	expiresAt := time.Now().UTC().Add(30 * time.Minute)
	if err := s.repo.StorePasswordResetToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		s.log.Errorf("failed to store password reset token: %v", err)
		return errors.New("internal server error")
	}

	resetURL := fmt.Sprintf("%s/auth?mode=reset&token=%s", strings.TrimRight(webBaseURL(), "/"), token)
	if s.emailSender == nil {
		return ErrEmailDeliveryFailure
	}
	if err := s.emailSender.SendPasswordReset(ctx, user.Email, resetURL); err != nil {
		s.log.Errorf("failed to send password reset email: %v", err)
		return ErrEmailDeliveryFailure
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetPasswordRequest) (*AuthResponse, error) {
	if strings.TrimSpace(req.Token) == "" {
		return nil, ValidationError{Fields: []response.FieldError{{Field: "token", Message: "token is required"}}}
	}
	if err := validatePassword(strings.TrimSpace(req.Password)); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(req.Password)), bcrypt.DefaultCost)
	if err != nil {
		s.log.Errorf("failed to hash password during reset: %v", err)
		return nil, errors.New("internal server error")
	}

	user, err := s.repo.ConsumePasswordResetToken(ctx, hashToken(req.Token), string(hash))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidResetToken
		}
		s.log.Errorf("failed to consume password reset token: %v", err)
		return nil, errors.New("internal server error")
	}

	return buildAuthResponse(user)
}

func (s *service) HandleGoogleLogin(ctx context.Context, req LoginGoogleRequest) (*AuthResponse, error) {
	s.log.Info("processing google login request")

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		s.log.Error("GOOGLE_CLIENT_ID is not set in environment")
		return nil, errors.New("internal server error")
	}

	// 1. Verifikasi token ke Google
	payload, err := s.verifyGoogleToken(ctx, req.Token, clientID)
	if err != nil {
		s.log.Errorf("failed to validate google token: %v", err)
		return nil, ErrInvalidGoogleToken
	}

	emailClaim, ok := payload.Claims["email"].(string)
	if !ok || strings.TrimSpace(emailClaim) == "" {
		return nil, ErrInvalidGoogleToken
	}
	email := normalizeEmail(emailClaim)
	googleID := payload.Subject
	displayName := strings.TrimSpace(claimString(payload.Claims, "name"))
	if displayName == "" {
		displayName = deriveDisplayName(email)
	}
	avatarURL := strings.TrimSpace(claimString(payload.Claims, "picture"))
	if strings.TrimSpace(googleID) == "" {
		return nil, ErrInvalidGoogleToken
	}
	s.log.Infof("successfully verified user from google: %s", email)

	// ==========================================
	// LOGIKA DATABASE: CARI ATAU BUAT PENGGUNA BARU
	// ==========================================

	// Coba cari pengguna berdasarkan email
	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		// Jika error-nya karena data tidak ditemukan (user baru)
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Infof("user %s not found, initiating auto-registration", email)

			// Eksekusi Transaction untuk buat Tenant dan User baru
			user, err = s.repo.CreateTenantAndUser(ctx, CreateUserParams{
				Email:       email,
				DisplayName: stringPtr(displayName),
				AvatarURL:   optionalStringPtr(avatarURL),
				GoogleID:    stringPtr(googleID),
				TenantName:  deriveTenantName(email),
			})
			if err != nil {
				s.log.Errorf("failed to auto-register user: %v", err)
				return nil, errors.New("failed to register new user")
			}
			s.log.Infof("successfully registered new tenant and user for %s", email)

		} else {
			// Jika error-nya karena hal lain (koneksi putus, dll)
			s.log.Errorf("database error during user lookup: %v", err)
			return nil, errors.New("internal server error")
		}
	} else {
		if err := s.repo.SyncGoogleAccount(ctx, user.ID, googleID, stringPtr(displayName), optionalStringPtr(avatarURL)); err != nil {
			s.log.Errorf("failed to sync google account for user %s: %v", email, err)
			return nil, errors.New("internal server error")
		}
		user.GoogleID = stringPtr(googleID)
		if user.DisplayName == nil || strings.TrimSpace(*user.DisplayName) == "" {
			user.DisplayName = stringPtr(displayName)
		}
		if (user.AvatarURL == nil || strings.TrimSpace(*user.AvatarURL) == "") && avatarURL != "" {
			user.AvatarURL = stringPtr(avatarURL)
		}
		s.log.Infof("existing user %s logged in", email)
	}

	resp, err := buildAuthResponse(user)
	if err != nil {
		s.log.Errorf("failed to generate jwt: %v", err)
		return nil, errors.New("failed to generate authentication token")
	}

	return resp, nil
}

func (s *service) GetAccountProfile(ctx context.Context, userID string) (*AccountProfileResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		s.log.Errorf("database error during account profile lookup: %v", err)
		return nil, errors.New("internal server error")
	}

	return buildAccountProfileResponse(user), nil
}

func (s *service) UpdateAccountProfile(ctx context.Context, userID string, req UpdateAccountProfileRequest) (*AccountProfileResponse, error) {
	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		return nil, ValidationError{Fields: []response.FieldError{{Field: "display_name", Message: "display name is required"}}}
	}
	if len(displayName) > 80 {
		return nil, ValidationError{Fields: []response.FieldError{{Field: "display_name", Message: "display name must be at most 80 characters"}}}
	}

	var avatarURL *string
	if req.AvatarURL != nil {
		trimmed := strings.TrimSpace(*req.AvatarURL)
		avatarURL = &trimmed
		if trimmed == "" {
			avatarURL = nil
		}
	}

	user, err := s.repo.UpdateAccountProfile(ctx, userID, UpdateAccountProfileParams{
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	})
	if err != nil {
		s.log.Errorf("failed to update account profile: %v", err)
		return nil, errors.New("internal server error")
	}

	return buildAccountProfileResponse(user), nil
}

func buildAuthResponse(user *User) (*AuthResponse, error) {
	tokenString, err := utils.GenerateToken(user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: tokenString,
		Email:       user.Email,
		DisplayName: stringValue(user.DisplayName),
		AvatarURL:   stringValue(user.AvatarURL),
		Role:        user.Role,
	}, nil
}

func buildAccountProfileResponse(user *User) *AccountProfileResponse {
	return &AccountProfileResponse{
		Email:           user.Email,
		DisplayName:     stringValueOrDefault(user.DisplayName, deriveDisplayName(user.Email)),
		AvatarURL:       stringValue(user.AvatarURL),
		Role:            user.Role,
		GoogleConnected: user.GoogleID != nil && strings.TrimSpace(*user.GoogleID) != "",
	}
}

func validateEmail(email string) error {
	if email == "" {
		return ValidationError{Fields: []response.FieldError{{Field: "email", Message: "email is required"}}}
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ValidationError{Fields: []response.FieldError{{Field: "email", Message: "email is invalid"}}}
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return ValidationError{Fields: []response.FieldError{{Field: "password", Message: "password is required"}}}
	}
	if len(password) < 8 {
		return ValidationError{Fields: []response.FieldError{{Field: "password", Message: "password must be at least 8 characters"}}}
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") || !strings.ContainsAny(password, "0123456789") {
		return ValidationError{Fields: []response.FieldError{{Field: "password", Message: "password must contain letters and numbers"}}}
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func deriveTenantName(email string) string {
	localPart := strings.Split(email, "@")[0]
	localPart = strings.TrimSpace(strings.ReplaceAll(localPart, ".", " "))
	if localPart == "" {
		return "Masjid Baru"
	}
	return fmt.Sprintf("%s Mosque", strings.Title(localPart))
}

func deriveDisplayName(email string) string {
	localPart := strings.Split(email, "@")[0]
	replacer := strings.NewReplacer(".", " ", "-", " ", "_", " ")
	cleaned := strings.TrimSpace(replacer.Replace(localPart))
	if cleaned == "" {
		return "Admin"
	}
	words := strings.Fields(cleaned)
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}
	return strings.Join(words, " ")
}

func generateResetToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	raw := hex.EncodeToString(buf)
	return raw, hashToken(raw), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(raw)))
	return hex.EncodeToString(sum[:])
}

func webBaseURL() string {
	if value := strings.TrimSpace(os.Getenv("APP_BASE_URL")); value != "" {
		return value
	}
	return "http://localhost:3000"
}

func stringPtr(value string) *string {
	return &value
}

func optionalStringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func stringValueOrDefault(value *string, fallback string) string {
	trimmed := stringValue(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func claimString(claims map[string]interface{}, key string) string {
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}
