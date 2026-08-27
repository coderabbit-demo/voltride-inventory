# ⚡ voltride-inventory

Stock levels and reservations for the [VoltRide](https://github.com/coderabbit-demo/voltride-platform) e-bike store. Go, standard library only, in-memory data seeded at startup. Runs on **port 4003**.

This is a leaf service: it calls no one, but its wire format is consumed by [voltride-catalog](https://github.com/coderabbit-demo/voltride-catalog), [voltride-pricing](https://github.com/coderabbit-demo/voltride-pricing), and [voltride-orders](https://github.com/coderabbit-demo/voltride-orders). See `AGENTS.md` before changing any response shape.

## Endpoints

- `GET /health`
- `GET /api/stock/:productId` → `{ productId, stockCount, warehouse, restockEtaDays }`
- `POST /api/stock/batch` `{ productIds: [] }` → `{ items: [...] }`
- `POST /api/reservations` `{ orderId, items: [{productId, quantity}] }` → 201 `{ reservationId, status: "reserved", items: [...] }`, 409 on shortage
- `DELETE /api/reservations/:id` → releases the hold (checkout rollback)
- `POST /api/admin/reset` → reseed stock

## Run

```sh
go run .          # requires Go >= 1.22; PORT env var overrides 4003
```

To run the whole VoltRide system, use the scripts in [voltride-platform](https://github.com/coderabbit-demo/voltride-platform).
