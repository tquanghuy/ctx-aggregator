package aggregator

import (
	"context"
	"runtime"
	"sync"
)

var _ ContextAggregator[any] = new(concurrentAggregator[any])
var _ IConcurrentAggregator = new(concurrentAggregator[any])

// RegisterConcurrentContextAggregator register a concurrentAggregator pointer into context
// for collecting and aggregating data asynchronously from multiple goroutines.
// In order to use many aggregators in a project, please use different keys.
func RegisterConcurrentContextAggregator[T any](ctx context.Context, keys ...string) context.Context {
	agg := &concurrentAggregator[T]{
		m:     &sync.Mutex{},
		wg:    &sync.WaitGroup{},
		datas: make([]T, 0),
	}

	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

// RegisterConcurrentContextAggregatorWithCapacity register a concurrentAggregator pointer into context
// with a capacity hint for pre-allocation. This can improve performance when the expected
// number of items is known in advance, reducing memory allocations.
func RegisterConcurrentContextAggregatorWithCapacity[T any](ctx context.Context, capacity int, keys ...string) context.Context {
	agg := &concurrentAggregator[T]{
		m:     &sync.Mutex{},
		wg:    &sync.WaitGroup{},
		datas: make([]T, 0, capacity),
	}

	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

type concurrentAggregator[T any] struct {
	m     *sync.Mutex
	wg    *sync.WaitGroup
	datas []T
}

type IConcurrentAggregator interface {
	AddWait()
	Done()
}

// WaitContextFinalizer current cannot used due to
// context cannot passed into runtime.SetFinalizer
// fatal error: cannot pass *context.valueCtx to finalizer func()

func WaitContextFinalizer(ctx context.Context, keys ...string) context.Context {
	ctxKey := buildContextKey(keys...)
	aggVal := ctx.Value(ctxKey)
	if aggVal == nil {
		return ctx
	}
	agg, ok := aggVal.(IConcurrentAggregator)
	if !ok {
		return ctx
	}

	agg.AddWait()
	runtime.SetFinalizer(ctx, func() {
		agg.Done()
	})

	return ctx
}

func WaitFunc(ctx context.Context, keys ...string) (context.Context, func()) {
	ctxKey := buildContextKey(keys...)
	aggVal := ctx.Value(ctxKey)
	if aggVal == nil {
		return ctx, func() {}
	}
	agg, ok := aggVal.(IConcurrentAggregator)
	if !ok {
		return ctx, func() {}
	}

	agg.AddWait()
	return ctx, func() { agg.Done() }
}

func (a *concurrentAggregator[T]) Collect(data T) {
	a.m.Lock()
	defer a.m.Unlock()

	a.datas = append(a.datas, data)
}

func (a *concurrentAggregator[T]) Aggregate() []T {
	// Always call Wait before lock mutex for not cause deadlock
	// between syncgroup and mutex
	a.wg.Wait()

	a.m.Lock()
	defer a.m.Unlock()

	return a.datas
}

func (a *concurrentAggregator[T]) AddWait() {
	a.wg.Add(1)
}

func (a *concurrentAggregator[T]) Done() {
	a.wg.Done()
}
