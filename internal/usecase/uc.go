package usecase

import (
	"context"
	"time"

	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
)

type CurrencyAvgCalculator interface {
	ProductID() string
	ProcessHeartBeat(h entity.HeartBeat)
	ProcessTicker(t entity.Ticker) (avg float64, totalPoints int, err error)
}

type OutputWriter interface {
	Write(r entity.AverageResult) error
}

type AvgManager struct {
	doneCh chan interface{}

	outWriter OutputWriter
	// avgCurrencyCalculators holds the avg calculators.
	// the key is the product id
	avgCurrencyCalculators map[string]CurrencyAvgCalculator
}

func NewAvgManager(calculators []CurrencyAvgCalculator, o OutputWriter) *AvgManager {
	avgManager := &AvgManager{
		doneCh:                 make(chan interface{}, 1),
		outWriter:              o,
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
						err := a.outWriter.Write(entity.AverageResult{
							ProductID:   v.ProductID,
							Average:     avg,
							Timestamp:   time.Now(),
							TotalPoints: totalPoints,
						})
						if err != nil {
							log.GetLogger().Warning("cannot write to output: %+v", err)
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

// Shutdown close the avg manager.
func (a *AvgManager) Shutdown() {
	a.doneCh <- struct{}{}
}
