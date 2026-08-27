package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var store = NewStore()

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "inventory"})
}

// GET /api/stock/{productId}
func handleGetStock(w http.ResponseWriter, r *http.Request) {
	productID := strings.TrimPrefix(r.URL.Path, "/api/stock/")
	rec, ok := store.Get(productID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product_not_found"})
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

// POST /api/stock/batch
func handleBatchStock(w http.ResponseWriter, r *http.Request) {
	var req BatchStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_body"})
		return
	}
	writeJSON(w, http.StatusOK, BatchStockResponse{Items: store.GetBatch(req.ProductIDs)})
}

// POST /api/reservations
func handleCreateReservation(w http.ResponseWriter, r *http.Request) {
	var req ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_body"})
		return
	}
	res, stockErr := store.Reserve(req)
	if stockErr != nil {
		writeJSON(w, http.StatusConflict, stockErr)
		return
	}
	log.Printf("reserved %s for order %s", res.ReservationID, req.OrderID)
	writeJSON(w, http.StatusCreated, res)
}

// DELETE /api/reservations/{reservationId}
func handleReleaseReservation(w http.ResponseWriter, r *http.Request) {
	reservationID := strings.TrimPrefix(r.URL.Path, "/api/reservations/")
	if !store.Release(reservationID) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "reservation_not_found"})
		return
	}
	log.Printf("released %s", reservationID)
	writeJSON(w, http.StatusOK, map[string]string{"reservationId": reservationID, "status": "released"})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /api/stock/{productId}", handleGetStock)
	mux.HandleFunc("POST /api/stock/batch", handleBatchStock)
	mux.HandleFunc("POST /api/reservations", handleCreateReservation)
	mux.HandleFunc("DELETE /api/reservations/{reservationId}", handleReleaseReservation)
	mux.HandleFunc("POST /api/admin/reset", func(w http.ResponseWriter, _ *http.Request) {
		store.Reset()
		writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4003"
	}
	log.Printf("inventory service listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
