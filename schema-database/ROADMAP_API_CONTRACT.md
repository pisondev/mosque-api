# Roadmap Delivery & API Contract

Dokumen ini menjadi baseline implementasi panel dashboard management dan backend sampai level produksi matang (bukan hanya MVP), berdasarkan:
- `db/migrations/000001_init_core_schema.up.sql`
- `schema_database/muslim_serv_schema.sql` (sudah diturunkan ke migrasi `000002`)

## 1) Tahapan / Sprint Plan (hingga matang)

### Sprint 0 — Foundation Delivery
- Finalisasi migrasi `000001` + `000002` dan seed data minimum.
- Standarisasi struktur modul (controller, service, repository, dto) untuk domain utama.
- Standarisasi response envelope API (`status`, `message`, `data`, `meta`).
- Setup observability dasar: request-id, structured log, panic recovery, health endpoint.
- Output: fondasi backend stabil untuk pengembangan paralel frontend dashboard.

### Sprint 1 — Tenant, Auth, Domain Mapping, Profile
- Selesaikan auth + tenant context untuk semua endpoint dashboard.
- Modul domain mapping (`website_domains`) untuk subdomain/custom domain.
- Modul profil masjid (`masjid_profiles`) + validasi bisnis inti.
- Endpoint `tenant/me` final (bukan endpoint uji coba).
- Output: dashboard bisa login, baca tenant aktif, kelola identitas website.

### Sprint 2 — Konten Informasi (Posts, Tags) + Static Pages
- Modul `posts`, `tags`, `post_tags` end-to-end.
- Workflow draft → publish → archive.
- Static page memakai `posts.category=static_page`.
- Search, pagination, sorting untuk dashboard content.
- Output: CMS dashboard siap dipakai operasional harian.

### Sprint 3 — Jadwal Ibadah & Penugasan
- Modul `prayer_time_settings`, `prayer_times_daily`, `prayer_duties`, `special_days`.
- Job sinkronisasi jadwal harian + endpoint kalender ibadah.
- Validasi domain rule (`fardhu` wajib `prayer`, mode lokasi city/coordinates).
- Output: dashboard operasional ibadah berjalan penuh.

### Sprint 4 — Event, Galeri, Pengurus
- Modul `events`, `gallery_albums`, `gallery_items`, `management_members`.
- Upload media via object-storage signed URL flow.
- Modul publik feed untuk event dan galeri.
- Output: website publik aktif dengan konten kegiatan lengkap.

### Sprint 5 — Donasi, Link Sosial, Fitur Website
- Modul `donation_channels`, `social_links`, `external_links`, `feature_catalog`, `website_features`.
- Pengaturan visibilitas kanal publik.
- Penguatan kontrol role admin/superadmin per fitur sensitif.
- Output: kanal engagement dan monetisasi siap produksi.

### Sprint 6 — Hardening Produksi
- Audit security (authz, rate limit, payload validation, brute-force protection).
- Audit performa (index review, query tuning, endpoint SLA).
- Audit reliability (idempotency, retry policy, backup/restore drill).
- Contract test + integration test + smoke test CI.
- Output: release candidate siap go-live multi-tenant.

## 2) API Contract — Sprint 1 (Tenant/Auth/Domain/Profile)

Base path: `/api/v1`
Format response sukses:

```json
{
  "status": "success",
  "message": "ok",
  "data": {},
  "meta": {}
}
```

Format response error:

```json
{
  "status": "error",
  "message": "validation failed",
  "errors": [
    { "field": "hostname", "message": "invalid hostname" }
  ]
}
```

### 2.1 Auth

**POST** `/auth/google`
- Auth: Public
- Body:

```json
{
  "token": "google_id_token"
}
```

- 200 data:

```json
{
  "status": "success",
  "message": "login success",
  "data": {
    "access_token": "jwt_token",
    "email": "admin@contoh.id",
    "role": "admin"
  }
}
```

- 401 jika token Google tidak valid.

### 2.2 Tenant Context

**GET** `/tenant/me`
- Auth: Bearer JWT
- 200 data:

```json
{
  "status": "success",
  "message": "tenant context loaded",
  "data": {
    "tenant_id": "uuid",
    "email": "admin@contoh.id",
    "role": "admin"
  }
}
```

### 2.3 Domain Mapping (`website_domains`)

**GET** `/tenant/domains`
- Auth: Bearer JWT
- Query: `status`, `domain_type`, `page`, `limit`
- 200 data: daftar domain milik tenant.

**POST** `/tenant/domains`
- Auth: Bearer JWT
- Body:

```json
{
  "domain_type": "custom_domain",
  "hostname": "masjidcontoh.org"
}
```

- 201 data:

```json
{
  "id": 10,
  "domain_type": "custom_domain",
  "hostname": "masjidcontoh.org",
  "status": "pending",
  "verified_at": null
}
```

**PATCH** `/tenant/domains/{id}`
- Auth: Bearer JWT
- Body (opsional):

```json
{
  "status": "verifying"
}
```

**DELETE** `/tenant/domains/{id}`
- Auth: Bearer JWT
- 204 no content.

### 2.4 Profil Masjid (`masjid_profiles`)

**GET** `/tenant/profile`
- Auth: Bearer JWT
- 200 data: single object profil tenant.

**PUT** `/tenant/profile`
- Auth: Bearer JWT
- Body minimal:

```json
{
  "official_name": "Masjid Al Ikhlas",
  "kind": "masjid",
  "city": "Bandung",
  "address_full": "Jl. Contoh No. 1"
}
```

- Rule validasi utama:
  - `official_name` wajib.
  - `kind` salah satu enum `place_kind`.
  - `established_year` di rentang 1000..3000 bila diisi.
  - `capacity_male/female/floors_count` tidak boleh negatif.

## 3) API Contract — Sprint 2 (Posts/Tags/Static Page)

### 3.1 Tags (`tags`)

**GET** `/tenant/tags`
- Auth: Bearer JWT
- Query: `scope`, `search`, `page`, `limit`

**POST** `/tenant/tags`
- Auth: Bearer JWT
- Body:

```json
{
  "scope": "post",
  "name": "Ramadhan 1448"
}
```

- Slug dibangkitkan server-side, unik per `(tenant_id, scope, slug)`.

**PATCH** `/tenant/tags/{id}`
- Auth: Bearer JWT
- Body: `name`

**DELETE** `/tenant/tags/{id}`
- Auth: Bearer JWT
- 409 bila masih dipakai relasi `post_tags`.

### 3.2 Posts (`posts`)

**GET** `/tenant/posts`
- Auth: Bearer JWT
- Query:
  - `category`: `news_activity|announcement|reflection|static_page`
  - `status`: `draft|published|archived`
  - `search`, `page`, `limit`, `sort_by`, `sort_order`

**POST** `/tenant/posts`
- Auth: Bearer JWT
- Body:

```json
{
  "title": "Agenda Kajian Pekanan",
  "category": "news_activity",
  "excerpt": "Kajian rutin ba'da maghrib",
  "content_markdown": "## Materi\nDetail materi...",
  "thumbnail_url": "https://cdn.example.com/poster.jpg",
  "author_name": "Sekretariat DKM",
  "status": "draft",
  "show_on_homepage": true,
  "sort_order": 1,
  "tag_ids": [1, 3]
}
```

**GET** `/tenant/posts/{id}`
- Auth: Bearer JWT
- Data termasuk daftar tag.

**PUT** `/tenant/posts/{id}`
- Auth: Bearer JWT
- Body: field yang sama dengan create.

**PATCH** `/tenant/posts/{id}/status`
- Auth: Bearer JWT
- Body:

```json
{
  "status": "published",
  "published_at": "2026-03-23T08:00:00Z"
}
```

- Rule:
  - Jika status `published` dan `published_at` kosong, backend isi waktu sekarang.
  - `expired_at` tidak boleh lebih kecil dari `published_at`.

**DELETE** `/tenant/posts/{id}`
- Auth: Bearer JWT
- Soft delete tidak tersedia di schema saat ini, sehingga delete permanen.

### 3.3 Static Pages (di atas `posts`)

**GET** `/tenant/static-pages`
- Auth: Bearer JWT
- Implementasi: filter `posts.category=static_page`.

**PUT** `/tenant/static-pages/{slug}`
- Auth: Bearer JWT
- Upsert halaman statis per slug (contoh: `visi-misi`, `sejarah`).

## 4) API Contract — Sprint 3 (Jadwal Ibadah & Penugasan)

### 4.1 Prayer Time Settings (`prayer_time_settings`)

**GET** `/tenant/prayer-time-settings`
- Auth: Bearer JWT
- 200 data: konfigurasi tunggal tenant.

**PUT** `/tenant/prayer-time-settings`
- Auth: Bearer JWT
- Body:

```json
{
  "timezone": "Asia/Jakarta",
  "location_mode": "city",
  "city_name": "Bandung",
  "latitude": null,
  "longitude": null,
  "calc_method": "kemenag",
  "asr_madhhab": "syafii",
  "adj_subuh_min": 2,
  "adj_dzuhur_min": 0,
  "adj_ashar_min": 1,
  "adj_maghrib_min": 0,
  "adj_isya_min": 0
}
```

- Rule:
  - `location_mode=city` → `city_name` wajib.
  - `location_mode=coordinates` → `latitude` & `longitude` wajib.
  - Nilai offset menit boleh negatif, tetapi dibatasi rentang bisnis (misal `-60..60`).

### 4.2 Prayer Times Daily (`prayer_times_daily`)

**GET** `/tenant/prayer-times-daily`
- Auth: Bearer JWT
- Query: `from`, `to`, `page`, `limit`, `sort_by`, `sort_order`
- Default sort: `day_date asc`.

**POST** `/tenant/prayer-times-daily`
- Auth: Bearer JWT
- Body:

```json
{
  "day_date": "2026-03-23",
  "subuh_time": "04:35:00",
  "dzuhur_time": "11:58:00",
  "ashar_time": "15:21:00",
  "maghrib_time": "18:02:00",
  "isya_time": "19:12:00",
  "sunrise_time": "05:48:00",
  "dhuha_time": "06:15:00",
  "source_label": "api.kemenag",
  "fetched_at": "2026-03-22T23:10:00Z"
}
```

**GET** `/tenant/prayer-times-daily/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/prayer-times-daily/{id}`
- Auth: Bearer JWT
- Body: sama dengan create.

**DELETE** `/tenant/prayer-times-daily/{id}`
- Auth: Bearer JWT

- Rule:
  - Unik per `(tenant_id, day_date)`, conflict → 409.
  - Semua waktu memakai timezone tenant.

### 4.3 Prayer Duties (`prayer_duties`)

**GET** `/tenant/prayer-duties`
- Auth: Bearer JWT
- Query: `from`, `to`, `category`, `prayer`, `page`, `limit`

**POST** `/tenant/prayer-duties`
- Auth: Bearer JWT
- Body:

```json
{
  "category": "fardhu",
  "duty_date": "2026-03-23",
  "prayer": "maghrib",
  "imam_name": "Ust. Ahmad",
  "muadzin_name": "Saepul",
  "khatib_name": null,
  "bilal_name": null,
  "note": "Petugas cadangan: Budi"
}
```

**GET** `/tenant/prayer-duties/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/prayer-duties/{id}`
- Auth: Bearer JWT

**DELETE** `/tenant/prayer-duties/{id}`
- Auth: Bearer JWT

- Rule:
  - Jika `category=fardhu`, field `prayer` wajib.
  - `duty_date` wajib format `YYYY-MM-DD`.

### 4.4 Special Days (`special_days`)

**GET** `/tenant/special-days`
- Auth: Bearer JWT
- Query: `year`, `kind`, `from`, `to`, `page`, `limit`

**POST** `/tenant/special-days`
- Auth: Bearer JWT
- Body:

```json
{
  "kind": "idul_fitri",
  "title": "Shalat Idul Fitri 1447 H",
  "day_date": "2026-03-21",
  "location_note": "Lapangan utama",
  "start_time": "06:30:00",
  "imam_name": "KH. Hasan",
  "khatib_name": "Dr. Fulan",
  "muadzin_name": "Rizki",
  "note": "Datang 30 menit lebih awal"
}
```

**GET** `/tenant/special-days/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/special-days/{id}`
- Auth: Bearer JWT

**DELETE** `/tenant/special-days/{id}`
- Auth: Bearer JWT

- Rule:
  - Unik per `(tenant_id, kind, day_date)`, conflict → 409.
  - `kind` harus enum: `idul_fitri|idul_adha|other`.

### 4.5 Calendar Aggregate

**GET** `/tenant/prayer-calendar`
- Auth: Bearer JWT
- Query: `from`, `to`
- Data: agregasi `prayer_times_daily` + `prayer_duties` + `special_days` untuk rendering kalender dashboard.

## 5) API Contract — Sprint 4 (Event, Galeri, Pengurus)

### 5.1 Events (`events`)

**GET** `/tenant/events`
- Auth: Bearer JWT
- Query: `status`, `category`, `from`, `to`, `search`, `page`, `limit`, `sort_by`, `sort_order`

**POST** `/tenant/events`
- Auth: Bearer JWT
- Body:

```json
{
  "title": "Kajian Ahad Pagi",
  "slug": "kajian-ahad-pagi",
  "category": "kajian",
  "description_markdown": "Materi tafsir surat pilihan",
  "status": "draft",
  "time_mode": "exact_time",
  "after_prayer": null,
  "start_date": "2026-04-05",
  "end_date": "2026-04-05",
  "start_time": "07:00:00",
  "end_time": "09:00:00",
  "location_text": "Aula utama",
  "speaker_name": "Ust. Fulan",
  "is_featured": true,
  "show_on_homepage": true,
  "capacity": 300,
  "fee_amount": 0
}
```

**GET** `/tenant/events/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/events/{id}`
- Auth: Bearer JWT

**PATCH** `/tenant/events/{id}/status`
- Auth: Bearer JWT
- Body: `{ "status": "published" }`

**DELETE** `/tenant/events/{id}`
- Auth: Bearer JWT

- Rule:
  - `end_date >= start_date`.
  - `time_mode=exact_time` wajib `start_time`.
  - `time_mode=after_prayer` wajib `after_prayer`.

### 5.2 Gallery Albums (`gallery_albums`)

**GET** `/tenant/gallery/albums`
- Auth: Bearer JWT
- Query: `status`, `search`, `page`, `limit`

**POST** `/tenant/gallery/albums`
- Auth: Bearer JWT
- Body:

```json
{
  "title": "Ramadhan 1447 H",
  "slug": "ramadhan-1447",
  "description": "Dokumentasi kegiatan Ramadhan",
  "cover_image_url": "https://cdn.example.com/albums/ramadhan-cover.jpg",
  "status": "published",
  "start_date": "2026-03-01",
  "end_date": "2026-03-30",
  "show_on_homepage": true
}
```

**GET** `/tenant/gallery/albums/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/gallery/albums/{id}`
- Auth: Bearer JWT

**DELETE** `/tenant/gallery/albums/{id}`
- Auth: Bearer JWT

### 5.3 Gallery Items (`gallery_items`)

**GET** `/tenant/gallery/items`
- Auth: Bearer JWT
- Query: `album_id`, `media_type`, `search`, `page`, `limit`

**POST** `/tenant/gallery/items`
- Auth: Bearer JWT
- Body:

```json
{
  "album_id": 12,
  "title": "Buka Puasa Bersama",
  "slug": "buka-puasa-bersama-1",
  "description": "Kegiatan buka puasa jamaah",
  "media_type": "image",
  "media_url": "https://cdn.example.com/gallery/item-001.jpg",
  "thumbnail_url": "https://cdn.example.com/gallery/item-001-thumb.jpg",
  "taken_at": "2026-03-14T17:45:00Z",
  "status": "published",
  "show_on_homepage": true,
  "sort_order": 1,
  "tag_ids": [2, 6]
}
```

**GET** `/tenant/gallery/items/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/gallery/items/{id}`
- Auth: Bearer JWT

**DELETE** `/tenant/gallery/items/{id}`
- Auth: Bearer JWT

### 5.4 Management Members (`management_members`)

**GET** `/tenant/management-members`
- Auth: Bearer JWT
- Query: `is_active`, `search`, `page`, `limit`

**POST** `/tenant/management-members`
- Auth: Bearer JWT
- Body:

```json
{
  "name": "Ahmad Fauzi",
  "position": "Ketua DKM",
  "photo_url": "https://cdn.example.com/management/ketua.jpg",
  "bio": "Periode 2026-2030",
  "phone": "08123456789",
  "email": "ketua@masjidcontoh.id",
  "sort_order": 1,
  "is_active": true
}
```

**GET** `/tenant/management-members/{id}`
- Auth: Bearer JWT

**PUT** `/tenant/management-members/{id}`
- Auth: Bearer JWT

**DELETE** `/tenant/management-members/{id}`
- Auth: Bearer JWT

### 5.5 Public Feed (Website Publik)

**GET** `/public/{hostname}/events`
**GET** `/public/{hostname}/gallery/albums`
**GET** `/public/{hostname}/gallery/items`
**GET** `/public/{hostname}/management-members`
- Auth: Public
- Rule: hanya data `published/active` dan visibility publik.

## 6) API Contract — Sprint 5 (Donasi, Link Sosial, Fitur Website)

### 6.1 Donation Channels (`donation_channels`)

**GET** `/tenant/donation-channels`
- Auth: Bearer JWT

**POST** `/tenant/donation-channels`
- Auth: Bearer JWT
- Body:

```json
{
  "channel_type": "bank_account",
  "provider_name": "Bank Syariah Indonesia",
  "account_name": "DKM Al Ikhlas",
  "account_number": "1234567890",
  "qris_image_url": null,
  "note": "Transfer dengan kode unik 001",
  "is_active": true,
  "sort_order": 1
}
```

**PUT** `/tenant/donation-channels/{id}`
- Auth: Bearer JWT

**DELETE** `/tenant/donation-channels/{id}`
- Auth: Bearer JWT

- Rule:
  - `channel_type=bank_account` wajib `account_name` + `account_number`.
  - `channel_type=qris` wajib `qris_image_url`.

### 6.2 Social Links (`social_links`)

**GET** `/tenant/social-links`
**POST** `/tenant/social-links`
**PUT** `/tenant/social-links/{id}`
**DELETE** `/tenant/social-links/{id}`
- Auth: Bearer JWT
- Field utama: `platform`, `label`, `url`, `is_active`, `sort_order`.

### 6.3 External Links (`external_links`)

**GET** `/tenant/external-links`
**POST** `/tenant/external-links`
**PUT** `/tenant/external-links/{id}`
**DELETE** `/tenant/external-links/{id}`
- Auth: Bearer JWT
- Field utama: `title`, `url`, `link_type`, `visibility`, `is_active`, `sort_order`.

### 6.4 Feature Catalog & Website Features

**GET** `/tenant/feature-catalog`
- Auth: Bearer JWT
- Data: daftar master fitur global.

**GET** `/tenant/website-features`
- Auth: Bearer JWT
- Data: fitur tenant + status aktif.

**PUT** `/tenant/website-features/{feature_id}`
- Auth: Bearer JWT
- Body:

```json
{
  "is_enabled": true,
  "config_json": {
    "show_on_homepage": true
  }
}
```

**PATCH** `/tenant/website-features:bulk`
- Auth: Bearer JWT
- Body: array update beberapa fitur sekaligus.

### 6.5 Public Endpoint

**GET** `/public/{hostname}/donation-channels`
**GET** `/public/{hostname}/social-links`
**GET** `/public/{hostname}/external-links`
- Auth: Public
- Rule: filter hanya `is_active=true` dan `visibility=public`.

## 7) API Contract — Sprint 6 (Hardening Produksi)

### 7.1 Security
- Rate limit per IP + per token untuk endpoint auth dan endpoint tulis.
- Payload size limit per endpoint (khusus upload link/media lebih ketat).
- Validation whitelist untuk `sort_by`, `sort_order`, enum, dan field filter.
- Penguatan JWT: `issuer`, `audience`, expiry, dan rotasi secret berkala.

### 7.2 Reliability
- Idempotency key untuk endpoint create sensitif (donation channel, event, domain).
- Retry policy internal untuk operasi network eksternal (signed URL/object storage).
- Standarisasi error code: `400`, `401`, `403`, `404`, `409`, `422`, `500`.

### 7.3 Performance
- Review index semua query list/filter utama (`tenant_id`, `status`, `day_date`, `start_date`).
- Batasi default pagination (`limit` default 10, max 100).
- Monitoring slow query dan target SLA endpoint dashboard.

### 7.4 Testing & Release Gate
- Contract test otomatis untuk semua endpoint sprint.
- Integration test minimal happy-path + validation-path + authz-path.
- Smoke test pasca migrasi: auth, tenant context, CRUD sample lintas modul.

## 8) Definisi Teknis Dokumentasi API

- Sumber kontrak: dokumen ini.
- Sumber turunan implementasi:
  - OpenAPI `openapi.yaml`
  - Collection test API (Postman/Bruno)
  - Contract test otomatis di CI
- Aturan versioning:
  - Breaking change: buat versi endpoint baru (`/api/v2`)
  - Non-breaking change: append field opsional tanpa mengganti struktur existing

## 9) Eksekusi Migrasi (tanpa Docker, lintas OS)

Prinsip lingkungan lokal:
- Prioritas non-Docker jika PostgreSQL lokal tersedia dan stabil.
- Gunakan Docker hanya jika instalasi PostgreSQL lokal tidak tersedia.

Migrasi tersedia:
- `db/migrations/000001_init_core_schema.up.sql`
- `db/migrations/000002_muslim_serv_schema.up.sql`
- `db/migrations/000003_seed_master_data.up.sql`

Perintah Linux/macOS (pakai Makefile):

```bash
make migrate-up DB_URL='postgres://root:secretpassword@localhost:5432/mosque_saas?sslmode=disable'
```

Perintah Windows PowerShell:

```powershell
$env:DB_URL='postgres://root:secretpassword@localhost:5432/mosque_saas?sslmode=disable'
go run -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path db/migrations -database "$env:DB_URL" -verbose up
```

Rollback 1 step:

```powershell
go run -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path db/migrations -database "$env:DB_URL" -verbose down 1
```

## 10) Definition of Done per Sprint

- Semua endpoint sprint punya kontrak request/response final.
- Semua endpoint sprint punya integration test minimal happy-path + validation-path.
- Tidak ada query lintas tenant tanpa filter `tenant_id`.
- Dashboard page untuk modul sprint dapat create-read-update-delete sesuai kontrak.
