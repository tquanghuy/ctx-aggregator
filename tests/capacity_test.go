package aggregator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func TestBaseContextAggregatorWithCapacity(t *testing.T) {
	ctx := context.Background()
	capacity := 100

	// Register with capacity
	ctx = aggregator.RegisterBaseContextAggregatorWithCapacity[int](ctx, capacity)

	// Collect items
	for i := 0; i < 50; i++ {
		err := aggregator.Collect(ctx, i)
		assert.NoError(t, err)
	}

	// Aggregate
	results, err := aggregator.Aggregate[int](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 50)

	// Verify values
	for i := 0; i < 50; i++ {
		assert.Equal(t, i, results[i])
	}
}

func TestConcurrentContextAggregatorWithCapacity(t *testing.T) {
	ctx := context.Background()
	capacity := 100

	// Register with capacity
	ctx = aggregator.RegisterConcurrentContextAggregatorWithCapacity[string](ctx, capacity)

	// Collect items
	for i := 0; i < 50; i++ {
		err := aggregator.Collect(ctx, "item")
		assert.NoError(t, err)
	}

	// Aggregate
	results, err := aggregator.Aggregate[string](ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 50)
}

func BenchmarkWithoutCapacity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		ctx = aggregator.RegisterBaseContextAggregator[int](ctx)

		for j := 0; j < 1000; j++ {
			_ = aggregator.Collect(ctx, j)
		}

		_, _ = aggregator.Aggregate[int](ctx)
	}
}

func BenchmarkWithCapacity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		ctx = aggregator.RegisterBaseContextAggregatorWithCapacity[int](ctx, 1000)

		for j := 0; j < 1000; j++ {
			_ = aggregator.Collect(ctx, j)
		}

		_, _ = aggregator.Aggregate[int](ctx)
	}
}
