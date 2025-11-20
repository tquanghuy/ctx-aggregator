# Examples

This directory contains example programs demonstrating different use cases of the `ctx-aggregator` library.

## Running the Examples

Each example is a standalone Go program that can be run using `go run`:

### Basic Sequential Aggregation

Demonstrates simple data collection through multiple layers of an application.

```bash
cd examples/basic
go run main.go
```

This example shows:
- Registering a base aggregator
- Collecting data sequentially through different layers
- Aggregating all collected data at the end

### Concurrent Aggregation

Shows thread-safe data collection from multiple goroutines.

```bash
cd examples/concurrent
go run main.go
```

This example demonstrates:
- Using the concurrent aggregator for thread-safe operations
- Collecting data from multiple goroutines simultaneously
- Proper synchronization with `sync.WaitGroup`

### Multiple Aggregators

Illustrates using multiple aggregators in the same context for different types of data.

```bash
cd examples/multiple
go run main.go
```

This example shows:
- Registering multiple aggregators with custom keys
- Collecting different types of data (errors, warnings, info) separately
- Aggregating each type independently

## Running All Examples

You can run all examples at once using the Makefile from the project root:

```bash
cd ..
make examples
```

## Integration with Your Project

To use `ctx-aggregator` in your own project:

1. Install the library:
   ```bash
   go get github.com/t-quanghuy/ctx-aggregator
   ```

2. Import it in your code:
   ```go
   import "github.com/t-quanghuy/ctx-aggregator"
   ```

3. Choose the appropriate aggregator type:
   - Use `RegisterBaseContextAggregator` for sequential operations
   - Use `RegisterConcurrentContextAggregator` for concurrent operations
   - Use `RegisterBaseContextAggregatorWithKey` or `RegisterConcurrentContextAggregatorWithKey` for multiple aggregators

For more details, see the [main README](../README.md) and the [GoDoc documentation](https://pkg.go.dev/github.com/t-quanghuy/ctx-aggregator).
