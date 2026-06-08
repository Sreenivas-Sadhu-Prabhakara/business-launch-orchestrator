---
sidebar_position: 6
title: API reference
---

# API reference

The UI is just a client of this API — you can drive an entire launch from the
terminal. Base URL defaults to `http://localhost:8080`.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/healthz` | Liveness check |
| `GET`  | `/api/v1/countries` | All countries + their step plans |
| `GET`  | `/api/v1/countries/{code}/plan` | Plan for `IN` / `PH` / `US` |
| `POST` | `/api/v1/businesses` | Create a launch → `{business, steps}` |
| `GET`  | `/api/v1/businesses` | List recent launches |
| `GET`  | `/api/v1/businesses/{id}` | Get one launch + its steps |
| `GET`  | `/api/v1/businesses/{id}/steps` | Steps only |
| `POST` | `/api/v1/businesses/{id}/advance` | Execute the next pending step |
| `POST` | `/api/v1/businesses/{id}/run` | Execute all remaining steps |

## Create a launch

```bash
curl -s -X POST http://localhost:8080/api/v1/businesses \
  -H 'Content-Type: application/json' \
  -d '{
    "country": "IN",
    "entity_type": "Private Limited Company",
    "legal_name": "Acme Technologies",
    "founder_name": "Jane Doe",
    "founder_email": "jane@acme.com",
    "founder_id_number": "ABCDE1234F",
    "address": { "city": "Bengaluru", "state": "Karnataka", "postal_code": "560001", "country": "IN" }
  }' | jq
```

## Run the whole pipeline

```bash
# capture the id
BIZ=$(curl -s -X POST http://localhost:8080/api/v1/businesses \
  -H 'Content-Type: application/json' \
  -d '{"country":"US","entity_type":"LLC","legal_name":"Globex LLC","founder_name":"John Roe","address":{"state":"Delaware","country":"US"}}' \
  | jq -r '.business.id')

# run all 11 steps to completion
curl -s -X POST http://localhost:8080/api/v1/businesses/$BIZ/run \
  | jq '.business.status, (.steps[] | {seq, title, status, external_ref})'
```

## Step one at a time

```bash
curl -s -X POST http://localhost:8080/api/v1/businesses/$BIZ/advance | jq
```

## Response shape

```json
{
  "business": { "id": "…", "country": "IN", "status": "completed", "...": "…" },
  "steps": [
    {
      "seq": 5,
      "step_type": "ip_check",
      "provider": "RDAP + IP India",
      "mode": "live",
      "status": "completed",
      "external_ref": "TM-8746792",
      "response": { "domains": { "acme.com": "registered" }, "trademark": { "...": "…" } }
    }
  ]
}
```
