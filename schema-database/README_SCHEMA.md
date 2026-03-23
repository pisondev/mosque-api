# Muslim-Serv (Masjid Website Management) — PostgreSQL Schema Documentation

Dokumen ini menjelaskan struktur database (PostgreSQL) yang sudah disesuaikan dengan fondasi migrasi core:

- Core tenant/auth ada di `db/migrations/000001_init_core_schema.up.sql`
- Tabel core: `tenants` dan `users` (UUID)
- Semua fitur masjid tambahan berada di `schema_database/muslim_serv_schema.sql`
- Semua data spesifik website/masjid memakai `tenant_id UUID` → FK ke `tenants(id)`
- Penyimpanan media memakai URL object storage
- Konten artikel/halaman memakai Markdown

---

## 1. Prinsip Multi-Tenant

Semua data domain bisnis (profil, jadwal, event, konten, galeri, donasi, dll) disimpan per tenant menggunakan `tenant_id`.

Alur query:
- Resolve request dari hostname/subdomain melalui tabel `website_domains`
- Dapatkan `tenant_id`
- Semua query fitur wajib difilter berdasarkan `tenant_id`

---

## 2. Core Tables (Migrasi 000001)

### 2.1 `tenants`
**Tujuan:** representasi tenant SaaS (1 tenant = 1 website masjid).

Kolom inti:
- `id UUID PRIMARY KEY`
- `name`
- `subdomain` (unik)
- `status` (`pending/active/suspended`)
- `created_at`, `updated_at`

### 2.2 `users`
**Tujuan:** akun pengelola tenant.

Kolom inti:
- `id UUID PRIMARY KEY`
- `tenant_id UUID` FK → `tenants(id)`
- `email` (unik)
- `password_hash` (nullable jika OAuth)
- `google_id` (unik, nullable)
- `role` (`superadmin/admin`)
- `created_at`, `updated_at`

Catatan:
- Tabel subscription/payment belum disediakan di core.

---

## 3. Domain Mapping

### 3.1 `website_domains`
**Tujuan:** mapping domain/subdomain ke tenant.

Kolom penting:
- `tenant_id UUID` FK
- `domain_type` (`subdomain/custom_domain`)
- `hostname` (unik global)
- `status` (`pending/verifying/active/disabled`)
- `verified_at`

---

## 4. Profil Masjid

### 4.1 `masjid_profiles`
**Tujuan:** identitas dan informasi dasar masjid per tenant (`UNIQUE tenant_id`).

Kolom wajib minimum:
- `tenant_id`
- `official_name`
- `kind` (`masjid/musala/surau/langgar`)

Catatan:
- Sejarah/visi-misi disimpan sebagai `posts` kategori `static_page`.

---

## 5. Jadwal Sholat

### 5.1 `prayer_time_settings`
Konfigurasi jadwal sholat per tenant:
- `tenant_id` unik
- mode lokasi `city/coordinates`
- penyesuaian menit (`adj_*_min`)
- validasi `CHECK` untuk kondisi city/coordinates

### 5.2 `prayer_times_daily`
Cache jadwal harian:
- FK `tenant_id`
- unik `(tenant_id, day_date)`
- simpan waktu wajib + opsional sunrise/dhuha

---

## 6. Penugasan Petugas

### 6.1 `prayer_duties`
Input petugas sholat/jumat/tarawih/id per tenant.

Aturan:
- Jika `category='fardhu'`, maka `prayer` wajib diisi (`CHECK`)
- Semua nama petugas disimpan sebagai string

---

## 7. Hari Raya / Special Days

### 7.1 `special_days`
Konfigurasi hari khusus:
- FK `tenant_id`
- `kind` (`idul_fitri/idul_adha/other`)
- unik `(tenant_id, kind, day_date)`

---

## 8. Agenda / Kegiatan

### 8.1 `events`
Kalender kegiatan masjid per tenant.

Aturan penting:
- `time_mode='exact_time'` wajib `start_time`
- `time_mode='after_prayer'` wajib `after_prayer`
- `end_date >= start_date`

---

## 9. Pengurus / DKM

### 9.1 `management_members`
Daftar pengurus per tenant untuk tampilan publik/internal.

---

## 10. Fasilitas & Layanan

### 10.1 `feature_catalog`
Master global fasilitas/layanan.

### 10.2 `website_features`
Relasi fitur ke tenant:
- FK `tenant_id`
- FK `feature_id`
- unik `(tenant_id, feature_id)`

---

## 11. Tags

### 11.1 `tags`
Tag per tenant dengan scope:
- `post`, `gallery`, `event`
- unik `(tenant_id, scope, slug)`

### 11.2 Relasi tag
- `post_tags` (post ↔ tag)
- `gallery_item_tags` (gallery item ↔ tag)

---

## 12. Konten Informasi

### 12.1 `posts`
Konten berita/pengumuman/refleksi/halaman statis:
- FK `tenant_id`
- `content_markdown`
- unik `(tenant_id, slug)`
- validasi `expired_at >= published_at` jika keduanya terisi

---

## 13. Galeri

### 13.1 `gallery_albums`
Album galeri per tenant.

### 13.2 `gallery_items`
Item foto/video per tenant, opsional terhubung album.

---

## 14. Donasi

### 14.1 `donation_channels`
Kanal donasi bank/QRIS per tenant.

Aturan:
- `bank_account` wajib `bank_name`, `account_number`, `account_holder_name`
- `qris` wajib `qris_image_url`

---

## 15. Sosial Media & Link Lain

### 15.1 `social_links`
Link sosial media per tenant.

### 15.2 `external_links`
Link eksternal tambahan per tenant.

---

## 16. Data Minimal Onboarding

Urutan minimal agar website bisa tayang:
1) Buat `tenants`
2) Buat `users` untuk tenant terkait
3) Buat `website_domains` untuk hostname akses publik
4) Buat `masjid_profiles` (minimal `official_name` + `kind`)
5) Opsional: isi `prayer_time_settings`
6) Opsional: isi `posts` kategori `static_page`

---

## 17. Catatan Teknis

- Semua timestamp menggunakan `TIMESTAMPTZ`.
- Banyak validasi domain dilakukan di level aplikasi.
- Reserved words subdomain tetap divalidasi di aplikasi.
