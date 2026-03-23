package community

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service interface {
	ListEvents(ctx context.Context, tenantID string, q EventListQuery) ([]EventResponse, int64, error)
	CreateEvent(ctx context.Context, tenantID string, req EventPayload) (*EventResponse, error)
	GetEvent(ctx context.Context, tenantID string, id int64) (*EventResponse, error)
	UpdateEvent(ctx context.Context, tenantID string, id int64, req EventPayload) error
	UpdateEventStatus(ctx context.Context, tenantID string, id int64, status string) error
	DeleteEvent(ctx context.Context, tenantID string, id int64) error
	ListPublicEvents(ctx context.Context, hostname string, page, limit int) ([]EventResponse, int64, error)
	ListGalleryAlbums(ctx context.Context, tenantID string, q BaseListQuery) ([]GalleryAlbumResponse, int64, error)
	CreateGalleryAlbum(ctx context.Context, tenantID string, req GalleryAlbumPayload) (*GalleryAlbumResponse, error)
	GetGalleryAlbum(ctx context.Context, tenantID string, id int64) (*GalleryAlbumResponse, error)
	UpdateGalleryAlbum(ctx context.Context, tenantID string, id int64, req GalleryAlbumPayload) error
	DeleteGalleryAlbum(ctx context.Context, tenantID string, id int64) error
	ListPublicGalleryAlbums(ctx context.Context, hostname string, page, limit int) ([]GalleryAlbumResponse, int64, error)
	ListGalleryItems(ctx context.Context, tenantID string, q BaseListQuery, albumID string) ([]GalleryItemResponse, int64, error)
	CreateGalleryItem(ctx context.Context, tenantID string, req GalleryItemPayload) (*GalleryItemResponse, error)
	GetGalleryItem(ctx context.Context, tenantID string, id int64) (*GalleryItemResponse, error)
	UpdateGalleryItem(ctx context.Context, tenantID string, id int64, req GalleryItemPayload) error
	DeleteGalleryItem(ctx context.Context, tenantID string, id int64) error
	ListPublicGalleryItems(ctx context.Context, hostname string, page, limit int) ([]GalleryItemResponse, int64, error)
	ListManagementMembers(ctx context.Context, tenantID string, q BaseListQuery) ([]ManagementMemberResponse, int64, error)
	CreateManagementMember(ctx context.Context, tenantID string, req ManagementMemberPayload) (*ManagementMemberResponse, error)
	GetManagementMember(ctx context.Context, tenantID string, id int64) (*ManagementMemberResponse, error)
	UpdateManagementMember(ctx context.Context, tenantID string, id int64, req ManagementMemberPayload) error
	DeleteManagementMember(ctx context.Context, tenantID string, id int64) error
	ListPublicManagementMembers(ctx context.Context, hostname string, page, limit int) ([]ManagementMemberResponse, int64, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) ListEvents(ctx context.Context, tenantID string, q EventListQuery) ([]EventResponse, int64, error) {
	return s.repo.ListEvents(ctx, tenantID, q)
}

func (s *service) CreateEvent(ctx context.Context, tenantID string, req EventPayload) (*EventResponse, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = normalizeEventCategory(req.Category)
	req.Status = normalizeEventStatus(req.Status)
	req.TimeMode = strings.ToLower(strings.TrimSpace(req.TimeMode))
	if req.Title == "" || req.Category == "" || req.StartDate == "" || req.TimeMode == "" || req.Status == "" {
		return nil, ErrValidation
	}
	if req.TimeMode == "exact_time" && (req.StartTime == nil || strings.TrimSpace(*req.StartTime) == "") {
		return nil, ErrValidation
	}
	if req.TimeMode == "after_prayer" && (req.AfterPrayer == nil || strings.TrimSpace(*req.AfterPrayer) == "") {
		return nil, ErrValidation
	}
	if req.RepeatWeekdays == nil {
		req.RepeatWeekdays = []int16{}
	}
	return s.repo.CreateEvent(ctx, tenantID, req)
}

func (s *service) GetEvent(ctx context.Context, tenantID string, id int64) (*EventResponse, error) {
	return s.repo.GetEvent(ctx, tenantID, id)
}

func (s *service) UpdateEvent(ctx context.Context, tenantID string, id int64, req EventPayload) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = normalizeEventCategory(req.Category)
	req.Status = normalizeEventStatus(req.Status)
	req.TimeMode = strings.ToLower(strings.TrimSpace(req.TimeMode))
	if req.Title == "" || req.Category == "" || req.StartDate == "" || req.TimeMode == "" || req.Status == "" {
		return ErrValidation
	}
	if req.TimeMode == "exact_time" && (req.StartTime == nil || strings.TrimSpace(*req.StartTime) == "") {
		return ErrValidation
	}
	if req.TimeMode == "after_prayer" && (req.AfterPrayer == nil || strings.TrimSpace(*req.AfterPrayer) == "") {
		return ErrValidation
	}
	if req.RepeatWeekdays == nil {
		req.RepeatWeekdays = []int16{}
	}
	return s.repo.UpdateEvent(ctx, tenantID, id, req)
}

func (s *service) UpdateEventStatus(ctx context.Context, tenantID string, id int64, status string) error {
	status = normalizeEventStatus(status)
	if status == "" {
		return ErrValidation
	}
	return s.repo.UpdateEventStatus(ctx, tenantID, id, status)
}

func (s *service) DeleteEvent(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteEvent(ctx, tenantID, id)
}

func (s *service) ListPublicEvents(ctx context.Context, hostname string, page, limit int) ([]EventResponse, int64, error) {
	return s.repo.ListPublicEvents(ctx, hostname, page, limit)
}

func (s *service) ListGalleryAlbums(ctx context.Context, tenantID string, q BaseListQuery) ([]GalleryAlbumResponse, int64, error) {
	return s.repo.ListGalleryAlbums(ctx, tenantID, q)
}

func (s *service) CreateGalleryAlbum(ctx context.Context, tenantID string, req GalleryAlbumPayload) (*GalleryAlbumResponse, error) {
	return s.repo.CreateGalleryAlbum(ctx, tenantID, req)
}

func (s *service) GetGalleryAlbum(ctx context.Context, tenantID string, id int64) (*GalleryAlbumResponse, error) {
	return s.repo.GetGalleryAlbum(ctx, tenantID, id)
}

func (s *service) UpdateGalleryAlbum(ctx context.Context, tenantID string, id int64, req GalleryAlbumPayload) error {
	return s.repo.UpdateGalleryAlbum(ctx, tenantID, id, req)
}

func (s *service) DeleteGalleryAlbum(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteGalleryAlbum(ctx, tenantID, id)
}

func (s *service) ListPublicGalleryAlbums(ctx context.Context, hostname string, page, limit int) ([]GalleryAlbumResponse, int64, error) {
	return s.repo.ListPublicGalleryAlbums(ctx, hostname, page, limit)
}

func (s *service) ListGalleryItems(ctx context.Context, tenantID string, q BaseListQuery, albumID string) ([]GalleryItemResponse, int64, error) {
	return s.repo.ListGalleryItems(ctx, tenantID, q, albumID)
}

func (s *service) CreateGalleryItem(ctx context.Context, tenantID string, req GalleryItemPayload) (*GalleryItemResponse, error) {
	return s.repo.CreateGalleryItem(ctx, tenantID, req)
}

func (s *service) GetGalleryItem(ctx context.Context, tenantID string, id int64) (*GalleryItemResponse, error) {
	return s.repo.GetGalleryItem(ctx, tenantID, id)
}

func (s *service) UpdateGalleryItem(ctx context.Context, tenantID string, id int64, req GalleryItemPayload) error {
	return s.repo.UpdateGalleryItem(ctx, tenantID, id, req)
}

func (s *service) DeleteGalleryItem(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteGalleryItem(ctx, tenantID, id)
}

func (s *service) ListPublicGalleryItems(ctx context.Context, hostname string, page, limit int) ([]GalleryItemResponse, int64, error) {
	return s.repo.ListPublicGalleryItems(ctx, hostname, page, limit)
}

func (s *service) ListManagementMembers(ctx context.Context, tenantID string, q BaseListQuery) ([]ManagementMemberResponse, int64, error) {
	return s.repo.ListManagementMembers(ctx, tenantID, q)
}

func (s *service) CreateManagementMember(ctx context.Context, tenantID string, req ManagementMemberPayload) (*ManagementMemberResponse, error) {
	return s.repo.CreateManagementMember(ctx, tenantID, req)
}

func (s *service) GetManagementMember(ctx context.Context, tenantID string, id int64) (*ManagementMemberResponse, error) {
	return s.repo.GetManagementMember(ctx, tenantID, id)
}

func (s *service) UpdateManagementMember(ctx context.Context, tenantID string, id int64, req ManagementMemberPayload) error {
	return s.repo.UpdateManagementMember(ctx, tenantID, id, req)
}

func (s *service) DeleteManagementMember(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteManagementMember(ctx, tenantID, id)
}

func (s *service) ListPublicManagementMembers(ctx context.Context, hostname string, page, limit int) ([]ManagementMemberResponse, int64, error) {
	return s.repo.ListPublicManagementMembers(ctx, hostname, page, limit)
}

func normalizeEventCategory(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	switch s {
	case "kajian":
		return "kajian_rutin"
	case "tabligh":
		return "tabligh_akbar"
	case "rapat":
		return "rapat_pengurus"
	case "sosial":
		return "kegiatan_sosial"
	}
	return s
}

func normalizeEventStatus(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	switch s {
	case "draft":
		return "upcoming"
	case "published":
		return "ongoing"
	case "archived":
		return "finished"
	}
	return s
}
