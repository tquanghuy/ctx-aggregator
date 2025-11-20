package main

import (
	"context"
	"fmt"
	"log"

	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func main() {
	fmt.Println("=== Multiple Aggregators Example ===")
	fmt.Println()

	// Create a context with multiple aggregators for different types of data
	ctx := context.Background()

	// Register separate aggregators using different keys
	// Each key creates a separate aggregator instance
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx, "errors")
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx, "warnings")
	ctx = aggregator.RegisterBaseContextAggregator[string](ctx, "info")

	// Simulate processing with different types of messages
	ctx = processStep1(ctx)
	ctx = processStep2(ctx)
	ctx = processStep3(ctx)

	// Aggregate and display results from each aggregator
	displayResults(ctx)
}

func processStep1(ctx context.Context) context.Context {
	fmt.Println("Step 1: Initializing...")
	aggregator.Collect(ctx, "Step 1: System initialized", "info")
	return ctx
}

func processStep2(ctx context.Context) context.Context {
	fmt.Println("Step 2: Processing data...")
	aggregator.Collect(ctx, "Step 2: Data processing started", "info")
	aggregator.Collect(ctx, "Step 2: Deprecated API used", "warnings")
	return ctx
}

func processStep3(ctx context.Context) context.Context {
	fmt.Println("Step 3: Finalizing...")
	aggregator.Collect(ctx, "Step 3: Finalization complete", "info")
	aggregator.Collect(ctx, "Step 3: Failed to save cache", "errors")
	return ctx
}

func displayResults(ctx context.Context) {
	fmt.Println("\n=== Results ===")

	// Display errors
	errors, err := aggregator.Aggregate[string](ctx, "errors")
	if err != nil {
		log.Printf("Failed to aggregate errors: %v", err)
	} else {
		fmt.Printf("\nErrors (%d):\n", len(errors))
		for i, e := range errors {
			fmt.Printf("  %d. %s\n", i+1, e)
		}
	}

	// Display warnings
	warnings, err := aggregator.Aggregate[string](ctx, "warnings")
	if err != nil {
		log.Printf("Failed to aggregate warnings: %v", err)
	} else {
		fmt.Printf("\nWarnings (%d):\n", len(warnings))
		for i, w := range warnings {
			fmt.Printf("  %d. %s\n", i+1, w)
		}
	}

	// Display info messages
	infos, err := aggregator.Aggregate[string](ctx, "info")
	if err != nil {
		log.Printf("Failed to aggregate info: %v", err)
	} else {
		fmt.Printf("\nInfo Messages (%d):\n", len(infos))
		for i, info := range infos {
			fmt.Printf("  %d. %s\n", i+1, info)
		}
	}
}
