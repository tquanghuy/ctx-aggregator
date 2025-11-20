# Design Documentation

## Overview

`ctx-aggregator` is a Go library that enables implicit data collection across application layers using `context.Context`. This design document explains the architectural decisions, implementation details, and usage patterns.

## Motivation

### The Problem

In complex applications, you often need to collect data from multiple layers or concurrent operations:

- **Logging/Tracing**: Collecting log entries or trace spans across function calls
- **Error Collection**: Gathering multiple errors from parallel operations
- **Metrics**: Accumulating performance metrics throughout request processing
- **Audit Trails**: Recording actions taken during request handling

Traditional approaches require:
1. Explicit parameter passing (verbose, couples layers)
2. Global state (not thread-safe, testing difficulties)
3. Thread-local storage (not idiomatic in Go)

### The Solution

`ctx-aggregator` leverages Go's `context.Context` to provide:
- **Implicit passing**: No need to modify function signatures
- **Type safety**: Uses Go generics for compile-time type checking
- **Thread safety**: Concurrent aggregator for parallel operations
- **Multiple aggregators**: Support for different data types in the same context

## Architecture

### Core Components

#### 1. Aggregator Interface

```go
type Aggregator[T any] interface {
    Collect(data T) error
    Aggregate() ([]T, error)
}
```

All aggregators implement this interface, providing a consistent API regardless of the underlying implementation.

#### 2. Base Aggregator

**Use Case**: Sequential data collection in single-threaded contexts.

**Implementation**:
- Simple slice-based storage
- No synchronization overhead
- Minimal memory footprint

**When to Use**:
- Linear request processing
- Single-threaded operations
- Performance-critical paths where concurrency is not needed

#### 3. Concurrent Aggregator

**Use Case**: Thread-safe data collection from multiple goroutines.

**Implementation**:
- `sync.Mutex` for thread-safe access
- `sync.WaitGroup` for goroutine coordination (optional)
- Slice-based storage protected by mutex

**When to Use**:
- Parallel processing (fan-out patterns)
- Concurrent worker pools
- Any scenario with multiple goroutines collecting data

#### 4. Streaming Aggregator

**Use Case**: Real-time data processing as items are collected.

**Implementation**:
- Callback-based execution
- Synchronous callback invocation
- Optional thread safety (ConcurrentStreamingAggregator)

**When to Use**:
- Real-time logging/monitoring
- Event-driven architectures
- When immediate processing is required vs batch processing

### Context Integration

The library uses `context.Context` to store aggregator instances:

```go
type contextKey string

func RegisterBaseContextAggregator[T any](ctx context.Context) context.Context {
    key := contextKey(fmt.Sprintf("aggregator_%s", reflect.TypeOf((*T)(nil)).Elem()))
    agg := NewBaseAggregator[T]()
    return context.WithValue(ctx, key, agg)
}
```

**Key Design Decisions**:

1. **Type-based keys**: Default keys are derived from the type parameter, allowing one aggregator per type
2. **Custom keys**: `WithKey` variants support multiple aggregators of the same type
3. **Immutability**: Following context conventions, registration returns a new context

## Thread Safety

### Base Aggregator

**Not thread-safe** by design:
- No mutex overhead
- Faster for sequential operations
- Caller responsible for synchronization if used concurrently

### Concurrent Aggregator

**Thread-safe** implementation:
- All operations protected by `sync.Mutex`
- Safe for concurrent `Collect()` calls
- `Aggregate()` safely reads while collection may continue

**Synchronization Pattern**:
```

### Advanced Operations

#### Filtering and Transformation

The library supports functional operations during aggregation:

- **Filtering**: `AggregateWithFilter` uses a predicate to select items.
- **Transformation**: `AggregateWithTransform` converts items to a different type.
- **Combined**: `AggregateWithFilterAndTransform` performs both in a single pass.

These operations are performed during the aggregation phase, keeping the collection phase fast and simple.

### Performance Optimizations

#### Capacity Hints

When the expected number of items is known, use `WithCapacity` variants to pre-allocate memory:

```go
// Allocates underlying slice with capacity 100
ctx = aggregator.RegisterConcurrentContextAggregatorWithCapacity[int](ctx, 100)
```

This avoids slice resizing overhead during collection.go
type ConcurrentAggregator[T any] struct {
    mu   sync.Mutex
    data []T
    wg   sync.WaitGroup
}

func (a *ConcurrentAggregator[T]) Collect(data T) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.data = append(a.data, data)
    return nil
}
```

## Performance Characteristics

### Memory

- **Base Aggregator**: `O(n)` where n = number of collected items
- **Concurrent Aggregator**: `O(n)` + mutex overhead (~8 bytes)
- **Context overhead**: One pointer per aggregator in context chain

### Time Complexity

| Operation | Base Aggregator | Concurrent Aggregator |
|-----------|----------------|----------------------|
| Collect   | O(1) amortized | O(1) amortized + lock |
| Aggregate | O(n)          | O(n) + lock          |

### Benchmarks

Typical performance (on modern hardware):

- **Base Collect**: ~10-20 ns/op
- **Concurrent Collect**: ~50-100 ns/op (includes mutex)
- **Aggregate**: ~1-5 ns per item

The concurrent aggregator's overhead is minimal compared to typical goroutine scheduling costs (1-2 μs).

## Usage Patterns

### Pattern 1: Error Collection

```go
ctx = aggregator.RegisterBaseContextAggregator[error](ctx)

// Collect errors from multiple operations
for _, item := range items {
    if err := processItem(ctx, item); err != nil {
        aggregator.Collect(ctx, err)
    }
}

// Check if any errors occurred
errors, _ := aggregator.Aggregate[error](ctx)
if len(errors) > 0 {
    return fmt.Errorf("multiple errors: %v", errors)
}
```

### Pattern 2: Concurrent Metrics

```go
ctx = aggregator.RegisterConcurrentContextAggregator[Metric](ctx)

var wg sync.WaitGroup
for _, task := range tasks {
    wg.Add(1)
    go func(t Task) {
        defer wg.Done()
        metric := processTask(ctx, t)
        aggregator.Collect(ctx, metric)
    }(task)
}
wg.Wait()

metrics, _ := aggregator.Aggregate[Metric](ctx)
```

### Pattern 3: Multiple Aggregators

```go
const (
    errorsKey   = "errors"
    warningsKey = "warnings"
)

ctx = aggregator.RegisterBaseContextAggregatorWithKey[string](ctx, errorsKey)
ctx = aggregator.RegisterBaseContextAggregatorWithKey[string](ctx, warningsKey)

// Collect different types
aggregator.CollectWithKey(ctx, errorsKey, "critical error")
aggregator.CollectWithKey(ctx, warningsKey, "deprecated API")

// Retrieve separately
errors, _ := aggregator.AggregateWithKey[string](ctx, errorsKey)
warnings, _ := aggregator.AggregateWithKey[string](ctx, warningsKey)
```

## Design Trade-offs

### Why Context?

**Pros**:
- Idiomatic Go pattern
- Already passed through most application layers
- Supports cancellation and deadlines
- Immutable value semantics

**Cons**:
- Slight performance overhead (pointer chasing)
- Can be misused (context should not be stored in structs)
- Type assertions required internally

**Decision**: The benefits of implicit passing outweigh the minimal overhead for most use cases.

### Why Not Channels?

Channels are great for streaming data but:
- Require explicit goroutine management
- Need buffering decisions
- More complex error handling
- Overkill for simple collection

Aggregators are simpler for the "collect now, aggregate later" pattern.

### Why Generics?

**Before Go 1.18**: Would require `interface{}` and type assertions
**With Generics**: Compile-time type safety, better performance

The generic implementation provides type safety without runtime overhead.

## Testing Considerations

### Unit Testing

```go
func TestWithAggregator(t *testing.T) {
    ctx := aggregator.RegisterBaseContextAggregator[string](context.Background())
    
    aggregator.Collect(ctx, "test")
    
    results, err := aggregator.Aggregate[string](ctx)
    assert.NoError(t, err)
    assert.Equal(t, []string{"test"}, results)
}
```

### Concurrent Testing

```go
func TestConcurrent(t *testing.T) {
    ctx := aggregator.RegisterConcurrentContextAggregator[int](context.Background())
    
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(val int) {
            defer wg.Done()
            aggregator.Collect(ctx, val)
        }(i)
    }
    wg.Wait()
    
    results, _ := aggregator.Aggregate[int](ctx)
    assert.Len(t, results, 100)
}
```

## Conclusion

`ctx-aggregator` provides a simple, type-safe way to collect data across application layers using Go's context pattern. The design prioritizes:

- **Simplicity**: Easy to understand and use
- **Performance**: Minimal overhead for common cases
- **Safety**: Thread-safe when needed, type-safe always
- **Flexibility**: Multiple aggregators, custom keys, streaming support

For most use cases, the base aggregator is sufficient. Use the concurrent aggregator when collecting from multiple goroutines, and streaming aggregator for real-time processing.
