package aggregator_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func TestConcurrentContextAggregator_CollectNotFoundAggregator(t *testing.T) {
	err := funcBaseCollecInt32(context.Background())
	assert.Equal(t, err, aggregator.ErrNotFoundAggregator)
}

func TestConcurrentContextAggregator_CollectInvalidType(t *testing.T) {
	ctx := aggregator.RegisterConcurrentContextAggregator[int](context.Background())
	err := funcBaseCollecInt32(ctx)
	assert.Equal(t, err, aggregator.ErrInvalidType)
}

func TestConcurrentContextAggregator_AggregateNotFoundAggregator(t *testing.T) {
	result, err := aggregator.Aggregate[int](context.Background())
	assert.Nil(t, result)
	assert.Equal(t, err, aggregator.ErrNotFoundAggregator)
}

func TestConcurrentContextAggregator_AggregateInvalidType(t *testing.T) {
	ctx := aggregator.RegisterConcurrentContextAggregator[int32](context.Background())
	err := funcBaseCollecInt32(ctx)
	assert.Nil(t, err)

	result, err := aggregator.Aggregate[int](ctx)
	assert.Nil(t, result)
	assert.Equal(t, err, aggregator.ErrInvalidType)
}

func TestConcurrentContextAggregator_SuccessNoKey(t *testing.T) {
	ctx := aggregator.RegisterConcurrentContextAggregator[int32](context.Background())

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

func TestConcurrentContextAggregator_SuccessBuildKey(t *testing.T) {
	key := "test"
	ctx := aggregator.RegisterConcurrentContextAggregator[int32](context.Background(), key)

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

func TestConcurrentContextAggregator_SuccessMultipleAggregator(t *testing.T) {
	key1 := "test1"
	key2 := "test2"
	ctx := context.Background()

	ctx = aggregator.RegisterConcurrentContextAggregator[int32](ctx, key1)
	ctx = aggregator.RegisterConcurrentContextAggregator[string](ctx, key2)

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

func TestConcurrentContextAggregator_SuccessConcurrent(t *testing.T) {
	key := "test"
	ctx := context.Background()

	ctx = aggregator.RegisterConcurrentContextAggregator[int32](ctx, key)

	var (
		err1, err2, err3 error
		wg               sync.WaitGroup
	)

	wg.Add(3)
	// Collect first element
	go func() {
		err1 = funcBaseCollecInt32WithVal(ctx, 1, key)
		wg.Done()
	}()

	// Collect second element
	go func() {
		err2 = funcBaseCollecInt32WithVal(ctx, 2, key)
		wg.Done()
	}()

	// Collect first element
	go func() {
		err3 = funcBaseCollecInt32WithVal(ctx, 3, key)
		wg.Done()
	}()

	wg.Wait()

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)

	result, err := aggregator.Aggregate[int32](ctx, key)
	assert.ElementsMatch(t, result, []int32{1, 2, 3})
	assert.Nil(t, err)
}

// func TestConcurrentContextAggregator_SuccessConcurrent_WaitFunc(t *testing.T) {
// 	key := "test"
// 	ctx := context.Background()

// 	ctx = RegisterConcurrentContextAggregator[int32](ctx, key)

// 	f := func(ctx context.Context, val int32, keys ...string) {
// 		ctx, fn := WaitFunc(ctx, keys...)
// 		defer fn()

// 		err := Collect(ctx, val, keys...)
// 		assert.Nil(t, err)
// 	}

// 	// Collect first element
// 	go f(ctx, 1, key)

// 	// Collect second element
// 	go f(ctx, 2, key)

// 	// Collect first element
// 	go f(ctx, 3, key)

// 	result, err := Aggregate[int32](ctx, key)
// 	assert.ElementsMatch(t, result, []int32{1, 2, 3})
// 	assert.Nil(t, err)
// }

// func TestConcurrentContextAggregator_SuccessConcurrent_WaitContextFinalizer(t *testing.T) {
// 	key := "test"
// 	ctx := context.Background()

// 	ctx = RegisterConcurrentContextAggregator[int32](ctx, key)

// 	// Collect first element
// 	go func() {
// 		err := funcBaseCollecInt32WithVal(WaitContextFinalizer(ctx, key), 1, key)
// 		assert.Nil(t, err)
// 	}()

// 	// Collect second element
// 	go func() {
// 		err := funcBaseCollecInt32WithVal(WaitContextFinalizer(ctx, key), 2, key)
// 		assert.Nil(t, err)
// 	}()

// 	// Collect first element
// 	go func() {
// 		err := funcBaseCollecInt32WithVal(WaitContextFinalizer(ctx, key), 3, key)
// 		assert.Nil(t, err)
// 	}()

// 	result, err := Aggregate[int32](ctx, key)
// 	assert.ElementsMatch(t, result, []int32{1, 2, 3})
// 	assert.Nil(t, err)
// }
