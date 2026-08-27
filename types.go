package main

// StockRecord is the wire format for a product's stock level.
// Consumed by: catalog (stock badges), pricing (discount/surcharge rules),
// orders (delivery estimates). Renaming a JSON tag here breaks all three.
type StockRecord struct {
	ProductID     string `json:"productId"`
	StockCount    int    `json:"availableUnits"`
	Warehouse     string `json:"warehouse"`
	RestockEtaDays int   `json:"restockEtaDays"`
}

type BatchStockRequest struct {
	ProductIDs []string `json:"productIds"`
}

type BatchStockResponse struct {
	Items []StockRecord `json:"items"`
}

type ReservationItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type ReservationRequest struct {
	OrderID string            `json:"orderId"`
	Items   []ReservationItem `json:"items"`
}

type ReservedItem struct {
	ProductID      string `json:"productId"`
	Quantity       int    `json:"quantity"`
	RemainingStock int    `json:"remainingStock"`
}

// ReservationResponse.Status is "reserved" on success; orders checks this
// exact string before proceeding with checkout.
type ReservationResponse struct {
	ReservationID string         `json:"reservationId"`
	Status        string         `json:"status"`
	Items         []ReservedItem `json:"items"`
}

type InsufficientStockError struct {
	Error      string `json:"error"`
	ProductID  string `json:"productId"`
	Requested  int    `json:"requested"`
	StockCount int    `json:"availableUnits"`
}
