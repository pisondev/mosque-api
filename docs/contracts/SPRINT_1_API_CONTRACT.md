# Sprint 1 API Contract

## Scope
- Auth
- Tenant context
- Domain mapping
- Profil masjid

## Base Path
- `/api/v1`

## Response Envelope

Sukses:

```json
{
  "status": "success",
  "message": "ok",
  "data": {},
  "meta": {}
}
```

Error:

```json
{
  "status": "error",
  "message": "validation failed",
  "errors": [
    { "field": "hostname", "message": "invalid hostname" }
  ]
}
```

## Endpoints

### Auth
- `POST /auth/google` (Public)
  - Body: `{ "token": "google_id_token" }`
  - 401 jika token Google tidak valid.

### Tenant Context
- `GET /tenant/me` (Bearer JWT)

### Domain Mapping (`website_domains`)
- `GET /tenant/domains`
- `POST /tenant/domains`
- `PATCH /tenant/domains/{id}`
- `DELETE /tenant/domains/{id}`

### Profil Masjid (`masjid_profiles`)
- `GET /tenant/profile`
- `PUT /tenant/profile`

## Validations
- `official_name` wajib.
- `kind` harus enum `place_kind`.
- `established_year` dalam rentang 1000..3000 jika diisi.
- `capacity_male`, `capacity_female`, `floors_count` tidak boleh negatif.
