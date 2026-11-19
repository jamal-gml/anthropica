---
name: seo-audit
description: Audit a website's technical SEO (crawlability, indexation, site architecture, Core Web Vitals, structured data, mobile, international) and produce prioritised fixes. Use when reviewing a site or template for search performance, before a migration or redesign, or when traffic drops.
---

# Technical SEO Audit

Work from evidence: page source, HTTP responses, sitemaps, robots.txt, Search Console or crawler exports. State which data sources were available.

## Checklist
**Crawl and index**
- `robots.txt` doesn't block important sections or assets (CSS and JS). XML sitemap is valid, lists only canonical 200 URLs, and is referenced in robots.txt.
- Status codes: no internal links to 3xx/4xx/5xx. Redirect chains collapsed to a single 301. Soft 404s identified.
- `noindex`, `canonical` and `hreflang` are consistent and not contradictory. Each indexable page self-canonicalises.
- Duplicate content handled (parameters, trailing slashes, http/https, www/non-www).

**Architecture and links**
- Important pages within about three clicks of the home page. Descriptive internal anchor text. No orphan pages.
- Clean, readable URLs. Breadcrumbs where hierarchy exists.

**Rendering and performance**
- Core content and links present in the rendered HTML (check JavaScript-rendered sites).
- Core Web Vitals (field data first): LCP ≤ 2.5 s, INP ≤ 200 ms, CLS ≤ 0.1. Identify the element or script responsible.
- Images: modern formats, explicit dimensions, lazy-loading below the fold, descriptive alt text.

**On-page templates**
- Unique `<title>` (about 50–60 chars) and meta description per page. One H1 and a logical heading order.
- Structured data (JSON-LD) valid for the page type (Organization, Product, Article, FAQ, Breadcrumb). Never mark up content not visible on the page.

**Mobile and international**
- Responsive, readable without zoom, tap targets spaced. `hreflang` pairs are reciprocal and include `x-default`.

**Migrations**
- 1:1 redirect map from old to new URLs, preserved metadata and structured data, updated sitemaps, and post-launch monitoring of indexation and rankings.

## Output
| Priority | Issue | Affected URLs or templates | Evidence | Fix | Expected impact |
|---|---|---|---|---|---|

Prioritise by impact × effort; template-level fixes first. Never promise specific rankings.
