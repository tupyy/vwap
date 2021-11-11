package compute_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/usecase"
)

func TestCurrencyAvgCalculator(t *testing.T) {
	c := usecase.NewCurrencyAvgCalculator("ID", 3)

	c.ProcessHeartBeat(entity.HeartBeat{
		Sequence:  1,
		ProductID: "ID",
	})

	avg, err := c.ProcessTicker(entity.Ticker{
		Sequence:  1,
		ProductID: "ID",
		Price:     1,
		Volume:    1,
		Timestamp: time.Now(),
	})

	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, float64(1), avg, "avg should be 1")

	avg, err = c.ProcessTicker(entity.Ticker{
		Sequence:  0,
		ProductID: "ID",
		Price:     1,
		Volume:    1,
		Timestamp: time.Now(),
	})

	assert.NotNil(t, err, "should have a error")
	assert.ErrorIs(t, err, usecase.ErrSequenceNotIncreasing, "should have err seq not increasing")
}
