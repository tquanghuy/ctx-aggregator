package aggregator

import (
	"context"
)

var _ ContextAggregator[any] = new(baseAggregator[any])

// RegisterBaseContextAggregator register a baseAggregator pointer into context
// for collecting and aggregating data sequentially without any asynchronous
// lock. In order to use many aggregators in a project, please use different keys.
func RegisterBaseContextAggregator[T any](ctx context.Context, keys ...string) context.Context {
	agg := &baseAggregator[T]{
		datas: make([]T, 0),
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

type baseAggregator[T any] struct {
	datas []T
}

func (a *baseAggregator[T]) Collect(data T) {
	a.datas = append(a.datas, data)
}

func (a *baseAggregator[T]) Aggregate() []T {
	return a.datas
}
