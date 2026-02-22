# PLATFORM-2864: Refactor secrets vault client caching layer

**Status:** In Progress · **Priority:** Medium
**Sprint:** Sprint 24 · **Story Points:** 5
**Reporter:** Suresh Kumar (Infra Lead) · **Assignee:** You (Intern)
**Due:** End of sprint (Friday)
**Labels:** `backend`, `infrastructure`, `golang`, `refactor`
**Task Type:** Code Maintenance

---

## Description

The secrets vault client **works correctly** — all tests pass. The code needs cleanup: the cache has no size limit (memory leak risk), error handling swallows errors silently, and the locking strategy is too broad.

## What Needs Improvement

- `// TODO (code review):` comments mark specific issues
- Cache grows without bound — add max size with LRU eviction
- Errors logged but not propagated (caller never knows fetch failed)
- Single mutex locks entire cache during any operation
- Cache TTL is hardcoded magic number
- No metrics/observability

## Acceptance Criteria

- [ ] All `// TODO (code review):` items addressed
- [ ] Cache limited to configurable max size
- [ ] Errors propagated appropriately
- [ ] Cache TTL configurable via constructor
- [ ] All existing tests still pass
