package auth

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/auth/credentials/idtoken"
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

	// 1. Verifikasi token ke server Google (secara kriptografis)
	payload, err := idtoken.Validate(ctx, req.Token, clientID)
	if err != nil {
		s.log.Errorf("failed to validate google token: %v", err)
		return nil, errors.New("invalid google token")
	}

	// 2. Ekstrak data dari token Google yang sudah valid
	email := payload.Claims["email"].(string)
	googleID := payload.Subject

	s.log.Infof("successfully verified user: %s (Google ID: %s)", email, googleID)

	// TODO Selanjutnya:
	// - Cek apakah email ini ada di database via s.repo.FindByEmail
	// - Jika tidak ada, buat User (dan Tenant) baru
	// - Generate JWT buatan kita sendiri dan kembalikan ke Next.js

	// Untuk sementara, kita kembalikan respons sukses dulu
	return &AuthResponse{
		AccessToken: "jwt_token_sementara",
		Email:       email,
		Role:        "admin",
	}, nil
}
