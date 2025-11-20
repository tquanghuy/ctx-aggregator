package aggregator

import (
	"context"
	"errors"
	"strings"
)

type contextKey string

const (
	contextAggregatorContextKey contextKey = "ctxAggCtxKey"
)

// ERRORS
var (
	ErrNotFoundAggregator = errors.New("not found aggregator")
	ErrInvalidType        = errors.New("invalid type of aggregator")
)

// FilterFunc is a predicate function that returns true if the item should be included
type FilterFunc[T any] func(T) bool

// TransformFunc transforms an item of type T to type R
type TransformFunc[T any, R any] func(T) R

type ContextAggregator[T any] interface {
	Collect(data T)
	Aggregate() []T
}

func Collect[T any](ctx context.Context, data T, keys ...string) error {
	agg, err := extractAggregator[T](ctx, keys...)
	if err != nil {
		return err
	}

	agg.Collect(data)
	return nil
}

func Aggregate[T any](ctx context.Context, keys ...string) ([]T, error) {
	agg, err := extractAggregator[T](ctx, keys...)
	if err != nil {
		return nil, err
	}

	return agg.Aggregate(), nil
}

// AggregateWithFilter aggregates only items that match the filter predicate
func AggregateWithFilter[T any](ctx context.Context, filter FilterFunc[T], keys ...string) ([]T, error) {
	agg, err := extractAggregator[T](ctx, keys...)
	if err != nil {
		return nil, err
	}

	allItems := agg.Aggregate()
	filtered := make([]T, 0, len(allItems))
	for _, item := range allItems {
		if filter(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

// AggregateWithTransform aggregates and transforms items from type T to type R
func AggregateWithTransform[T any, R any](ctx context.Context, transform TransformFunc[T, R], keys ...string) ([]R, error) {
	agg, err := extractAggregator[T](ctx, keys...)
	if err != nil {
		return nil, err
	}

	allItems := agg.Aggregate()
	transformed := make([]R, 0, len(allItems))
	for _, item := range allItems {
		transformed = append(transformed, transform(item))
	}

	return transformed, nil
}

// AggregateWithFilterAndTransform filters and transforms items in a single pass
func AggregateWithFilterAndTransform[T any, R any](ctx context.Context, filter FilterFunc[T], transform TransformFunc[T, R], keys ...string) ([]R, error) {
	agg, err := extractAggregator[T](ctx, keys...)
	if err != nil {
		return nil, err
	}

	allItems := agg.Aggregate()
	result := make([]R, 0, len(allItems))
	for _, item := range allItems {
		if filter(item) {
			result = append(result, transform(item))
		}
	}

	return result, nil
}

// buildContextKey builds context key from default context key and input keys
func buildContextKey(keys ...string) contextKey {
	if len(keys) == 0 {
		return contextAggregatorContextKey
	}

	keys = append([]string{string(contextAggregatorContextKey)}, keys...)
	return contextKey(strings.Join(keys, "_"))
}

func extractAggregator[T any](ctx context.Context, keys ...string) (ContextAggregator[T], error) {
	ctxKey := buildContextKey(keys...)
	aggVal := ctx.Value(ctxKey)
	if aggVal == nil {
		return nil, ErrNotFoundAggregator
	}

	agg, ok := aggVal.(ContextAggregator[T])
	if !ok {
		return nil, ErrInvalidType
	}

	return agg, nil
}
