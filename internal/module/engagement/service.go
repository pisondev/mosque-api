package engagement

import (
	"context"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service interface {
	ListStaticPaymentMethods(ctx context.Context, tenantID string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error)
	CreateStaticPaymentMethod(ctx context.Context, tenantID string, req StaticPaymentMethodPayload) (*StaticPaymentMethodResponse, error)
	GetStaticPaymentMethod(ctx context.Context, tenantID string, id int64) (*StaticPaymentMethodResponse, error)
	UpdateStaticPaymentMethod(ctx context.Context, tenantID string, id int64, req StaticPaymentMethodPayload) error
	DeleteStaticPaymentMethod(ctx context.Context, tenantID string, id int64) error
	ListPublicStaticPaymentMethods(ctx context.Context, hostname string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error)
	ListSocialLinks(ctx context.Context, tenantID string, q ListQuery) ([]SocialLinkResponse, int64, error)
	CreateSocialLink(ctx context.Context, tenantID string, req SocialLinkPayload) (*SocialLinkResponse, error)
	GetSocialLink(ctx context.Context, tenantID string, id int64) (*SocialLinkResponse, error)
	UpdateSocialLink(ctx context.Context, tenantID string, id int64, req SocialLinkPayload) error
	DeleteSocialLink(ctx context.Context, tenantID string, id int64) error
	ListPublicSocialLinks(ctx context.Context, hostname string, q ListQuery) ([]SocialLinkResponse, int64, error)
	ListExternalLinks(ctx context.Context, tenantID string, q ListQuery) ([]ExternalLinkResponse, int64, error)
	CreateExternalLink(ctx context.Context, tenantID string, req ExternalLinkPayload) (*ExternalLinkResponse, error)
	GetExternalLink(ctx context.Context, tenantID string, id int64) (*ExternalLinkResponse, error)
	UpdateExternalLink(ctx context.Context, tenantID string, id int64, req ExternalLinkPayload) error
	DeleteExternalLink(ctx context.Context, tenantID string, id int64) error
	ListPublicExternalLinks(ctx context.Context, hostname string, q ListQuery) ([]ExternalLinkResponse, int64, error)
	ListFeatureCatalog(ctx context.Context) ([]FeatureCatalogResponse, error)
	ListWebsiteFeatures(ctx context.Context, tenantID string) ([]WebsiteFeatureResponse, error)
	UpsertWebsiteFeature(ctx context.Context, tenantID string, featureID int64, req WebsiteFeatureUpdateRequest) error
	BulkUpsertWebsiteFeatures(ctx context.Context, tenantID string, items []WebsiteFeatureBulkItem) error
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

var labelRegex = regexp.MustCompile(`^[A-Za-z -]+$`)
var digitsOnlyRegex = regexp.MustCompile(`^[0-9]+$`)

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) ListStaticPaymentMethods(ctx context.Context, tenantID string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error) {
	return s.repo.ListStaticPaymentMethods(ctx, tenantID, q)
}
func (s *service) CreateStaticPaymentMethod(ctx context.Context, tenantID string, req StaticPaymentMethodPayload) (*StaticPaymentMethodResponse, error) {
	if err := validateStaticPaymentMethod(&req); err != nil {
		return nil, err
	}
	return s.repo.CreateStaticPaymentMethod(ctx, tenantID, req)
}
func (s *service) GetStaticPaymentMethod(ctx context.Context, tenantID string, id int64) (*StaticPaymentMethodResponse, error) {
	return s.repo.GetStaticPaymentMethod(ctx, tenantID, id)
}
func (s *service) UpdateStaticPaymentMethod(ctx context.Context, tenantID string, id int64, req StaticPaymentMethodPayload) error {
	if err := validateStaticPaymentMethod(&req); err != nil {
		return err
	}
	return s.repo.UpdateStaticPaymentMethod(ctx, tenantID, id, req)
}
func (s *service) DeleteStaticPaymentMethod(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteStaticPaymentMethod(ctx, tenantID, id)
}
func (s *service) ListPublicStaticPaymentMethods(ctx context.Context, hostname string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error) {
	return s.repo.ListPublicStaticPaymentMethods(ctx, hostname, q)
}
func (s *service) ListSocialLinks(ctx context.Context, tenantID string, q ListQuery) ([]SocialLinkResponse, int64, error) {
	return s.repo.ListSocialLinks(ctx, tenantID, q)
}
func (s *service) CreateSocialLink(ctx context.Context, tenantID string, req SocialLinkPayload) (*SocialLinkResponse, error) {
	return s.repo.CreateSocialLink(ctx, tenantID, req)
}
func (s *service) GetSocialLink(ctx context.Context, tenantID string, id int64) (*SocialLinkResponse, error) {
	return s.repo.GetSocialLink(ctx, tenantID, id)
}
func (s *service) UpdateSocialLink(ctx context.Context, tenantID string, id int64, req SocialLinkPayload) error {
	return s.repo.UpdateSocialLink(ctx, tenantID, id, req)
}
func (s *service) DeleteSocialLink(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteSocialLink(ctx, tenantID, id)
}
func (s *service) ListPublicSocialLinks(ctx context.Context, hostname string, q ListQuery) ([]SocialLinkResponse, int64, error) {
	return s.repo.ListPublicSocialLinks(ctx, hostname, q)
}
func (s *service) ListExternalLinks(ctx context.Context, tenantID string, q ListQuery) ([]ExternalLinkResponse, int64, error) {
	return s.repo.ListExternalLinks(ctx, tenantID, q)
}
func (s *service) CreateExternalLink(ctx context.Context, tenantID string, req ExternalLinkPayload) (*ExternalLinkResponse, error) {
	return s.repo.CreateExternalLink(ctx, tenantID, req)
}
func (s *service) GetExternalLink(ctx context.Context, tenantID string, id int64) (*ExternalLinkResponse, error) {
	return s.repo.GetExternalLink(ctx, tenantID, id)
}
func (s *service) UpdateExternalLink(ctx context.Context, tenantID string, id int64, req ExternalLinkPayload) error {
	return s.repo.UpdateExternalLink(ctx, tenantID, id, req)
}
func (s *service) DeleteExternalLink(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteExternalLink(ctx, tenantID, id)
}
func (s *service) ListPublicExternalLinks(ctx context.Context, hostname string, q ListQuery) ([]ExternalLinkResponse, int64, error) {
	return s.repo.ListPublicExternalLinks(ctx, hostname, q)
}
func (s *service) ListFeatureCatalog(ctx context.Context) ([]FeatureCatalogResponse, error) {
	return s.repo.ListFeatureCatalog(ctx)
}
func (s *service) ListWebsiteFeatures(ctx context.Context, tenantID string) ([]WebsiteFeatureResponse, error) {
	return s.repo.ListWebsiteFeatures(ctx, tenantID)
}
func (s *service) UpsertWebsiteFeature(ctx context.Context, tenantID string, featureID int64, req WebsiteFeatureUpdateRequest) error {
	return s.repo.UpsertWebsiteFeature(ctx, tenantID, featureID, req)
}
func (s *service) BulkUpsertWebsiteFeatures(ctx context.Context, tenantID string, items []WebsiteFeatureBulkItem) error {
	return s.repo.BulkUpsertWebsiteFeatures(ctx, tenantID, items)
}

func validateStaticPaymentMethod(req *StaticPaymentMethodPayload) error {
	req.Label = strings.TrimSpace(req.Label)
	if req.Label == "" || len(req.Label) > 25 || !labelRegex.MatchString(req.Label) {
		return ErrValidation
	}
	if req.Description != nil && len(strings.TrimSpace(*req.Description)) > 250 {
		return ErrValidation
	}
	if req.ChannelType == "bank_account" {
		if req.BankName == nil || req.AccountNumber == nil || req.AccountHolderName == nil {
			return ErrValidation
		}
		bankName := strings.TrimSpace(*req.BankName)
		holder := strings.TrimSpace(*req.AccountHolderName)
		accountNo := strings.TrimSpace(*req.AccountNumber)
		if bankName == "" || len(bankName) > 25 || !labelRegex.MatchString(bankName) {
			return ErrValidation
		}
		if holder == "" || len(holder) > 25 || !labelRegex.MatchString(holder) {
			return ErrValidation
		}
		if accountNo == "" || !digitsOnlyRegex.MatchString(accountNo) {
			return ErrValidation
		}
		*req.BankName = bankName
		*req.AccountHolderName = holder
		*req.AccountNumber = accountNo
	}
	if req.ChannelType == "qris" {
		if req.QrisImageURL == nil || strings.TrimSpace(*req.QrisImageURL) == "" {
			return ErrValidation
		}
		qris := strings.TrimSpace(*req.QrisImageURL)
		if len(qris) > 1000000 {
			return ErrValidation
		}
		*req.QrisImageURL = qris
	}
	return nil
}
