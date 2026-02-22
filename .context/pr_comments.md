# PR #415 Review — Vault Client (by Suresh Kumar)

## Reviewer: Anil Sharma — Feb 16, 2026

---

### `vaultClient.go`

> **Line 48** — Cache TTL check:  
> `time.Since(entry.cachedAt) < cs.ttl` — this means the entry is always considered fresh. Should be `>` to detect expiry.

> **Line 72** — Retry backoff:  
> `backoff = backoff * 1` — backoff never increases. Should multiply by 2 for exponential backoff.

### `secretCache.go`

> Clean implementation. The mutex usage is correct.

---

**Suresh Kumar** — Feb 17, 2026

> Both are silly bugs. The TTL one means secrets never refresh, which is a security risk.
