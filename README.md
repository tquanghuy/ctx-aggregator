# ctx-aggregator

[![Go Reference](https://pkg.go.dev/badge/github.com/t-quanghuy/ctx-aggregator.svg)](https://pkg.go.dev/github.com/t-quanghuy/ctx-aggregator)
[![Go](https://github.com/tquanghuy/ctx-aggregator/actions/workflows/go.yml/badge.svg)](https://github.com/tquanghuy/ctx-aggregator/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`ctx-aggregator` is a Go library for collecting and aggregating data via `context.Context`. It supports both sequential and concurrent aggregation patterns, making it useful for gathering data across different layers of an application or concurrent goroutines without explicit parameter passing.

## Features

- **Context-based Aggregation**: Pass data implicitly through `context.Context`.
- **Sequential Aggregation**: Simple collection for linear flows with minimal overhead.
- **Concurrent Aggregation**: Thread-safe collection for concurrent operations with mutex protection.
- **Streaming Aggregation**: Process data in real-time with callbacks during collection.
- **Concurrent Streaming**: Thread-safe real-time processing with callback support.
- **Advanced Operations**: Filter and transform data during aggregation without intermediate allocations.
- **Performance Optimization**: Capacity hints for pre-allocating memory when item count is predictable.
- **Type Safety**: Uses Go generics (Go 1.18+) for compile-time type safety.
- **Multiple Aggregators**: Support for multiple independent aggregators in the same context using unique keys.

## Installation

```bash
go get github.com/t-quanghuy/ctx-aggregator
```

## Quick Start

### Sequential Aggregation

```go
package main

import (
	"context"
	"fmt"
	"github.com/t-quanghuy/ctx-aggregator"
)

func main() {
	ctx := context.Background()
	
	// Register the aggregator in the context
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx)

	// Collect data
	_ = aggregator.Collect(ctx, "hello")
	_ = aggregator.Collect(ctx, "world")

	// Aggregate data
	results, _ := aggregator.Aggregate[string](ctx)
	fmt.Println(results) // Output: [hello world]
}
```

### Concurrent Aggregation

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"github.com/t-quanghuy/ctx-aggregator"
)

func main() {
	ctx := context.Background()
	
	// Register the concurrent aggregator
	ctx = aggregator.RegisterConcurrentContextAggregator[int](ctx)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = aggregator.Collect(ctx, val)
		}(i)
	}
	wg.Wait()

	results, _ := aggregator.Aggregate[int](ctx)
	fmt.Println(results) // Output: [0 1 2 3 4] (order may vary)
}
```

### Multiple Aggregators

Use different keys to maintain multiple aggregators in the same context:

```go
ctx = aggregator.RegisterBaseContextAggregator[string](ctx, "errors")
ctx = aggregator.RegisterBaseContextAggregator[string](ctx, "warnings")

aggregator.Collect(ctx, "critical error", "errors")
aggregator.Collect(ctx, "deprecated API", "warnings")

errors, _ := aggregator.Aggregate[string](ctx, "errors")
warnings, _ := aggregator.Aggregate[string](ctx, "warnings")
```

### Streaming Aggregation

Process data immediately as it is collected using callbacks:

```go
// Define a callback function
callback := func(data string) {
    fmt.Printf("Received: %s\n", data)
}

// Register streaming aggregator
ctx = aggregator.RegisterStreamingAggregator[string](ctx, callback)

// Data is processed immediately
aggregator.Collect(ctx, "stream item 1")
aggregator.Collect(ctx, "stream item 2")
```

### Concurrent Streaming Aggregation

Thread-safe streaming aggregation for concurrent collection with real-time callbacks:

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"github.com/t-quanghuy/ctx-aggregator"
)

func main() {
	ctx := context.Background()
	
	// Define a thread-safe callback
	callback := func(data int) {
		fmt.Printf("Processing: %d\n", data)
	}
	
	// Register concurrent streaming aggregator
	ctx = aggregator.RegisterConcurrentStreamingAggregator[int](ctx, callback)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = aggregator.Collect(ctx, val)
		}(i)
	}
	wg.Wait()

	results, _ := aggregator.Aggregate[int](ctx)
	fmt.Println(results) // All items collected and callbacks executed
}
```

### Advanced Aggregation

#### Filtering

Retrieve only items that match a specific condition:

```go
// Filter only even numbers
filter := func(i int) bool {
    return i%2 == 0
}

results, _ := aggregator.AggregateWithFilter(ctx, filter)
```

#### Transformation

Transform items during aggregation:

```go
// Transform int to string
transform := func(i int) string {
    return fmt.Sprintf("Value: %d", i)
}

results, _ := aggregator.AggregateWithTransform(ctx, transform)
```

#### Capacity Hints

Optimize performance by pre-allocating memory when the expected number of items is known:

```go
// Pre-allocate slice with capacity 100
ctx = aggregator.RegisterConcurrentContextAggregatorWithCapacity[int](ctx, 100)
```

#### Synchronization with Concurrent Operations

For concurrent aggregators, use `WaitFunc` to ensure all goroutines complete before retrieving results:

```go
ctx = aggregator.RegisterConcurrentContextAggregator[int](ctx)

// Get wait function to track goroutine completion
ctx, waitDone := aggregator.WaitFunc(ctx)

var wg sync.WaitGroup
for i := 0; i < 5; i++ {
	wg.Add(1)
	go func(val int) {
		defer wg.Done()
		defer waitDone()
		aggregator.Collect(ctx, val)
	}(i)
}
wg.Wait()

// All data is collected and ready
results, _ := aggregator.Aggregate[int](ctx)
```

## Examples

For more detailed examples, see the [examples](./examples) directory:

- **[Basic](./examples/basic)**: Sequential aggregation through multiple layers
- **[Concurrent](./examples/concurrent)**: Thread-safe collection from multiple goroutines
- **[Multiple Aggregators](./examples/multiple)**: Using multiple aggregators for different data types

Run all examples:
```bash
make examples
```

## Documentation

- **[GoDoc](https://pkg.go.dev/github.com/t-quanghuy/ctx-aggregator)**: API reference
- **[Design Documentation](./docs/design.md)**: Architecture, performance, and usage patterns
- **[Examples](./examples/README.md)**: Detailed example programs

## Project Structure

```
ctx-aggregator/
├── aggregator.go              # Main aggregator interface and utility functions
├── base.go                    # Sequential aggregator implementation
├── concurrent.go              # Thread-safe aggregator implementation
├── streaming.go               # Streaming aggregators with callback support
├── doc.go                     # Package documentation
├── *_test.go                  # Comprehensive test files
├── examples/                  # Example programs
│   ├── README.md              # Examples documentation
│   ├── basic/                 # Sequential aggregation example
│   ├── concurrent/            # Concurrent aggregation example
│   └── multiple/              # Multiple aggregators example
├── docs/                      # Project documentation
│   └── design.md              # Design, architecture, and performance docs
├── vendor/                    # Vendored dependencies (testify, yaml)
├── .github/workflows/         # CI/CD workflows
├── Makefile                   # Build and development tasks
├── go.mod & go.sum            # Go module definition
├── CONTRIBUTING.md            # Contribution guidelines
├── LICENSE                    # MIT License
└── README.md                  # This file
```

## Development

### Running Tests

```bash
make test
```

### Running Linter

```bash
make lint
```

### Generate Coverage Report

```bash
make coverage
```

### Run Benchmarks

```bash
make bench
```

### Run All Checks

```bash
make check
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## License

MIT - See [LICENSE](./LICENSE) for details.
