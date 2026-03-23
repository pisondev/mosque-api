# Sprint 4 API Contract

## Scope
- Events
- Gallery albums
- Gallery items
- Management members
- Public feed endpoints

## Base Path
- Tenant: `/api/v1/tenant`
- Public: `/api/v1/public/{hostname}`

## Tenant Endpoints
- `GET /events`
- `POST /events`
- `GET /events/{id}`
- `PUT /events/{id}`
- `PATCH /events/{id}/status`
- `DELETE /events/{id}`
- `GET /gallery/albums`
- `POST /gallery/albums`
- `GET /gallery/albums/{id}`
- `PUT /gallery/albums/{id}`
- `DELETE /gallery/albums/{id}`
- `GET /gallery/items`
- `POST /gallery/items`
- `GET /gallery/items/{id}`
- `PUT /gallery/items/{id}`
- `DELETE /gallery/items/{id}`
- `GET /management-members`
- `POST /management-members`
- `GET /management-members/{id}`
- `PUT /management-members/{id}`
- `DELETE /management-members/{id}`

## Public Endpoints
- `GET /events`
- `GET /gallery/albums`
- `GET /gallery/items`
- `GET /management-members`
