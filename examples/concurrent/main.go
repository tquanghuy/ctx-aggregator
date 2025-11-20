package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	aggregator "github.com/t-quanghuy/ctx-aggregator"
)

func main() {
	fmt.Println("=== Concurrent Aggregation Example ===")
	fmt.Println()

	// Create a context with a concurrent aggregator
	ctx := context.Background()
	ctx = aggregator.RegisterConcurrentContextAggregator[string](ctx)

	// Simulate multiple concurrent workers
	var wg sync.WaitGroup
	numWorkers := 5

	fmt.Printf("Starting %d concurrent workers...\n", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i+1, &wg)
	}

	// Wait for all workers to complete
	wg.Wait()

	// Aggregate all collected data
	results, err := aggregator.Aggregate[string](ctx)
	if err != nil {
		log.Fatalf("Failed to aggregate: %v", err)
	}

	fmt.Printf("\nCollected %d results from concurrent workers:\n", len(results))
	for i, result := range results {
		fmt.Printf("  %d. %s\n", i+1, result)
	}
}

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate some work
	time.Sleep(time.Duration(id*10) * time.Millisecond)

	message := fmt.Sprintf("Worker %d completed task", id)
	if err := aggregator.Collect(ctx, message); err != nil {
		log.Printf("Worker %d error: %v", id, err)
	} else {
		fmt.Printf("Worker %d: collected data\n", id)
	}
}
