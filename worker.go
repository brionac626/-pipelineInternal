package pipelineinternal

import "context"

// Worker is an executer of a pipe line
type Worker interface {
	Run(ctx context.Context) (context.Context, error)
	Revert(ctx context.Context) (context.Context, error)
}
