# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-11-20

### Added

- **Initial stable release** of ctx-aggregator
- **Base Aggregator**: Sequential aggregation for linear data collection flows
- **Concurrent Aggregator**: Thread-safe aggregation using mutex protection for concurrent goroutines
- **Streaming Aggregator**: Real-time data processing with callback support during collection
- **Concurrent Streaming Aggregator**: Thread-safe streaming with callback support for concurrent operations
- **Advanced Operations**:
  - `AggregateWithFilter()`: Filter collected items based on predicates
  - `AggregateWithTransform()`: Transform items during aggregation
- **Performance Features**:
  - Capacity hints for pre-allocation (`RegisterBaseContextAggregatorWithCapacity`, `RegisterConcurrentContextAggregatorWithCapacity`)
  - Synchronization utilities (`WaitFunc`) for concurrent operations
- **Multiple Aggregators**: Support for multiple independent aggregators using unique keys
- **Type Safety**: Full support for Go generics (Go 1.18+)
- **Comprehensive Documentation**:
  - README with quick start guides and examples
  - Design documentation covering architecture and performance patterns
  - Example programs for sequential, concurrent, and multiple aggregator scenarios
- **Extensive Test Coverage**:
  - Base aggregator tests
  - Concurrent aggregator tests
  - Capacity hint tests
  - Filter and transform tests
  - Streaming aggregator tests
- **CI/CD Pipeline**: GitHub Actions workflow for automated testing and linting
- **Developer Tools**: Makefile with targets for testing, linting, coverage, and benchmarking

### Features

- Context-based aggregation without explicit parameter passing
- Sequential aggregation with minimal overhead
- Thread-safe concurrent aggregation with mutex protection
- Real-time streaming with callback support
- Flexible filtering and transformation operations
- Memory optimization through capacity hints
- Multiple independent aggregators in a single context
- Full type safety with Go generics

### Documentation

- Comprehensive README with examples
- Design documentation with architecture details
- Example programs demonstrating usage patterns
- API documentation via GoDoc

