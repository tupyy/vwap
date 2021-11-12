package compute

import (
	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
)

type node struct {
	prev  *node
	value entity.DataPoint
}

type stack struct {
	root *node
	size int
}

func newStack() *stack {
	return &stack{}
}

func (s *stack) Push(p entity.DataPoint) {
	if s.root == nil {
		s.root = &node{value: p}
		s.size++

		log.GetLogger().Trace("set root to %+v", p)

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

	log.GetLogger().Trace("element popped %+v. new size: %d", root, s.size)

	return root
}

func (s *stack) Size() int {
	return s.size
}

func (s *stack) pushEnd(n *node, p entity.DataPoint) {
	if n.prev == nil {
		n.prev = &node{value: p}
		s.size++

		return
	}

	s.pushEnd(n.prev, p)
}
