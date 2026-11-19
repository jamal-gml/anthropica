---
name: security-review
description: Review application changes for authentication, authorisation, input validation, secrets, injection and common security weaknesses. Use when a change touches authentication, user input, secrets, file handling or external integrations.
---

# Security Review

## Review scope
- Authentication and session handling
- Authorisation and tenant/resource boundaries
- Input validation, output encoding and injection
- Secrets and sensitive logging
- File uploads and path handling
- CORS, CSRF, rate limiting and security headers where applicable
- Dependency and configuration risks

## Method
1. Establish the changed attack surface.
2. Trace trust boundaries and privilege checks.
3. Verify findings against code and configuration.
4. Classify severity and explain realistic impact.
5. Recommend minimal remediations.
6. Re-run relevant tests or scanners if available.

Do not claim a system is secure based on a limited review. Do not exploit production systems.
