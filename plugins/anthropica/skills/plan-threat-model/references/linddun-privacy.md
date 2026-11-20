# LINDDUN privacy threats

Run over the same DFD when personal data is processed. Record threats in the same table with category prefix `P-` (e.g. `P-Li`). Privacy impact is rated for the data subject, not only the organisation.

## Categories and prompts
| Category | Question | Typical mitigations |
|---|---|---|
| **L**inking | Can records or actions be linked to each other across contexts (sessions, products, datasets)? | Separate identifiers per context, pseudonymous IDs, avoid joining analytics with identity |
| **I**dentifying | Can a person be identified from data believed anonymous (IP, device ID, quasi-identifiers like postcode + birth date + gender)? | Aggregation, k-anonymity thresholds, truncation, drop IPs after N days |
| **N**on-repudiation (privacy sense) | Is a user unable to deny an action when deniability matters (whistleblowing, sensitive browsing)? | Minimise per-user logging of sensitive actions, avoid signed receipts where not needed |
| **D**etecting | Can someone infer that a person is in a dataset or used a service (membership inference, error messages, timing)? | Uniform responses, no enumeration, access control on existence |
| **D**ata disclosure | Is more personal data collected, shared, retained or exposed than necessary? | Data minimisation, field-level encryption, retention limits, purpose-scoped access |
| **U**nawareness / unintervenability | Do users know what is collected and can they access, correct, export, delete or object? | Clear notices, consent records, self-service DSAR, deletion that reaches backups/replicas/vendors |
| **N**on-compliance | Does processing break law, policy or contract (lawful basis, cross-border transfer, retention, DPIA requirement)? | Records of processing, DPIA, transfer mechanisms, vendor DPAs, retention enforcement jobs |

## Per-element application
- External entity (data subject): Unawareness, Non-compliance.
- Data flow: Linking, Identifying, Detecting, Data disclosure (third-party SDKs, analytics beacons, logs).
- Data store: all categories; check retention and deletion paths explicitly.
- Process: Linking, Identifying, Data disclosure, Non-compliance (purpose creep).

## Checks that are often missed
- Analytics, error-tracking and session-replay SDKs sending form contents, URLs with tokens, or IPs to third parties.
- Logs, traces and data warehouses holding PII with longer retention than the primary store.
- Deletion requests not propagating to backups, search indexes, caches, ML training sets and vendors.
- Free-text fields (support tickets, notes, LLM prompts) collecting special-category data nobody planned for.
- Cross-border transfer via a vendor's sub-processors or region defaults.
- Children's data, location trails, biometrics and health data: treat as high impact by default.

## When to recommend a formal DPIA
Large-scale processing of special-category data, systematic monitoring of public spaces or behaviour, profiling with legal or similar effects, new technology on personal data (including AI decision-making), or combining datasets from different sources. State the trigger; the DPIA itself is owned by the privacy/legal function.
