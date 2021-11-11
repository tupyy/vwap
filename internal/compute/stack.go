package compute

import (
	"github.com/tupyy/vwap/internal/entity"
)

type node struct {
	prev  *node
	value entity.VolumePoint
}

type stack struct {
	root *node
	size int
}

func newStack() *stack {
	return &stack{}
}

func (s *stack) Push(p entity.VolumePoint) {
	if s.root == nil {
		s.root = &node{value: p}
		s.size++

		return
	}

	s.pushEnd(s.root, p)
}

func (s *stack) Pop() *node {
	if s.root == nil {
		return nil
	}

	root := s.root
	s.root = s.root.prev
	s.size--

	return root
}

func (s *stack) Size() int {
	return s.size
}

func (s *stack) pushEnd(n *node, p entity.VolumePoint) {
	if n.prev == nil {
		n.prev = &node{value: p}
		s.size++

		return
	}

	s.pushEnd(n.prev, p)
}
