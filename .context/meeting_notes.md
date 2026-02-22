# Meeting Notes — Sprint 24 Standup

**Date:** Feb 18, 2026  
**Attendees:** Anil (Security), Suresh, Intern

---

## Secrets Vault Client

- **Anil:** We're migrating to HashiCorp Vault. Need a Go client that all services can use. @Intern, PLATFORM-2864 is yours.

- **Suresh:** The client fetches secrets fine but the cache TTL comparison is backwards — it caches forever. Also the retry backoff multiplies by 1 instead of doubling.

## Action Items

- [ ] @Intern — Fix vault client (PLATFORM-2864)
- [ ] @Suresh — Integration test with real vault instance
