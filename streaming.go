package aggregator

import (
	"context"
	"sync"
)

var _ ContextAggregator[any] = new(streamingAggregator[any])
var _ ContextAggregator[any] = new(concurrentStreamingAggregator[any])

// CollectCallback is a function that is called whenever an item is collected
type CollectCallback[T any] func(T)

// RegisterStreamingAggregator registers a streaming aggregator that calls a callback
// function for each collected item. The callback is invoked synchronously during collection.
// Use this for sequential data collection with real-time processing.
func RegisterStreamingAggregator[T any](ctx context.Context, callback CollectCallback[T], keys ...string) context.Context {
	agg := &streamingAggregator[T]{
		datas:    make([]T, 0),
		callback: callback,
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

// RegisterStreamingAggregatorWithCapacity registers a streaming aggregator with capacity hint
func RegisterStreamingAggregatorWithCapacity[T any](ctx context.Context, capacity int, callback CollectCallback[T], keys ...string) context.Context {
	agg := &streamingAggregator[T]{
		datas:    make([]T, 0, capacity),
		callback: callback,
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

// RegisterConcurrentStreamingAggregator registers a thread-safe streaming aggregator
// that calls a callback function for each collected item. The callback is invoked
// synchronously during collection with mutex protection.
func RegisterConcurrentStreamingAggregator[T any](ctx context.Context, callback CollectCallback[T], keys ...string) context.Context {
	agg := &concurrentStreamingAggregator[T]{
		m:        &sync.Mutex{},
		wg:       &sync.WaitGroup{},
		datas:    make([]T, 0),
		callback: callback,
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

// RegisterConcurrentStreamingAggregatorWithCapacity registers a thread-safe streaming aggregator with capacity hint
func RegisterConcurrentStreamingAggregatorWithCapacity[T any](ctx context.Context, capacity int, callback CollectCallback[T], keys ...string) context.Context {
	agg := &concurrentStreamingAggregator[T]{
		m:        &sync.Mutex{},
		wg:       &sync.WaitGroup{},
		datas:    make([]T, 0, capacity),
		callback: callback,
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

// streamingAggregator is a sequential aggregator with callback support
type streamingAggregator[T any] struct {
	datas    []T
	callback CollectCallback[T]
}

func (a *streamingAggregator[T]) Collect(data T) {
	// Call callback first (with panic recovery)
	if a.callback != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Silently recover from callback panics to prevent disrupting collection
				}
			}()
			a.callback(data)
		}()
	}

	// Store data for later aggregation
	a.datas = append(a.datas, data)
}

func (a *streamingAggregator[T]) Aggregate() []T {
	return a.datas
}

// concurrentStreamingAggregator is a thread-safe aggregator with callback support
type concurrentStreamingAggregator[T any] struct {
	m        *sync.Mutex
	wg       *sync.WaitGroup
	datas    []T
	callback CollectCallback[T]
}

func (a *concurrentStreamingAggregator[T]) Collect(data T) {
	a.m.Lock()
	defer a.m.Unlock()

	// Call callback first (with panic recovery)
	if a.callback != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Silently recover from callback panics to prevent disrupting collection
				}
			}()
			a.callback(data)
		}()
	}

	// Store data for later aggregation
	a.datas = append(a.datas, data)
}

func (a *concurrentStreamingAggregator[T]) Aggregate() []T {
	// Always call Wait before lock mutex for not cause deadlock
	a.wg.Wait()

	a.m.Lock()
	defer a.m.Unlock()

	return a.datas
}

func (a *concurrentStreamingAggregator[T]) AddWait() {
	a.wg.Add(1)
}

func (a *concurrentStreamingAggregator[T]) Done() {
	a.wg.Done()
}
