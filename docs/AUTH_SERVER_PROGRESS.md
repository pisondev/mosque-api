# Auth Server Progress

## Ringkasan

Implementasi autentikasi email/password dan reset password telah ditambahkan ke `mosque-api` pada branch `feat/auth-server`.

## File Ditambahkan

- `db/migrations/000009_add_password_reset_auth.up.sql`
- `db/migrations/000009_add_password_reset_auth.down.sql`
- `internal/module/auth/mailer.go`
- `internal/module/auth/auth_integration_test.go`

## File Diupdate

- `internal/module/auth/controller.go`
- `internal/module/auth/dto.go`
- `internal/module/auth/entity.go`
- `internal/module/auth/repository.go`
- `internal/module/auth/service.go`
- `internal/router/router.go`
- `internal/router/hardening_integration_test.go`
- `scripts/jwtgen/main.go`

## File Dihapus

- Tidak ada.

## Perubahan Utama

- Menambahkan endpoint publik auth:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/forgot-password`
  - `POST /api/v1/auth/reset-password`
  - `POST /api/v1/auth/google`
- Menambahkan penyimpanan token reset password di database.
- Menambahkan pengiriman email reset password via Resend (`RESEND_API_KEY`, opsional `RESEND_FROM_EMAIL`).
- Menjaga alur Google login lama, sekaligus bisa menautkan `google_id` ke akun email yang sudah ada.
- Menambahkan integration test auth end-to-end di level route/service.

## Alur Testing Yang Dijalankan

### 1. Database

Command:

```powershell
docker compose up -d postgres
migrate -path db/migrations -database "postgres://root:secretpassword@localhost:5435/mosque_saas?sslmode=disable" -verbose up
```

Hasil:

- Container database berjalan.
- Migration `000009_add_password_reset_auth` sukses diaplikasikan.

### 2. Automated Test Backend

Command:

```powershell
$env:DB_URL='postgres://root:secretpassword@localhost:5435/mosque_saas?sslmode=disable'; go test ./...
```

Hasil:

- Lulus.
- Test auth baru memverifikasi alur register -> login -> forgot password -> reset password -> login ulang.

## Catatan Operasional

- Forgot password akan mengirim email asli jika `RESEND_API_KEY` tersedia.
- URL reset password dibentuk dari `APP_BASE_URL`, default ke `http://localhost:3000`.
