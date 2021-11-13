package entity

import "time"

// Ticker represent the json ticker message from coinbase.
type Ticker struct {
	Sequence  int64     `json:"sequence"`
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price,string"`
	Volume    float64   `json:"last_size,string"`
	Timestamp time.Time `json:"time"`
}

type Ticker2 struct {
	Sequence  int64     `json:"sequence"`
	ProductID string    `json:"product_id"`
	Price     string    `json:"price"`
	Volume    string    `json:"last_size"`
	Timestamp time.Time `json:"time"`
}

// Heartbeat message.
type HeartBeat struct {
	ProductID string `json:"product_id"`
	Sequence  int64  `json:"sequence"`
}
