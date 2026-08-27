package main

import (
	"fmt"
	"sync"
)

type reservation struct {
	id    string
	items []ReservationItem
}

type Store struct {
	mu           sync.Mutex
	stock        map[string]*StockRecord
	reservations map[string]*reservation
	nextResID    int
}

func seedStock() map[string]*StockRecord {
	return map[string]*StockRecord{
		"volt-vaquero":       {ProductID: "volt-vaquero", StockCount: 12, Warehouse: "PDX-1", RestockEtaDays: 5},
		"thunder-pedal":      {ProductID: "thunder-pedal", StockCount: 7, Warehouse: "PDX-1", RestockEtaDays: 7},
		"watt-wanderer":      {ProductID: "watt-wanderer", StockCount: 25, Warehouse: "SLC-2", RestockEtaDays: 4},
		"ampere-al-fresco":   {ProductID: "ampere-al-fresco", StockCount: 2, Warehouse: "PDX-1", RestockEtaDays: 10},
		"circuit-breaker-xl": {ProductID: "circuit-breaker-xl", StockCount: 0, Warehouse: "SLC-2", RestockEtaDays: 14},
		"sparrow-glide":      {ProductID: "sparrow-glide", StockCount: 18, Warehouse: "SLC-2", RestockEtaDays: 3},
		"joule-junior":       {ProductID: "joule-junior", StockCount: 9, Warehouse: "PDX-1", RestockEtaDays: 6},
	}
}

func NewStore() *Store {
	return &Store{
		stock:        seedStock(),
		reservations: map[string]*reservation{},
	}
}

func (s *Store) Get(productID string) (StockRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.stock[productID]
	if !ok {
		return StockRecord{}, false
	}
	return *rec, true
}

func (s *Store) GetBatch(productIDs []string) []StockRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := []StockRecord{}
	for _, id := range productIDs {
		if rec, ok := s.stock[id]; ok {
			items = append(items, *rec)
		}
	}
	return items
}

// Reserve atomically decrements stock for every item, or fails without
// touching anything if any single item is short.
func (s *Store) Reserve(req ReservationRequest) (ReservationResponse, *InsufficientStockError) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range req.Items {
		rec, ok := s.stock[item.ProductID]
		if !ok {
			return ReservationResponse{}, &InsufficientStockError{
				Error: "insufficient_stock", ProductID: item.ProductID,
				Requested: item.Quantity, StockCount: 0,
			}
		}
		if rec.StockCount < item.Quantity {
			return ReservationResponse{}, &InsufficientStockError{
				Error: "insufficient_stock", ProductID: item.ProductID,
				Requested: item.Quantity, StockCount: rec.StockCount,
			}
		}
	}

	s.nextResID++
	res := &reservation{id: fmt.Sprintf("res-%04x", s.nextResID), items: req.Items}
	s.reservations[res.id] = res

	reserved := make([]ReservedItem, 0, len(req.Items))
	for _, item := range req.Items {
		rec := s.stock[item.ProductID]
		rec.StockCount -= item.Quantity
		reserved = append(reserved, ReservedItem{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			RemainingStock: rec.StockCount,
		})
	}

	return ReservationResponse{ReservationID: res.id, Status: "reserved", Items: reserved}, nil
}

// Release returns a reservation's stock to the pool (checkout rollback).
func (s *Store) Release(reservationID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, ok := s.reservations[reservationID]
	if !ok {
		return false
	}
	for _, item := range res.items {
		if rec, exists := s.stock[item.ProductID]; exists {
			rec.StockCount += item.Quantity
		}
	}
	delete(s.reservations, reservationID)
	return true
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stock = seedStock()
	s.reservations = map[string]*reservation{}
}
