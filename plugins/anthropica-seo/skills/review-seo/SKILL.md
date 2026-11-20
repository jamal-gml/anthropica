---
name: review-seo
description: Audit technical SEO from evidence (headers, robots.txt, sitemaps, rendered HTML, crawl and search-console exports) - crawlability, indexation, canonicals, redirects, rendering, Core Web Vitals, structured data, hreflang, migrations. Use when reviewing a site or template, before/after a migration or redesign, or when organic traffic or indexing drops. Not for writing content (use `create-seo-content`).
---

# Technical SEO Audit

Find what stops search engines from crawling, rendering, indexing and understanding the site, prove each issue with evidence, and rank fixes by impact.

## Procedure
1. **Scope.** Domains/subdomains, page templates (home, category, product, article, search, faceted, paginated, localised), site size, CMS/framework and rendering mode (static, server-rendered, client-rendered), recent changes (releases, migrations, CMS updates) and the question being answered (traffic drop? pre-launch? general health?).
2. **Collect evidence** and record which sources were available: robots.txt, XML sitemaps, `curl` headers and HTML for a sample URL per template, rendered DOM, a full crawl export, search-console exports (performance, page-indexing reasons, sitemaps, Core Web Vitals, manual actions), field CWV data, server logs for large sites. Read `references/evidence-commands.md` for the exact commands.
3. **Triage a traffic drop first** if that's the trigger: confirm it's organic (not tracking/consent change), date it, compare against releases and known algorithm updates, split by template/directory/country/device, and check indexing counts and manual actions before anything else.
4. **Audit by area** using `references/checklist.md`: crawl and index, status codes and redirects, canonicals and duplicates, architecture and internal links, rendering, performance, on-page templates, structured data, mobile, international, security.
5. **Generalise to templates.** One broken product page is usually every product page: confirm on 3+ URLs of the same template and estimate the affected URL count from the crawl/sitemap.
6. **Prioritise** (Rules) and write fixes a developer can implement and verify.
7. **Migrations/redesigns:** read `references/migrations.md` when URLs, domains, CMS, rendering or site structure are changing.
8. **Define verification** for each fix: the command or report that will show it worked, and when to re-check.

## Rules
**Thresholds and facts that matter**
- Core Web Vitals, judged on field data at the 75th percentile over 28 days: LCP <= 2.5 s, INP <= 200 ms, CLS <= 0.1. Lab tools diagnose; field data decides. INP replaced FID.
- Redirects: permanent moves use 301/308, one hop. Chains >1 hop and 302s for permanent moves are findings. Internal links should point at final 200 URLs.
- XML sitemap: <= 50,000 URLs and <= 50 MB uncompressed per file (use an index beyond that); only canonical, indexable, 200 URLs; accurate `lastmod` or none.
- Important pages reachable within ~3 clicks of the home page; orphan pages (in sitemap/analytics, no internal links) are findings.

**Gotchas**
- `robots.txt` Disallow blocks crawling, not indexing: a blocked URL can still be indexed from links, and its `noindex` will never be seen. To deindex: allow crawl + `noindex`, then block later if needed.
- `noindex` + `canonical` to another URL send conflicting signals; pick one.
- Canonical is a hint: it is overridden when sitemap, internal links, hreflang and redirects point elsewhere. All signals must agree.
- Paginated pages self-canonicalise (page 2 is not canonical to page 1); keep crawlable `<a href>` links between pages.
- Faceted/filter URLs, internal search results, calendars and session IDs create crawl traps: control with link design, `noindex`, or robots rules depending on whether they have search value.
- JavaScript sites: links must be `<a href="...">`, not click handlers; critical content and metadata (title, canonical, robots, hreflang) should be in the initial server HTML. Compare raw HTML vs rendered DOM.
- Soft 404s: "not found" or empty pages returning 200. Error pages must return 404/410; temporary outages 503 with `Retry-After`.
- Non-HTML files (PDFs, images) are controlled with the `X-Robots-Tag` header, not meta tags.
- A leftover `noindex`, `Disallow: /` or basic-auth from staging is the classic launch-day disaster: check first after any release.
- hreflang: reciprocal, self-referencing, valid language(-region) codes (`en-GB`, not `en-UK`), pointing at canonical 200 URLs, plus `x-default`.
- Structured data must describe content visible on the page; invalid or misleading markup can lose rich results or trigger manual action.
- Large sites (roughly 100k+ URLs) or "discovered but not indexed" backlogs: analyse server logs for crawl waste before blaming content.

**Priority**
| Priority | Meaning | Examples |
|---|---|---|
| P0 | Blocks indexing/ranking site-wide or for a key template; fix now | Site-wide noindex, robots blocking key sections, 5xx on templates, canonicals to wrong domain, mass 404 after migration |
| P1 | Large measurable loss on important templates | Redirect chains on linked URLs, JS-only links/content, CWV failing on high-traffic templates, hreflang broken on main markets |
| P2 | Moderate or limited-scope | Duplicate titles, orphan pages, missing breadcrumbs schema, sitemap includes non-canonical URLs |
| P3 | Hygiene | Minor heading order, alt text gaps on decorative images |
Order within a priority by affected traffic x effort; template fixes before page-by-page fixes.

## Safety gates
- Read-only by default: audits use GET/HEAD requests at a polite rate (identify the crawler, respect robots.txt, throttle; don't crawl production at a rate that could degrade it).
- Ask before changing robots.txt, redirects, canonicals, sitemaps, noindex or DNS, submitting removal requests, or changing search-console settings - these can deindex a site quickly.
- Don't crawl behind login, or use credentials, unless the user provides scoped access for that purpose.
- Never promise rankings or traffic figures; describe expected impact qualitatively, with evidence.
- Label each finding as verified (evidence attached) or suspected (what would confirm it).

## Output
```
Scope: <domains, templates, size>   Trigger: <question>
Data sources: <robots, sitemaps, curl samples, rendered DOM, crawl export, search-console exports, CWV field data, logs> (missing: <...>)
Summary: <3-5 lines: biggest risks and what to do first>

| Priority | Issue | Affected templates / URL count | Evidence | Fix | Verify by | Expected impact |
|---|---|---|---|---|---|---|

Quick wins (< 1 day): ...
Needs developer/platform work: ...
Open questions / data needed: ...
Re-check plan: <what to re-run and when>
```

---

**Phase:** Review. **Next:** last phase; feed findings back into Build: `/anthropica-seo:build-frontend`. Phase map: `/anthropica-seo:workflow`.
