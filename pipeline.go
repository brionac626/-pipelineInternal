package pipelineinternal

import (
	"context"
	"errors"
)

// PipeNode is a node have prev and next worker node
type PipeNode struct {
	prevNode *PipeNode
	nextNode *PipeNode
	worker   Worker
}

// NewPipeNode create a new pipe line node with context
func NewPipeNode() *PipeNode {
	return &PipeNode{}
}

// SetPrevNode set a prev pipe line node from current node
func (n *PipeNode) SetPrevNode(prev *PipeNode) {
	n.prevNode = prev
}

// GetPrevNode get pipe's prev node
func (n *PipeNode) GetPrevNode() *PipeNode {
	return n.prevNode
}

// SetNextNode set a next pipe line node from current node
func (n *PipeNode) SetNextNode(next *PipeNode) {
	n.nextNode = next
}

// GetNextNode get pipe's next node
func (n *PipeNode) GetNextNode() *PipeNode {
	return n.nextNode
}

// SetNodes set both prev and next node at same time
func (n *PipeNode) SetNodes(prev, next *PipeNode) {
	n.prevNode = prev
	n.nextNode = next
}

// SetWork set current pipe's work
func (n *PipeNode) SetWork(w Worker) {
	n.worker = w
}

var (
	// ErrFristNode return when there are no node before current node
	ErrFristNode = errors.New("frist node")
	// ErrLastNode return when there are no node after current node
	ErrLastNode = errors.New("last node")
	// ErrEmptyNode return when parameter is nil
	ErrEmptyNode = errors.New("empty node")
)

// Pipe is a pipe line main executor
type Pipe struct {
	ctx context.Context
}

// NewPipe create a new pipe
func NewPipe(ctx context.Context) *Pipe {
	return &Pipe{ctx: ctx}
}

// Exec execute worker's work forward
func (p *Pipe) Exec(n *PipeNode) error {
	if n != nil {
		ctx, err := n.worker.Run(p.ctx)
		if err != nil {
			return err
		}

		if n.GetNextNode() == nil {
			return ErrLastNode
		}

		p.ctx = ctx

		return p.Exec(n.GetNextNode())
	}

	return ErrEmptyNode
}

// Revert execute worker's work backward
func (p *Pipe) Revert(n *PipeNode) error {
	if n != nil {
		ctx, err := n.worker.Revert(p.ctx)
		if err != nil {
			return err
		}

		if n.GetPrevNode() == nil {
			return ErrFristNode
		}

		p.ctx = ctx

		return p.Revert(n.GetPrevNode())
	}

	return ErrEmptyNode
}
