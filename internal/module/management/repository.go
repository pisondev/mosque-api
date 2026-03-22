package management

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetTenantContext(ctx context.Context, tenantID, email string) (map[string]interface{}, error)
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
	UpdateTenantSetup(ctx context.Context, tenantID, name, subdomain string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetTenantContext(ctx context.Context, tenantID, email string) (map[string]interface{}, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tenants WHERE id=$1)`, tenantID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}
	return map[string]interface{}{
		"tenant_id": tenantID,
		"email":     email,
	}, nil
}

func (r *repository) ListDomains(ctx context.Context, tenantID string, q DomainListQuery) ([]DomainResponse, int64, error) {
	where := []string{"tenant_id = $1"}
	args := []interface{}{tenantID}
	argN := 2

	if q.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argN))
		args = append(args, q.Status)
		argN++
	}
	if q.DomainType != "" {
		where = append(where, fmt.Sprintf("domain_type = $%d", argN))
		args = append(args, q.DomainType)
		argN++
	}

	whereSQL := strings.Join(where, " AND ")
	countSQL := "SELECT COUNT(*) FROM website_domains WHERE " + whereSQL
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id, domain_type::text, hostname, status::text, verified_at
		FROM website_domains WHERE %s ORDER BY id DESC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []DomainResponse
	for rows.Next() {
		var item DomainResponse
		if err := rows.Scan(&item.ID, &item.DomainType, &item.Hostname, &item.Status, &item.VerifiedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (r *repository) CreateDomain(ctx context.Context, tenantID string, req CreateDomainRequest) (*DomainResponse, error) {
	var item DomainResponse
	err := r.db.QueryRow(ctx, `INSERT INTO website_domains (tenant_id, domain_type, hostname)
		VALUES ($1,$2,$3)
		RETURNING id, domain_type::text, hostname, status::text, verified_at`,
		tenantID, req.DomainType, req.Hostname).
		Scan(&item.ID, &item.DomainType, &item.Hostname, &item.Status, &item.VerifiedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) UpdateDomain(ctx context.Context, tenantID string, id int64, req UpdateDomainRequest) error {
	tag, err := r.db.Exec(ctx, `UPDATE website_domains SET status=$1, updated_at=now() WHERE id=$2 AND tenant_id=$3`, req.Status, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteDomain(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM website_domains WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) GetProfile(ctx context.Context, tenantID string) (*ProfileResponse, error) {
	var item ProfileResponse
	err := r.db.QueryRow(ctx, `SELECT official_name, kind::text, COALESCE(short_name,''), COALESCE(city,''), COALESCE(address_full,''), COALESCE(phone_whatsapp,''), COALESCE(email,'')
		FROM masjid_profiles WHERE tenant_id=$1`, tenantID).
		Scan(&item.OfficialName, &item.Kind, &item.ShortName, &item.City, &item.AddressFull, &item.PhoneWA, &item.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) UpsertProfile(ctx context.Context, tenantID string, req ProfileRequest) (*ProfileResponse, error) {
	var item ProfileResponse
	err := r.db.QueryRow(ctx, `INSERT INTO masjid_profiles (tenant_id, official_name, kind, short_name, city, address_full, phone_whatsapp, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (tenant_id) DO UPDATE SET
		official_name=EXCLUDED.official_name,
		kind=EXCLUDED.kind,
		short_name=EXCLUDED.short_name,
		city=EXCLUDED.city,
		address_full=EXCLUDED.address_full,
		phone_whatsapp=EXCLUDED.phone_whatsapp,
		email=EXCLUDED.email,
		updated_at=now()
		RETURNING official_name, kind::text, COALESCE(short_name,''), COALESCE(city,''), COALESCE(address_full,''), COALESCE(phone_whatsapp,''), COALESCE(email,'')`,
		tenantID, req.OfficialName, req.Kind, nullable(req.ShortName), nullable(req.City), nullable(req.AddressFull), nullable(req.PhoneWA), nullable(req.Email)).
		Scan(&item.OfficialName, &item.Kind, &item.ShortName, &item.City, &item.AddressFull, &item.PhoneWA, &item.Email)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) ListTags(ctx context.Context, tenantID, scope, search string, page, limit int) ([]TagResponse, int64, error) {
	where := []string{"tenant_id = $1"}
	args := []interface{}{tenantID}
	argN := 2

	if scope != "" {
		where = append(where, fmt.Sprintf("scope = $%d", argN))
		args = append(args, scope)
		argN++
	}
	if search != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR slug ILIKE $%d)", argN, argN))
		args = append(args, "%"+search+"%")
		argN++
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM tags WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	args = append(args, limit, offset)
	sql := fmt.Sprintf(`SELECT id, scope::text, name, slug FROM tags WHERE %s ORDER BY id DESC LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []TagResponse
	for rows.Next() {
		var item TagResponse
		if err := rows.Scan(&item.ID, &item.Scope, &item.Name, &item.Slug); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (r *repository) CreateTag(ctx context.Context, tenantID string, req CreateTagRequest) (*TagResponse, error) {
	var item TagResponse
	err := r.db.QueryRow(ctx, `INSERT INTO tags (tenant_id, scope, name, slug)
		VALUES ($1,$2,$3,$4)
		RETURNING id, scope::text, name, slug`, tenantID, req.Scope, req.Name, slugify(req.Name)).
		Scan(&item.ID, &item.Scope, &item.Name, &item.Slug)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) UpdateTag(ctx context.Context, tenantID string, id int64, req UpdateTagRequest) error {
	tag, err := r.db.Exec(ctx, `UPDATE tags SET name=$1, slug=$2, updated_at=now() WHERE id=$3 AND tenant_id=$4`, req.Name, slugify(req.Name), id, tenantID)
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

func (r *repository) DeleteTag(ctx context.Context, tenantID string, id int64) error {
	var used bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(
		SELECT 1 FROM post_tags pt
		JOIN posts p ON p.id = pt.post_id
		WHERE pt.tag_id=$1 AND p.tenant_id=$2
	)`, id, tenantID).Scan(&used)
	if err != nil {
		return err
	}
	if used {
		return ErrTagInUse
	}
	tag, err := r.db.Exec(ctx, `DELETE FROM tags WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPosts(ctx context.Context, tenantID string, q PostListQuery) ([]PostResponse, int64, error) {
	where := []string{"tenant_id = $1"}
	args := []interface{}{tenantID}
	argN := 2
	if q.Category != "" {
		where = append(where, fmt.Sprintf("category = $%d", argN))
		args = append(args, q.Category)
		argN++
	}
	if q.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argN))
		args = append(args, q.Status)
		argN++
	}
	if q.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR excerpt ILIKE $%d)", argN, argN))
		args = append(args, "%"+q.Search+"%")
		argN++
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM posts WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	sortBy := "created_at"
	if q.SortBy == "published_at" || q.SortBy == "sort_order" || q.SortBy == "title" {
		sortBy = q.SortBy
	}
	sortOrder := "DESC"
	if strings.ToUpper(q.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)
	sql := fmt.Sprintf(`SELECT id, title, slug, category::text, COALESCE(excerpt,''), content_markdown, COALESCE(thumbnail_url,''), COALESCE(author_name,''), published_at, expired_at, status::text, show_on_homepage, sort_order
		FROM posts WHERE %s ORDER BY %s %s LIMIT $%d OFFSET $%d`, whereSQL, sortBy, sortOrder, argN, argN+1)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []PostResponse
	for rows.Next() {
		var item PostResponse
		if err := rows.Scan(&item.ID, &item.Title, &item.Slug, &item.Category, &item.Excerpt, &item.ContentMarkdown, &item.ThumbnailURL, &item.AuthorName, &item.PublishedAt, &item.ExpiredAt, &item.Status, &item.ShowOnHomepage, &item.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (r *repository) CreatePost(ctx context.Context, tenantID string, req PostPayload) (*PostResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var id int64
	var slug string
	err = tx.QueryRow(ctx, `INSERT INTO posts (tenant_id, title, slug, category, excerpt, content_markdown, thumbnail_url, author_name, published_at, expired_at, status, show_on_homepage, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, slug`,
		tenantID, req.Title, slugify(req.Title), req.Category, nullable(req.Excerpt), req.ContentMarkdown, nullable(req.ThumbnailURL), nullable(req.AuthorName), req.PublishedAt, req.ExpiredAt, req.Status, req.ShowOnHomepage, req.SortOrder).
		Scan(&id, &slug)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, err
	}

	// DI DALAM FUNC CreatePost
	for _, tagID := range req.TagIDs {
		// PERBAIKAN: Gunakan SELECT untuk memastikan tag_id tersebut milik tenant_id ini
		if _, err = tx.Exec(ctx, `
            INSERT INTO post_tags (post_id, tag_id) 
            SELECT $1, id FROM tags WHERE id=$2 AND tenant_id=$3 
            ON CONFLICT DO NOTHING`, id, tagID, tenantID); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetPost(ctx, tenantID, id)
}

func (r *repository) GetPost(ctx context.Context, tenantID string, id int64) (*PostResponse, error) {
	var item PostResponse
	err := r.db.QueryRow(ctx, `SELECT id, title, slug, category::text, COALESCE(excerpt,''), content_markdown, COALESCE(thumbnail_url,''), COALESCE(author_name,''), published_at, expired_at, status::text, show_on_homepage, sort_order
		FROM posts WHERE id=$1 AND tenant_id=$2`, id, tenantID).
		Scan(&item.ID, &item.Title, &item.Slug, &item.Category, &item.Excerpt, &item.ContentMarkdown, &item.ThumbnailURL, &item.AuthorName, &item.PublishedAt, &item.ExpiredAt, &item.Status, &item.ShowOnHomepage, &item.SortOrder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT t.id, t.scope::text, t.name, t.slug
		FROM post_tags pt JOIN tags t ON t.id=pt.tag_id WHERE pt.post_id=$1 ORDER BY t.id ASC`, item.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag TagResponse
		if err := rows.Scan(&tag.ID, &tag.Scope, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}
		item.Tags = append(item.Tags, tag)
	}
	return &item, nil
}

func (r *repository) UpdatePost(ctx context.Context, tenantID string, id int64, req PostPayload) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `UPDATE posts SET
		title=$1, slug=$2, category=$3, excerpt=$4, content_markdown=$5, thumbnail_url=$6, author_name=$7,
		published_at=$8, expired_at=$9, status=$10, show_on_homepage=$11, sort_order=$12, updated_at=now()
		WHERE id=$13 AND tenant_id=$14`,
		req.Title, slugify(req.Title), req.Category, nullable(req.Excerpt), req.ContentMarkdown, nullable(req.ThumbnailURL), nullable(req.AuthorName),
		req.PublishedAt, req.ExpiredAt, req.Status, req.ShowOnHomepage, req.SortOrder, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	if _, err = tx.Exec(ctx, `DELETE FROM post_tags WHERE post_id=$1`, id); err != nil {
		return err
	}
	for _, tagID := range req.TagIDs {
		// PERBAIKAN: Keamanan yang sama diterapkan saat update
		if _, err = tx.Exec(ctx, `
            INSERT INTO post_tags (post_id, tag_id) 
            SELECT $1, id FROM tags WHERE id=$2 AND tenant_id=$3 
            ON CONFLICT DO NOTHING`, id, tagID, tenantID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) UpdatePostStatus(ctx context.Context, tenantID string, id int64, req UpdatePostStatusRequest) error {
	publishedAt := req.PublishedAt
	if req.Status == "published" && publishedAt == nil {
		now := time.Now().UTC()
		publishedAt = &now
	}
	tag, err := r.db.Exec(ctx, `UPDATE posts SET status=$1, published_at=COALESCE($2,published_at), updated_at=now() WHERE id=$3 AND tenant_id=$4`,
		req.Status, publishedAt, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeletePost(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM posts WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListStaticPages(ctx context.Context, tenantID string) ([]PostResponse, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, slug, category::text, COALESCE(excerpt,''), content_markdown, COALESCE(thumbnail_url,''), COALESCE(author_name,''), published_at, expired_at, status::text, show_on_homepage, sort_order
		FROM posts WHERE tenant_id=$1 AND category='static_page' ORDER BY sort_order ASC, id ASC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PostResponse
	for rows.Next() {
		var item PostResponse
		if err := rows.Scan(&item.ID, &item.Title, &item.Slug, &item.Category, &item.Excerpt, &item.ContentMarkdown, &item.ThumbnailURL, &item.AuthorName, &item.PublishedAt, &item.ExpiredAt, &item.Status, &item.ShowOnHomepage, &item.SortOrder); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *repository) UpsertStaticPageBySlug(ctx context.Context, tenantID, slug string, req StaticPagePayload) (*PostResponse, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO posts (tenant_id, title, slug, category, excerpt, content_markdown, status, show_on_homepage, sort_order)
		VALUES ($1,$2,$3,'static_page',$4,$5,'published',false,0)
		ON CONFLICT (tenant_id, slug) DO UPDATE SET
			title=EXCLUDED.title,
			excerpt=EXCLUDED.excerpt,
			content_markdown=EXCLUDED.content_markdown,
			updated_at=now()
		RETURNING id`, tenantID, req.Title, slug, nullable(req.Excerpt), req.ContentMarkdown).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetPost(ctx, tenantID, id)
}

func (r *repository) UpdateTenantSetup(ctx context.Context, tenantID, name, subdomain string) error {
	// Kita ubah namanya, subdomainnya, dan ubah statusnya dari 'pending' menjadi 'active'
	tag, err := r.db.Exec(ctx, `UPDATE tenants SET name=$1, subdomain=$2, status='active', updated_at=now() WHERE id=$3`, name, subdomain, tenantID)
	if err != nil {
		// Tangkap error jika subdomain sudah dipakai orang lain
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

// ==========================================
// UTILITY FUNCTIONS
// ==========================================

// nullable mengubah string kosong menjadi nil agar masuk ke DB sebagai NULL (bukan string kosong "")
func nullable(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// slugify membuat URL-friendly string (Contoh: "Kajian Rutin" -> "kajian-rutin")
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	// Opsional: Jika ingin lebih ketat, kamu bisa tambahkan regex untuk menghapus karakter spesial di sini nanti
	return s
}
