package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID string) (*User, error)
	CreateTenantAndUser(ctx context.Context, params CreateUserParams) (*User, error)
	SyncGoogleAccount(ctx context.Context, userID, googleID string, displayName, avatarURL *string) error
	UpdateAccountProfile(ctx context.Context, userID string, params UpdateAccountProfileParams) (*User, error)
	StorePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	ConsumePasswordResetToken(ctx context.Context, tokenHash, passwordHash string) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, tenant_id, email, display_name, avatar_url, password_hash, password_set_at, google_id, role, created_at, updated_at FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.DisplayName, &user.AvatarURL, &user.PasswordHash, &user.PasswordSetAt,
		&user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err // Akan mengembalikan pgx.ErrNoRows jika tidak ditemukan
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, userID string) (*User, error) {
	query := `SELECT id, tenant_id, email, display_name, avatar_url, password_hash, password_set_at, google_id, role, created_at, updated_at FROM users WHERE id = $1::uuid`

	var user User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.DisplayName, &user.AvatarURL, &user.PasswordHash, &user.PasswordSetAt,
		&user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateTenantAndUser(ctx context.Context, params CreateUserParams) (*User, error) {
	// 1. Mulai Transaksi Database
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Pastikan rollback otomatis dipanggil jika terjadi panic/error sebelum Commit
	defer tx.Rollback(ctx)

	var tenantID string
	// Karena subdomain harus UNIQUE, kita buat subdomain sementara (klien bisa menggantinya nanti di dashboard)
	tempSubdomain := fmt.Sprintf("tenant-%d", time.Now().UnixNano())
	tenantName := params.TenantName
	if tenantName == "" {
		tenantName = "Masjid Baru"
	}

	// 2. Insert ke tabel tenants
	tenantQuery := `INSERT INTO tenants (name, subdomain, status) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx, tenantQuery, tenantName, tempSubdomain, "pending").Scan(&tenantID)
	if err != nil {
		return nil, err
	}

	// 3. Insert ke tabel users, sambungkan dengan tenantID
	var user User
	userQuery := `
		INSERT INTO users (tenant_id, email, display_name, avatar_url, password_hash, password_set_at, google_id, role) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, tenant_id, email, display_name, avatar_url, password_hash, password_set_at, google_id, role, created_at, updated_at`

	var passwordSetAt *time.Time
	if params.PasswordHash != nil {
		now := time.Now().UTC()
		passwordSetAt = &now
	}

	err = tx.QueryRow(ctx, userQuery, tenantID, params.Email, params.DisplayName, params.AvatarURL, params.PasswordHash, passwordSetAt, params.GoogleID, "admin").Scan(
		&user.ID, &user.TenantID, &user.Email, &user.DisplayName, &user.AvatarURL, &user.PasswordHash, &user.PasswordSetAt, &user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 4. Jika semua lancar, kunci perubahan ke database (Commit)
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) SyncGoogleAccount(ctx context.Context, userID, googleID string, displayName, avatarURL *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET google_id = $2,
			display_name = CASE
				WHEN COALESCE(display_name, '') = '' AND COALESCE($3, '') <> '' THEN $3
				ELSE display_name
			END,
			avatar_url = CASE
				WHEN COALESCE(avatar_url, '') = '' AND COALESCE($4, '') <> '' THEN $4
				ELSE avatar_url
			END,
			updated_at = now()
		WHERE id = $1::uuid`, userID, googleID, displayName, avatarURL)
	return err
}

func (r *repository) UpdateAccountProfile(ctx context.Context, userID string, params UpdateAccountProfileParams) (*User, error) {
	query := `
		UPDATE users
		SET display_name = $2,
			avatar_url = $3,
			updated_at = now()
		WHERE id = $1::uuid
		RETURNING id, tenant_id, email, display_name, avatar_url, password_hash, password_set_at, google_id, role, created_at, updated_at`

	var user User
	err := r.db.QueryRow(ctx, query, userID, params.DisplayName, params.AvatarURL).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.DisplayName, &user.AvatarURL, &user.PasswordHash, &user.PasswordSetAt,
		&user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) StorePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM password_reset_tokens WHERE user_id = $1::uuid AND used_at IS NULL`, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1::uuid, $2, $3)`, userID, tokenHash, expiresAt); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) ConsumePasswordResetToken(ctx context.Context, tokenHash, passwordHash string) (*User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT u.id, u.tenant_id, u.email, u.display_name, u.avatar_url, u.password_hash, u.password_set_at, u.google_id, u.role, u.created_at, u.updated_at
		FROM password_reset_tokens prt
		JOIN users u ON u.id = prt.user_id
		WHERE prt.token_hash = $1 AND prt.used_at IS NULL AND prt.expires_at > now()
		FOR UPDATE OF prt, u`

	var user User
	if err := tx.QueryRow(ctx, query, tokenHash).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.DisplayName, &user.AvatarURL, &user.PasswordHash, &user.PasswordSetAt,
		&user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET password_hash = $2, password_set_at = now(), updated_at = now() WHERE id = $1::uuid`, user.ID, passwordHash); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE password_reset_tokens SET used_at = now() WHERE token_hash = $1`, tokenHash); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM password_reset_tokens WHERE user_id = $1::uuid AND token_hash <> $2`, user.ID, tokenHash); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	user.PasswordHash = &passwordHash
	now := time.Now().UTC()
	user.PasswordSetAt = &now
	return &user, nil
}
