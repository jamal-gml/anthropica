# Framework notes

Idioms and gotchas only. Follow the project's existing patterns and versions first.

## React
- Function components and hooks. Rules of hooks: top level only, never conditional.
- Do not mirror props or derived data into state; compute during render (`useMemo` only if profiling shows cost).
- `useEffect` is for syncing with external systems (subscriptions, DOM APIs, non-React widgets). Data fetching belongs in the data library or framework loader; event-driven logic belongs in event handlers.
- Effects with async work must ignore results after unmount or param change (AbortController or an `ignore` flag); React 18+ StrictMode runs effects twice in development to surface missing cleanup.
- Keys: stable IDs. Changing a `key` intentionally resets component state.
- Context re-renders all consumers on value change: split contexts, memoise the value object, or use a store with selectors.
- Forms: controlled inputs for simple forms; a form library (React Hook Form, Formik, TanStack Form) with a schema (zod etc.) for complex ones; share the schema with the server where possible.
- Server state: TanStack Query, SWR, RTK Query or framework loaders; set `staleTime`, invalidate on mutation, use `queryKey` arrays including all params.
- Error boundaries around route segments and risky widgets; they do not catch errors in event handlers or async code.
- `dangerouslySetInnerHTML` only with sanitised input.

## Next.js (App Router)
- Components are Server Components by default; add `'use client'` only at the leaf that needs state, effects or browser APIs. Server-only code (DB, secrets) must never be imported into client components; use the `server-only` package to enforce.
- Fetch on the server where possible; understand the caching behaviour of `fetch` and route segment config for the project's Next.js version, and set `revalidate`/`cache` explicitly.
- Server Actions are public HTTP endpoints: validate input and check authorisation inside every action.
- `NEXT_PUBLIC_*` env vars are inlined into the client bundle.
- Use `next/image` (dimensions required), `next/font`, `next/link`; `loading.tsx` and `error.tsx` for route states; `generateMetadata` for SEO.
- Middleware runs on every matched request: keep it light; do not rely on it as the only authorisation check.

## Vue 3 / Nuxt
- Composition API with `<script setup>`; `ref` for primitives, `reactive` for objects; do not destructure reactive objects without `toRefs` (loses reactivity).
- `computed` for derived state; `watch` only for side effects; clean up in `onUnmounted` or via the watcher's cleanup callback.
- Pinia for shared client state; server data via Nuxt `useFetch`/`useAsyncData` (keys must be unique and stable) or TanStack Query for Vue.
- `v-html` only with sanitised input; `v-for` always with `:key`; never `v-if` and `v-for` on the same element.
- Nuxt: runtime config `public` keys ship to the client; server routes in `server/api` need their own validation and auth.

## Angular
- Standalone components, signals for local state where the version supports them; RxJS for streams. Always unsubscribe (`takeUntilDestroyed`, `async` pipe).
- `OnPush` change detection for performance; immutable inputs.
- Reactive forms with typed `FormGroup` and validators; display errors via `aria-describedby`.
- Angular sanitises `[innerHTML]` but `bypassSecurityTrust*` disables that: treat every use as a security review item.
- HTTP interceptors for auth headers and error mapping; `HttpClient` with typed responses.
- Lazy-load routes with `loadComponent`/`loadChildren`.

## Svelte / SvelteKit
- Svelte 5 runes (`$state`, `$derived`, `$effect`) or Svelte 4 stores depending on version; do not mix styles within a component.
- `{@html}` only with sanitised input.
- SvelteKit `load` functions: `+page.server.ts` for server-only data and secrets; `+page.ts` runs on both sides. Form actions need validation and authorisation like any endpoint.
- `$env/static/public` and `$env/dynamic/public` are exposed to the browser; `$env/*/private` are not.

## Styling
- Follow the existing approach (CSS modules, utility classes, CSS-in-JS, preprocessors); do not introduce a second one.
- Use design tokens/variables for colour, spacing, typography; no magic numbers (`build-design-system` owns tokens).
- Logical properties (`margin-inline`, `padding-block`) for RTL support; `rem` for type; container queries where component-level responsiveness is needed.
- Avoid `!important` and deep selectors into third-party components; scope overrides.

## Testing stack
- Unit/component: Vitest or Jest + Testing Library (`@testing-library/react|vue|angular|svelte`) + `@testing-library/user-event` + `jest-dom` matchers.
- Network: MSW handlers shared between tests and local development.
- E2E: Playwright (role-based locators `getByRole`, auto-waiting; avoid fixed sleeps); keep to critical journeys.
- Visual regression only for stable, shared components; review diffs, do not auto-accept.
- Avoid snapshot tests of large component trees; assert on behaviour and accessible output.
