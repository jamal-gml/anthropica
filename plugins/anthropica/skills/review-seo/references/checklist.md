# Technical SEO checklist

Check each item on at least one URL per template, then confirm scale from the crawl. Record evidence (command output, export row, screenshot) next to every failed item.

## Crawl and index
- [ ] robots.txt returns 200 (or 404), does not block important sections, CSS, JS, images or APIs needed for rendering; lists absolute `Sitemap:` URLs.
- [ ] No unintended site-wide or template-wide `noindex` (meta or `X-Robots-Tag`); no basic-auth or IP block on production.
- [ ] XML sitemaps valid, within size limits, list only canonical, indexable 200 URLs; split by type (pages, products, articles) to read indexing per type; submitted in the search console.
- [ ] Submitted vs indexed ratio per sitemap investigated where low.
- [ ] Page-indexing exclusion reasons reviewed; large "crawled/discovered but not indexed" buckets explained (quality, duplication, crawl budget).
- [ ] Indexed page count roughly matches the intended indexable count (no index bloat from parameters, tags, search pages).

## Status codes and redirects
- [ ] No internal links to 3xx, 4xx or 5xx URLs.
- [ ] Permanent redirects are 301/308, single hop, to the relevant equivalent (not all to the home page).
- [ ] http->https and host (www/non-www) normalised in one hop; trailing slash and case handled consistently.
- [ ] Missing pages return 404/410; no soft 404s (200 with "not found" or empty listings).
- [ ] Out-of-stock/expired items handled deliberately (keep with alternatives, 301 to successor, or 404/410 if gone for good).
- [ ] 5xx and timeouts rare in crawl stats; maintenance uses 503 + `Retry-After`.

## Canonicals and duplicates
- [ ] Every indexable page has an absolute, self-referencing canonical to its preferred 200 URL.
- [ ] Canonical, sitemap, internal links, hreflang and redirects all point at the same URL version.
- [ ] Parameter variants (sort, filter, tracking, session) canonicalise to the clean URL or are kept out of crawl.
- [ ] No canonicals to redirected, noindexed, 404 or cross-domain URLs unless intended (syndication).
- [ ] Near-duplicate templates (thin category pages, tag pages, printer versions) consolidated or noindexed.
- [ ] Search-engine-chosen canonical matches the declared one (URL inspection) for key pages.

## Architecture and internal links
- [ ] Key pages within ~3 clicks of the home page; crawl depth distribution reviewed.
- [ ] No orphans: URLs in sitemap/analytics/backlink reports with zero internal inlinks.
- [ ] Navigation and related-content links are crawlable `<a href>` with descriptive anchor text.
- [ ] Breadcrumbs reflect hierarchy (and `BreadcrumbList` schema).
- [ ] Pagination uses crawlable links; paginated pages self-canonical.
- [ ] Faceted navigation: only facets with search demand are indexable; others noindex or not linked crawlably; no infinite combinations.
- [ ] URLs readable, lowercase, hyphenated, stable; no session IDs.
- [ ] Nofollow on internal links not used to "sculpt" (it wastes link equity).

## Rendering (JavaScript)
- [ ] Main content, links, title, meta robots, canonical and hreflang present in raw server HTML (SSR/SSG/prerender) or reliably after render.
- [ ] Rendered DOM matches what users see; no content behind interactions (tabs that load on click, infinite scroll without paginated URLs).
- [ ] No resources needed for render blocked by robots.txt; no render-time errors in the console.
- [ ] Client-side routes return correct status codes server-side (a missing route returns 404, not a 200 shell).
- [ ] Lazy-loaded content uses native lazy-loading or intersection-based loading that works without scroll events.

## Performance (Core Web Vitals)
- [ ] Field p75 per template: LCP <= 2.5 s, INP <= 200 ms, CLS <= 0.1 (mobile first).
- [ ] LCP: element identified; image preloaded/prioritised, not lazy-loaded, served in modern format at the right size; TTFB under ~0.8 s; render-blocking CSS/JS minimised.
- [ ] INP: long tasks (>50 ms) on interaction identified; heavy third-party scripts deferred; event handlers yield to the main thread.
- [ ] CLS: images/video/ads/embeds have reserved dimensions; fonts use `font-display` with fallback metrics; no content injected above existing content.
- [ ] Caching and compression (Brotli/gzip), HTTP/2 or HTTP/3, CDN for static assets.
- [ ] Hand deep performance work to `test-performance` / frontend owners.

## On-page templates
- [ ] Unique, descriptive `<title>` (~50-60 chars) and meta description per indexable page; no boilerplate duplicates across the template.
- [ ] One H1 describing the page; logical heading hierarchy.
- [ ] Image alt text meaningful for content images, empty for decorative.
- [ ] Main content not thin or duplicated across many URLs (template text outweighing unique text).
- [ ] Hand content quality issues to `create-seo-content`.

## Structured data
- [ ] JSON-LD valid (schema validator, rich-results test) with no errors on key templates.
- [ ] Types match the page (`Organization`, `WebSite`, `BreadcrumbList`, `Product`+`Offer`, `Article`, `VideoObject`, `LocalBusiness`, `Event`, `JobPosting`).
- [ ] Marked-up values match visible content (price, availability, ratings, dates); no self-serving or fabricated reviews.
- [ ] Rich-result eligibility changes acknowledged (FAQ/HowTo rich results are restricted); don't recommend markup only for a feature that no longer shows.

## Mobile
- [ ] Responsive, `<meta name="viewport" content="width=device-width, initial-scale=1">`.
- [ ] Mobile and desktop versions have the same content, links, structured data and metadata (mobile-first indexing).
- [ ] Text readable without zoom, tap targets spaced, no intrusive interstitials on entry.

## International
- [ ] hreflang via HTML, header or sitemap (one method), reciprocal, self-referencing, valid codes, `x-default` set.
- [ ] hreflang targets are canonical, indexable 200 URLs.
- [ ] No automatic IP/language redirects that hide versions from crawlers; offer a switcher instead.
- [ ] Localised content actually translated (not just currency swap) and canonicalised per locale.
- [ ] URL structure consistent (ccTLD, subdomain or subdirectory) with clear targeting.

## Security and hygiene
- [ ] HTTPS everywhere, valid certificate, no mixed content; HSTS where appropriate.
- [ ] No hacked/spam pages indexed (`site:` searches for spammy terms, security issues report).
- [ ] Staging/dev hosts not indexable (auth, not just robots.txt).
- [ ] Manual actions report clean.
