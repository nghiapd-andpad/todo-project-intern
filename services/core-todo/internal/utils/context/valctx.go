package context

import (
	"context"
	"time"
)

// ValueOnly returns a context that keeps only the parent's values.
func ValueOnly(ctx context.Context) context.Context {
	return valueOnlyContext{ctx}
}

type valueOnlyContext struct {
	parent context.Context
}

func (v valueOnlyContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (v valueOnlyContext) Done() <-chan struct{} {
	return nil
}

func (v valueOnlyContext) Err() error {
	return nil
}

func (v valueOnlyContext) Value(key interface{}) interface{} {
	return v.parent.Value(key)
}
