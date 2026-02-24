# url2md vs markdown.new Benchmark Report

**Date**: 2026-02-24
**url2md version**: v0.1.0 (Go, static layer, frontmatter enabled)
**markdown.new**: Cloudflare Workers AI

---

## Test Matrix

| # | Type | URL |
|---|------|-----|
| 1 | Simple Page | `https://example.com` |
| 2 | English Blog | `https://go.dev/blog/go1.23` |
| 3 | Chinese Article | `https://sspai.com/post/88053` |
| 4 | Tech Documentation | `https://docs.python.org/3/tutorial/classes.html` |
| 5 | News/Blog | `https://blog.cloudflare.com/markdown-for-agents` |
| 6 | Wikipedia | `https://en.wikipedia.org/wiki/Markdown` |

---

## Results Summary

### Raw Data

| Page Type | url2md Tokens | url2md Chars | url2md FM | markdown.new Tokens | markdown.new Chars | markdown.new FM |
|-----------|---------------|--------------|-----------|---------------------|--------------------|-----------------|
| Simple Page | 32 | 199 | Yes | 42 | 167 | No |
| English Blog | 986 | 5,272 | Yes | 1,359 | 5,436 | Yes |
| Chinese Article | 6,598 | 14,635 | Yes | 2,296 | 6,579 | Yes |
| Tech Docs | 7,078 | 39,639 | Yes | 10,392 | 41,501 | Yes |
| News/Blog (CF) | 1,886 | 11,278 | Yes | 3,697 | 14,166 | Yes |
| Wikipedia | 3,706 | 38,658 | Yes | 17,174 | 68,255 | Yes |

### Token Efficiency

| Page Type | url2md | markdown.new | Ratio |
|-----------|--------|--------------|-------|
| Simple Page | 32 | 42 | **1.3x fewer** |
| English Blog | 986 | 1,359 | **1.4x fewer** |
| Chinese Article | 6,598 | 2,296 | md.new truncated content |
| Tech Docs | 7,078 | 10,392 | **1.5x fewer** |
| News/Blog (CF) | 1,886 | 3,697 | **2.0x fewer** |
| Wikipedia | 3,706 | 17,174 | **4.6x fewer** |

url2md produces **fewer tokens** in 5 out of 6 tests (Chinese article excluded due to markdown.new truncation). Average reduction: **2.2x more token-efficient**.

### Content Quality

| Page Type | url2md | markdown.new | Winner |
|-----------|--------|--------------|--------|
| **Simple Page** | Frontmatter + `# Title` + clean body | No frontmatter, `# Title` + body | **url2md** |
| **English Blog** | Full article + frontmatter + title | Full article + frontmatter + author/date | Tie |
| **Chinese Article** | **Complete** (14,635 chars) | Truncated (6,579 chars) | **url2md** |
| **Tech Docs** | Full content, **no nav noise** | Full content, includes nav breadcrumbs | **url2md** |
| **News/Blog (CF)** | Clean article + frontmatter | Full article + frontmatter + author | Tie |
| **Wikipedia** | Clean article + references | Full page with nav/tables/templates — noisy | **url2md** |

**Score: url2md wins 4, ties 2, loses 0.**

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

## Recommendations

1. **Improve GitHub/SPA extraction** — detect README content, handle JS-heavy pages via Layer 3
2. **Add response caching** — in-memory or Redis cache for HTTP server mode
3. **Deploy to edge** — Fly.io or similar for lower latency
