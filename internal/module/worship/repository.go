package worship

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
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

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetPrayerTimeSettings(ctx context.Context, tenantID string) (*PrayerTimeSettingsResponse, error) {
	var out PrayerTimeSettingsResponse
	err := r.db.QueryRow(ctx, `SELECT id, timezone, location_mode::text, city_name, latitude, longitude, calc_method, asr_madhhab,
		adj_subuh_min, adj_dzuhur_min, adj_ashar_min, adj_maghrib_min, adj_isya_min
		FROM prayer_time_settings WHERE tenant_id=$1`, tenantID).
		Scan(&out.ID, &out.Timezone, &out.LocationMode, &out.CityName, &out.Latitude, &out.Longitude, &out.CalcMethod, &out.AsrMadhhab,
			&out.AdjSubuhMin, &out.AdjDzuhurMin, &out.AdjAsharMin, &out.AdjMaghribMin, &out.AdjIsyaMin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpsertPrayerTimeSettings(ctx context.Context, tenantID string, req PrayerTimeSettingsRequest) (*PrayerTimeSettingsResponse, error) {
	var out PrayerTimeSettingsResponse
	err := r.db.QueryRow(ctx, `INSERT INTO prayer_time_settings
		(tenant_id, timezone, location_mode, city_name, latitude, longitude, calc_method, asr_madhhab,
		adj_subuh_min, adj_dzuhur_min, adj_ashar_min, adj_maghrib_min, adj_isya_min)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (tenant_id) DO UPDATE SET
		timezone=EXCLUDED.timezone, location_mode=EXCLUDED.location_mode, city_name=EXCLUDED.city_name,
		latitude=EXCLUDED.latitude, longitude=EXCLUDED.longitude, calc_method=EXCLUDED.calc_method,
		asr_madhhab=EXCLUDED.asr_madhhab, adj_subuh_min=EXCLUDED.adj_subuh_min, adj_dzuhur_min=EXCLUDED.adj_dzuhur_min,
		adj_ashar_min=EXCLUDED.adj_ashar_min, adj_maghrib_min=EXCLUDED.adj_maghrib_min, adj_isya_min=EXCLUDED.adj_isya_min, updated_at=now()
		RETURNING id, timezone, location_mode::text, city_name, latitude, longitude, calc_method, asr_madhhab,
		adj_subuh_min, adj_dzuhur_min, adj_ashar_min, adj_maghrib_min, adj_isya_min`,
		tenantID, req.Timezone, req.LocationMode, req.CityName, req.Latitude, req.Longitude, req.CalcMethod, req.AsrMadhhab,
		req.AdjSubuhMin, req.AdjDzuhurMin, req.AdjAsharMin, req.AdjMaghribMin, req.AdjIsyaMin).
		Scan(&out.ID, &out.Timezone, &out.LocationMode, &out.CityName, &out.Latitude, &out.Longitude, &out.CalcMethod, &out.AsrMadhhab,
			&out.AdjSubuhMin, &out.AdjDzuhurMin, &out.AdjAsharMin, &out.AdjMaghribMin, &out.AdjIsyaMin)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) ListPrayerTimesDaily(ctx context.Context, tenantID string, q PrayerTimesDailyQuery) ([]PrayerTimesDailyResponse, int64, error) {
	where := []string{"tenant_id = $1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.From != "" {
		where = append(where, fmt.Sprintf("day_date >= $%d::date", argN))
		args = append(args, q.From)
		argN++
	}
	if q.To != "" {
		where = append(where, fmt.Sprintf("day_date <= $%d::date", argN))
		args = append(args, q.To)
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM prayer_times_daily WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id, to_char(day_date,'YYYY-MM-DD'),
		to_char(subuh_time,'HH24:MI:SS'), to_char(dzuhur_time,'HH24:MI:SS'), to_char(ashar_time,'HH24:MI:SS'),
		to_char(maghrib_time,'HH24:MI:SS'), to_char(isya_time,'HH24:MI:SS'),
		CASE WHEN sunrise_time IS NULL THEN NULL ELSE to_char(sunrise_time,'HH24:MI:SS') END,
		CASE WHEN dhuha_time IS NULL THEN NULL ELSE to_char(dhuha_time,'HH24:MI:SS') END,
		source_label,
		CASE WHEN fetched_at IS NULL THEN NULL ELSE to_char(fetched_at at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"') END
		FROM prayer_times_daily WHERE %s ORDER BY day_date ASC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []PrayerTimesDailyResponse
	for rows.Next() {
		var it PrayerTimesDailyResponse
		if err := rows.Scan(&it.ID, &it.DayDate, &it.SubuhTime, &it.DzuhurTime, &it.AsharTime, &it.MaghribTime, &it.IsyaTime, &it.SunriseTime, &it.DhuhaTime, &it.SourceLabel, &it.FetchedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreatePrayerTimesDaily(ctx context.Context, tenantID string, req PrayerTimesDailyPayload) (*PrayerTimesDailyResponse, error) {
	var out PrayerTimesDailyResponse
	err := r.db.QueryRow(ctx, `INSERT INTO prayer_times_daily
		(tenant_id, day_date, subuh_time, dzuhur_time, ashar_time, maghrib_time, isya_time, sunrise_time, dhuha_time, source_label, fetched_at)
		VALUES ($1,$2::date,$3::time,$4::time,$5::time,$6::time,$7::time,$8::time,$9::time,$10,$11::timestamptz)
		RETURNING id, to_char(day_date,'YYYY-MM-DD'),
		to_char(subuh_time,'HH24:MI:SS'), to_char(dzuhur_time,'HH24:MI:SS'), to_char(ashar_time,'HH24:MI:SS'),
		to_char(maghrib_time,'HH24:MI:SS'), to_char(isya_time,'HH24:MI:SS'),
		CASE WHEN sunrise_time IS NULL THEN NULL ELSE to_char(sunrise_time,'HH24:MI:SS') END,
		CASE WHEN dhuha_time IS NULL THEN NULL ELSE to_char(dhuha_time,'HH24:MI:SS') END,
		source_label,
		CASE WHEN fetched_at IS NULL THEN NULL ELSE to_char(fetched_at at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"') END`,
		tenantID, req.DayDate, req.SubuhTime, req.DzuhurTime, req.AsharTime, req.MaghribTime, req.IsyaTime, req.SunriseTime, req.DhuhaTime, req.SourceLabel, req.FetchedAt).
		Scan(&out.ID, &out.DayDate, &out.SubuhTime, &out.DzuhurTime, &out.AsharTime, &out.MaghribTime, &out.IsyaTime, &out.SunriseTime, &out.DhuhaTime, &out.SourceLabel, &out.FetchedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetPrayerTimesDaily(ctx context.Context, tenantID string, id int64) (*PrayerTimesDailyResponse, error) {
	var out PrayerTimesDailyResponse
	err := r.db.QueryRow(ctx, `SELECT id, to_char(day_date,'YYYY-MM-DD'),
		to_char(subuh_time,'HH24:MI:SS'), to_char(dzuhur_time,'HH24:MI:SS'), to_char(ashar_time,'HH24:MI:SS'),
		to_char(maghrib_time,'HH24:MI:SS'), to_char(isya_time,'HH24:MI:SS'),
		CASE WHEN sunrise_time IS NULL THEN NULL ELSE to_char(sunrise_time,'HH24:MI:SS') END,
		CASE WHEN dhuha_time IS NULL THEN NULL ELSE to_char(dhuha_time,'HH24:MI:SS') END,
		source_label,
		CASE WHEN fetched_at IS NULL THEN NULL ELSE to_char(fetched_at at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"') END
		FROM prayer_times_daily WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.DayDate, &out.SubuhTime, &out.DzuhurTime, &out.AsharTime, &out.MaghribTime, &out.IsyaTime, &out.SunriseTime, &out.DhuhaTime, &out.SourceLabel, &out.FetchedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdatePrayerTimesDaily(ctx context.Context, tenantID string, id int64, req PrayerTimesDailyPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE prayer_times_daily SET
		day_date=$1::date, subuh_time=$2::time, dzuhur_time=$3::time, ashar_time=$4::time, maghrib_time=$5::time, isya_time=$6::time,
		sunrise_time=$7::time, dhuha_time=$8::time, source_label=$9, fetched_at=$10::timestamptz, updated_at=now()
		WHERE tenant_id=$11 AND id=$12`,
		req.DayDate, req.SubuhTime, req.DzuhurTime, req.AsharTime, req.MaghribTime, req.IsyaTime, req.SunriseTime, req.DhuhaTime, req.SourceLabel, req.FetchedAt, tenantID, id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeletePrayerTimesDaily(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM prayer_times_daily WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPrayerDuties(ctx context.Context, tenantID string, q PrayerDutiesQuery) ([]PrayerDutyResponse, int64, error) {
	where := []string{"tenant_id = $1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.From != "" {
		where = append(where, fmt.Sprintf("duty_date >= $%d::date", argN))
		args = append(args, q.From)
		argN++
	}
	if q.To != "" {
		where = append(where, fmt.Sprintf("duty_date <= $%d::date", argN))
		args = append(args, q.To)
		argN++
	}
	if q.Category != "" {
		where = append(where, fmt.Sprintf("category = $%d", argN))
		args = append(args, q.Category)
		argN++
	}
	if q.Prayer != "" {
		where = append(where, fmt.Sprintf("prayer = $%d", argN))
		args = append(args, q.Prayer)
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM prayer_duties WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id, category::text, to_char(duty_date,'YYYY-MM-DD'),
		CASE WHEN prayer IS NULL THEN NULL ELSE prayer::text END, khatib_name, imam_name, muadzin_name,
		CASE WHEN first_adhan_time IS NULL THEN NULL ELSE to_char(first_adhan_time,'HH24:MI:SS') END,
		CASE WHEN khutbah_start_time IS NULL THEN NULL ELSE to_char(khutbah_start_time,'HH24:MI:SS') END,
		khutbah_topic, note
		FROM prayer_duties WHERE %s ORDER BY duty_date ASC, id ASC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []PrayerDutyResponse
	for rows.Next() {
		var it PrayerDutyResponse
		if err := rows.Scan(&it.ID, &it.Category, &it.DutyDate, &it.Prayer, &it.KhatibName, &it.ImamName, &it.MuadzinName, &it.FirstAdhanTime, &it.KhutbahStartTime, &it.KhutbahTopic, &it.Note); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreatePrayerDuty(ctx context.Context, tenantID string, req PrayerDutyPayload) (*PrayerDutyResponse, error) {
	var out PrayerDutyResponse
	err := r.db.QueryRow(ctx, `INSERT INTO prayer_duties
		(tenant_id, category, duty_date, prayer, khatib_name, imam_name, muadzin_name, first_adhan_time, khutbah_start_time, khutbah_topic, note)
		VALUES ($1,$2,$3::date,$4,$5,$6,$7,$8::time,$9::time,$10,$11)
		RETURNING id, category::text, to_char(duty_date,'YYYY-MM-DD'), CASE WHEN prayer IS NULL THEN NULL ELSE prayer::text END,
		khatib_name, imam_name, muadzin_name,
		CASE WHEN first_adhan_time IS NULL THEN NULL ELSE to_char(first_adhan_time,'HH24:MI:SS') END,
		CASE WHEN khutbah_start_time IS NULL THEN NULL ELSE to_char(khutbah_start_time,'HH24:MI:SS') END,
		khutbah_topic, note`,
		tenantID, req.Category, req.DutyDate, req.Prayer, req.KhatibName, req.ImamName, req.MuadzinName, req.FirstAdhanTime, req.KhutbahStartTime, req.KhutbahTopic, req.Note).
		Scan(&out.ID, &out.Category, &out.DutyDate, &out.Prayer, &out.KhatibName, &out.ImamName, &out.MuadzinName, &out.FirstAdhanTime, &out.KhutbahStartTime, &out.KhutbahTopic, &out.Note)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetPrayerDuty(ctx context.Context, tenantID string, id int64) (*PrayerDutyResponse, error) {
	var out PrayerDutyResponse
	err := r.db.QueryRow(ctx, `SELECT id, category::text, to_char(duty_date,'YYYY-MM-DD'),
		CASE WHEN prayer IS NULL THEN NULL ELSE prayer::text END, khatib_name, imam_name, muadzin_name,
		CASE WHEN first_adhan_time IS NULL THEN NULL ELSE to_char(first_adhan_time,'HH24:MI:SS') END,
		CASE WHEN khutbah_start_time IS NULL THEN NULL ELSE to_char(khutbah_start_time,'HH24:MI:SS') END,
		khutbah_topic, note
		FROM prayer_duties WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.Category, &out.DutyDate, &out.Prayer, &out.KhatibName, &out.ImamName, &out.MuadzinName, &out.FirstAdhanTime, &out.KhutbahStartTime, &out.KhutbahTopic, &out.Note)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdatePrayerDuty(ctx context.Context, tenantID string, id int64, req PrayerDutyPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE prayer_duties SET
		category=$1, duty_date=$2::date, prayer=$3, khatib_name=$4, imam_name=$5, muadzin_name=$6,
		first_adhan_time=$7::time, khutbah_start_time=$8::time, khutbah_topic=$9, note=$10, updated_at=now()
		WHERE tenant_id=$11 AND id=$12`,
		req.Category, req.DutyDate, req.Prayer, req.KhatibName, req.ImamName, req.MuadzinName, req.FirstAdhanTime, req.KhutbahStartTime, req.KhutbahTopic, req.Note, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeletePrayerDuty(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM prayer_duties WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListSpecialDays(ctx context.Context, tenantID string, q SpecialDaysQuery) ([]SpecialDayResponse, int64, error) {
	where := []string{"tenant_id = $1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.Year != "" {
		where = append(where, fmt.Sprintf("extract(year from day_date) = $%d::int", argN))
		args = append(args, q.Year)
		argN++
	}
	if q.Kind != "" {
		where = append(where, fmt.Sprintf("kind = $%d", argN))
		args = append(args, q.Kind)
		argN++
	}
	if q.From != "" {
		where = append(where, fmt.Sprintf("day_date >= $%d::date", argN))
		args = append(args, q.From)
		argN++
	}
	if q.To != "" {
		where = append(where, fmt.Sprintf("day_date <= $%d::date", argN))
		args = append(args, q.To)
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM special_days WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id, kind::text, title, to_char(day_date,'YYYY-MM-DD'), location_note,
		CASE WHEN start_time IS NULL THEN NULL ELSE to_char(start_time,'HH24:MI:SS') END, note, imam_name, khatib_name, muadzin_name
		FROM special_days WHERE %s ORDER BY day_date ASC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []SpecialDayResponse
	for rows.Next() {
		var it SpecialDayResponse
		if err := rows.Scan(&it.ID, &it.Kind, &it.Title, &it.DayDate, &it.LocationNote, &it.StartTime, &it.Note, &it.ImamName, &it.KhatibName, &it.MuadzinName); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateSpecialDay(ctx context.Context, tenantID string, req SpecialDayPayload) (*SpecialDayResponse, error) {
	var out SpecialDayResponse
	err := r.db.QueryRow(ctx, `INSERT INTO special_days
		(tenant_id, kind, title, day_date, location_note, start_time, note, imam_name, khatib_name, muadzin_name)
		VALUES ($1,$2,$3,$4::date,$5,$6::time,$7,$8,$9,$10)
		RETURNING id, kind::text, title, to_char(day_date,'YYYY-MM-DD'), location_note,
		CASE WHEN start_time IS NULL THEN NULL ELSE to_char(start_time,'HH24:MI:SS') END, note, imam_name, khatib_name, muadzin_name`,
		tenantID, req.Kind, req.Title, req.DayDate, req.LocationNote, req.StartTime, req.Note, req.ImamName, req.KhatibName, req.MuadzinName).
		Scan(&out.ID, &out.Kind, &out.Title, &out.DayDate, &out.LocationNote, &out.StartTime, &out.Note, &out.ImamName, &out.KhatibName, &out.MuadzinName)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetSpecialDay(ctx context.Context, tenantID string, id int64) (*SpecialDayResponse, error) {
	var out SpecialDayResponse
	err := r.db.QueryRow(ctx, `SELECT id, kind::text, title, to_char(day_date,'YYYY-MM-DD'), location_note,
		CASE WHEN start_time IS NULL THEN NULL ELSE to_char(start_time,'HH24:MI:SS') END, note, imam_name, khatib_name, muadzin_name
		FROM special_days WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.Kind, &out.Title, &out.DayDate, &out.LocationNote, &out.StartTime, &out.Note, &out.ImamName, &out.KhatibName, &out.MuadzinName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateSpecialDay(ctx context.Context, tenantID string, id int64, req SpecialDayPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE special_days SET
		kind=$1, title=$2, day_date=$3::date, location_note=$4, start_time=$5::time, note=$6, imam_name=$7, khatib_name=$8, muadzin_name=$9, updated_at=now()
		WHERE tenant_id=$10 AND id=$11`,
		req.Kind, req.Title, req.DayDate, req.LocationNote, req.StartTime, req.Note, req.ImamName, req.KhatibName, req.MuadzinName, tenantID, id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteSpecialDay(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM special_days WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) GetPrayerCalendar(ctx context.Context, tenantID, from, to string) (map[string]interface{}, error) {
	qDaily := PrayerTimesDailyQuery{From: from, To: to, Page: 1, Limit: 1000}
	daily, _, err := r.ListPrayerTimesDaily(ctx, tenantID, qDaily)
	if err != nil {
		return nil, err
	}
	qDuty := PrayerDutiesQuery{From: from, To: to, Page: 1, Limit: 1000}
	duties, _, err := r.ListPrayerDuties(ctx, tenantID, qDuty)
	if err != nil {
		return nil, err
	}
	qSpecial := SpecialDaysQuery{From: from, To: to, Page: 1, Limit: 1000}
	special, _, err := r.ListSpecialDays(ctx, tenantID, qSpecial)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"prayer_times_daily": daily,
		"prayer_duties":      duties,
		"special_days":       special,
	}, nil
}
