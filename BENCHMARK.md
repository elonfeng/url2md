# url2md vs markdown.new Benchmark Report

**Date**: 2026-02-24
**url2md version**: v0.1.0 (Go, static layer, frontmatter enabled)
**markdown.new**: Cloudflare Workers AI

---

## Test Matrix

### Web Pages

| # | Type | URL |
|---|------|-----|
| 1 | Simple Page | `https://example.com` |
| 2 | English Blog | `https://go.dev/blog/go1.23` |
| 3 | Chinese Article | `https://sspai.com/post/88053` |
| 4 | Tech Documentation | `https://docs.python.org/3/tutorial/classes.html` |
| 5 | News/Blog | `https://blog.cloudflare.com/markdown-for-agents` |
| 6 | Wikipedia | `https://en.wikipedia.org/wiki/Markdown` |

### File Types

| # | Type | URL |
|---|------|-----|
| 7 | PDF | `https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf` |
| 8 | DOCX | `https://calibre-ebook.com/downloads/demos/demo.docx` |
| 9 | PNG (Image) | `https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png` |
| 10 | CSV | `https://people.sc.fsu.edu/~jburkardt/data/csv/addresses.csv` |
| 11 | XLSX | `https://go.microsoft.com/fwlink/?LinkID=521962` |

---

## Web Page Results

### Token Efficiency

| Page Type | url2md | markdown.new | Ratio |
|-----------|--------|--------------|-------|
| Simple Page | 32 | 42 | **1.3x fewer** |
| English Blog | 986 | 1,359 | **1.4x fewer** |
| Chinese Article | 6,598 | 2,296 | md.new truncated content |
| Tech Docs | 7,078 | 10,392 | **1.5x fewer** |
| News/Blog (CF) | 1,886 | 3,697 | **2.0x fewer** |
| Wikipedia | 3,706 | 17,174 | **4.6x fewer** |

url2md produces **fewer tokens** in 5 out of 6 web page tests. Average reduction: **2.2x more token-efficient**.

### Content Quality

| Page Type | url2md | markdown.new | Winner |
|-----------|--------|--------------|--------|
| **Simple Page** | Frontmatter + `# Title` + clean body | No frontmatter, body only | **url2md** |
| **English Blog** | Full article + frontmatter + title | Full article + frontmatter + author/date | Tie |
| **Chinese Article** | **Complete** (14,635 chars) | Truncated (6,579 chars) | **url2md** |
| **Tech Docs** | Full content, **no nav noise** | Full content, includes nav breadcrumbs | **url2md** |
| **News/Blog (CF)** | Clean article + frontmatter | Full article + frontmatter + author | Tie |
| **Wikipedia** | Clean article + references | Full page with nav/tables/templates — noisy | **url2md** |

**Web page score: url2md wins 4, ties 2, loses 0.**

---

## File Type Results

| File Type | url2md | markdown.new | Winner |
|-----------|--------|--------------|--------|
| **PDF** | FAIL (outputs raw binary) | Extracts text + metadata (83 tokens) | **markdown.new** |
| **DOCX** | FAIL (outputs raw binary) | Full document content (2,836 tokens) | **markdown.new** |
| **PNG** | FAIL (outputs raw binary) | AI image description (81 tokens) | **markdown.new** |
| **CSV** | FAIL (outputs raw binary) | Wraps in markdown code block | **markdown.new** |
| **XLSX** | FAIL (outputs raw binary) | FAIL (outputs raw binary) | Tie |

**File type score: url2md wins 0, ties 1, loses 4.**

### File Type Detail

**PDF** — markdown.new extracts structured content via Workers AI:
```
# dummy.pdf
## Metadata
- Author=Evangelos Vlachogiannis
- Creator=Writer
## Content
Dummy PDF file
```

**DOCX** — markdown.new converts full document structure (headings, paragraphs, lists, tables):
```
# demo.docx
Demonstration of DOCX support in calibre
This document demonstrates the ability of the calibre DOCX Input plugin...
```

**PNG** — markdown.new uses AI vision to describe the image:
```
The image displays the Google logo. The logo is composed of four overlapping,
rounded shapes in the colors red, yellow, green, and blue...
```

**CSV** — markdown.new wraps raw CSV in a code block:
```
# addresses.csv
\```csv
John,Doe,120 jefferson st.,Riverside, NJ, 08075
Jack,McGinnis,220 hobo Av.,Phila, PA,09119
\```
```

**XLSX** — Both tools fail. markdown.new outputs raw PK binary, url2md same.

---

## Overall Summary

### Combined Score

| Category | url2md Wins | Ties | markdown.new Wins |
|----------|-------------|------|-------------------|
| Web Pages (6 tests) | **4** | 2 | 0 |
| File Types (5 tests) | 0 | 1 | **4** |
| **Total (11 tests)** | **4** | **3** | **4** |

### Where Each Tool Excels

**url2md is better for**:
- Web page → Markdown conversion (cleaner, fewer tokens)
- CJK/Chinese content (complete extraction vs truncation)
- Noise removal (readability-based, strips nav/ads/UI)
- Self-hosted / offline use
- Customizable pipeline

**markdown.new is better for**:
- Non-HTML file conversion (PDF, DOCX, images)
- AI-powered image description (OCR/vision)
- Zero-setup usage (SaaS, no deployment needed)

---

## Detailed Analysis

### 1. Noise Removal

**url2md** uses go-readability + goquery selector-based cleaning:
- Removes navigation bars, headers, footers, sidebars
- Removes cookie banners, ad containers, popup modals
- Removes script/style/noscript/iframe/svg tags

**markdown.new** does raw HTML-to-markdown with less aggressive cleaning, preserving navigation breadcrumbs, login prompts, and UI chrome.

| Page | url2md Noise | markdown.new Noise |
|------|-------------|-------------------|
| Python Docs | None — pure tutorial content | Navigation bar, theme picker, breadcrumbs |
| Wikipedia | Minimal — article + references | Jump-to-content, edit links, nav tables |

### 2. Metadata & Frontmatter

| Feature | url2md | markdown.new |
|---------|--------|--------------|
| YAML frontmatter | `--frontmatter` (default on) | Sometimes (inconsistent) |
| Title heading | Auto `# Title` when missing | Auto `# Title` |
| Description | Yes | Yes |
| OG Image | Yes (when available) | Yes |
| Author/Date | Not extracted | Sometimes present |
| API JSON fields | title, description, metadata, tokens | title, tokens |

### 3. CJK (Chinese) Content

url2md significantly outperforms on Chinese content:
- **url2md**: 14,635 chars — complete article with all sections
- **markdown.new**: 6,579 chars — truncated, missing later sections

go-readability handles CJK text extraction better than Cloudflare Workers AI's converter.

### 4. Anti-Bot Handling

Both tools fail on anti-bot protected sites (Reuters, TechCrunch). url2md's Layer 3 (headless Chrome) could potentially bypass some protections but was not benchmarked here.

---

## Feature Comparison

| Feature | url2md | markdown.new |
|---------|--------|--------------|
| Web page → Markdown | Yes | Yes |
| PDF → Markdown | **No** | Yes |
| DOCX → Markdown | **No** | Yes |
| Image → Description | **No** | Yes (AI vision) |
| CSV → Markdown | **No** | Yes (code block) |
| XLSX → Markdown | No | No |
| YAML frontmatter | Yes (default on) | Inconsistent |
| Auto `# Title` heading | Yes | Yes |
| Noise removal | Better (readability-based) | Weaker (preserves UI chrome) |
| Token efficiency | **1.3x-4.6x fewer** | Higher token count |
| CJK support | Complete extraction | Truncation issues |
| Self-hosted | Yes | No (Cloudflare only) |
| Open source | Yes | No |
| CLI tool | Yes | No |
| HTTP API | Yes | Yes |
| Go SDK | Yes (`import pkg`) | No |
| Headless Chrome fallback | Yes (Layer 3) | Yes (Browser Rendering API) |
| Customizable cleaning rules | Yes (goquery selectors) | No |

---

## Roadmap for url2md

1. **Add PDF support** — integrate Go PDF parsing library (e.g. pdfcpu, unipdf)
2. **Add DOCX support** — integrate Go DOCX parser (e.g. fumiama/go-docx)
3. **Add CSV/XLSX support** — table → markdown table conversion
4. **Add image OCR** — optional Tesseract or API-based vision
5. **Improve GitHub/SPA extraction** — detect README content, handle JS-heavy pages
6. **Add response caching** — in-memory or Redis cache for HTTP server mode
