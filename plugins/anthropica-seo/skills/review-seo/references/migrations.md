# Migrations and redesigns

Applies to domain moves, HTTPS/host changes, URL restructures, CMS/platform changes, rendering changes (SSR <-> CSR), site merges and redesigns. Any of these can cut organic traffic for weeks or permanently if signals are lost.

## Before (weeks ahead)
1. **Benchmark:** export 16 months of search-console performance by page and query, analytics landing-page sessions/conversions, indexed counts, top-linked URLs (backlink export), and CWV per template. Save them; they are the comparison baseline.
2. **URL inventory:** union of crawl of the live site, all XML sitemaps, analytics landing pages (12+ months), search-console pages with impressions, and URLs with external backlinks. Dedupe.
3. **Redirect map:** every old URL -> its most relevant new URL, 1:1 where possible, as a spreadsheet (old, new, type, owner). Rules:
   - 301/308 single hop; no chains with existing redirects (flatten old redirects to point at new final URLs).
   - Removed content without an equivalent -> 404/410, not the home page.
   - Pattern rules (regex) tested against the full inventory, not samples.
   - Keep parameters that change content; drop tracking parameters.
4. **Parity check on staging:** titles, meta descriptions, H1s, canonicals, hreflang, structured data, internal links, body content and indexability match or improve on old pages. Crawl staging and diff against the old-site crawl.
5. **Staging protection:** block staging with authentication (not only robots.txt/noindex, which tend to leak to production).
6. **Rendering:** if moving to a JavaScript framework, confirm content and links exist in server HTML; compare raw vs rendered.
7. **Plan the launch window:** low-traffic period, not before a peak season; team available for 2 weeks after; rollback plan.
8. **Search console:** verify the new property (and domain property) before launch.

## Launch day
- [ ] Remove staging auth/noindex/`Disallow: /` on the new production site; confirm robots.txt and meta robots on key templates.
- [ ] Deploy redirects; run the full old-URL list through curl (status + final URL) and fix any non-301 or wrong destinations.
- [ ] Canonicals on new pages point at new URLs (not old, not staging host).
- [ ] New XML sitemaps live and submitted; optionally keep a temporary sitemap of old URLs to speed discovery of redirects.
- [ ] Domain move: use the search console's change-of-address tool; keep the old domain registered and redirecting.
- [ ] Update internal links, navigation, hreflang, structured data and assets (images, PDFs) to new URLs - don't rely on redirects internally.
- [ ] Update external profiles and key backlinks you control; analytics and tag configuration on the new site.

```bash
# Verify the redirect map: old URL -> expected new URL (tab-separated file: old<TAB>expected)
while IFS=$'\t' read -r old expected; do
  read -r code final < <(curl -s -o /dev/null -L -w '%{http_code} %{url_effective}\n' "$old")
  [ "$final" = "$expected" ] && [ "$code" = "200" ] || echo "FAIL $old -> $final ($code), expected $expected"
  sleep 0.3
done < redirect-map.tsv
```

## After (8-12 weeks)
- Daily for the first two weeks: 404s and 5xx (crawl stats, logs), redirect errors, page-indexing reasons, sitemap indexed counts, top pages' traffic.
- Weekly: organic clicks/impressions vs benchmark by template, old URLs dropping out and new URLs being indexed, rankings for priority queries, CWV.
- Fix new 404s with real backlinks or traffic by adding redirects.
- Expect temporary fluctuation; a sustained drop beyond ~4-6 weeks needs investigation (lost content, changed internal linking, rendering, canonical mistakes).
- Keep redirects for at least a year, ideally permanently; keep the old domain.

## Common failure modes
| Symptom | Likely cause |
|---|---|
| Sharp drop, pages "excluded by noindex" | Staging noindex shipped |
| Old URLs still indexed weeks later | Redirects 302, chained, or blocked by robots.txt |
| New pages "duplicate, different canonical chosen" | Canonicals point at old/staging URLs, or content duplicated across hosts |
| Category traffic lost | Internal linking/navigation changed; pages now deeper or orphaned |
| Rich results disappeared | Structured data not carried over |
| International traffic mixed up | hreflang not updated to new URLs |
| Many soft 404s | Old URLs redirected to the home page or generic category |
