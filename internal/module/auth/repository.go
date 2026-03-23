package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateTenantAndUser(ctx context.Context, email, googleID string) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, tenant_id, email, password_hash, google_id, role, created_at, updated_at FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err // Akan mengembalikan pgx.ErrNoRows jika tidak ditemukan
	}

	return &user, nil
}

// Implementasi fungsi Transaction
func (r *repository) CreateTenantAndUser(ctx context.Context, email, googleID string) (*User, error) {
	// 1. Mulai Transaksi Database
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Pastikan rollback otomatis dipanggil jika terjadi panic/error sebelum Commit
	defer tx.Rollback(ctx)

	var tenantID string
	// Karena subdomain harus UNIQUE, kita buat subdomain sementara (klien bisa menggantinya nanti di dashboard)
	tempSubdomain := fmt.Sprintf("tenant-%d", time.Now().Unix())

	// 2. Insert ke tabel tenants
	tenantQuery := `INSERT INTO tenants (name, subdomain, status) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx, tenantQuery, "Toko Baru", tempSubdomain, "pending").Scan(&tenantID)
	if err != nil {
		return nil, err
	}

	// 3. Insert ke tabel users, sambungkan dengan tenantID
	var user User
	userQuery := `
		INSERT INTO users (tenant_id, email, google_id, role) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, tenant_id, email, google_id, role, created_at, updated_at`

	err = tx.QueryRow(ctx, userQuery, tenantID, email, googleID, "admin").Scan(
		&user.ID, &user.TenantID, &user.Email, &user.GoogleID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
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
