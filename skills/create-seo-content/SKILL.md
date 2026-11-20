---
name: create-seo-content
description: Plan, write and optimise search-led content - intent, SERP analysis, keyword clusters, briefs, E-E-A-T, on-page and schema, refreshing or pruning pages. Use when writing a content brief, article, comparison or landing page for search, planning topic clusters, or fixing cannibalisation or decaying pages. Not for technical audits (use `review-seo`) or ad/email copy (use `create-marketing-copy`).
---

# SEO Content

Create the page that best satisfies the searcher's intent, so it deserves to rank - written for people first, then made easy for search engines to understand.

## Procedure
1. **Define the job.** Page purpose, audience, business goal (lead, sale, signup, support deflection), primary query, and the one action the reader should take next.
2. **Classify intent** (table in Rules). If the results page mixes intents, pick the dominant one or plan separate pages.
3. **Analyse the SERP** for the primary query (top ~10 results plus features). Record: page types ranking (guide, list, product, category, tool, video, forum), format and angle, depth and sections covered, freshness, SERP features (featured snippet, people-also-ask, video, images, local pack), and what is missing or weak. If no live SERP or tool data is available, say so and mark everything inferred as an assumption.
4. **Cluster keywords.** Group queries that the same page should target: if two queries share roughly 3+ of the same top-10 URLs, they belong on one page; if not, separate pages. Map each cluster to one existing or new URL. Check the site first (`site:` search, Search Console queries by page) so you don't create a competitor to an existing page.
5. **Write the brief.** Read `references/content-brief.md` for the full template and SERP worksheet. The angle must state the information gain: what this page adds that current results lack (original data, first-hand testing, expert input, examples, tools, better structure).
6. **Draft.** Answer the core question within the first ~100 words, then expand. Use headings that match real sub-questions (people-also-ask, support tickets, sales calls, forums). Cover entities and subtopics the intent requires; skip filler sections that exist only because competitors have them.
7. **On-page optimisation** (Rules). Write title, meta description, H1, slug, image alt text, schema.
8. **Internal links.** Add contextual links from 3-5 relevant, well-linked existing pages to the new page, and from the new page to its pillar and siblings. Descriptive anchors.
9. **Helpful-content check** (Rules) and claims check. Fix before handing over.
10. **Maintenance.** Set a review date. Read `references/refresh-and-prune.md` when refreshing decaying pages, merging cannibalising pages or pruning thin content.

## Rules
**Intent**
| Intent | Query signals | Page that wins |
|---|---|---|
| Informational | how, what, why, guide, examples, ideas | Guide, explainer, tutorial, definition with examples |
| Commercial investigation | best, top, vs, review, alternatives, comparison | Comparison, list with criteria, review based on real use |
| Transactional | buy, price, discount, near me, template, download, signup | Product, category, pricing, tool or download page |
| Navigational | brand or product name, login | The official page; don't fight for other brands' navigational terms |
| Local | city/area names, near me | Location page with real address, hours, reviews |
Match the format the SERP rewards (a product query is rarely won by a blog post); differentiate on usefulness, not on format.

**On-page**
- Title tag: unique, primary topic near the start, ~50-60 chars, specific benefit or qualifier (year only if the content is genuinely maintained for it), brand at the end if space.
- Meta description: ~120-160 chars, what the reader gets and why click; not a keyword list. Search engines may rewrite it.
- One H1 that matches the intent (can differ from the title). H2/H3 in logical order; headings readable as a summary of the page.
- Slug: short, lowercase, hyphenated, no dates or stop-word strings; never change a ranking URL without a 301.
- Natural vocabulary of the topic (synonyms, entities, related terms). No keyword density targets, no stuffing, no hidden text, no exact-match anchor spam.
- Featured-snippet friendly: a 40-60 word direct answer under a question heading, or a list/table where the SERP shows one.
- Images: original where possible, descriptive filename and alt text (describe the image, don't stuff keywords), compressed, with dimensions.
- Outbound links to primary sources for claims; cite data with date.
- No word-count targets: length follows intent. Cover the topic completely, then stop.

**Schema (JSON-LD, only for content visible on the page)**
| Page | Type |
|---|---|
| Article/blog | `Article` or `BlogPosting` with author, datePublished, dateModified |
| Product | `Product` with `Offer`, real `AggregateRating` only if reviews are on-page |
| Organisation/about | `Organization` with logo and sameAs |
| Navigation | `BreadcrumbList` |
| Video | `VideoObject` |
| FAQ / HowTo | Only if the page really is that; major engines have restricted these rich results, so don't promise them |
Validate with a rich-results or schema validator before publishing.

**E-E-A-T and trust**
- Experience: first-hand detail - original photos, screenshots, test results, "we tried X for 3 months".
- Expertise: named author with relevant credentials and an author page; expert reviewer for health, finance, legal, safety (YMYL) topics.
- Authoritativeness: cite primary sources; earn links with original data rather than asking for them.
- Trust: accurate, dated ("Updated <date>" only when content materially changed), transparent about affiliate links and sponsorship, easy contact and about pages.

**Helpful-content check (all must be yes)**
- Would the target reader be satisfied without going back to search?
- Does it add information beyond what already ranks?
- Is it written for this site's audience, not only for search traffic outside its expertise?
- Is every claim sourced or from verifiable first-hand experience?
- Would you publish it if search engines didn't exist?
- Not mass-produced: no near-duplicate city/keyword-swap pages, no doorway pages, no auto-generated text without human editing and unique value.

## Safety gates
- Never invent search volumes, difficulty scores, rankings, traffic numbers or competitor data; use tool exports the user provides or label estimates.
- Never fabricate expertise: no fake authors, credentials, reviews, test results, quotes or "we tested" claims that didn't happen.
- Never promise rankings, positions or traffic outcomes.
- Ask before publishing, changing live URLs, deleting/noindexing pages, or bulk-editing titles across a site.
- YMYL topics: flag the need for qualified expert review before publication.

## Output
```
Primary query: <query>  Intent: <type>  Target URL: <new or existing>
Cluster: <secondary queries>  Cannibalisation check: <result>
SERP summary: <dominant format, angle, gaps>  Data sources: <tools/exports or "none - assumptions">
Information gain: <what this page adds>

Title (<n> chars): ...
Meta description (<n> chars): ...
H1: ...   Slug: /...
Outline: H2/H3 list with the question each answers
<Draft or optimised content>
Internal links: from <page> anchor "<text>" -> this page; this page -> <pages>
Schema: <type, key properties>
E-E-A-T: author, reviewer, sources, first-hand evidence
[VERIFY] items and assumptions
Review date: <date>
```
