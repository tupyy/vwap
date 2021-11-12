package compute_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/vwap/internal/compute"
	"github.com/tupyy/vwap/internal/entity"
)

func TestCalculator(t *testing.T) {
	calc := compute.NewCalculator(3)

	calc.Add(entity.DataPoint{Value: 1, Volume: 1})
	calc.Add(entity.DataPoint{Value: 2, Volume: 2})

	// avg := (1*1 + 2*2) / 3 = 1.666
	avg, totalPoints := calc.ComputeAverage()

	assert.Equal(t, float64(5)/float64(3), avg, "expect avg = 1.6667")
	assert.Equal(t, 2, totalPoints, "expect 2 computation points")

	calc.Add(entity.DataPoint{Value: 1, Volume: 1})

	avg, totalPoints = calc.ComputeAverage()
	assert.Equal(t, float64(6)/float64(4), avg, "expected avg = 1.5")
	assert.Equal(t, 3, totalPoints, "expect 3 computation points")

	// avg = (2*2 + 1*1 + 2*2)/ 5
	calc.Add(entity.DataPoint{Value: 2, Volume: 2})

	avg, totalPoints = calc.ComputeAverage()
	assert.Equal(t, float64(9)/float64(5), avg, "expected avg = 1")
	assert.Equal(t, 3, totalPoints, "expect 3 computation points")
}
