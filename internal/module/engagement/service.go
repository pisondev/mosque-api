package engagement

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service interface {
	ListDonationChannels(ctx context.Context, tenantID string, q ListQuery) ([]DonationChannelResponse, int64, error)
	CreateDonationChannel(ctx context.Context, tenantID string, req DonationChannelPayload) (*DonationChannelResponse, error)
	GetDonationChannel(ctx context.Context, tenantID string, id int64) (*DonationChannelResponse, error)
	UpdateDonationChannel(ctx context.Context, tenantID string, id int64, req DonationChannelPayload) error
	DeleteDonationChannel(ctx context.Context, tenantID string, id int64) error
	ListPublicDonationChannels(ctx context.Context, hostname string, q ListQuery) ([]DonationChannelResponse, int64, error)
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

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) ListDonationChannels(ctx context.Context, tenantID string, q ListQuery) ([]DonationChannelResponse, int64, error) {
	return s.repo.ListDonationChannels(ctx, tenantID, q)
}
func (s *service) CreateDonationChannel(ctx context.Context, tenantID string, req DonationChannelPayload) (*DonationChannelResponse, error) {
	return s.repo.CreateDonationChannel(ctx, tenantID, req)
}
func (s *service) GetDonationChannel(ctx context.Context, tenantID string, id int64) (*DonationChannelResponse, error) {
	return s.repo.GetDonationChannel(ctx, tenantID, id)
}
func (s *service) UpdateDonationChannel(ctx context.Context, tenantID string, id int64, req DonationChannelPayload) error {
	return s.repo.UpdateDonationChannel(ctx, tenantID, id, req)
}
func (s *service) DeleteDonationChannel(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteDonationChannel(ctx, tenantID, id)
}
func (s *service) ListPublicDonationChannels(ctx context.Context, hostname string, q ListQuery) ([]DonationChannelResponse, int64, error) {
	return s.repo.ListPublicDonationChannels(ctx, hostname, q)
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
