---
name: database-engineering
description: Review or change SQL schemas, migrations, indexes, transactions and data-access code. Use when writing migrations, changing a schema, tuning queries or reviewing data-access code.
---

# Database Engineering

1. Identify database engine, migration tool and schema source of truth.
2. Inspect the relevant tables, relations, constraints and query paths.
3. Determine compatibility and data migration implications.
4. Prefer reversible, additive migrations.
5. Use parameterised queries and transactions where appropriate.
6. Check indexes against actual query patterns.
7. Add migration and repository tests where available.
8. Provide rollback considerations.

Never run destructive or production data operations without explicit approval.
