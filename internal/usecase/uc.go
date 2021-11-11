package usecase

import (
	"context"
	"time"

	"github.com/tupyy/vwap/internal/entity"
)

type CurrencyAvgCalculator interface {
	ProductID() string
	ProcessHeartBeat(h entity.HeartBeat)
	ProcessTicker(t entity.Ticker) (avg float64, totalPoints int, err error)
}

type AvgManager struct {
	doneCh   chan interface{}
	outputCh chan<- entity.AverageResult

	// avgCurrencyCalculators holds the avg calculators.
	// the key is the product id
	avgCurrencyCalculators map[string]CurrencyAvgCalculator
}

func NewAvgManager(calculators []CurrencyAvgCalculator) *AvgManager {
	avgManager := &AvgManager{
		doneCh:                 make(chan interface{}, 1),
		outputCh:               make(chan entity.AverageResult),
		avgCurrencyCalculators: make(map[string]CurrencyAvgCalculator),
	}

	for _, c := range calculators {
		avgManager.avgCurrencyCalculators[c.ProductID()] = c
	}

	return avgManager
}

// Start starts the avg manager.
// It receive an input channel and a context.
// From input channel reads Ticker and HeartBeat messages.
func (a *AvgManager) Start(ctx context.Context, inputCh <-chan interface{}) {
	go func() {
		for {
			select {
			case msg := <-inputCh:
				switch v := msg.(type) {
				case entity.HeartBeat:
					c := a.avgCurrencyCalculators[v.ProductID]
					c.ProcessHeartBeat(v)
				case entity.Ticker:
					c := a.avgCurrencyCalculators[v.ProductID]

					avg, totalPoints, err := c.ProcessTicker(v)
					if err != nil {
						// log it
					} else {
						a.outputCh <- entity.AverageResult{
							ProductID:   v.ProductID,
							Average:     avg,
							Timestamp:   time.Now(),
							TotalPoints: totalPoints,
						}
					}
				}
			case <-a.doneCh:
				return
			case <-ctx.Done():
				// ctx canceled
				return
			}
		}
	}()
}

// OutputChannel returns the channel on which the results will be written.
func (a *AvgManager) OutputChannel() chan<- entity.AverageResult {
	return a.outputCh
}

// Shutdown close the avg manager.
func (a *AvgManager) Shutdown() {
	a.doneCh <- struct{}{}
}
