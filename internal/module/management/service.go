package management

import (
	"context"
	"regexp"
	"strings"

	"github.com/pisondev/mosque-api/internal/constant"
	"github.com/sirupsen/logrus"
)

type Service interface {
	GetTenantMe(ctx context.Context, tenantID, email, role string) (map[string]interface{}, error)
	ListDomains(ctx context.Context, tenantID string, q DomainListQuery) ([]DomainResponse, int64, error)
	CreateDomain(ctx context.Context, tenantID string, req CreateDomainRequest) (*DomainResponse, error)
	UpdateDomain(ctx context.Context, tenantID string, id int64, req UpdateDomainRequest) error
	DeleteDomain(ctx context.Context, tenantID string, id int64) error
	GetProfile(ctx context.Context, tenantID string) (*ProfileResponse, error)
	UpsertProfile(ctx context.Context, tenantID string, req ProfileRequest) (*ProfileResponse, error)
	ListTags(ctx context.Context, tenantID, scope, search string, page, limit int) ([]TagResponse, int64, error)
	CreateTag(ctx context.Context, tenantID string, req CreateTagRequest) (*TagResponse, error)
	UpdateTag(ctx context.Context, tenantID string, id int64, req UpdateTagRequest) error
	DeleteTag(ctx context.Context, tenantID string, id int64) error
	ListPosts(ctx context.Context, tenantID string, q PostListQuery) ([]PostResponse, int64, error)
	CreatePost(ctx context.Context, tenantID string, req PostPayload) (*PostResponse, error)
	GetPost(ctx context.Context, tenantID string, id int64) (*PostResponse, error)
	UpdatePost(ctx context.Context, tenantID string, id int64, req PostPayload) error
	UpdatePostStatus(ctx context.Context, tenantID string, id int64, req UpdatePostStatusRequest) error
	DeletePost(ctx context.Context, tenantID string, id int64) error
	ListStaticPages(ctx context.Context, tenantID string) ([]PostResponse, error)
	UpsertStaticPageBySlug(ctx context.Context, tenantID, slug string, req StaticPagePayload) (*PostResponse, error)
	SetupTenant(ctx context.Context, tenantID, name, subdomain string) error

	GetBillingStatus(ctx context.Context, tenantID string) (*BillingStatusResponse, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

var setupSubdomainRegex = regexp.MustCompile(`^[a-z-]{1,10}$`)
var alphaSpaceHyphenRegex = regexp.MustCompile(`^[A-Za-z -]+$`)
var digitsOnlyRegex = regexp.MustCompile(`^[0-9]+$`)

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) GetTenantMe(ctx context.Context, tenantID, email, role string) (map[string]interface{}, error) {
	data, err := s.repo.GetTenantContext(ctx, tenantID, email)
	if err != nil {
		return nil, err
	}
	data["role"] = role
	return data, nil
}

func (s *service) ListDomains(ctx context.Context, tenantID string, q DomainListQuery) ([]DomainResponse, int64, error) {
	return s.repo.ListDomains(ctx, tenantID, q)
}

func (s *service) CreateDomain(ctx context.Context, tenantID string, req CreateDomainRequest) (*DomainResponse, error) {
	req.Hostname = strings.ToLower(strings.TrimSpace(req.Hostname))
	return s.repo.CreateDomain(ctx, tenantID, req)
}

func (s *service) UpdateDomain(ctx context.Context, tenantID string, id int64, req UpdateDomainRequest) error {
	return s.repo.UpdateDomain(ctx, tenantID, id, req)
}

func (s *service) DeleteDomain(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteDomain(ctx, tenantID, id)
}

func (s *service) GetProfile(ctx context.Context, tenantID string) (*ProfileResponse, error) {
	return s.repo.GetProfile(ctx, tenantID)
}

func (s *service) UpsertProfile(ctx context.Context, tenantID string, req ProfileRequest) (*ProfileResponse, error) {
	req.OfficialName = strings.TrimSpace(req.OfficialName)
	req.ShortName = strings.TrimSpace(req.ShortName)
	req.ShortDesc = strings.TrimSpace(req.ShortDesc)
	req.Province = strings.TrimSpace(req.Province)
	req.City = strings.TrimSpace(req.City)
	req.AddressFull = strings.TrimSpace(req.AddressFull)
	req.PhoneWA = strings.TrimSpace(req.PhoneWA)
	req.HeaderImageURL = strings.TrimSpace(req.HeaderImageURL)

	if req.OfficialName == "" || len(req.OfficialName) > 25 || !alphaSpaceHyphenRegex.MatchString(req.OfficialName) {
		return nil, ErrValidation
	}
	if req.ShortName != "" {
		if len(req.ShortName) > 25 || !alphaSpaceHyphenRegex.MatchString(req.ShortName) {
			return nil, ErrValidation
		}
	}
	if req.PhoneWA != "" && !digitsOnlyRegex.MatchString(req.PhoneWA) {
		return nil, ErrValidation
	}
	if len(req.AddressFull) > 250 || len(req.ShortDesc) > 250 {
		return nil, ErrValidation
	}

	return s.repo.UpsertProfile(ctx, tenantID, req)
}

func (s *service) ListTags(ctx context.Context, tenantID, scope, search string, page, limit int) ([]TagResponse, int64, error) {
	return s.repo.ListTags(ctx, tenantID, scope, search, page, limit)
}

func (s *service) CreateTag(ctx context.Context, tenantID string, req CreateTagRequest) (*TagResponse, error) {
	return s.repo.CreateTag(ctx, tenantID, req)
}

func (s *service) UpdateTag(ctx context.Context, tenantID string, id int64, req UpdateTagRequest) error {
	return s.repo.UpdateTag(ctx, tenantID, id, req)
}

func (s *service) DeleteTag(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteTag(ctx, tenantID, id)
}

func (s *service) ListPosts(ctx context.Context, tenantID string, q PostListQuery) ([]PostResponse, int64, error) {
	return s.repo.ListPosts(ctx, tenantID, q)
}

func (s *service) CreatePost(ctx context.Context, tenantID string, req PostPayload) (*PostResponse, error) {
	return s.repo.CreatePost(ctx, tenantID, req)
}

func (s *service) GetPost(ctx context.Context, tenantID string, id int64) (*PostResponse, error) {
	return s.repo.GetPost(ctx, tenantID, id)
}

func (s *service) UpdatePost(ctx context.Context, tenantID string, id int64, req PostPayload) error {
	return s.repo.UpdatePost(ctx, tenantID, id, req)
}

func (s *service) UpdatePostStatus(ctx context.Context, tenantID string, id int64, req UpdatePostStatusRequest) error {
	return s.repo.UpdatePostStatus(ctx, tenantID, id, req)
}

func (s *service) DeletePost(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeletePost(ctx, tenantID, id)
}

func (s *service) ListStaticPages(ctx context.Context, tenantID string) ([]PostResponse, error) {
	return s.repo.ListStaticPages(ctx, tenantID)
}

func (s *service) UpsertStaticPageBySlug(ctx context.Context, tenantID, slug string, req StaticPagePayload) (*PostResponse, error) {
	return s.repo.UpsertStaticPageBySlug(ctx, tenantID, slug, req)
}

func (s *service) SetupTenant(ctx context.Context, tenantID, name, subdomain string) error {
	name = strings.TrimSpace(name)
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if name == "" || subdomain == "" {
		return ErrValidation
	}
	if !setupSubdomainRegex.MatchString(subdomain) {
		return ErrValidation
	}

	return s.repo.UpdateTenantSetup(ctx, tenantID, name, subdomain)
}

func (s *service) GetBillingStatus(ctx context.Context, tenantID string) (*BillingStatusResponse, error) {
	// 1. Ambil data mentah dari DB
	raw, err := s.repo.GetRawBillingData(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. Cocokkan dengan Kamus Plan (Fallback ke Free jika plan tidak valid)
	planDetail, exists := constant.SubscriptionPlans[raw.SubscriptionPlan]
	if !exists {
		planDetail = constant.SubscriptionPlans[constant.PlanFree]
		raw.SubscriptionPlan = constant.PlanFree
	}

	// 3. Rakit Response sesuai Kontrak Frontend
	res := &BillingStatusResponse{
		SubscriptionPlan: raw.SubscriptionPlan,
		ActiveTemplate:   raw.ActiveTemplate,
		Storage: StorageInfo{
			LimitMB: planDetail.StorageLimitMB,
			UsedMB:  raw.StorageUsedMB,
		},
		FeaturesUnlocked:      planDetail.FeaturesUnlocked,
		AttributionEnabled:    planDetail.AttributionEnabled,
		PlatformFeePercentage: planDetail.PlatformFeePercentage,
	}

	return res, nil
}
