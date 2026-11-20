# Channel limits and truncation points

Limits change and differ by platform, placement and device. Treat these as typical working values; check the current platform spec before final delivery, and front-load the message before the truncation point. Count characters, not words; emoji and non-Latin scripts may count double.

## Email
| Field | Hard limit | Visible / recommended |
|---|---|---|
| Subject line | None practical | 30-50 chars (mobile shows ~30-40); key words first |
| Preheader / preview text | None practical | 40-100 chars; complements the subject, never repeats it; pad to stop body text leaking in |
| From name | - | Recognisable, under ~25 chars |
| Body | - | One idea; primary CTA visible without scrolling on mobile; works with images off |

## Search ads (responsive text ads, common format)
| Field | Limit |
|---|---|
| Headline | 30 chars each; supply many distinct ones (they are combined automatically) |
| Description | 90 chars each |
| Display URL path | 15 chars per path segment |
- Every headline must make sense next to any other (no headline that depends on another). Pin only when legally required.
- Include the core keyword theme in at least a few headlines; include a CTA and a proof point among the set.

## Social and feed ads
| Field | Typical guidance |
|---|---|
| Feed ad primary text | ~125 chars visible before "more"; hard limits much higher |
| Feed ad headline | ~25-40 chars to avoid truncation |
| Feed ad link description | ~25-30 chars; often hidden on mobile |
| Professional-network sponsored intro text | ~150 chars visible; hard limit ~600 |
| Professional-network sponsored headline | ~70 chars recommended; hard limit ~200 |
| Video ads | Hook in first 1-3 seconds; design for sound off with captions |

## Organic social
| Format | Hard limit | Visible before truncation |
|---|---|---|
| Microblog post | 280 chars (standard accounts) | Whole post; links consume a fixed count |
| Professional-network post | ~3,000 chars | ~140-210 chars before "see more" |
| Photo-sharing caption | ~2,200 chars | ~125 chars; max ~30 hashtags, 3-5 is better |
| Short-video caption | Varies (thousands) | ~100 chars; on-screen text matters more |
| Video title (long-form video sites) | ~100 chars | ~60-70 chars |

## SMS and push
| Format | Limit |
|---|---|
| SMS (GSM-7 alphabet) | 160 chars per segment; 153 when concatenated |
| SMS (any emoji or non-GSM char -> UCS-2) | 70 chars per segment; 67 when concatenated |
| SMS must include | Sender identity, opt-out instruction (for example "Reply STOP to opt out") where required |
| Push title | ~30-50 chars visible |
| Push body | ~90-150 chars visible depending on OS and lock screen |

One emoji or curly quote can switch an SMS to UCS-2 and more than double its cost; use straight quotes.

## App stores
| Field | Typical limit |
|---|---|
| App name/title | 30 chars |
| Subtitle / short description | 30-80 chars |
| Keyword field (where offered) | 100 chars, comma-separated, no spaces |
| Long description | ~4,000 chars; first 1-3 lines visible before "more" |

## Web snippets
| Field | Guidance |
|---|---|
| Title tag | ~50-60 chars (pixel-width truncation, ~600 px) |
| Meta description | ~120-160 chars; may be rewritten by the search engine |
| Open Graph title / description | ~60 / ~110 chars |
| Button/CTA text | 2-5 words |
| Display/banner ad | 5-8 words of headline max; logo, offer, CTA |
