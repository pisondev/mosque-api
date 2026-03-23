# Sprint 5 API Contract

## Scope
- Donation channels
- Social links
- External links
- Feature catalog
- Website features

## Base Path
- Tenant: `/api/v1/tenant`
- Public: `/api/v1/public/{hostname}`

## Tenant Endpoints
- `GET /donation-channels`
- `POST /donation-channels`
- `GET /donation-channels/{id}`
- `PUT /donation-channels/{id}`
- `DELETE /donation-channels/{id}`
- `GET /social-links`
- `POST /social-links`
- `GET /social-links/{id}`
- `PUT /social-links/{id}`
- `DELETE /social-links/{id}`
- `GET /external-links`
- `POST /external-links`
- `GET /external-links/{id}`
- `PUT /external-links/{id}`
- `DELETE /external-links/{id}`
- `GET /feature-catalog`
- `GET /website-features`
- `PUT /website-features/{feature_id}`
- `PATCH /website-features/bulk`

## Public Endpoints
- `GET /donation-channels`
- `GET /social-links`
- `GET /external-links`

## Validations
- `channel_type=bank_account` mewajibkan `bank_name`, `account_number`, `account_holder_name`.
- `channel_type=qris` mewajibkan `qris_image_url`.
