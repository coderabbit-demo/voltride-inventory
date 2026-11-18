package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateReservationReturnsReservedStatus(t *testing.T) {
	previousStore := store
	store = NewStore()
	t.Cleanup(func() { store = previousStore })

	req := httptest.NewRequest(http.MethodPost, "/api/reservations", strings.NewReader(`{"orderId":"order-1","items":[{"productId":"volt-vaquero","quantity":1}]}`))
	response := httptest.NewRecorder()
	handleCreateReservation(response, req)

	if response.Code != http.StatusCreated {
		t.Fatalf("status code = %d, want %d; body: %s", response.Code, http.StatusCreated, response.Body.String())
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Status != "reserved" {
		t.Errorf("reservation status = %q, want %q", body.Status, "reserved")
	}
}
