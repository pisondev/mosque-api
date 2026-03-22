package auth

import (
	"context"

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
	// TODO: Implementasi logika verifikasi token Google dan pembuatan JWT di sini
	return &AuthResponse{}, nil
}
