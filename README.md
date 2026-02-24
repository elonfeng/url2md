# url2md

Convert web pages to clean, LLM-friendly Markdown.

## Features

- **Three-layer fallback pipeline**: Content negotiation → Static fetch → Headless Chrome
- **Smart extraction**: Readability-based article extraction with noise removal
- **15 file types**: PDF, DOCX, XLSX, XLS, ODT, CSV, JSON, XML, HTML, TXT, MD, PNG, JPG, SVG, WEBP
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

```bash
# start server
url2md serve --port 8080

# GET
curl "http://localhost:8080/https://example.com"

# POST
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

Full API documentation: [API.md](API.md)

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

Image URLs (PNG, JPEG, GIF, WEBP, SVG) can be described using **Cloudflare Workers AI** vision model (`@cf/meta/llama-3.2-11b-vision-instruct`). This feature is **optional** — without credentials, url2md falls back to image metadata + embed.

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

Tested against [markdown.new](https://markdown.new) (Cloudflare Workers AI) across 16 test cases. Full report: [BENCHMARK.md](BENCHMARK.md).

| Category | url2md Wins | Ties | markdown.new Wins |
|----------|-------------|------|-------------------|
| Web Pages (6 tests) | **4** | 2 | 0 |
| File Types (10 tests) | **2** | 8 | 0 |
| **Total (16 tests)** | **6** | **10** | **0** |

**Key advantages**:
- **1.3x-4.6x fewer tokens** on web pages (avg 2.2x more efficient)
- **Complete CJK extraction** (markdown.new truncates Chinese content)
- **Better noise removal** — strips nav/ads/UI chrome via readability
- **CSV → markdown table** (vs code block), **XLSX → markdown table** (markdown.new fails)
- **Self-hosted**, open source, customizable pipeline

## License

MIT
