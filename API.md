# API Documentation

## Start Server

```bash
url2md serve --port 8080
```

## Endpoints

### `GET /{url}`

Convert a URL to Markdown via GET request.

```bash
curl "http://localhost:8080/https://example.com"
```

Query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `method` | string | `auto` | Conversion method: `auto`, `negotiate`, `static`, `browser` |
| `retain_images` | bool | `false` | Keep image tags in output |
| `retain_links` | bool | `true` | Keep hyperlinks in output |
| `enable_browser` | bool | `false` | Enable headless Chrome fallback |
| `frontmatter` | bool | `true` | Prepend YAML frontmatter |

Example:

```bash
curl "http://localhost:8080/https://go.dev/blog/go1.23?retain_images=true&frontmatter=false"
```

### `POST /`

Convert a URL to Markdown via POST request.

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "method": "auto",
    "retain_images": false,
    "retain_links": true,
    "frontmatter": true
  }'
```

Request body:

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `url` | string | **yes** | — | URL to convert |
| `method` | string | no | `auto` | Conversion method |
| `retain_images` | bool | no | `false` | Keep image tags |
| `retain_links` | bool | no | `true` | Keep hyperlinks |
| `frontmatter` | bool | no | `true` | Prepend YAML frontmatter |

### Response

Both GET and POST return the same JSON response:

```json
{
  "url": "https://example.com",
  "title": "Example Domain",
  "description": "This domain is for use in illustrative examples.",
  "markdown": "---\ntitle: Example Domain\n---\n\n# Example Domain\n\nThis domain is for use in illustrative examples...",
  "token_count": 32,
  "method": "static",
  "metadata": {
    "og:title": "Example Domain"
  },
  "fetch_ms": 215,
  "convert_ms": 3
}
```

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Final URL (after redirects) |
| `title` | string | Page title |
| `description` | string | Page description |
| `markdown` | string | Converted Markdown content |
| `token_count` | int | Estimated token count |
| `method` | string | Which layer succeeded (`negotiate`, `static`, `browser`) |
| `metadata` | object | Extracted Open Graph / meta tags |
| `fetch_ms` | int | Fetch duration in milliseconds |
| `convert_ms` | int | Conversion duration in milliseconds |

Response headers:

| Header | Description |
|--------|-------------|
| `X-Markdown-Tokens` | Estimated token count |
| `X-Convert-Method` | Which layer succeeded |
| `X-Fetch-Time` | Fetch duration |

### Error Response

```json
{
  "error": "HTTP 404"
}
```

### `GET /health`

Health check endpoint.

```json
{"status": "ok"}
```
