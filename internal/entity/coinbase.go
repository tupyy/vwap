package entity

import "time"

// Ticker represent the json ticker message from coinbase.
// nolint: tagliatelle
type Ticker struct {
	Sequence  int64     `json:"sequence"`
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price,string"`
	Volume    float64   `json:"last_size,string"`
	Timestamp time.Time `json:"time"`
}

// Heartbeat message.
// nolint: tagliatelle
type HeartBeat struct {
	ProductID string `json:"product_id"`
	Sequence  int64  `json:"sequence"`
}
