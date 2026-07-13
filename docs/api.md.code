# SNISID API Reference

Base URL: `https://<host>/api/v1`

All endpoints require `Authorization: Bearer <token>` unless noted.

## Health & Status

| Method | Path               | Description                  |
|--------|--------------------|------------------------------|
| GET    | `/health`          | Liveness check               |
| GET    | `/health/ready`    | Readiness check (deps ready) |
| GET    | `/status`          | Full system status           |

## Data Records

| Method | Path                | Description                        |
|--------|---------------------|------------------------------------|
| GET    | `/records`          | List records (paginated)           |
| POST   | `/records`          | Create a new record                |
| GET    | `/records/:id`      | Get a single record by ID          |
| PUT    | `/records/:id`      | Update an existing record          |
| DELETE | `/records/:id`      | Soft-delete a record               |

### Query Parameters (GET /records)
- `page` (int, default 1)
- `limit` (int, default 50, max 500)
- `sort` (string, e.g. `-created_at`)
- `filter` (JSON-encoded filter object)

### Request Body (POST/PUT /records)
```json
{
  "type": "observation",
  "payload": { ... },
  "tags": ["field", "urgent"]
}
```

## Sync

| Method | Path                  | Description                          |
|--------|-----------------------|--------------------------------------|
| POST   | `/sync/upload`        | Upload offline data file             |
| GET    | `/sync/status`        | Sync agent status & last sync time   |
| GET    | `/sync/conflicts`     | List unresolved sync conflicts       |

## Administration

| Method | Path                  | Description                          |
|--------|-----------------------|--------------------------------------|
| GET    | `/admin/users`        | List users                           |
| POST   | `/admin/users`        | Create user                          |
| DELETE | `/admin/users/:id`    | Remove user                          |
| GET    | `/admin/logs`         | Recent application logs              |
| POST   | `/admin/restart`      | Restart core service                 |

## Error Responses
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Field 'type' is required",
    "details": { "field": "type" }
  }
}
```

HTTP status codes: 200 (OK), 201 (Created), 400 (Bad Request), 401 (Unauthorized), 403 (Forbidden), 404 (Not Found), 409 (Conflict), 500 (Internal Server Error).
