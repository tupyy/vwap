package manager_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/manager"
)

func TestAvgManager(t *testing.T) {
	writerMock := &outputWriter{}
	pairMock := &pairMockCalculator{}

	avgM := manager.NewAvgManager(writerMock)
	avgM.AddAvgCalculator("id", pairMock)

	inputCh := make(chan interface{})
	avgM.Start(context.Background(), inputCh)

	// add one heartbeat and one ticker
	inputCh <- entity.HeartBeat{ProductID: "id", Sequence: 1}
	inputCh <- entity.Ticker{ProductID: "id", Price: 1}

	<-time.After(1 * time.Second)
	assert.Equal(t, 1, pairMock.HeartbeatCallCount, "should have one heart beat")
	assert.Equal(t, 1, pairMock.TickerCallCount, "should have one ticker")

	// assert if the avg was wrote on writer
	assert.Equal(t, 1, writerMock.WriteCallCount, "should have one call")
	assert.Equal(t, float64(1), writerMock.Avg, "should have avg = 1")

	// push one message of another product
	inputCh <- entity.HeartBeat{ProductID: "unkown_product", Sequence: 1}
	inputCh <- entity.Ticker{ProductID: "unkown_product", Price: 1}

	<-time.After(1 * time.Second)
	assert.Equal(t, 1, pairMock.HeartbeatCallCount, "should have one heart beat")
	assert.Equal(t, 1, pairMock.TickerCallCount, "should have one ticker")

	avgM.Shutdown()

	// goroutine should be closed now
	doneCh := make(chan interface{}, 1)
	go func() {
		inputCh <- entity.HeartBeat{ProductID: "id", Sequence: 1}
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		t.Error("goroutine still running")
	case <-time.After(3 * time.Second):
		// we are good. it is closed
	}
}

/***************
	Mocks
***************/

type pairMockCalculator struct {
	Seq int64
	// counts how many heartbeat method was called
	HeartbeatCallCount int
	// counts how many times ticker method was called
	TickerCallCount int
}

func (p *pairMockCalculator) ProcessHeartBeat(h entity.HeartBeat) {
	p.Seq = h.Sequence
	p.HeartbeatCallCount++
}

func (p *pairMockCalculator) ProcessTicker(t entity.Ticker) (avg float64, totalPoints int, err error) {
	p.TickerCallCount++

	if t.Sequence == 10 {
		return 0, 0, errors.New("ticker error")
	}

	return t.Price, 1, nil
}

type outputWriter struct {
	WriteCallCount int
	Avg            float64
}

func (o *outputWriter) Write(r entity.AverageResult) error {
	o.Avg = r.Average
	o.WriteCallCount++

	return nil
}
