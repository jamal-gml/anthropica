# Evidence commands

All read-only. Throttle loops (`sleep 1`) on production. On Windows run these in Git Bash, or call `curl.exe` explicitly (in Windows PowerShell 5.1 `curl` is an alias for a different command). Replace `https://example.com` with the audited site.

## Headers, status and redirects
```bash
# Status and headers for one URL (look for status, Location, X-Robots-Tag, Cache-Control, Content-Type, Link rel=canonical)
curl -sI https://example.com/page

# Full redirect chain: one line per hop
curl -sIL https://example.com/old-page | grep -iE '^(HTTP/|location:)'

# Final status, final URL and hop count
curl -sIL -o /dev/null -w '%{http_code} %{num_redirects} %{url_effective}\n' https://example.com/old-page

# Protocol/host variants should all 301 once to the canonical origin
for u in http://example.com https://example.com http://www.example.com https://www.example.com; do
  curl -s -o /dev/null -w "$u -> %{http_code} %{redirect_url}\n" "$u"; done

# Soft 404 check: a nonsense URL must return 404/410, not 200
curl -s -o /dev/null -w '%{http_code}\n' https://example.com/this-should-not-exist-12345

# Response as a search crawler (compare to a browser UA to spot cloaking or bot blocking by a WAF/CDN)
curl -sI -A '<crawler user-agent string from the search engine docs>' https://example.com/page
```

## robots.txt
```bash
curl -s https://example.com/robots.txt
curl -s -o /dev/null -w '%{http_code}\n' https://example.com/robots.txt   # 5xx can halt crawling; 404 = everything allowed
```
Check: `Disallow: /` under `User-agent: *`; blocked CSS/JS or API paths needed for rendering; `Sitemap:` lines present and absolute; rules per subdomain (each host has its own robots.txt).

## Sitemaps
```bash
# List URLs from a sitemap (or sitemap index)
curl -s https://example.com/sitemap.xml | grep -o '<loc>[^<]*' | sed 's/<loc>//' > sitemap-urls.txt
wc -l sitemap-urls.txt

# Status of every sitemap URL (throttled); anything not 200 is a finding
while read -r u; do curl -s -o /dev/null -w "%{http_code} $u\n" "$u"; sleep 0.5; done < sitemap-urls.txt | grep -v '^200'

# Gzipped sitemaps
curl -s https://example.com/sitemap-1.xml.gz | gunzip | grep -c '<loc>'
```

## Indexing signals in HTML
```bash
curl -s https://example.com/page > raw.html
grep -ioE '<title>[^<]*</title>' raw.html
grep -ioE '<meta[^>]+name="?(robots|description)"?[^>]*>' raw.html
grep -ioE '<link[^>]+rel="?(canonical|alternate)"?[^>]*>' raw.html      # canonical and hreflang
grep -ioE '<h1[^>]*>.*</h1>' raw.html
grep -c '<a ' raw.html                                                   # crawlable links in raw HTML
grep -o '<script type="application/ld+json">.*</script>' raw.html        # JSON-LD (may span lines; open the file if empty)
```

## Rendered vs raw HTML (JavaScript sites)
```bash
# Rendered DOM with a headless browser
chromium --headless=new --dump-dom https://example.com/page > rendered.html
# or: npx playwright@latest ... (script that loads the page and prints page.content())

# Compare what exists only after rendering
diff <(grep -oE '<a [^>]*href="[^"]*"' raw.html | sort -u) <(grep -oE '<a [^>]*href="[^"]*"' rendered.html | sort -u)
```
Content, links, title, canonical or robots present only in `rendered.html` depend on rendering: higher risk, slower indexing.

## Performance
```bash
# Lab diagnostics (not ranking data)
npx lighthouse https://example.com/page --only-categories=performance,seo --form-factor=mobile --output=json --output-path=lh.json --chrome-flags="--headless"
```
Field data: the search console's Core Web Vitals report (grouped by template), the public field-data API (CrUX) or PageSpeed Insights for URL/origin p75 values. Report which source and date.

## Search console exports to request
| Export | Use |
|---|---|
| Performance: pages and queries, 16 months, by device/country | Locate drops by template, query, market; compare like-for-like periods |
| Page indexing report with reasons and example URLs | Why pages are excluded: noindex, redirect, 404, soft 404, duplicate/canonical chosen differently, crawled or discovered but not indexed |
| Sitemaps report | Submitted vs indexed per sitemap |
| URL inspection (single URLs) | Rendered HTML, chosen canonical, last crawl, blocked resources |
| Crawl stats | Response codes, host availability, crawl rate changes after releases |
| Core Web Vitals report | Failing URL groups by template |
| Manual actions and security issues | Always check during a drop |

## Crawl export (site crawler)
Columns to request: URL, status, indexability + reason, canonical, meta robots, X-Robots-Tag, title, title length, meta description, H1, word count, crawl depth, inlinks, outlinks, response time, hreflang, structured data errors. Crawl in "list mode" for sitemap URLs and old-URL lists; compare crawl vs sitemap vs analytics URLs to find orphans.

## Server logs (large sites)
```bash
# Hits by search crawlers per status code (combined log format; verify crawler IPs via reverse DNS before trusting the UA)
grep -i '<crawler UA token>' access.log | awk '{print $9}' | sort | uniq -c | sort -rn
# Most-crawled paths (spot parameter/facet waste)
grep -i '<crawler UA token>' access.log | awk '{print $7}' | sed 's/?.*//' | sort | uniq -c | sort -rn | head -50
```
