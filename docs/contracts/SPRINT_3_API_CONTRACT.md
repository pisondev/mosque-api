# Sprint 3 API Contract

## Scope
- Prayer time settings
- Prayer times daily
- Prayer duties
- Special days
- Prayer calendar aggregate

## Base Path
- `/api/v1/tenant`

## Endpoints
- `GET /prayer-time-settings`
- `PUT /prayer-time-settings`
- `GET /prayer-times-daily`
- `POST /prayer-times-daily`
- `GET /prayer-times-daily/{id}`
- `PUT /prayer-times-daily/{id}`
- `DELETE /prayer-times-daily/{id}`
- `GET /prayer-duties`
- `POST /prayer-duties`
- `GET /prayer-duties/{id}`
- `PUT /prayer-duties/{id}`
- `DELETE /prayer-duties/{id}`
- `GET /special-days`
- `POST /special-days`
- `GET /special-days/{id}`
- `PUT /special-days/{id}`
- `DELETE /special-days/{id}`
- `GET /prayer-calendar`

## Validations
- `location_mode=city` mewajibkan `city_name`.
- `location_mode=coordinates` mewajibkan `latitude` dan `longitude`.
- `category=fardhu` mewajibkan `prayer`.
- Konflik unik:
  - `(tenant_id, day_date)` pada `prayer_times_daily`.
  - `(tenant_id, kind, day_date)` pada `special_days`.
