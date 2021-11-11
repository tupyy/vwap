package compute

import "github.com/tupyy/vwap/internal/entity"

type Calculator struct {
	// stack holds the points
	stack *stack
	// totalVolume is the sum of all volumes in the stack
	totalVolume float64
	// valueVolumeSum is the Sum(point.value * point.volume) of all points in the stack
	valueVolumeSum float64
	// size is the max number of points used in calculation
	size int
}

func New(size int) *Calculator {
	return &Calculator{
		stack: newStack(),
		size:  size,
	}
}

func (c *Calculator) Add(p entity.VolumePoint) {
	if c.stack.Size() == c.size {
		lastPoint := c.stack.Pop()

		// substract the volume of the lastpoint from totalVolume
		c.totalVolume -= lastPoint.value.Volume

		// substract the product value*volume of the lastpoint from valueVolumeSum
		c.valueVolumeSum -= lastPoint.value.Value * lastPoint.value.Volume
	}

	c.stack.Push(p)

	// recompute totalVolume and valueVolumeSum
	c.totalVolume += p.Volume
	c.valueVolumeSum += p.Value * p.Volume
}

func (c *Calculator) ComputeAverage() (avg float64, totalPoints int) {
	return c.valueVolumeSum / c.totalVolume, c.stack.Size()
}
