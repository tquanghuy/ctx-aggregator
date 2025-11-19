# ctx-aggregator

[![Go Reference](https://pkg.go.dev/badge/github.com/t-quanghuy/ctx-aggregator.svg)](https://pkg.go.dev/github.com/t-quanghuy/ctx-aggregator)
[![CI](https://github.com/t-quanghuy/ctx-aggregator/actions/workflows/ci.yml/badge.svg)](https://github.com/t-quanghuy/ctx-aggregator/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/t-quanghuy/ctx-aggregator)](https://goreportcard.com/report/github.com/t-quanghuy/ctx-aggregator)

`ctx-aggregator` is a Go library for collecting and aggregating data via `context.Context`. It supports both sequential and concurrent aggregation patterns, making it useful for gathering data across different layers of an application or concurrent goroutines without explicit parameter passing.

## Features

- **Context-based Aggregation**: Pass data implicitly through `context.Context`.
- **Sequential Aggregation**: Simple collection for linear flows.
- **Concurrent Aggregation**: Thread-safe collection for concurrent operations using `sync.Mutex` and `sync.WaitGroup`.
- **Type Safety**: Uses Go generics for type-safe data collection.
- **Multiple Aggregators**: Support for multiple aggregators in the same context using unique keys.

## Installation

```bash
go get github.com/t-quanghuy/ctx-aggregator
```

## Usage

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
			// AddWait/Done is handled internally if using the aggregator's wait mechanism,
			// but here we just use Collect safely.
			// Note: For strict synchronization ensuring all goroutines finish before Aggregate,
			// you should manage the goroutine lifecycle or use the aggregator's wait group if exposed/extended.
			// In this simple example, we wait for our own WaitGroup.
			_ = aggregator.Collect(ctx, val)
		}(i)
	}
	wg.Wait()

	results, _ := aggregator.Aggregate[int](ctx)
	fmt.Println(results) // Output: [0 1 2 3 4] (order may vary)
}
```

## License

MIT
