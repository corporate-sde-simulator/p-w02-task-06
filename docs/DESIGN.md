# ADR-023: Secrets Caching â€” TTL vs On-Demand

**Date:**  
**Status:** Accepted  
**Authors:** Anil Sharma, Suresh Kumar

## Decision

Use **TTL-based local caching** for secrets with a default 5-minute expiry.

## Rationale

- Reduces vault API calls from every request to every 5 minutes per secret
- Acceptable staleness window given our rotation schedule (weekly)
- Local cache is faster than network call to vault

## Consequences

- Secret rotations take up to 5 minutes to propagate
- Must handle cache invalidation on explicit rotation events
- Memory footprint grows with number of secrets cached
