package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestTenant struct {
	ID       string
	Hostname string
}

func OpenTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://root:secretpassword@localhost:5432/mosque_saas?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("skip integration test: db not available: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("skip integration test: db not reachable: %v", err)
	}
	return pool
}

func SeedTenant(t *testing.T, db *pgxpool.Pool) TestTenant {
	t.Helper()
	ctx := context.Background()
	suffix := time.Now().UnixNano()
	name := fmt.Sprintf("tenant-it-%d", suffix)
	subdomain := fmt.Sprintf("it-%d", suffix)
	hostname := fmt.Sprintf("%s.local", subdomain)
	var tenantID string
	err := db.QueryRow(ctx, `INSERT INTO tenants (name, subdomain, status) VALUES ($1,$2,'active') RETURNING id::text`, name, subdomain).Scan(&tenantID)
	if err != nil {
		t.Fatalf("seed tenant failed: %v", err)
	}
	_, err = db.Exec(ctx, `INSERT INTO website_domains (tenant_id, domain_type, hostname, status, verified_at) VALUES ($1,'custom_domain',$2,'active',now())`, tenantID, hostname)
	if err != nil {
		_, _ = db.Exec(ctx, `DELETE FROM tenants WHERE id=$1::uuid`, tenantID)
		t.Fatalf("seed website domain failed: %v", err)
	}
	return TestTenant{ID: tenantID, Hostname: hostname}
}

func CleanupTenant(t *testing.T, db *pgxpool.Pool, tenantID string) {
	t.Helper()
	_, _ = db.Exec(context.Background(), `DELETE FROM tenants WHERE id=$1::uuid`, tenantID)
}

func IssueTestJWT(t *testing.T, secret, tenantID string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id":   "00000000-0000-0000-0000-000000000001",
		"tenant_id": tenantID,
		"email":     "integration@test.local",
		"role":      "admin",
		"exp":       time.Now().Add(2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("issue jwt failed: %v", err)
	}
	return signed
}
