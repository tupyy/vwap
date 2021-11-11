package entity

import "time"

type AverageResult struct {
	// ProductID -- id of the product
	ProductID string
	// Timestamp -- timestamp of the calculation
	Timestamp time.Time
	// Average -- actual value of the average
	Average float64
	// TotalPoints -- number of points used in calculation
	TotalPoints int
}
