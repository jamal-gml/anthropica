---
name: backend-development
description: Implement backend services, business logic, API routes, validation and error handling.
disable-model-invocation: true
---

# Backend Development

# Backend Development

1. Trace the existing request path from route to service to persistence.
2. Preserve established dependency injection and error-handling patterns.
3. Validate input at trust boundaries.
4. Enforce authentication and resource-level authorisation.
5. Keep business logic out of transport handlers where the project separates them.
6. Add focused tests for success, invalid input and denied access.
7. Run targeted tests, lint and type checks.
8. Report API contract changes and affected consumers.

Do not add infrastructure or dependencies for a one-off requirement.
