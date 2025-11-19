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
