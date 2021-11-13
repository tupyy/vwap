package manager

import (
	"context"
	"time"

	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
)

type CurrencyAvgCalculator interface {
	ProcessHeartBeat(h entity.HeartBeat)
	ProcessTicker(t entity.Ticker) (avg float64, totalPoints int, err error)
}

type OutputWriter interface {
	Write(r entity.AverageResult) error
}

type AvgManager struct {
	doneCh chan chan interface{}

	outWriter OutputWriter
	// avgCurrencyCalculators holds the avg calculators.
	// the key is the product id
	avgCurrencyCalculators map[string]CurrencyAvgCalculator
}

func NewAvgManager(o OutputWriter) *AvgManager {
	avgManager := &AvgManager{
		doneCh:                 make(chan chan interface{}, 1),
		outWriter:              o,
		avgCurrencyCalculators: make(map[string]CurrencyAvgCalculator),
	}

	return avgManager
}

func (a *AvgManager) AddAvgCalculator(productID string, c CurrencyAvgCalculator) {
	a.avgCurrencyCalculators[productID] = c
}

// Start starts the avg manager.
// It receive an input channel and a context.
// From input channel reads Ticker and HeartBeat messages.
// nolint: gocognit
func (a *AvgManager) Start(ctx context.Context, inputCh <-chan interface{}) {
	logger := log.GetLogger()

	go func() {
		for {
			select {
			case msg := <-inputCh:
				switch v := msg.(type) {
				case entity.HeartBeat:
					logger.Debugf("heart beat received: %+v", v)

					c := a.avgCurrencyCalculators[v.ProductID]
					c.ProcessHeartBeat(v)
				case entity.Ticker:
					logger.Debugf("ticker received: %+v", v)

					c, found := a.avgCurrencyCalculators[v.ProductID]
					if !found {
						logger.Errorf("received ticker for a product that does not exists: %s", v.ProductID)

						continue
					}

					avg, totalPoints, err := c.ProcessTicker(v)
					if err != nil {
						logger.Errorf("cannot compute average: %+v", err)
					} else {
						err := a.outWriter.Write(entity.AverageResult{
							ProductID:   v.ProductID,
							Average:     avg,
							Timestamp:   time.Now(),
							TotalPoints: totalPoints,
						})
						if err != nil {
							log.GetLogger().Warningf("cannot write to output: %+v", err)
						}
					}
				default:
					log.GetLogger().Warningf("cannot cast received message: %+v", msg)
				}
			case retCh := <-a.doneCh:
				retCh <- struct{}{}
				return
			case err := <-ctx.Done():
				logger.Errorf("ctx canceled: %+v. exit", err)
				return
			}
		}
	}()
}

// Shutdown close the avg manager.
// Block until goroutine returned.
func (a *AvgManager) Shutdown() {
	retCh := make(chan interface{})

	a.doneCh <- retCh
	<-retCh

	log.GetLogger().Infof("avg manager closed")
}
