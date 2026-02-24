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

### Image Understanding

Image URLs are processed with AI-powered description via **Cloudflare Workers AI** (free tier: 10,000 neurons/day). See [.env.example](.env.example) for setup:

```bash
export CLOUDFLARE_ACCOUNT_ID="your-account-id"
export CLOUDFLARE_API_TOKEN="your-api-token"
url2md https://example.com/photo.png
```

Uses `@cf/meta/llama-3.2-11b-vision-instruct` to generate natural language descriptions. Without env vars, falls back to metadata + image embed.

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

### HTTP Server

```bash
# start server
url2md serve --port 8080
```

**GET request:**
```
GET /https://example.com?method=auto&retain_images=false
```

**POST request:**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","method":"auto"}'
```

**Response headers:**
- `X-Markdown-Tokens` — estimated token count
- `X-Convert-Method` — which layer succeeded
- `X-Fetch-Time` — fetch duration

### Docker

```bash
docker build -t url2md .
docker run -p 8080:8080 url2md
```

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
