package aggregator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func funcBaseCollecInt32(ctx context.Context, keys ...string) error {
	return aggregator.Collect(ctx, int32(0), keys...)
}

func funcBaseCollecInt32WithVal(ctx context.Context, val int32, keys ...string) error {
	return aggregator.Collect(ctx, val, keys...)
}

func funcBaseCollecString(ctx context.Context, keys ...string) error {
	return aggregator.Collect(ctx, "", keys...)
}

func TestBaseContextAggregator_CollectNotFoundAggregator(t *testing.T) {
	err := funcBaseCollecInt32(context.Background())
	assert.Equal(t, err, aggregator.ErrNotFoundAggregator)
}

func TestBaseContextAggregator_CollectInvalidType(t *testing.T) {
	ctx := aggregator.RegisterBaseContextAggregator[int](context.Background())
	err := funcBaseCollecInt32(ctx)
	assert.Equal(t, err, aggregator.ErrInvalidType)
}

func TestBaseContextAggregator_AggregateNotFoundAggregator(t *testing.T) {
	result, err := aggregator.Aggregate[int](context.Background())
	assert.Nil(t, result)
	assert.Equal(t, err, aggregator.ErrNotFoundAggregator)
}

func TestBaseContextAggregator_AggregateInvalidType(t *testing.T) {
	ctx := aggregator.RegisterBaseContextAggregator[int32](context.Background())
	err := funcBaseCollecInt32(ctx)
	assert.Nil(t, err)

	result, err := aggregator.Aggregate[int](ctx)
	assert.Nil(t, result)
	assert.Equal(t, err, aggregator.ErrInvalidType)
}

func TestBaseContextAggregator_SuccessNoKey(t *testing.T) {
	ctx := aggregator.RegisterBaseContextAggregator[int32](context.Background())

	// Collect first element
	err := funcBaseCollecInt32(ctx)
	assert.Nil(t, err)

	// Collect secomd element
	err = funcBaseCollecInt32(ctx)
	assert.Nil(t, err)

	result, err := aggregator.Aggregate[int32](ctx)
	assert.Equal(t, result, []int32{0, 0})
	assert.Nil(t, err)
}

func TestBaseContextAggregator_SuccessBuildKey(t *testing.T) {
	key := "test"
	ctx := aggregator.RegisterBaseContextAggregator[int32](context.Background(), key)

	// Collect first element
	err := funcBaseCollecInt32(ctx, key)
	assert.Nil(t, err)

	// Collect secomd element
	err = funcBaseCollecInt32(ctx, key)
	assert.Nil(t, err)

	result, err := aggregator.Aggregate[int32](ctx, key)
	assert.Equal(t, result, []int32{0, 0})
	assert.Nil(t, err)
}

func TestBaseContextAggregator_SuccessMultipleAggregator(t *testing.T) {
	key1 := "test1"
	key2 := "test2"
	ctx := context.Background()

	ctx = aggregator.RegisterBaseContextAggregator[int32](ctx, key1)
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx, key2)

	// Collect first element
	err := funcBaseCollecInt32(ctx, key1)
	assert.Nil(t, err)

	// Collect secomd element
	err = funcBaseCollecInt32(ctx, key1)
	assert.Nil(t, err)

	// Collect first element
	err = funcBaseCollecString(ctx, key2)
	assert.Nil(t, err)

	// Collect secomd element
	err = funcBaseCollecString(ctx, key2)
	assert.Nil(t, err)

	result1, err := aggregator.Aggregate[int32](ctx, key1)
	assert.Equal(t, result1, []int32{0, 0})
	assert.Nil(t, err)

	result2, err := aggregator.Aggregate[string](ctx, key2)
	assert.Equal(t, result2, []string{"", ""})
	assert.Nil(t, err)
}
