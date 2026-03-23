package auth

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/jackc/pgx/v5"
	"github.com/pisondev/mosque-api/internal/utils"
	"github.com/sirupsen/logrus"
)

type Service interface {
	HandleGoogleLogin(ctx context.Context, req LoginGoogleRequest) (*AuthResponse, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) HandleGoogleLogin(ctx context.Context, req LoginGoogleRequest) (*AuthResponse, error) {
	s.log.Info("processing google login request")

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		s.log.Error("GOOGLE_CLIENT_ID is not set in environment")
		return nil, errors.New("internal server error")
	}

	// 1. Verifikasi token ke Google
	payload, err := idtoken.Validate(ctx, req.Token, clientID)
	if err != nil {
		s.log.Errorf("failed to validate google token: %v", err)
		return nil, errors.New("invalid google token")
	}

	email := payload.Claims["email"].(string)
	googleID := payload.Subject
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
			user, err = s.repo.CreateTenantAndUser(ctx, email, googleID)
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
		s.log.Infof("existing user %s logged in", email)
	}

	// ==========================================
	// GENERATE JWT DENGAN DATA ASLI DARI DATABASE
	// ==========================================

	tokenString, err := utils.GenerateToken(user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		s.log.Errorf("failed to generate jwt: %v", err)
		return nil, errors.New("failed to generate authentication token")
	}

	return &AuthResponse{
		AccessToken: tokenString,
		Email:       user.Email,
		Role:        user.Role,
	}, nil
}
