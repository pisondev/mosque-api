package community

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
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

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) ListEvents(ctx context.Context, tenantID string, q EventListQuery) ([]EventResponse, int64, error) {
	where := []string{"tenant_id=$1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.Status != "" {
		where = append(where, fmt.Sprintf("status=$%d", argN))
		args = append(args, q.Status)
		argN++
	}
	if q.Category != "" {
		where = append(where, fmt.Sprintf("category=$%d", argN))
		args = append(args, q.Category)
		argN++
	}
	if q.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR COALESCE(description,'') ILIKE $%d)", argN, argN))
		args = append(args, "%"+q.Search+"%")
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM events WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id,title,category::text,speaker_name,person_in_charge,description,note_public,note_internal,
		to_char(start_date,'YYYY-MM-DD'),
		CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date,'YYYY-MM-DD') END,
		time_mode::text,
		CASE WHEN start_time IS NULL THEN NULL ELSE to_char(start_time,'HH24:MI:SS') END,
		CASE WHEN end_time IS NULL THEN NULL ELSE to_char(end_time,'HH24:MI:SS') END,
		CASE WHEN after_prayer IS NULL THEN NULL ELSE after_prayer::text END,
		after_prayer_offset_min, repeat_pattern, COALESCE(repeat_weekdays,'{}'),
		CASE WHEN audience IS NULL THEN NULL ELSE audience::text END,
		capacity, fee_type, CASE WHEN fee_amount IS NULL THEN NULL ELSE fee_amount::text END,
		contact_phone, location_inside, location_outside, status::text, poster_image_url
		FROM events WHERE %s ORDER BY start_date DESC,id DESC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []EventResponse
	for rows.Next() {
		var it EventResponse
		if err := rows.Scan(&it.ID, &it.Title, &it.Category, &it.SpeakerName, &it.PersonInCharge, &it.Description, &it.NotePublic, &it.NoteInternal,
			&it.StartDate, &it.EndDate, &it.TimeMode, &it.StartTime, &it.EndTime, &it.AfterPrayer, &it.AfterPrayerOffsetMin,
			&it.RepeatPattern, &it.RepeatWeekdays, &it.Audience, &it.Capacity, &it.FeeType, &it.FeeAmount, &it.ContactPhone, &it.LocationInside, &it.LocationOutside, &it.Status, &it.PosterImageURL); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateEvent(ctx context.Context, tenantID string, req EventPayload) (*EventResponse, error) {
	var out EventResponse
	err := r.db.QueryRow(ctx, `INSERT INTO events
		(tenant_id,title,category,speaker_name,person_in_charge,description,note_public,note_internal,start_date,end_date,time_mode,start_time,end_time,after_prayer,after_prayer_offset_min,repeat_pattern,repeat_weekdays,audience,capacity,fee_type,fee_amount,contact_phone,location_inside,location_outside,status,poster_image_url)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::date,$10::date,$11,$12::time,$13::time,$14,$15,$16,$17,$18,$19,$20,$21::numeric,$22,$23,$24,$25,$26)
		RETURNING id,title,category::text,speaker_name,person_in_charge,description,note_public,note_internal,
		to_char(start_date,'YYYY-MM-DD'),
		CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date,'YYYY-MM-DD') END,
		time_mode::text,
		CASE WHEN start_time IS NULL THEN NULL ELSE to_char(start_time,'HH24:MI:SS') END,
		CASE WHEN end_time IS NULL THEN NULL ELSE to_char(end_time,'HH24:MI:SS') END,
		CASE WHEN after_prayer IS NULL THEN NULL ELSE after_prayer::text END,
		after_prayer_offset_min, repeat_pattern, COALESCE(repeat_weekdays,'{}'),
		CASE WHEN audience IS NULL THEN NULL ELSE audience::text END,
		capacity, fee_type, CASE WHEN fee_amount IS NULL THEN NULL ELSE fee_amount::text END,
		contact_phone, location_inside, location_outside, status::text, poster_image_url`,
		tenantID, req.Title, req.Category, nullableString(req.SpeakerName), nullableString(req.PersonInCharge), nullableString(req.Description), nullableString(req.NotePublic), nullableString(req.NoteInternal), req.StartDate, nullableString(req.EndDate), req.TimeMode, nullableString(req.StartTime), nullableString(req.EndTime), nullableString(req.AfterPrayer), nullableInt16(req.AfterPrayerOffsetMin), nullableString(req.RepeatPattern), req.RepeatWeekdays, nullableString(req.Audience), nullableInt(req.Capacity), nullableString(req.FeeType), nullableString(req.FeeAmount), nullableString(req.ContactPhone), nullableString(req.LocationInside), nullableString(req.LocationOutside), req.Status, nullableString(req.PosterImageURL)).
		Scan(&out.ID, &out.Title, &out.Category, &out.SpeakerName, &out.PersonInCharge, &out.Description, &out.NotePublic, &out.NoteInternal,
			&out.StartDate, &out.EndDate, &out.TimeMode, &out.StartTime, &out.EndTime, &out.AfterPrayer, &out.AfterPrayerOffsetMin,
			&out.RepeatPattern, &out.RepeatWeekdays, &out.Audience, &out.Capacity, &out.FeeType, &out.FeeAmount, &out.ContactPhone, &out.LocationInside, &out.LocationOutside, &out.Status, &out.PosterImageURL)
	if err != nil {
		return nil, mapPgError(err)
	}
	return &out, nil
}

func (r *repository) GetEvent(ctx context.Context, tenantID string, id int64) (*EventResponse, error) {
	var out EventResponse
	err := r.db.QueryRow(ctx, `SELECT id,title,category::text,speaker_name,person_in_charge,description,note_public,note_internal,
		to_char(start_date,'YYYY-MM-DD'),
		CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date,'YYYY-MM-DD') END,
		time_mode::text,
		CASE WHEN start_time IS NULL THEN NULL ELSE to_char(start_time,'HH24:MI:SS') END,
		CASE WHEN end_time IS NULL THEN NULL ELSE to_char(end_time,'HH24:MI:SS') END,
		CASE WHEN after_prayer IS NULL THEN NULL ELSE after_prayer::text END,
		after_prayer_offset_min, repeat_pattern, COALESCE(repeat_weekdays,'{}'),
		CASE WHEN audience IS NULL THEN NULL ELSE audience::text END,
		capacity, fee_type, CASE WHEN fee_amount IS NULL THEN NULL ELSE fee_amount::text END,
		contact_phone, location_inside, location_outside, status::text, poster_image_url
		FROM events WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.Title, &out.Category, &out.SpeakerName, &out.PersonInCharge, &out.Description, &out.NotePublic, &out.NoteInternal,
			&out.StartDate, &out.EndDate, &out.TimeMode, &out.StartTime, &out.EndTime, &out.AfterPrayer, &out.AfterPrayerOffsetMin,
			&out.RepeatPattern, &out.RepeatWeekdays, &out.Audience, &out.Capacity, &out.FeeType, &out.FeeAmount, &out.ContactPhone, &out.LocationInside, &out.LocationOutside, &out.Status, &out.PosterImageURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateEvent(ctx context.Context, tenantID string, id int64, req EventPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE events SET
		title=$1,category=$2,speaker_name=$3,person_in_charge=$4,description=$5,note_public=$6,note_internal=$7,
		start_date=$8::date,end_date=$9::date,time_mode=$10,start_time=$11::time,end_time=$12::time,after_prayer=$13,after_prayer_offset_min=$14,
		repeat_pattern=$15,repeat_weekdays=$16,audience=$17,capacity=$18,fee_type=$19,fee_amount=$20::numeric,contact_phone=$21,location_inside=$22,location_outside=$23,status=$24,poster_image_url=$25,updated_at=now()
		WHERE tenant_id=$26 AND id=$27`,
		req.Title, req.Category, nullableString(req.SpeakerName), nullableString(req.PersonInCharge), nullableString(req.Description), nullableString(req.NotePublic), nullableString(req.NoteInternal), req.StartDate, nullableString(req.EndDate), req.TimeMode, nullableString(req.StartTime), nullableString(req.EndTime), nullableString(req.AfterPrayer), nullableInt16(req.AfterPrayerOffsetMin), nullableString(req.RepeatPattern), req.RepeatWeekdays, nullableString(req.Audience), nullableInt(req.Capacity), nullableString(req.FeeType), nullableString(req.FeeAmount), nullableString(req.ContactPhone), nullableString(req.LocationInside), nullableString(req.LocationOutside), req.Status, nullableString(req.PosterImageURL), tenantID, id)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) UpdateEventStatus(ctx context.Context, tenantID string, id int64, status string) error {
	tag, err := r.db.Exec(ctx, `UPDATE events SET status=$1, updated_at=now() WHERE tenant_id=$2 AND id=$3`, status, tenantID, id)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteEvent(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM events WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicEvents(ctx context.Context, hostname string, page, limit int) ([]EventResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	return r.ListEvents(ctx, tenantID, EventListQuery{Status: "upcoming", Page: page, Limit: limit})
}

func (r *repository) ListGalleryAlbums(ctx context.Context, tenantID string, q BaseListQuery) ([]GalleryAlbumResponse, int64, error) {
	where := []string{"tenant_id=$1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.Search != "" {
		where = append(where, fmt.Sprintf("title ILIKE $%d", argN))
		args = append(args, "%"+q.Search+"%")
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM gallery_albums WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id,title,description,
		CASE WHEN start_date IS NULL THEN NULL ELSE to_char(start_date,'YYYY-MM-DD') END,
		CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date,'YYYY-MM-DD') END,
		media_kind::text
		FROM gallery_albums WHERE %s ORDER BY id DESC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []GalleryAlbumResponse
	for rows.Next() {
		var it GalleryAlbumResponse
		if err := rows.Scan(&it.ID, &it.Title, &it.Description, &it.StartDate, &it.EndDate, &it.MediaKind); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateGalleryAlbum(ctx context.Context, tenantID string, req GalleryAlbumPayload) (*GalleryAlbumResponse, error) {
	var out GalleryAlbumResponse
	err := r.db.QueryRow(ctx, `INSERT INTO gallery_albums (tenant_id,title,description,start_date,end_date,media_kind)
		VALUES ($1,$2,$3,$4::date,$5::date,$6)
		RETURNING id,title,description,
		CASE WHEN start_date IS NULL THEN NULL ELSE to_char(start_date,'YYYY-MM-DD') END,
		CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date,'YYYY-MM-DD') END, media_kind::text`,
		tenantID, req.Title, req.Description, req.StartDate, req.EndDate, req.MediaKind).
		Scan(&out.ID, &out.Title, &out.Description, &out.StartDate, &out.EndDate, &out.MediaKind)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetGalleryAlbum(ctx context.Context, tenantID string, id int64) (*GalleryAlbumResponse, error) {
	var out GalleryAlbumResponse
	err := r.db.QueryRow(ctx, `SELECT id,title,description,
		CASE WHEN start_date IS NULL THEN NULL ELSE to_char(start_date,'YYYY-MM-DD') END,
		CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date,'YYYY-MM-DD') END, media_kind::text
		FROM gallery_albums WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.Title, &out.Description, &out.StartDate, &out.EndDate, &out.MediaKind)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateGalleryAlbum(ctx context.Context, tenantID string, id int64, req GalleryAlbumPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE gallery_albums SET title=$1,description=$2,start_date=$3::date,end_date=$4::date,media_kind=$5,updated_at=now() WHERE tenant_id=$6 AND id=$7`,
		req.Title, req.Description, req.StartDate, req.EndDate, req.MediaKind, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteGalleryAlbum(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM gallery_albums WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicGalleryAlbums(ctx context.Context, hostname string, page, limit int) ([]GalleryAlbumResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	return r.ListGalleryAlbums(ctx, tenantID, BaseListQuery{Page: page, Limit: limit})
}

func (r *repository) ListGalleryItems(ctx context.Context, tenantID string, q BaseListQuery, albumID string) ([]GalleryItemResponse, int64, error) {
	where := []string{"tenant_id=$1"}
	args := []interface{}{tenantID}
	argN := 2
	if albumID != "" {
		where = append(where, fmt.Sprintf("album_id=$%d", argN))
		args = append(args, albumID)
		argN++
	}
	if q.Search != "" {
		where = append(where, fmt.Sprintf("COALESCE(caption,'') ILIKE $%d", argN))
		args = append(args, "%"+q.Search+"%")
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM gallery_items WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id, album_id, media_type::text, media_url, caption,
		CASE WHEN taken_at IS NULL THEN NULL ELSE to_char(taken_at at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"') END,
		location_note, is_highlight, sort_order
		FROM gallery_items WHERE %s ORDER BY sort_order ASC, id DESC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []GalleryItemResponse
	for rows.Next() {
		var it GalleryItemResponse
		if err := rows.Scan(&it.ID, &it.AlbumID, &it.MediaType, &it.MediaURL, &it.Caption, &it.TakenAt, &it.LocationNote, &it.IsHighlight, &it.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateGalleryItem(ctx context.Context, tenantID string, req GalleryItemPayload) (*GalleryItemResponse, error) {
	var out GalleryItemResponse
	err := r.db.QueryRow(ctx, `INSERT INTO gallery_items
		(tenant_id,album_id,media_type,media_url,caption,taken_at,location_note,is_highlight,sort_order)
		VALUES ($1,$2,$3,$4,$5,$6::timestamptz,$7,$8,$9)
		RETURNING id,album_id,media_type::text,media_url,caption,
		CASE WHEN taken_at IS NULL THEN NULL ELSE to_char(taken_at at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"') END,
		location_note,is_highlight,sort_order`,
		tenantID, req.AlbumID, req.MediaType, req.MediaURL, req.Caption, req.TakenAt, req.LocationNote, req.IsHighlight, req.SortOrder).
		Scan(&out.ID, &out.AlbumID, &out.MediaType, &out.MediaURL, &out.Caption, &out.TakenAt, &out.LocationNote, &out.IsHighlight, &out.SortOrder)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetGalleryItem(ctx context.Context, tenantID string, id int64) (*GalleryItemResponse, error) {
	var out GalleryItemResponse
	err := r.db.QueryRow(ctx, `SELECT id,album_id,media_type::text,media_url,caption,
		CASE WHEN taken_at IS NULL THEN NULL ELSE to_char(taken_at at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS"Z"') END,
		location_note,is_highlight,sort_order
		FROM gallery_items WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.AlbumID, &out.MediaType, &out.MediaURL, &out.Caption, &out.TakenAt, &out.LocationNote, &out.IsHighlight, &out.SortOrder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateGalleryItem(ctx context.Context, tenantID string, id int64, req GalleryItemPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE gallery_items SET album_id=$1,media_type=$2,media_url=$3,caption=$4,taken_at=$5::timestamptz,location_note=$6,is_highlight=$7,sort_order=$8,updated_at=now() WHERE tenant_id=$9 AND id=$10`,
		req.AlbumID, req.MediaType, req.MediaURL, req.Caption, req.TakenAt, req.LocationNote, req.IsHighlight, req.SortOrder, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteGalleryItem(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM gallery_items WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicGalleryItems(ctx context.Context, hostname string, page, limit int) ([]GalleryItemResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	return r.ListGalleryItems(ctx, tenantID, BaseListQuery{Page: page, Limit: limit}, "")
}

func (r *repository) ListManagementMembers(ctx context.Context, tenantID string, q BaseListQuery) ([]ManagementMemberResponse, int64, error) {
	where := []string{"tenant_id=$1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.Search != "" {
		where = append(where, fmt.Sprintf("(full_name ILIKE $%d OR role_title ILIKE $%d)", argN, argN))
		args = append(args, "%"+q.Search+"%")
		argN++
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM management_members WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id,full_name,role_title,phone_whatsapp,profile_image_url,show_public,sort_order
		FROM management_members WHERE %s ORDER BY sort_order ASC, id DESC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ManagementMemberResponse
	for rows.Next() {
		var it ManagementMemberResponse
		if err := rows.Scan(&it.ID, &it.FullName, &it.RoleTitle, &it.PhoneWhatsapp, &it.ProfileImageURL, &it.ShowPublic, &it.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateManagementMember(ctx context.Context, tenantID string, req ManagementMemberPayload) (*ManagementMemberResponse, error) {
	var out ManagementMemberResponse
	err := r.db.QueryRow(ctx, `INSERT INTO management_members
		(tenant_id,full_name,role_title,phone_whatsapp,profile_image_url,show_public,sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id,full_name,role_title,phone_whatsapp,profile_image_url,show_public,sort_order`,
		tenantID, req.FullName, req.RoleTitle, req.PhoneWhatsapp, req.ProfileImageURL, req.ShowPublic, req.SortOrder).
		Scan(&out.ID, &out.FullName, &out.RoleTitle, &out.PhoneWhatsapp, &out.ProfileImageURL, &out.ShowPublic, &out.SortOrder)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetManagementMember(ctx context.Context, tenantID string, id int64) (*ManagementMemberResponse, error) {
	var out ManagementMemberResponse
	err := r.db.QueryRow(ctx, `SELECT id,full_name,role_title,phone_whatsapp,profile_image_url,show_public,sort_order
		FROM management_members WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.FullName, &out.RoleTitle, &out.PhoneWhatsapp, &out.ProfileImageURL, &out.ShowPublic, &out.SortOrder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateManagementMember(ctx context.Context, tenantID string, id int64, req ManagementMemberPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE management_members SET full_name=$1,role_title=$2,phone_whatsapp=$3,profile_image_url=$4,show_public=$5,sort_order=$6,updated_at=now()
		WHERE tenant_id=$7 AND id=$8`, req.FullName, req.RoleTitle, req.PhoneWhatsapp, req.ProfileImageURL, req.ShowPublic, req.SortOrder, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteManagementMember(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM management_members WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicManagementMembers(ctx context.Context, hostname string, page, limit int) ([]ManagementMemberResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	where := []string{"tenant_id=$1", "show_public=true"}
	args := []interface{}{tenantID}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM management_members WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, `SELECT id,full_name,role_title,phone_whatsapp,profile_image_url,show_public,sort_order
		FROM management_members WHERE `+whereSQL+` ORDER BY sort_order ASC, id DESC LIMIT $2 OFFSET $3`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ManagementMemberResponse
	for rows.Next() {
		var it ManagementMemberResponse
		if err := rows.Scan(&it.ID, &it.FullName, &it.RoleTitle, &it.PhoneWhatsapp, &it.ProfileImageURL, &it.ShowPublic, &it.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) resolveTenantByHostname(ctx context.Context, hostname string) (string, error) {
	var tenantID string
	err := r.db.QueryRow(ctx, `SELECT tenant_id::text FROM website_domains WHERE hostname=$1 AND status='active' LIMIT 1`, hostname).Scan(&tenantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}
	return tenantID, nil
}

func mapPgError(err error) error {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}
	switch pgErr.Code {
	case "22P02", "22007", "23502", "23514":
		return ErrValidation
	}
	return err
}

func nullableString(v *string) interface{} {
	if v == nil {
		return nil
	}
	return *v
}

func nullableInt16(v *int16) interface{} {
	if v == nil {
		return nil
	}
	return *v
}

func nullableInt(v *int) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
