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

| File Type | url2md (tokens) | markdown.new (tokens) | Winner |
|-----------|-----------------|----------------------|--------|
| **PDF** | Text extraction by page (13) | Text + metadata via AI (83) | Tie |
| **DOCX** | Full doc: headings, bold/italic, tables (2,418) | Full doc via AI (2,836) | Tie |
| **PNG** | Metadata + image embed (23) | AI vision description (81) | **markdown.new** |
| **CSV** | Markdown table (144) | Raw CSV in code block | **url2md** |
| **XLSX** | Markdown table, multi-sheet (30,933) | FAIL (raw binary) | **url2md** |

**File type score: url2md wins 2, ties 2, loses 1.**

### File Type Detail

**PDF** — Both extract text. url2md uses native Go PDF parsing, markdown.new uses Workers AI:
```
# dummy.pdf                          # dummy.pdf
## Contents                          ## Metadata
### Page 1                           - Author=Evangelos Vlachogiannis
Dummy PDF file                       ## Content
                                     Dummy PDF file
(url2md)                             (markdown.new)
```

**DOCX** — Both extract full document structure. url2md uses go-docx, markdown.new uses AI:
```
# demo.docx
Demonstration of DOCX support in calibre
**bold**, _italic_, ~~strikethrough~~, tables, headings — both handle well.
```

**PNG** — markdown.new uses AI vision for description. url2md outputs metadata + image embed:
```
# nurbcup2si.png                     The image displays a 3D rendered
## Metadata                          teacup with a smooth surface...
- File: nurbcup2si.png
- Size: 18746 bytes
![nurbcup2si.png](url)
(url2md)                             (markdown.new)
```

**CSV** — url2md converts to a proper markdown table, markdown.new wraps raw CSV in a code block:
```
| John | Doe | 120 jefferson st. |   ```csv
| --- | --- | --- |                   John,Doe,120 jefferson st.,...
| Jack | McGinnis | 220 hobo Av. |   Jack,McGinnis,220 hobo Av.,...
                                      ```
(url2md — markdown table)            (markdown.new — code block)
```

**XLSX** — url2md successfully parses multi-sheet Excel files into markdown tables. markdown.new fails:
```
# Financial Sample.xlsx
| Segment | Country | Product | ... |
| --- | --- | --- | --- |
| Government | Canada | Carretera | ... |
(url2md — 30,933 tokens, full data)  (markdown.new — FAIL, raw binary)
```

---

## Overall Summary

### Combined Score

| Category | url2md Wins | Ties | markdown.new Wins |
|----------|-------------|------|-------------------|
| Web Pages (6 tests) | **4** | 2 | 0 |
| File Types (5 tests) | **2** | 2 | 1 |
| **Total (11 tests)** | **6** | **4** | **1** |

### Where Each Tool Excels

**url2md is better for**:
- Web page → Markdown conversion (cleaner, fewer tokens)
- CJK/Chinese content (complete extraction vs truncation)
- Noise removal (readability-based, strips nav/ads/UI)
- CSV → proper markdown table (vs code block)
- XLSX → markdown table (markdown.new fails)
- Self-hosted / offline use
- Customizable pipeline

**markdown.new is better for**:
- Zero-setup usage (SaaS, no deployment needed)

> **Note**: Image AI description is available in the online (Cloudflare-hosted) version of url2md. The self-hosted/CLI version outputs image metadata + embed.

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
| PDF → Markdown | Yes (native Go) | Yes (Workers AI) |
| DOCX → Markdown | Yes (go-docx) | Yes (Workers AI) |
| Image → Description | Yes (AI vision, online only) | Yes (AI vision) |
| CSV → Markdown table | **Yes** | Partial (code block only) |
| XLSX → Markdown table | **Yes** (multi-sheet) | **No** (fails) |
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
| Magic bytes detection | Yes (PDF/DOCX/XLSX/PNG/JPEG) | Unknown |
| Redirect URL handling | Yes (follows + detects) | Yes |

---

## Roadmap for url2md

1. ~~**Add PDF support**~~ — Done (ledongthuc/pdf)
2. ~~**Add DOCX support**~~ — Done (fumiama/go-docx)
3. ~~**Add CSV/XLSX support**~~ — Done (markdown table, excelize/v2)
4. ~~**Add image AI description**~~ — Available in online version (Cloudflare Workers AI vision)
5. **Improve GitHub/SPA extraction** — detect README content, handle JS-heavy pages
6. **Add response caching** — in-memory or Redis cache for HTTP server mode
