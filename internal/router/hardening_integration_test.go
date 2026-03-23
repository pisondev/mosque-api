package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pisondev/mosque-api/internal/testutil"
	"github.com/sirupsen/logrus"
)

func TestHardeningNegativeIntegration(t *testing.T) {
	db := testutil.OpenTestDB(t)
	defer db.Close()

	tenant := testutil.SeedTenant(t, db)
	defer testutil.CleanupTenant(t, db, tenant.ID)

	t.Setenv("JWT_SECRET", "hardening-secret")
	app := fiber.New()
	SetupRoutes(app, db, logrus.New())
	token := testutil.IssueTestJWT(t, "hardening-secret", tenant.ID)

	status, _ := doJSONRequest(t, app, http.MethodGet, "/api/v1/tenant/me", "", nil)
	if status != http.StatusUnauthorized {
		t.Fatalf("missing auth expected 401 got %d", status)
	}

	status, _ = doJSONRequest(t, app, http.MethodGet, "/api/v1/tenant/me", "invalid-token", nil)
	if status != http.StatusUnauthorized {
		t.Fatalf("invalid token expected 401 got %d", status)
	}

	status, body := doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/prayer-times-daily", token, map[string]interface{}{})
	if status != http.StatusBadRequest {
		t.Fatalf("validation expected 400 got %d body=%v", status, body)
	}

	dailyPayload := map[string]interface{}{
		"day_date":     "2040-01-01",
		"subuh_time":   "04:35:00",
		"dzuhur_time":  "11:58:00",
		"ashar_time":   "15:21:00",
		"maghrib_time": "18:02:00",
		"isya_time":    "19:12:00",
	}
	status, _ = doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/prayer-times-daily", token, dailyPayload)
	if status != http.StatusCreated {
		t.Fatalf("first create prayer_times_daily expected 201 got %d", status)
	}
	status, _ = doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/prayer-times-daily", token, dailyPayload)
	if status != http.StatusConflict {
		t.Fatalf("duplicate prayer_times_daily expected 409 got %d", status)
	}

	specialPayload := map[string]interface{}{
		"kind":     "other",
		"title":    "Hari Konflik",
		"day_date": "2040-01-02",
	}
	status, _ = doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/special-days", token, specialPayload)
	if status != http.StatusCreated {
		t.Fatalf("first create special_day expected 201 got %d", status)
	}
	status, _ = doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/special-days", token, specialPayload)
	if status != http.StatusConflict {
		t.Fatalf("duplicate special_day expected 409 got %d", status)
	}

	domainPayload := map[string]interface{}{
		"domain_type": "custom_domain",
		"hostname":    "conflict-hardening.local",
	}
	status, _ = doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/domains", token, domainPayload)
	if status != http.StatusCreated {
		t.Fatalf("first create domain expected 201 got %d", status)
	}
	status, _ = doJSONRequest(t, app, http.MethodPost, "/api/v1/tenant/domains", token, domainPayload)
	if status != http.StatusConflict {
		t.Fatalf("duplicate domain expected 409 got %d", status)
	}
}

func TestHardeningListPerformanceBaseline(t *testing.T) {
	db := testutil.OpenTestDB(t)
	defer db.Close()

	tenant := testutil.SeedTenant(t, db)
	defer testutil.CleanupTenant(t, db, tenant.ID)

	t.Setenv("JWT_SECRET", "hardening-secret")
	app := fiber.New()
	SetupRoutes(app, db, logrus.New())
	token := testutil.IssueTestJWT(t, "hardening-secret", tenant.ID)

	if err := seedListBaselineData(context.Background(), db, tenant.ID); err != nil {
		t.Fatalf("seed baseline data failed: %v", err)
	}

	type endpointCase struct {
		name      string
		path      string
		authToken string
		max       time.Duration
	}

	cases := []endpointCase{
		{name: "tenant domains", path: "/api/v1/tenant/domains?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant posts", path: "/api/v1/tenant/posts?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant prayer times", path: "/api/v1/tenant/prayer-times-daily?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant events", path: "/api/v1/tenant/events?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant gallery items", path: "/api/v1/tenant/gallery/items?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant donation channels", path: "/api/v1/tenant/donation-channels?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant social links", path: "/api/v1/tenant/social-links?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant external links", path: "/api/v1/tenant/external-links?page=1&limit=50", authToken: token, max: 2 * time.Second},
		{name: "tenant feature catalog", path: "/api/v1/tenant/feature-catalog", authToken: token, max: 2 * time.Second},
		{name: "tenant website features", path: "/api/v1/tenant/website-features", authToken: token, max: 2 * time.Second},
		{name: "public events", path: fmt.Sprintf("/api/v1/public/%s/events?page=1&limit=50", tenant.Hostname), max: 2 * time.Second},
		{name: "public gallery items", path: fmt.Sprintf("/api/v1/public/%s/gallery/items?page=1&limit=50", tenant.Hostname), max: 2 * time.Second},
		{name: "public donation channels", path: fmt.Sprintf("/api/v1/public/%s/donation-channels?page=1&limit=50", tenant.Hostname), max: 2 * time.Second},
		{name: "public external links", path: fmt.Sprintf("/api/v1/public/%s/external-links?page=1&limit=50", tenant.Hostname), max: 2 * time.Second},
	}

	for _, tc := range cases {
		start := time.Now()
		status, body := doJSONRequest(t, app, http.MethodGet, tc.path, tc.authToken, nil)
		elapsed := time.Since(start)
		if status != http.StatusOK {
			t.Fatalf("%s expected 200 got %d body=%v", tc.name, status, body)
		}
		if body["status"] != "success" {
			t.Fatalf("%s expected success status body=%v", tc.name, body)
		}
		if elapsed > tc.max {
			t.Fatalf("%s baseline performance exceeded: %s > %s", tc.name, elapsed, tc.max)
		}
	}
}

func seedListBaselineData(ctx context.Context, db *pgxpool.Pool, tenantID string) error {
	for i := 0; i < 40; i++ {
		hostname := fmt.Sprintf("d-%d-%d.local", time.Now().UnixNano(), i)
		if _, err := db.Exec(ctx, `INSERT INTO website_domains (tenant_id,domain_type,hostname,status,verified_at) VALUES ($1,'custom_domain',$2,'active',now())`, tenantID, hostname); err != nil {
			return err
		}
		day := fmt.Sprintf("2041-01-%02d", (i%28)+1)
		if _, err := db.Exec(ctx, `INSERT INTO prayer_times_daily (tenant_id,day_date,subuh_time,dzuhur_time,ashar_time,maghrib_time,isya_time,source_label)
			VALUES ($1,$2::date,'04:30:00','12:00:00','15:00:00','18:00:00','19:00:00','baseline')
			ON CONFLICT (tenant_id,day_date) DO NOTHING`, tenantID, day); err != nil {
			return err
		}
		if _, err := db.Exec(ctx, `INSERT INTO events (tenant_id,title,category,start_date,time_mode,start_time,status)
			VALUES ($1,$2,'kajian_rutin',$3::date,'exact_time','07:00:00','upcoming')`, tenantID, fmt.Sprintf("Event %d", i), day); err != nil {
			return err
		}
		if _, err := db.Exec(ctx, `INSERT INTO posts (tenant_id,title,slug,category,content_markdown,status)
			VALUES ($1,$2,$3,'announcement','konten','published')`, tenantID, fmt.Sprintf("Post %d", i), fmt.Sprintf("post-%d-%d", time.Now().UnixNano(), i)); err != nil {
			return err
		}
		if _, err := db.Exec(ctx, `INSERT INTO donation_channels (tenant_id,channel_type,label,bank_name,account_number,account_holder_name,sort_order,is_public)
			VALUES ($1,'bank_account',$2,'BSI','12345','DKM',$3,true)`, tenantID, fmt.Sprintf("Donasi %d", i), i); err != nil {
			return err
		}
		if _, err := db.Exec(ctx, `INSERT INTO social_links (tenant_id,platform,url,sort_order)
			VALUES ($1,'instagram',$2,$3)`, tenantID, fmt.Sprintf("https://instagram.com/%d", i), i); err != nil {
			return err
		}
		if _, err := db.Exec(ctx, `INSERT INTO external_links (tenant_id,link_type,label,url,visibility,sort_order)
			VALUES ($1,'registration',$2,$3,'public',$4)`, tenantID, fmt.Sprintf("Link %d", i), fmt.Sprintf("https://example.com/%d", i), i); err != nil {
			return err
		}
	}

	var albumID int64
	if err := db.QueryRow(ctx, `INSERT INTO gallery_albums (tenant_id,title,media_kind) VALUES ($1,'Album Baseline','photo') RETURNING id`, tenantID).Scan(&albumID); err != nil {
		return err
	}
	for i := 0; i < 40; i++ {
		if _, err := db.Exec(ctx, `INSERT INTO gallery_items (tenant_id,album_id,media_type,media_url,sort_order)
			VALUES ($1,$2,'image',$3,$4)`, tenantID, albumID, fmt.Sprintf("https://cdn.example.com/%d.jpg", i), i); err != nil {
			return err
		}
	}

	if _, err := db.Exec(ctx, `INSERT INTO management_members (tenant_id,full_name,role_title,show_public,sort_order)
		VALUES ($1,'Ketua Baseline','Ketua',true,1)`, tenantID); err != nil {
		return err
	}

	_, err := db.Exec(ctx, `INSERT INTO website_features (tenant_id,feature_id,enabled,is_active,detail,note)
		SELECT $1,id,true,true,'baseline','baseline'
		FROM feature_catalog
		ON CONFLICT (tenant_id,feature_id) DO NOTHING`, tenantID)
	if err != nil {
		return err
	}
	return nil
}

func doJSONRequest(t *testing.T, app *fiber.App, method, path, token string, payload interface{}) (int, map[string]interface{}) {
	t.Helper()
	var bodyReader io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal payload failed: %v", err)
		}
		bodyReader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	var body map[string]interface{}
	_ = json.Unmarshal(rawBody, &body)
	return resp.StatusCode, body
}
