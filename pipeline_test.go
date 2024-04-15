package pipelineinternal

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestPipe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p1 := NewPipeNode()
	p2 := NewPipeNode()
	p3 := NewPipeNode()

	// work 1 first work
	p1.SetNextNode(p2)
	p1.SetWork(&p1Worker{})

	// work 2
	p2.SetPrevNode(p1)
	p2.SetNextNode(p3)
	p2.SetWork(&p2Worker{})

	// work 3 last work
	p3.SetPrevNode(p2)
	p3.SetWork(&p3Worker{})

	p := NewPipe(ctx)
	if err := p.Exec(p1); err != nil {
		if !errors.Is(err, ErrLastNode) {
			t.Error(err)
		}
	}

	if err := p.Revert(p3); err != nil {
		if !errors.Is(err, ErrFirstNode) {
			t.Error(err)
		}
	}
}

type p1Worker struct {
	Name string
}

func (p1 *p1Worker) Run(ctx context.Context) (context.Context, error) {
	fmt.Println("p1")
	p1.Name = "worker 1"
	ctx = context.WithValue(ctx, p1Worker{}, p1)
	return ctx, nil
}

func (p1 *p1Worker) Revert(ctx context.Context) (context.Context, error) {
	fmt.Println("1p")
	return ctx, nil
}

type p2Worker struct {
	Name string
}

func (p2 *p2Worker) Run(ctx context.Context) (context.Context, error) {
	fmt.Println("p2")

	p1, ok := ctx.Value(p1Worker{}).(*p1Worker)
	if !ok {
		return context.TODO(), errors.New("p1 worker error")
	}

	fmt.Println("worker in ctx", p1.Name)

	p2.Name = "worker 2"
	ctx = context.WithValue(ctx, p2Worker{}, p2)

	return ctx, nil
}

func (p2 *p2Worker) Revert(ctx context.Context) (context.Context, error) {
	fmt.Println("2p")
	return ctx, nil
}

type p3Worker struct {
	Name string
}

func (p3 *p3Worker) Run(ctx context.Context) (context.Context, error) {
	fmt.Println("p3")

	p2, ok := ctx.Value(p2Worker{}).(*p2Worker)
	if !ok {
		return context.TODO(), errors.New("p2 worker error")
	}

	fmt.Println("worker in ctx", p2.Name)

	return ctx, nil
}

func (p3 *p3Worker) Revert(ctx context.Context) (context.Context, error) {
	fmt.Println("3p")
	return ctx, nil
}
