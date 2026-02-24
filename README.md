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

Currently, image URLs output metadata (filename, size, content-type) and a markdown image embed. AI-powered image description is planned via a pluggable vision provider:

1. Detect image type via magic bytes (PNG, JPEG, GIF, WEBP)
2. Base64-encode the image data
3. Send to a vision model API (Cloudflare Workers AI `@cf/meta/llama-3.2-11b-vision-instruct`, OpenAI GPT-4o, or Claude) with the prompt: *"Describe this image in detail for a markdown document"*
4. Return the AI-generated description as markdown text

This will be configurable via `--vision-provider` flag and `VISION_API_KEY` environment variable. Contributions welcome.

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

## License

MIT
