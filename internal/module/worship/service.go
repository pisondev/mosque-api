package worship

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service interface {
	GetPrayerTimeSettings(ctx context.Context, tenantID string) (*PrayerTimeSettingsResponse, error)
	UpsertPrayerTimeSettings(ctx context.Context, tenantID string, req PrayerTimeSettingsRequest) (*PrayerTimeSettingsResponse, error)
	ListPrayerTimesDaily(ctx context.Context, tenantID string, q PrayerTimesDailyQuery) ([]PrayerTimesDailyResponse, int64, error)
	CreatePrayerTimesDaily(ctx context.Context, tenantID string, req PrayerTimesDailyPayload) (*PrayerTimesDailyResponse, error)
	GetPrayerTimesDaily(ctx context.Context, tenantID string, id int64) (*PrayerTimesDailyResponse, error)
	UpdatePrayerTimesDaily(ctx context.Context, tenantID string, id int64, req PrayerTimesDailyPayload) error
	DeletePrayerTimesDaily(ctx context.Context, tenantID string, id int64) error
	ListPrayerDuties(ctx context.Context, tenantID string, q PrayerDutiesQuery) ([]PrayerDutyResponse, int64, error)
	CreatePrayerDuty(ctx context.Context, tenantID string, req PrayerDutyPayload) (*PrayerDutyResponse, error)
	GetPrayerDuty(ctx context.Context, tenantID string, id int64) (*PrayerDutyResponse, error)
	UpdatePrayerDuty(ctx context.Context, tenantID string, id int64, req PrayerDutyPayload) error
	DeletePrayerDuty(ctx context.Context, tenantID string, id int64) error
	ListSpecialDays(ctx context.Context, tenantID string, q SpecialDaysQuery) ([]SpecialDayResponse, int64, error)
	CreateSpecialDay(ctx context.Context, tenantID string, req SpecialDayPayload) (*SpecialDayResponse, error)
	GetSpecialDay(ctx context.Context, tenantID string, id int64) (*SpecialDayResponse, error)
	UpdateSpecialDay(ctx context.Context, tenantID string, id int64, req SpecialDayPayload) error
	DeleteSpecialDay(ctx context.Context, tenantID string, id int64) error
	GetPrayerCalendar(ctx context.Context, tenantID, from, to string) (map[string]interface{}, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) GetPrayerTimeSettings(ctx context.Context, tenantID string) (*PrayerTimeSettingsResponse, error) {
	return s.repo.GetPrayerTimeSettings(ctx, tenantID)
}

func (s *service) UpsertPrayerTimeSettings(ctx context.Context, tenantID string, req PrayerTimeSettingsRequest) (*PrayerTimeSettingsResponse, error) {
	req.Timezone = strings.TrimSpace(req.Timezone)
	req.LocationMode = strings.TrimSpace(req.LocationMode)
	return s.repo.UpsertPrayerTimeSettings(ctx, tenantID, req)
}

func (s *service) ListPrayerTimesDaily(ctx context.Context, tenantID string, q PrayerTimesDailyQuery) ([]PrayerTimesDailyResponse, int64, error) {
	return s.repo.ListPrayerTimesDaily(ctx, tenantID, q)
}

func (s *service) CreatePrayerTimesDaily(ctx context.Context, tenantID string, req PrayerTimesDailyPayload) (*PrayerTimesDailyResponse, error) {
	return s.repo.CreatePrayerTimesDaily(ctx, tenantID, req)
}

func (s *service) GetPrayerTimesDaily(ctx context.Context, tenantID string, id int64) (*PrayerTimesDailyResponse, error) {
	return s.repo.GetPrayerTimesDaily(ctx, tenantID, id)
}

func (s *service) UpdatePrayerTimesDaily(ctx context.Context, tenantID string, id int64, req PrayerTimesDailyPayload) error {
	return s.repo.UpdatePrayerTimesDaily(ctx, tenantID, id, req)
}

func (s *service) DeletePrayerTimesDaily(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeletePrayerTimesDaily(ctx, tenantID, id)
}

func (s *service) ListPrayerDuties(ctx context.Context, tenantID string, q PrayerDutiesQuery) ([]PrayerDutyResponse, int64, error) {
	return s.repo.ListPrayerDuties(ctx, tenantID, q)
}

func (s *service) CreatePrayerDuty(ctx context.Context, tenantID string, req PrayerDutyPayload) (*PrayerDutyResponse, error) {
	return s.repo.CreatePrayerDuty(ctx, tenantID, req)
}

func (s *service) GetPrayerDuty(ctx context.Context, tenantID string, id int64) (*PrayerDutyResponse, error) {
	return s.repo.GetPrayerDuty(ctx, tenantID, id)
}

func (s *service) UpdatePrayerDuty(ctx context.Context, tenantID string, id int64, req PrayerDutyPayload) error {
	return s.repo.UpdatePrayerDuty(ctx, tenantID, id, req)
}

func (s *service) DeletePrayerDuty(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeletePrayerDuty(ctx, tenantID, id)
}

func (s *service) ListSpecialDays(ctx context.Context, tenantID string, q SpecialDaysQuery) ([]SpecialDayResponse, int64, error) {
	return s.repo.ListSpecialDays(ctx, tenantID, q)
}

func (s *service) CreateSpecialDay(ctx context.Context, tenantID string, req SpecialDayPayload) (*SpecialDayResponse, error) {
	return s.repo.CreateSpecialDay(ctx, tenantID, req)
}

func (s *service) GetSpecialDay(ctx context.Context, tenantID string, id int64) (*SpecialDayResponse, error) {
	return s.repo.GetSpecialDay(ctx, tenantID, id)
}

func (s *service) UpdateSpecialDay(ctx context.Context, tenantID string, id int64, req SpecialDayPayload) error {
	return s.repo.UpdateSpecialDay(ctx, tenantID, id, req)
}

func (s *service) DeleteSpecialDay(ctx context.Context, tenantID string, id int64) error {
	return s.repo.DeleteSpecialDay(ctx, tenantID, id)
}

func (s *service) GetPrayerCalendar(ctx context.Context, tenantID, from, to string) (map[string]interface{}, error) {
	return s.repo.GetPrayerCalendar(ctx, tenantID, from, to)
}
