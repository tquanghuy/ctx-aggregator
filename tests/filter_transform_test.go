package aggregator_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func TestAggregateWithFilter(t *testing.T) {
	ctx := context.Background()
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx)

	// Collect mixed data
	_ = aggregator.Collect(ctx, "ERROR: something went wrong")
	_ = aggregator.Collect(ctx, "INFO: normal operation")
	_ = aggregator.Collect(ctx, "ERROR: another error")
	_ = aggregator.Collect(ctx, "WARN: warning message")

	// Filter only errors
	errors, err := aggregator.AggregateWithFilter(ctx, func(s string) bool {
		return strings.HasPrefix(s, "ERROR:")
	})

	assert.NoError(t, err)
	assert.Len(t, errors, 2)
	assert.Equal(t, "ERROR: something went wrong", errors[0])
	assert.Equal(t, "ERROR: another error", errors[1])
}

func TestAggregateWithTransform(t *testing.T) {
	ctx := context.Background()
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx)

	// Collect strings
	_ = aggregator.Collect(ctx, "hello")
	_ = aggregator.Collect(ctx, "world")
	_ = aggregator.Collect(ctx, "test")

	// Transform to lengths
	lengths, err := aggregator.AggregateWithTransform[string, int](ctx, func(s string) int {
		return len(s)
	})

	assert.NoError(t, err)
	assert.Len(t, lengths, 3)
	assert.Equal(t, 5, lengths[0]) // "hello"
	assert.Equal(t, 5, lengths[1]) // "world"
	assert.Equal(t, 4, lengths[2]) // "test"
}

func TestAggregateWithFilterAndTransform(t *testing.T) {
	ctx := context.Background()
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx)

	// Collect mixed data
	_ = aggregator.Collect(ctx, "ERROR: short")
	_ = aggregator.Collect(ctx, "INFO: information")
	_ = aggregator.Collect(ctx, "ERROR: this is a longer error message")
	_ = aggregator.Collect(ctx, "WARN: warning")

	// Filter errors and transform to lengths
	errorLengths, err := aggregator.AggregateWithFilterAndTransform(
		ctx,
		func(s string) bool {
			return strings.HasPrefix(s, "ERROR:")
		},
		func(s string) int {
			return len(s)
		},
	)

	assert.NoError(t, err)
	assert.Len(t, errorLengths, 2)
	assert.Equal(t, 12, errorLengths[0]) // "ERROR: short"
	assert.Equal(t, 37, errorLengths[1]) // "ERROR: this is a longer error message"
}

func TestFilterWithConcurrentAggregator(t *testing.T) {
	ctx := context.Background()
	ctx = aggregator.RegisterConcurrentContextAggregator[int](ctx)

	// Collect numbers
	for i := 0; i < 10; i++ {
		_ = aggregator.Collect(ctx, i)
	}

	// Filter even numbers
	evens, err := aggregator.AggregateWithFilter(ctx, func(n int) bool {
		return n%2 == 0
	})

	assert.NoError(t, err)
	assert.Len(t, evens, 5)
}

func TestTransformWithDifferentTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	ctx := context.Background()
	ctx = aggregator.RegisterBaseContextAggregator[User](ctx)

	// Collect users
	_ = aggregator.Collect(ctx, User{ID: 1, Name: "Alice"})
	_ = aggregator.Collect(ctx, User{ID: 2, Name: "Bob"})
	_ = aggregator.Collect(ctx, User{ID: 3, Name: "Charlie"})

	// Transform to IDs
	ids, err := aggregator.AggregateWithTransform[User, int](ctx, func(u User) int {
		return u.ID
	})

	assert.NoError(t, err)
	assert.Len(t, ids, 3)
	assert.Equal(t, []int{1, 2, 3}, ids)

	// Transform to names
	names, err := aggregator.AggregateWithTransform[User, string](ctx, func(u User) string {
		return u.Name
	})

	assert.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, names)
}

func TestEmptyFilterResult(t *testing.T) {
	ctx := context.Background()
	ctx = aggregator.RegisterBaseContextAggregator[int](ctx)

	// Collect some numbers
	_ = aggregator.Collect(ctx, 1)
	_ = aggregator.Collect(ctx, 2)
	_ = aggregator.Collect(ctx, 3)

	// Filter for numbers > 10 (none match)
	results, err := aggregator.AggregateWithFilter(ctx, func(n int) bool {
		return n > 10
	})

	assert.NoError(t, err)
	assert.Len(t, results, 0)
	assert.NotNil(t, results) // Should be empty slice, not nil
}
