package compute

import (
	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
)

type Calculator struct {
	// stack holds the points
	stack *stack
	// totalVolume is the sum of all volumes in the stack
	totalVolume float64
	// valueVolumeSum is the Sum(point.value * point.volume) of all points in the stack
	valueVolumeSum float64
	// maxSize is the max number of points used in calculation
	maxSize int
}

func NewCalculator(size int) *Calculator {
	return &Calculator{
		stack:   newStack(),
		maxSize: size,
	}
}

func (c *Calculator) Add(p entity.VolumePoint) {
	if c.stack.Size() == c.maxSize {
		poppedPoint := c.stack.Pop()

		// substract the volume of the poppedPoint from totalVolume
		c.totalVolume -= poppedPoint.value.Volume

		// substract the product value*volume of the poppedPoint from valueVolumeSum
		c.valueVolumeSum -= poppedPoint.value.Value * poppedPoint.value.Volume
	}

	c.stack.Push(p)

	// recompute totalVolume and valueVolumeSum
	c.totalVolume += p.Volume
	c.valueVolumeSum += p.Value * p.Volume

	log.GetLogger().Debug("Total volume: %f ValueVolumeSum: %f Number of points: %d", c.totalVolume, c.valueVolumeSum, c.stack.Size())
}

func (c *Calculator) ComputeAverage() (avg float64, totalPoints int) {
	return c.valueVolumeSum / c.totalVolume, c.stack.Size()
}
