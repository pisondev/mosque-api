# Sprint 6 Hardening Contract

## Scope
- Security hardening
- Reliability hardening
- Performance hardening
- Contract and integration quality gate

## Mandatory Controls
- Rate limiting untuk endpoint auth dan endpoint write.
- Payload validation ketat untuk seluruh query/body.
- Error mapping terstandar: `400`, `401`, `403`, `404`, `409`, `422`, `500`.
- Tenant scoping wajib di setiap query data domain.
- Pagination default dan max limit untuk endpoint list.

## Verification Gate
- Contract test seluruh endpoint sprint 1-5.
- Integration test minimal happy-path, validation-path, dan authz-path.
- Smoke test pasca migrasi dan startup service.
