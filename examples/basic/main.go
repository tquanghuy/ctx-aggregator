package main

import (
	"context"
	"fmt"
	"log"

	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func main() {
	fmt.Println("=== Basic Sequential Aggregation Example ===")
	fmt.Println()

	// Create a context with a base aggregator
	ctx := context.Background()
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx)

	// Simulate collecting data through different layers
	ctx = collectFromLayer1(ctx)
	ctx = collectFromLayer2(ctx)
	ctx = collectFromLayer3(ctx)

	// Aggregate all collected data
	results, err := aggregator.Aggregate[string](ctx)
	if err != nil {
		log.Fatalf("Failed to aggregate: %v", err)
	}

	fmt.Println("\nCollected data:")
	for i, result := range results {
		fmt.Printf("  %d. %s\n", i+1, result)
	}
}

func collectFromLayer1(ctx context.Context) context.Context {
	fmt.Println("Layer 1: Processing request...")
	if err := aggregator.Collect(ctx, "Layer 1: Request received"); err != nil {
		log.Printf("Warning: %v", err)
	}
	return ctx
}

func collectFromLayer2(ctx context.Context) context.Context {
	fmt.Println("Layer 2: Validating data...")
	if err := aggregator.Collect(ctx, "Layer 2: Validation passed"); err != nil {
		log.Printf("Warning: %v", err)
	}
	return ctx
}

func collectFromLayer3(ctx context.Context) context.Context {
	fmt.Println("Layer 3: Processing business logic...")
	if err := aggregator.Collect(ctx, "Layer 3: Business logic executed"); err != nil {
		log.Printf("Warning: %v", err)
	}
	return ctx
}
