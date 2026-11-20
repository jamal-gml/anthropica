# Refreshing, merging and pruning content

Work from data: Search Console (or equivalent) performance by page and query over 16 months where available, analytics sessions and conversions, backlink counts, and a crawl list of URLs. State the date range used.

## Find candidates
| Signal | How to see it | Likely action |
|---|---|---|
| Clicks/impressions down >20-30% year over year (same season) | Performance by page, compare periods | Refresh |
| Position slipping from top 3 to 4-10 | Query report for the page | Refresh, improve snippet answer, add information gain |
| High impressions, low CTR vs similar positions | Query + page, CTR by position | Rewrite title/meta to match intent |
| Two+ URLs alternating for the same query | Query report -> pages; position volatility | Merge or differentiate (cannibalisation) |
| Zero clicks and near-zero impressions for 12 months, no links, no conversions | Page report + backlink export + analytics | Merge, improve or prune |
| Outdated facts, prices, screenshots, product names | Content review | Refresh |

Seasonality and core-update volatility cause swings; compare like-for-like periods before acting.

## Decision table
| Traffic/impressions | Backlinks | Conversions or business need | Content quality | Action |
|---|---|---|---|---|
| Declining | Any | Any | Fixable | Refresh in place (same URL) |
| Low | Yes | Any | Overlaps a stronger page | Merge into the stronger page, 301 the old URL |
| Low | No | Yes (sales, support, legal) | Any | Keep; improve or noindex if thin |
| Low | No | No | Thin or obsolete | Remove: 410/404, or 301 only if a truly equivalent page exists |
| Cannibalising pair | Any | Any | Both decent | Differentiate intents clearly, or merge into the best performer |

- Never 301 everything to the home page; irrelevant redirects are treated as soft 404s.
- Remove internal links to pruned URLs and drop them from the sitemap.
- Prune in batches, then watch impact for 4-8 weeks before the next batch.

## Refresh checklist
- [ ] Re-check the SERP: has intent or dominant format shifted?
- [ ] Update facts, statistics (with new sources and dates), prices, screenshots, product/UI names.
- [ ] Fill subtopic gaps current top results cover that the page misses; add new information gain.
- [ ] Tighten the intro so the answer appears in the first ~100 words.
- [ ] Improve title/meta if CTR is below peers at the same position.
- [ ] Fix broken outbound and internal links; add links from newer related pages.
- [ ] Update `dateModified` (schema) and the visible updated date only for material changes - cosmetic date bumps are misleading.
- [ ] Keep the same URL; if the slug must change, 301 and update internal links.
- [ ] Request re-crawl via the search console's URL inspection after publishing (optional, single URLs).
- [ ] Note the change date in a content log to read impact later.

## Merging checklist
1. Pick the surviving URL: best rankings/links/conversions, cleanest slug.
2. Combine the unique, still-accurate content; drop duplicates.
3. 301 each retired URL to the survivor (single hop).
4. Update internal links to point straight at the survivor; update sitemap.
5. Monitor queries from the retired URLs moving to the survivor over 2-6 weeks.

## Measuring a refresh
- Compare 28 days after vs 28 days before (and vs same period last year), excluding the first ~1-2 weeks of recrawl noise.
- Track clicks, impressions, average position for the target cluster, and conversions.
- Annotate core algorithm updates in the timeline; they can mask or mimic the effect.
