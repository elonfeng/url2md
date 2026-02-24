# url2md vs markdown.new Benchmark Report

**Date**: 2026-02-24
**url2md version**: v0.1.0 (Go, static layer only, frontmatter enabled)
**markdown.new**: Cloudflare Workers AI

---

## Test Matrix

| # | Type | URL | Anti-Bot |
|---|------|-----|----------|
| 1 | Simple Page | `https://example.com` | No |
| 2 | English Blog | `https://go.dev/blog/go1.23` | No |
| 3 | Chinese Article | `https://sspai.com/post/88053` | No |
| 4 | Tech Documentation | `https://docs.python.org/3/tutorial/classes.html` | No |
| 5 | News/Blog | `https://blog.cloudflare.com/markdown-for-agents` | No |
| 6 | Wikipedia | `https://en.wikipedia.org/wiki/Markdown` | No |
| 7 | GitHub Repo | `https://github.com/chromedp/chromedp` | Partial |
| 8 | News (Reuters) | `https://www.reuters.com/technology/artificial-intelligence/` | Yes |

---

## Results Summary

### Token Efficiency

| Page Type | url2md Tokens | markdown.new Tokens | Ratio |
|-----------|---------------|---------------------|-------|
| Simple Page | 22 | 42 | **1.9x fewer** |
| English Blog | 942 | 1,359 | **1.4x fewer** |
| Chinese Article | 6,511 | 2,326 | md.new truncated |
| Tech Docs | 7,008 | 10,392 | **1.5x fewer** |
| News/Blog (CF) | 1,795 | 3,697 | **2.1x fewer** |
| Wikipedia | 3,693 | 17,174 | **4.6x fewer** |
| GitHub Repo | 44 | 4,805 | N/A (both poor) |
| News (Reuters) | FAIL | FAIL | - |

url2md produces **fewer tokens** in 5 out of 6 successful tests, averaging **2.3x more token-efficient**. For LLM consumption, fewer tokens = lower cost + better signal-to-noise ratio.

### Content Quality Comparison

| Page Type | url2md | markdown.new | Winner |
|-----------|--------|--------------|--------|
| **Simple Page** | YAML frontmatter + `# Title` + clean body | YAML frontmatter + `# Title` + clean body | Tie |
| **English Blog** | Full article + frontmatter + title heading | Full article + frontmatter + author/date | Tie |
| **Chinese Article** | **Complete article** (14,435 chars), all sections | Truncated (6,603 chars), content cut off | **url2md** |
| **Tech Docs** | Full tutorial content, **no navigation noise** | Full content, **includes navigation breadcrumbs** | **url2md** |
| **News/Blog (CF)** | Clean article body + frontmatter | Full article + frontmatter + author info | Tie |
| **Wikipedia** | Article body, clean extraction, references intact | **Full page with nav, tables, templates** — noisy | **url2md** |
| **GitHub Repo** | Only extracted description + license (too sparse) | Full page dump including UI chrome (too noisy) | Neither |
| **News (Reuters)** | Both fail — anti-bot protection | Both fail | Tie |

**url2md wins 3, ties 4, loses 0** (excluding mutual failures).

---

## Detailed Analysis

### 1. Noise Removal

**url2md** uses go-readability + goquery selector-based cleaning. This removes:
- Navigation bars, headers, footers, sidebars
- Cookie banners, ad containers, popup modals
- Script/style/noscript/iframe/svg tags

**markdown.new** (Cloudflare Workers AI) does raw HTML-to-markdown conversion with less aggressive cleaning. Result: it often preserves navigation breadcrumbs, login prompts, and UI chrome.

| Page | url2md Noise | markdown.new Noise |
|------|-------------|-------------------|
| Python Docs | None — pure tutorial content | Navigation bar, theme picker, breadcrumbs |
| Wikipedia | Minimal — article + references | Jump-to-content, edit links, nav tables |
| GitHub | Over-stripped (only description) | Full UI: sign-in prompts, fork/star buttons |

### 2. Metadata & Frontmatter

Both tools now produce YAML frontmatter:

| Feature | url2md | markdown.new |
|---------|--------|--------------|
| YAML frontmatter | `--frontmatter` (default on) | Always on |
| Title | title field + auto `# Title` | title field + `# Title` |
| Description | description field | description field |
| OG Image | image field (when available) | image field |
| Author/Date | Not extracted | Sometimes present |
| Output format | Both inline YAML + API JSON | Inline YAML only |

### 3. Anti-Bot Handling

Both tools fail equally on anti-bot protected sites (Reuters, TechCrunch). Neither has anti-detection capabilities in the default static fetch path. url2md's Layer 3 (headless Chrome) could potentially bypass some of these, but was not tested here.

### 4. CJK (Chinese) Content

url2md significantly outperforms on Chinese content:
- **url2md**: 14,435 chars — complete article with all sections
- **markdown.new**: 6,603 chars — truncated, missing later sections

go-readability handles CJK text extraction better than Cloudflare Workers AI's converter.

---

## Feature Comparison

| Feature | url2md | markdown.new |
|---------|--------|--------------|
| YAML frontmatter | Yes | Yes |
| Auto `# Title` heading | Yes | Yes |
| Noise removal quality | Better (readability-based) | Weaker (preserves UI chrome) |
| Token efficiency | 1.4x-4.6x fewer tokens | Higher token count |
| CJK support | Complete extraction | Truncation issues |
| Self-hosted | Yes | No (Cloudflare only) |
| Open source | Yes | No |
| CLI tool | Yes | No |
| HTTP API | Yes | Yes |
| Go SDK | Yes (import pkg) | No |
| Headless Chrome fallback | Yes (Layer 3) | Yes (Browser Rendering API) |
| Edge deployment | DIY (Fly.io, etc.) | Built-in (Cloudflare) |
| Customizable cleaning rules | Yes (goquery selectors) | No |

---

## Recommendations

1. **Improve GitHub extraction** — detect README content specifically
2. **Add response caching** — in-memory or Redis cache for the HTTP server mode
3. **Deploy to edge** — Fly.io or similar for latency parity
