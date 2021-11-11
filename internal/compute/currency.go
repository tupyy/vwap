package compute

import (
	"errors"
	"fmt"
	"time"

	"github.com/tupyy/vwap/internal/calculator"
	"github.com/tupyy/vwap/internal/entity"
)

// DefaultVolumeSize is the default size for average calculation.
const DefaultVolumeSize = 200

var (
	// ErrSequenceNotIncreasing means that the sequence of the receiving ticker message is inferior of the last seen sequence.
	// It means that the message arrive too late and is not taken into account.
	ErrSequenceNotIncreasing = errors.New("error sequence not increasing")
)

type currencyAvgCalculator struct {
	// c -- avg calculator
	calc *calculator.Calculator
	// heartBeatSequence -- holds the last received sequence
	heartBeatSequence int64
	// lastTimestamp -- holds the timestamp of the last message
	lastTimestamp time.Time
}

func NewCurrencyAvgCalculator(volumeSize int) *currencyAvgCalculator {
	return &currencyAvgCalculator{
		calc: calculator.New(volumeSize),
	}
}

// ProcessHeartBeat updates the lastSequence and last timestamp
func (c *currencyAvgCalculator) ProcessHeartBeat(h entity.HeartBeat) {
	c.heartBeatSequence = h.Sequence
}

func (c *currencyAvgCalculator) ProcessTicker(t entity.Ticker) (avg float64, totalPoints int, err error) {
	if t.Sequence < c.heartBeatSequence {
		return 0, 0, fmt.Errorf("%w received sequence: %d last sequence: %d", ErrSequenceNotIncreasing, t.Sequence, c.heartBeatSequence)
	}

	newPoint := entity.VolumePoint{
		Value:  t.Price,
		Volume: t.Volume,
	}

	// add the new point to calculator
	c.calc.Add(newPoint)

	avg, totalPoints = c.calc.ComputeAverage()

	return
}
