# Comprehensive Status Report — template-sqlc

**Date:** 2026-04-16 08:54
**Branch:** master (up to date with origin/master)
**Commits this session:** 8 (11a5df9 → 491884f)

---

## A. Brutally Honest Self-Reflection

### 1. What Did We Forget?

| # | Forgotten Issue | Impact |
|---|-----------------|--------|
| 1 | **`pkg/validation/` is NOT in git** — it's gitignored (`/pkg/` in `.gitignore`). Only `pkg/errors/errors.go` was force-added. `pkg/validation/validator.go` will be LOST on clone | 🔴 Critical — silent data loss |
| 2 | **No `cmd/` or `main.go`** — this is a template with no entry point. Nothing wires together monitoring, events, adapters, or servers | 🟡 By design (template), but limits testability |
| 3 | **`internal/monitoring/` is a FULL GHOST** — 17 Prometheus metrics, HTTP server, middleware — zero external callers. No code anywhere instantiates or uses it | 🔴 400+ lines of dead code |
| 4 | **Events are half-ghost** — `InMemoryEventPublisher` collects events into a slice, but no consumer processes them. Events vanish after publish. Only useful in tests | 🟡 Architecture incomplete |
| 5 | **Depguard references nonexistent `pkg/`** — `.golangci.yml:136` allows `github.com/LarsArtmann/template-sqlc/pkg` but most of pkg/ is gitignored | 🟡 Stale config |
| 6 | **`lint-sql` recipe is empty** — `justfile:86` has a recipe with no body | 🟢 Dead recipe |
| 7 | **No test coverage tooling** — justfile has no `-cover` commands, no coverage thresholds | 🟡 Quality gap |

### 2. What's Stupid That We Do Anyway?

| # | Stupid Thing | Why It's Stupid |
|---|-------------|----------------|
| 1 | **Two error systems** — `entities/errors.go` (6 types) + `pkg/errors/errors.go` (AppError with 26 codes). Neither integrates. Same name `ValidationError` in both | Callers must check both. Confusing. Split brain. |
| 2 | **Triple validation** — `pkg/validation/validator.go`, `adapters/validation/validation.go`, AND (formerly) `entities/errors.go:ValidateSearchQuery` all validate the same things with DIFFERENT limits | Different limits = bugs. maxSearchLimit=100 vs 1000, different reserved usernames |
| 3 | **Duplicate reserved usernames** — 4 entries in `entities/user.go` vs 36 entries in `pkg/validation/validator.go`. Neither references the other | "admin" blocked by one but not the other |
| 4 | **Duplicate email regex** — identical pattern compiled in `entities/user.go:73` AND `pkg/validation/validator.go:22` | DRY violation, wasted memory |
| 5 | **`pkg/` imports `internal/`** — `pkg/validation/validator.go` imports `internal/domain/entities` | Breaks Go convention. pkg/ should be independently importable |
| 6 | **`/pkg/` in `.gitignore`** — `pkg/validation/` is not tracked by git but exists on disk. Only `pkg/errors/` was force-added | Will lose code on clone |
| 7 | **`IDID` type name** — defined in `user.go` but is an event identifier. Should be `EventID` in an events file | Confusing naming, wrong file |
| 8 | **`fmt.Printf` in domain service** — 3 occurrences for "warning" logging. No structured logging | Anti-pattern in production code |
| 9 | **`fmt.Printf` in config/builder.go** — 2 occurrences for build output | Should use proper logging |

### 3. What Could We Have Done Better?

| # | Better Approach | What We Did Instead |
|---|----------------|---------------------|
| 1 | Unify error systems before adding more error types | Two parallel systems grew organically |
| 2 | Single validation package from the start | Three packages with different limits |
| 3 | Move `pkg/validation/` to `internal/validation/` first, then fix | Configured linters around the violation |
| 4 | Make monitoring `Used` before making it `Good` | 400+ lines of polished ghost code |
| 5 | Event system should have consumers, not just publishers | Events vanish into in-memory void |

### 4. What Can We Still Improve?

See Section F: Top 25 Things To Do Next.

### 5. Did We Lie?

**Yes, by omission.** The previous session's status report claimed `pkg/validation/` exists and needs moving. It does exist on disk, but it's NOT in git — it's gitignored. The report didn't mention this. Anyone cloning the repo would get a broken build because `pkg/validation/` is missing.

Also: the monitoring package was described as having "lint issues to fix" when it should have been flagged as a ghost system. We were polishing dead code.

### 6. How Can We Be Less Stupid?

- **If code has zero callers, DELETE IT or wire it IN.** No ghost systems.
- **Single source of truth** — one error system, one validation package, one reserved usernames list
- **Everything in git or it doesn't exist** — fix `.gitignore` or move files to `internal/`
- **Name things what they ARE** — `IDID` → `EventID`, move to events package
- **Don't configure linters around architectural violations** — fix the violations

### 7. Ghost Systems Found

| # | System | Location | Status | Action |
|---|--------|----------|--------|--------|
| 1 | **Monitoring** | `internal/monitoring/` | ☠️ **FULL GHOST** — 0 external callers | Wire in (add to cmd/) or DELETE |
| 2 | **Events (publisher only)** | `internal/domain/events/` | ⚠️ **HALF GHOST** — published but never consumed | Add consumer or simplify |
| 3 | **Depguard `pkg` reference** | `.golangci.yml:136` | ⚠️ **GHOST REFERENCE** — points to gitignored dir | Update after pkg→internal move |
| 4 | **`lint-sql` recipe** | `justfile:86` | ⚠️ **EMPTY RECIPE** — does nothing | Implement or remove |
| 5 | **`internal/tests/e2e/`** | Directory | ⚠️ **EMPTY DIR** — no tests | Add E2E tests or remove dir |
| 6 | **MySQL/Postgres session repos** | `adapters/mysql/`, `adapters/postgres/` | ❌ **MISSING** — only SQLite has stubs | Implement or document as intentional |

### 8. Scope Creep Check

We ARE at risk of scope creep. The original request was "improve quality and deduplicate." We've been:
- ✅ Fixing real lint issues (good)
- ✅ Removing dead code (good)
- ⚠️ Polishing ghost monitoring code (bad — should wire in or delete)
- ⚠️ Not addressing the root architectural problems (dual error systems, triple validation)

**Recommendation:** Focus on integrating ghost systems OR deleting them. Stop polishing unused code.

### 9. Did We Remove Anything Useful?

- ✅ `ValidateSearchQuery` from `entities/errors.go` — duplicates `adapters/validation/`. Correct removal.
- ✅ `newResourceError` — unused. Correct removal.
- ✅ `postgres/error_helper.go` — unused. Correct removal.
- ⚠️ `IDID` type still exists but is poorly named. Not removed yet.
- ⚠️ Monitoring code not removed despite being a ghost system.

### 10. Split Brains

| # | Split Brain | Details |
|---|-----------|---------|
| 1 | **Error systems** | `entities/errors.go` has `ValidationError` struct; `pkg/errors/errors.go` has `AppError` with `ErrCodeValidation`. Same concept, different types, no bridge |
| 2 | **Validation packages** | `pkg/validation/` vs `adapters/validation/` — both validate, different APIs, different limits |
| 3 | **Reserved usernames** | 4 entries in entities vs 36 in pkg/validation — will block different usernames |
| 4 | **Email regex** | Compiled in two packages — if one is updated, the other won't be |
| 5 | **pkg/ gitignored vs tracked** | `pkg/errors/` is tracked (force-added), `pkg/validation/` is not — inconsistent |
| 6 | **Events publish but nobody listens** | Publisher exists, consumer doesn't — events go into void |

### 11. Test Status

| Metric | Value |
|--------|-------|
| Test files total | 3 |
| Test packages with tests | 3 (unit, integration, bdd) |
| Packages with NO tests | 14 (entities, services, adapters/*, mappers, converters, monitoring, pkg/*) |
| E2E tests | 0 (empty directory) |
| Test coverage commands | None in justfile |
| Coverage threshold | None |
| `t.Parallel()` usage | Missing (thelper warning) |

**Only 3 of 17 packages have any tests.** This is a critical quality gap.

---

## B. Work Fully Done ✅

| # | Task | Commit | Impact |
|---|------|--------|--------|
| 1 | Configure depguard for hexagonal architecture | a4c1590 | Enforces layer separation |
| 2 | Delete unused postgres/error_helper.go | 450c56a | Removes dead code |
| 3 | Replace `any` field types with proper interfaces in adapters | 56f9a0a | Type safety |
| 4 | Remove duplicate ParseBool/FormatBool from mappers | 2a52a18 | DRY |
| 5 | Unify SQLite stubs to embed NotImplementedRepository | 09b0366 | Pattern consistency |
| 6 | Fix gosec G602 slice index out of range | 6e8e9e9 | Security bug fix |
| 7 | Configure ireturn/tagliatelle/wrapcheck/exhaustruct linters | 11a5df9 | ~50 lint false positives eliminated |
| 8 | Add ReadHeaderTimeout to http.Server (G112) | 3c48c4b | Security fix |
| 9 | Remove unused newResourceError, fix err113, remove duplicate ValidateSearchQuery | 7a67ce5 | Dead code + lint fixes |
| 10 | Fix nolintlint: remove unused G101 from directive | aa00617 | Lint hygiene |
| 11 | Fix gci import ordering | d1b59c7 | Formatting |
| 12 | Implement HTTP duration recording, fix godox TODO, fix varnamelen | b1998d7 | Dead code → real code |
| 13 | Remove unused nolint directive in BDD test | 491884f | Lint hygiene |

**Lint issues reduced: ~200 → 153** (23.5% reduction)

---

## C. Work Partially Done ⚠️

| # | Task | Status | Blocker |
|---|------|--------|---------|
| 1 | Move `pkg/validation/` to `internal/validation/` | Not started — most impactful architectural fix | Needs import updates across codebase |
| 2 | Unify reserved usernames | Not started — data inconsistency bug | Need to decide: 4 entries or 36? |
| 3 | Replace `fmt.Printf` with `slog` | Not started | 3 occurrences in user_service.go |
| 4 | Fix godoclint/revive doc comments | Not started | ~54 issues |
| 5 | Extract magic numbers | Not started | 14 mnd issues |
| 6 | Split newMetrics (funlen) | Not started | 119-line function |
| 7 | Wire monitoring into application | Not started — GHOST SYSTEM | No cmd/main exists |
| 8 | Fix varnamelen | Not started | 5 issues in production code |

---

## D. Not Started ❌

| # | Task | Priority |
|---|------|----------|
| 1 | Unify dual error systems (entities/errors.go vs pkg/errors/errors.go) | 🔴 Critical |
| 2 | Move pkg/validation to internal/validation | 🔴 Critical |
| 3 | Fix .gitignore: pkg/validation is not tracked | 🔴 Critical |
| 4 | Wire monitoring into application OR delete ghost system | 🔴 Critical |
| 5 | Add event consumers or simplify event system | 🟡 Important |
| 6 | Rename IDID → EventID, move to events package | 🟡 Important |
| 7 | Add missing MySQL/Postgres session repos | 🟡 Important |
| 8 | Add test coverage tooling to justfile | 🟡 Important |
| 9 | Add unit tests for 14 untested packages | 🟡 Important |
| 10 | Remove/update depguard ghost reference to pkg/ | 🟢 Minor |
| 11 | Remove empty lint-sql recipe or implement it | 🟢 Minor |
| 12 | Remove empty e2e test directory or add tests | 🟢 Minor |
| 13 | Fix containedctx in integration test | 🟢 Minor |
| 14 | Fix inamedparam on interface methods | 🟢 Minor |
| 15 | Fix noinlineerr in user_service.go | 🟢 Minor |
| 16 | Fix gosec G101 false positive | 🟢 Minor |
| 17 | Fix tagliatelle snake_case JSON tags (33 issues) | 🟢 Minor |
| 18 | Add structured logging (slog) | 🟡 Important |
| 19 | Consider adding: gin, koanf, lo, do, OTEL, uniflow/errors | 🟢 Enhancement |
| 20 | Add main.go / cmd/ to wire application together | 🟡 Important |

---

## E. Totally Fucked Up 💀

| # | Issue | Severity | Details |
|---|-------|----------|---------|
| 1 | **`pkg/validation/` not in git** | 🔴 CRITICAL | `/pkg/` in `.gitignore`. Clone → broken build. `pkg/errors/` was force-added but `pkg/validation/` was not. |
| 2 | **Dual error systems with same type names** | 🔴 CRITICAL | `entities.ValidationError` ≠ `pkg/errors.AppError(ErrCodeValidation)`. Callers must check both. No bridge. |
| 3 | **Reserved usernames: 4 vs 36** | 🔴 BUG | Different users pass validation in one place but fail in another. Silent data inconsistency. |
| 4 | **Monitoring is 400+ lines of ghost code** | 🔴 WASTE | Zero callers. Polished but unused. |
| 5 | **Events have no consumers** | 🟡 HALF-DONE | Events published to in-memory slice, then ignored. |

---

## F. Top 25 Things To Do Next (Sorted by Impact × Urgency ÷ Effort)

| # | Task | Impact | Urgency | Effort | Score |
|---|------|--------|---------|--------|-------|
| 1 | Fix .gitignore: move pkg/validation to internal/validation (fixes git tracking AND architecture) | 10 | 10 | 3 | **33** |
| 2 | Unify reserved usernames: single source of truth in entities | 9 | 9 | 2 | **40** |
| 3 | Deduplicate email regex: keep one in entities, reference from validation | 8 | 8 | 1 | **64** |
| 4 | Replace fmt.Printf with slog in user_service.go | 7 | 8 | 1 | **56** |
| 5 | Rename IDID → EventID, move to events package | 7 | 7 | 1 | **49** |
| 6 | Fix tagliatelle: convert snake_case JSON tags to camelCase (33 issues) | 6 | 7 | 2 | **21** |
| 7 | Wire monitoring into cmd/ OR delete it (resolve ghost) | 9 | 6 | 3 | **18** |
| 8 | Fix containedctx in integration test | 4 | 5 | 1 | **20** |
| 9 | Fix inamedparam: add named params to interface methods | 5 | 5 | 1 | **25** |
| 10 | Fix noinlineerr: use plain assignment in user_service.go | 4 | 5 | 1 | **20** |
| 11 | Fix godoclint: correct doc comment formats (4 issues) | 4 | 4 | 1 | **16** |
| 12 | Fix gosec G101 false positive with nolint comment | 3 | 4 | 1 | **12** |
| 13 | Remove/update depguard ghost reference to pkg/ | 5 | 4 | 1 | **20** |
| 14 | Fix revive: add missing doc comments (~50 issues) | 4 | 3 | 3 | **4** |
| 15 | Extract magic numbers to named constants (14 mnd issues) | 5 | 3 | 2 | **7.5** |
| 16 | Split newMetrics function (funlen) | 4 | 3 | 2 | **6** |
| 17 | Fix varnamelen: rename short variables (5 issues) | 3 | 3 | 2 | **4.5** |
| 18 | Add event consumer interface or document events as test-only | 7 | 3 | 3 | **7** |
| 19 | Add test coverage commands to justfile | 6 | 4 | 2 | **12** |
| 20 | Remove empty lint-sql recipe or implement it | 2 | 2 | 1 | **4** |
| 21 | Remove empty e2e test directory | 1 | 2 | 1 | **2** |
| 22 | Unify dual error systems: bridge entities errors and pkg/errors | 9 | 2 | 5 | **3.6** |
| 23 | Add samber/lo, samber/do, koanf dependencies | 6 | 2 | 4 | **3** |
| 24 | Add unit tests for untested packages | 7 | 3 | 8 | **2.6** |
| 25 | Create cmd/ with main.go wiring everything together | 8 | 2 | 5 | **3.2** |

---

## G. Top #1 Question I Cannot Figure Out Myself

**Should the monitoring package be wired into the application (add cmd/main.go, start metrics server), or should it be deleted/simplified?**

Arguments for wiring in:
- It's a sqlc template — monitoring SQL queries IS valuable
- The Prometheus metrics server is well-implemented
- It would make the template actually runnable

Arguments for deleting:
- No `cmd/` or `main.go` exists currently — this is a library/template, not an app
- 400+ lines of unused code is worse than no code
- Events also have no consumers — we'd need to wire those too
- OpenTelemetry is the future (code even says `// DEPRECATED: prefer go.opentelemetry.io/otel`)

**My recommendation:** Create a minimal `cmd/template-sqlc/main.go` that wires monitoring + a simple HTTP server using gin. This gives the template a real runnable demo. But this needs a user decision.

---

## H. Current Lint Breakdown (153 issues)

| Linter | Count | Category |
|--------|-------|----------|
| revive | 50 | Doc comments |
| tagliatelle | 33 | JSON tag naming |
| wrapcheck | 14 | Error wrapping |
| mnd | 14 | Magic numbers |
| exhaustruct | 5 | Struct initialization |
| varnamelen | 5 | Variable naming |
| testifylint | 5 | Test assertions |
| inamedparam | 6 | Interface params |
| noinlineerr | 3 | Error handling style |
| forbidigo | 3 | Forbidden functions |
| godoclint | 4 | Doc comment format |
| gochecknoglobals | 2 | Package-level vars |
| ireturn | 2 | Interface returns |
| gosec | 1 | Security |
| funlen | 1 | Function length |
| interfacebloat | 1 | Interface size |
| containedctx | 1 | Context in struct |
| err113 | 1 | Dynamic errors |
| thelper | 1 | Test helper |
| staticcheck | 1 | Static analysis |

---

## I. Architecture Diagram

```
                    ┌─────────────┐
                    │   cmd/      │ ← DOES NOT EXIST
                    │  (main.go)  │
                    └──────┬──────┘
                           │ (no wiring)
            ┌──────────────┼──────────────┐
            ▼              ▼              ▼
    ┌───────────────┐ ┌─────────┐ ┌──────────────┐
    │   domain/     │ │ pkg/    │ │  monitoring/  │
    │  (entities,   │ │ (errors,│ │  (17 metrics, │ ← GHOST: 0 callers
    │   services,   │ │  valid.)│ │   HTTP server)│
    │   repos,      │ │         │ │               │
    │   events)     │ │ ⚠️ VIOL │ │               │
    └───────┬───────┘ │ pkg/    │ └──────────────┘
            │         │ imports │
            │         │ intern. │
            ▼         └────┬────┘
    ┌───────────────┐      │
    │  adapters/    │◄─────┘ (should NOT depend on domain this way)
    │ (sqlite,mysql,│
    │  postgres,    │
    │  validation,  │
    │  mappers,     │
    │  converters)  │
    └───────┬───────┘
            │
            ▼
    ┌───────────────┐
    │  db/          │
    │ (sqlc gen'd)  │
    └───────────────┘
```

---

## J. Test Coverage Summary

| Package | Has Tests | Test Type |
|---------|-----------|-----------|
| `internal/domain/entities` | ✅ (external) | Unit (in tests/unit/) |
| `internal/domain/services` | ✅ (external) | Integration + BDD |
| `internal/adapters/*` | ❌ | None |
| `internal/adapters/mappers` | ❌ | None |
| `internal/adapters/converters` | ❌ | None |
| `internal/monitoring` | ❌ | None (ghost) |
| `pkg/errors` | ❌ | None |
| `pkg/validation` | ❌ | None |
| `internal/tests/e2e` | ❌ | Empty dir |

---

_End of Status Report_
