package community

import (
	"context"
	"regexp"
	"strings"
	"time"

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

var memberNameRoleRegex = regexp.MustCompile(`^[A-Za-z -]+$`)
var memberDigitsOnlyRegex = regexp.MustCompile(`^[0-9]+$`)

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) ListEvents(ctx context.Context, tenantID string, q EventListQuery) ([]EventResponse, int64, error) {
	return s.repo.ListEvents(ctx, tenantID, q)
}

func (s *service) CreateEvent(ctx context.Context, tenantID string, req EventPayload) (*EventResponse, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = normalizeEventCategory(req.Category)
	req.TimeMode = strings.ToLower(strings.TrimSpace(req.TimeMode))
	if req.Title == "" || req.Category == "" || req.StartDate == "" || req.TimeMode == "" || req.EndDate == nil || strings.TrimSpace(*req.EndDate) == "" {
		return nil, ErrValidation
	}
	if req.TimeMode == "exact_time" && (req.StartTime == nil || strings.TrimSpace(*req.StartTime) == "" || req.EndTime == nil || strings.TrimSpace(*req.EndTime) == "") {
		return nil, ErrValidation
	}
	if req.TimeMode == "after_prayer" && (req.AfterPrayer == nil || strings.TrimSpace(*req.AfterPrayer) == "" || req.EndTime == nil || strings.TrimSpace(*req.EndTime) == "") {
		return nil, ErrValidation
	}
	if !isValidEventDateRange(req.StartDate, *req.EndDate) {
		return nil, ErrValidation
	}
	if !isValidOptionalPerson(req.SpeakerName) || !isValidOptionalPerson(req.PersonInCharge) || !isValidOptionalPhone(req.ContactPhone) || !isValidOptionalLongText(req.Description) || !isValidOptionalLongText(req.NoteInternal) || !isValidOptionalLongText(req.NotePublic) {
		return nil, ErrValidation
	}
	req.Status = deriveEventStatus(req)
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
	req.TimeMode = strings.ToLower(strings.TrimSpace(req.TimeMode))
	if req.Title == "" || req.Category == "" || req.StartDate == "" || req.TimeMode == "" || req.EndDate == nil || strings.TrimSpace(*req.EndDate) == "" {
		return ErrValidation
	}
	if req.TimeMode == "exact_time" && (req.StartTime == nil || strings.TrimSpace(*req.StartTime) == "" || req.EndTime == nil || strings.TrimSpace(*req.EndTime) == "") {
		return ErrValidation
	}
	if req.TimeMode == "after_prayer" && (req.AfterPrayer == nil || strings.TrimSpace(*req.AfterPrayer) == "" || req.EndTime == nil || strings.TrimSpace(*req.EndTime) == "") {
		return ErrValidation
	}
	if !isValidEventDateRange(req.StartDate, *req.EndDate) {
		return ErrValidation
	}
	if !isValidOptionalPerson(req.SpeakerName) || !isValidOptionalPerson(req.PersonInCharge) || !isValidOptionalPhone(req.ContactPhone) || !isValidOptionalLongText(req.Description) || !isValidOptionalLongText(req.NoteInternal) || !isValidOptionalLongText(req.NotePublic) {
		return ErrValidation
	}
	req.Status = deriveEventStatus(req)
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
	if err := validateManagementMemberPayload(&req); err != nil {
		return nil, err
	}
	return s.repo.CreateManagementMember(ctx, tenantID, req)
}

func (s *service) GetManagementMember(ctx context.Context, tenantID string, id int64) (*ManagementMemberResponse, error) {
	return s.repo.GetManagementMember(ctx, tenantID, id)
}

func (s *service) UpdateManagementMember(ctx context.Context, tenantID string, id int64, req ManagementMemberPayload) error {
	if err := validateManagementMemberPayload(&req); err != nil {
		return err
	}
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

func isValidEventDateRange(startDate, endDate string) bool {
	start, err := time.Parse("2006-01-02", strings.TrimSpace(startDate))
	if err != nil {
		return false
	}
	end, err := time.Parse("2006-01-02", strings.TrimSpace(endDate))
	if err != nil {
		return false
	}
	return !end.Before(start)
}

func isValidOptionalPerson(v *string) bool {
	if v == nil {
		return true
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return true
	}
	return len(s) <= 25 && memberNameRoleRegex.MatchString(s)
}

func isValidOptionalPhone(v *string) bool {
	if v == nil {
		return true
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return true
	}
	return memberDigitsOnlyRegex.MatchString(s)
}

func isValidOptionalLongText(v *string) bool {
	if v == nil {
		return true
	}
	return len(strings.TrimSpace(*v)) <= 250
}

func deriveEventStatus(req EventPayload) string {
	now := time.Now()
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return "upcoming"
	}
	end, err := time.Parse("2006-01-02", *req.EndDate)
	if err != nil {
		return "upcoming"
	}

	if req.StartTime != nil && strings.TrimSpace(*req.StartTime) != "" {
		if t, err := time.Parse("15:04:05", strings.TrimSpace(*req.StartTime)); err == nil {
			start = time.Date(start.Year(), start.Month(), start.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())
		}
	}
	if req.EndTime != nil && strings.TrimSpace(*req.EndTime) != "" {
		if t, err := time.Parse("15:04:05", strings.TrimSpace(*req.EndTime)); err == nil {
			end = time.Date(end.Year(), end.Month(), end.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())
		}
	}

	if now.Before(start) {
		return "upcoming"
	}
	if now.After(end) {
		return "finished"
	}
	return "ongoing"
}

func validateManagementMemberPayload(req *ManagementMemberPayload) error {
	req.FullName = strings.TrimSpace(req.FullName)
	req.RoleTitle = strings.TrimSpace(req.RoleTitle)

	if req.FullName == "" || len(req.FullName) > 25 || !memberNameRoleRegex.MatchString(req.FullName) {
		return ErrValidation
	}
	if req.RoleTitle == "" || len(req.RoleTitle) > 25 || !memberNameRoleRegex.MatchString(req.RoleTitle) {
		return ErrValidation
	}

	if req.PhoneWhatsapp != nil {
		v := strings.TrimSpace(*req.PhoneWhatsapp)
		if v == "" {
			req.PhoneWhatsapp = nil
		} else {
			if !memberDigitsOnlyRegex.MatchString(v) {
				return ErrValidation
			}
			req.PhoneWhatsapp = &v
		}
	}

	if req.ProfileImageURL != nil {
		v := strings.TrimSpace(*req.ProfileImageURL)
		if v == "" {
			req.ProfileImageURL = nil
		} else {
			if len(v) > 500 {
				return ErrValidation
			}
			req.ProfileImageURL = &v
		}
	}

	return nil
}
