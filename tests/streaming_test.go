package aggregator_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func TestStreamingAggregator(t *testing.T) {
	ctx := context.Background()
	var callbackCount int32

	// Register with callback
	ctx = aggregator.RegisterStreamingAggregator(ctx, func(s string) {
		atomic.AddInt32(&callbackCount, 1)
	})

	// Collect items
	_ = aggregator.Collect(ctx, "item1")
	_ = aggregator.Collect(ctx, "item2")
	_ = aggregator.Collect(ctx, "item3")

	// Verify callback was called
	assert.Equal(t, int32(3), atomic.LoadInt32(&callbackCount))

	// Verify data is still aggregated
	results, err := aggregator.Aggregate[string](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestStreamingAggregatorWithCapacity(t *testing.T) {
	ctx := context.Background()
	var callbackCount int32

	// Register with capacity and callback
	ctx = aggregator.RegisterStreamingAggregatorWithCapacity(ctx, 100, func(n int) {
		atomic.AddInt32(&callbackCount, 1)
	})

	// Collect items
	for i := 0; i < 50; i++ {
		_ = aggregator.Collect(ctx, i)
	}

	// Verify callback was called
	assert.Equal(t, int32(50), atomic.LoadInt32(&callbackCount))

	// Verify data is still aggregated
	results, err := aggregator.Aggregate[int](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 50)
}

func TestConcurrentStreamingAggregator(t *testing.T) {
	ctx := context.Background()
	var callbackCount int32

	// Register concurrent streaming aggregator
	ctx = aggregator.RegisterConcurrentStreamingAggregator(ctx, func(n int) {
		atomic.AddInt32(&callbackCount, 1)
	})

	// Collect from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = aggregator.Collect(ctx, val)
		}(i)
	}
	wg.Wait()

	// Verify callback was called for each item
	assert.Equal(t, int32(10), atomic.LoadInt32(&callbackCount))

	// Verify data is still aggregated
	results, err := aggregator.Aggregate[int](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 10)
}

func TestStreamingAggregatorCallbackPanic(t *testing.T) {
	ctx := context.Background()

	// Register with callback that panics
	ctx = aggregator.RegisterStreamingAggregator(ctx, func(s string) {
		panic("callback panic")
	})

	// Collect should not panic
	assert.NotPanics(t, func() {
		_ = aggregator.Collect(ctx, "item1")
		_ = aggregator.Collect(ctx, "item2")
	})

	// Data should still be collected despite callback panics
	results, err := aggregator.Aggregate[string](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestStreamingAggregatorWithProcessing(t *testing.T) {
	ctx := context.Background()
	var processedItems []string
	var mu sync.Mutex

	// Register with callback that processes items
	ctx = aggregator.RegisterStreamingAggregator(ctx, func(s string) {
		mu.Lock()
		defer mu.Unlock()
		processedItems = append(processedItems, "processed: "+s)
	})

	// Collect items
	_ = aggregator.Collect(ctx, "item1")
	_ = aggregator.Collect(ctx, "item2")
	_ = aggregator.Collect(ctx, "item3")

	// Verify processing
	mu.Lock()
	assert.Len(t, processedItems, 3)
	assert.Equal(t, "processed: item1", processedItems[0])
	mu.Unlock()

	// Verify original data is still available
	results, err := aggregator.Aggregate[string](ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"item1", "item2", "item3"}, results)
}

func TestConcurrentStreamingAggregatorWithCapacity(t *testing.T) {
	ctx := context.Background()
	var callbackCount int32

	// Register with capacity
	ctx = aggregator.RegisterConcurrentStreamingAggregatorWithCapacity(ctx, 100, func(n int) {
		atomic.AddInt32(&callbackCount, 1)
	})

	// Collect from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = aggregator.Collect(ctx, val)
		}(i)
	}
	wg.Wait()

	// Verify callback was called
	assert.Equal(t, int32(50), atomic.LoadInt32(&callbackCount))

	// Verify data is aggregated
	results, err := aggregator.Aggregate[int](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 50)
}

func TestStreamingWithFilterAndTransform(t *testing.T) {
	ctx := context.Background()
	var callbackItems []string
	var mu sync.Mutex

	// Register streaming aggregator
	ctx = aggregator.RegisterStreamingAggregator(ctx, func(s string) {
		mu.Lock()
		defer mu.Unlock()
		callbackItems = append(callbackItems, s)
	})

	// Collect mixed data
	_ = aggregator.Collect(ctx, "ERROR: error1")
	_ = aggregator.Collect(ctx, "INFO: info1")
	_ = aggregator.Collect(ctx, "ERROR: error2")

	// Verify all items were sent to callback
	mu.Lock()
	assert.Len(t, callbackItems, 3)
	mu.Unlock()

	// Filter and transform
	errorLengths, err := aggregator.AggregateWithFilterAndTransform(
		ctx,
		func(s string) bool {
			return len(s) > 10
		},
		func(s string) int {
			return len(s)
		},
	)

	assert.NoError(t, err)
	assert.Len(t, errorLengths, 3)
}
