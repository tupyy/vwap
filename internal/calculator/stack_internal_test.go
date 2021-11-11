package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/vwap/internal/entity"
)

func TestStack(t *testing.T) {
	s := newStack()

	s.Push(entity.VolumePoint{Value: 1, Volume: 1})
	s.Push(entity.VolumePoint{Value: 2, Volume: 2})

	assert.Equal(t, 2, s.Size(), "expected 2 elements in the stack")

	n := s.Pop()
	assert.Equal(t, float64(1), n.value.Value, "expect 1 as value")
	assert.Equal(t, 1, s.Size(), "expect size 1")

	n = s.Pop()
	assert.Equal(t, float64(2), n.value.Value, "expect 2 as value")
	assert.Equal(t, 0, s.Size(), "expect size 0")

	// pop again we should get nil
	noNode := s.Pop()
	assert.Nil(t, noNode, "expect nil")
}
