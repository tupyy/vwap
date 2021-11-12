package compute_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/vwap/internal/compute"
	"github.com/tupyy/vwap/internal/entity"
)

func TestCurrencyAvgCalculator(t *testing.T) {
	c := compute.NewAvgCalculator(3)

	c.ProcessHeartBeat(entity.HeartBeat{
		Sequence: 1,
	})

	avg, totalPoints, err := c.ProcessTicker(entity.Ticker{
		Sequence:  1,
		Price:     1,
		Volume:    1,
		Timestamp: time.Now(),
	})

	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, float64(1), avg, "avg should be 1")
	assert.Equal(t, 1, totalPoints, "total points should be 1")

	avg, totalPoints, err = c.ProcessTicker(entity.Ticker{
		Sequence:  0,
		Price:     1,
		Volume:    1,
		Timestamp: time.Now(),
	})

	assert.NotNil(t, err, "should have a error")
	assert.ErrorIs(t, err, compute.ErrSequenceNotIncreasing, "should have err seq not increasing")
}
