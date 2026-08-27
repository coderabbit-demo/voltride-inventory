# AGENTS.md — voltride-inventory

Part of VoltRide, a multi-repo microservices demo (see the `voltride-platform` repo for the system map). This service is a **leaf**: it calls no other service, but three sibling repos hand-maintain local copies of its wire format. There is **no shared types package anywhere in VoltRide** — contracts are duplicated per repo on purpose, and nothing must ever change that.

## Who consumes this service's contracts

| Contract | Consumer repo | Consumer file | Failure mode if changed |
|---|---|---|---|
| Stock record (`productId`, `stockCount`, `warehouse`, `restockEtaDays`) | voltride-catalog | `src/types.ts` + `src/clients/inventoryClient.ts` | `stockCount` reads `undefined` → every product shows "Backordered" |
| Stock record | voltride-pricing | `inventory_client.py` | `data.get("stockCount", 0)` → silently 0 → corrupt discounts/surcharges |
| Stock record | voltride-orders | `clients.go` (`StockRecord`) | Go decodes missing keys as 0 **silently** → inflated delivery estimates |
| Reservation response (`status: "reserved"` exact string) | voltride-orders | `clients.go` (`reserveStock` guard) | any other string → every checkout rolls back and 502s |
| 409 body (`error`, `productId`, `requested`, `stockCount`) | voltride-orders | `clients.go` (`InsufficientStock`) | broken out-of-stock UX |

**Changing any JSON tag, field type/unit, status string, endpoint path, or status code here is a breaking change for those repos.** They cannot be fixed in this PR — open coordinated PRs in each consumer repo and link them.

## Conventions

- In-memory store only, reseeded via `POST /api/admin/reset`.
- Standard library only; no external Go deps.
- Verify with: `go vet ./... && go build -o /dev/null .`, then `go run .` and curl the endpoints above.
