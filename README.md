# url2md

Convert web pages to clean, LLM-friendly Markdown.

## Features

- **Three-layer fallback pipeline**: Content negotiation → Static fetch → Headless Chrome
- **Smart extraction**: Readability-based article extraction with noise removal
- **File type support**: PDF, DOCX, XLSX, CSV, images (magic bytes detection)
- **YAML frontmatter**: Auto-generated title, description, og:image metadata
- **Token estimation**: Approximate token count with CJK support
- **Metadata extraction**: Title, description, Open Graph tags
- **Dual interface**: CLI tool + HTTP API server

## Install

```bash
go install github.com/elonfeng/url2md/cmd/url2md@latest
```

## Usage

### CLI

```bash
# basic conversion
url2md https://example.com

# save to file
url2md https://example.com -o output.md

# force specific method
url2md https://example.com --method static

# enable headless Chrome fallback
url2md https://example.com --browser

# retain images
url2md https://example.com --images

# batch convert
url2md batch https://example.com https://example.org
```

### HTTP API

Start the server:

```bash
url2md serve --port 8080
```

#### `GET /{url}`

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

#### `POST /`

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

#### Response

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

Error response:

```json
{
  "error": "HTTP 404"
}
```

#### `GET /health`

Health check endpoint.

```json
{"status": "ok"}
```

## Deploy

Anyone can deploy url2md as a self-hosted service. Image AI description is optional and requires your own Cloudflare credentials.

### Docker

```bash
docker build -t url2md .
docker run -p 8080:8080 url2md
```

With image AI description enabled:

```bash
docker run -p 8080:8080 \
  -e CLOUDFLARE_ACCOUNT_ID="your-account-id" \
  -e CLOUDFLARE_API_TOKEN="your-api-token" \
  url2md
```

### Image AI Description

Image URLs (PNG, JPEG, GIF, WEBP) can be described using **Cloudflare Workers AI** vision model (`@cf/meta/llama-3.2-11b-vision-instruct`). This feature is **optional** — without credentials, url2md falls back to image metadata + embed.

To enable, set two environment variables (see [.env.example](.env.example)):

```bash
export CLOUDFLARE_ACCOUNT_ID="your-account-id"
export CLOUDFLARE_API_TOKEN="your-api-token"
```

**How to get credentials:**

1. Sign up at [dash.cloudflare.com](https://dash.cloudflare.com/)
2. **Account ID**: Dashboard right sidebar → "Account ID"
3. **API Token**: My Profile → API Tokens → Create Token → template "Workers AI (Read)"

**Pricing**: Free tier includes 10,000 neurons/day (roughly dozens of image descriptions). Beyond that, $0.011 per 1,000 neurons. Each deployment uses its own credentials and quota. See [Workers AI Pricing](https://developers.cloudflare.com/workers-ai/platform/pricing/).

## Architecture

```
URL → [Layer 1: Content Negotiation]
       ↓ fail
      [Layer 2: Static HTTP + Readability + html-to-markdown]
       ↓ fail
      [Layer 3: Headless Chrome + Readability + html-to-markdown]
       ↓
      Clean Markdown + Metadata + Token Count
```

## Benchmark vs markdown.new

Tested against [markdown.new](https://markdown.new) (Cloudflare Workers AI) across 11 test cases. Full report: [BENCHMARK.md](BENCHMARK.md).

| Category | url2md Wins | Ties | markdown.new Wins |
|----------|-------------|------|-------------------|
| Web Pages (6 tests) | **4** | 2 | 0 |
| File Types (5 tests) | **2** | 3 | 0 |
| **Total (11 tests)** | **6** | **5** | **0** |

**Key advantages**:
- **1.3x-4.6x fewer tokens** on web pages (avg 2.2x more efficient)
- **Complete CJK extraction** (markdown.new truncates Chinese content)
- **Better noise removal** — strips nav/ads/UI chrome via readability
- **CSV → markdown table** (vs code block), **XLSX → markdown table** (markdown.new fails)
- **Self-hosted**, open source, customizable pipeline

## License

MIT
