# Sprint 2 API Contract

## Scope
- Tags
- Posts
- Static pages

## Base Path
- `/api/v1/tenant`

## Endpoints

### Tags (`tags`)
- `GET /tags`
- `POST /tags`
- `PATCH /tags/{id}`
- `DELETE /tags/{id}`

Rule:
- Slug dibangkitkan server-side.
- Unik per `(tenant_id, scope, slug)`.
- `DELETE` mengembalikan 409 bila tag masih dipakai relasi.

### Posts (`posts`)
- `GET /posts`
- `POST /posts`
- `GET /posts/{id}`
- `PUT /posts/{id}`
- `PATCH /posts/{id}/status`
- `DELETE /posts/{id}`

Query `GET /posts`:
- `category`: `news_activity|announcement|reflection|static_page`
- `status`: `draft|published|archived`
- `search`, `page`, `limit`, `sort_by`, `sort_order`

Rule:
- Jika status `published` dan `published_at` kosong, backend isi waktu sekarang.
- `expired_at` tidak boleh lebih kecil dari `published_at`.
- Delete permanen sesuai schema saat ini.

### Static Pages (berbasis `posts`)
- `GET /static-pages`
- `PUT /static-pages/{slug}`

Rule:
- `GET /static-pages` memfilter `posts.category=static_page`.
- Upsert halaman statis berdasarkan slug.
