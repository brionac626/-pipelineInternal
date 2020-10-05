package pipelineinternal

import (
	"context"
	"errors"
)

// PipeNode is a node have prev and next worker node
type PipeNode struct {
	// ctx      context.Context
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
	if n != nil {
		n.prevNode = prev
	}
}

// GetPrevNode get pipe's prev node
func (n *PipeNode) GetPrevNode() *PipeNode {
	return n.prevNode
}

// SetNextNode set a next pipe line node from current node
func (n *PipeNode) SetNextNode(next *PipeNode) {
	if n != nil {
		n.nextNode = next
	}
}

// GetNextNode get pipe's next node
func (n *PipeNode) GetNextNode() *PipeNode {
	return n.nextNode
}

// SetWork set current pipe's work
func (n *PipeNode) SetWork(w Worker) {
	if n != nil {
		n.worker = w
	}
}

// ErrLastNode return when there are no node left
var ErrLastNode = errors.New("last node")

// DefaultExec is a default method for a pipe line to run
func DefaultExec(ctx context.Context, n *PipeNode) error {
	if n != nil {
		ctx, err := n.worker.Run(ctx)
		if err != nil {
			return err
		}

		if n.GetNextNode() == nil {
			return ErrLastNode
		}

		return DefaultExec(ctx, n.GetNextNode())
	}

	return errors.New("empty node")
}
