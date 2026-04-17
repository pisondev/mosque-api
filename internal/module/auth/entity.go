package auth

import "time"

type User struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	Email         string     `json:"email"`
	PasswordHash  *string    `json:"-"`
	PasswordSetAt *time.Time `json:"password_set_at"`
	GoogleID      *string    `json:"google_id"`
	Role          string     `json:"role"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateUserParams struct {
	Email        string
	PasswordHash *string
	GoogleID     *string
	TenantName   string
}

type PasswordResetToken struct {
	UserID    string
	TenantID  string
	Email     string
	Role      string
	ExpiresAt time.Time
}
