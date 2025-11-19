// Batch Job Template
//
// This template provides a background job structure for building
// batch processing applications with the Canvus SDK. It includes:
// - Progress reporting
// - Error aggregation
// - Resumable operations with checkpointing
// - Concurrent processing with rate limiting
//
// Usage:
//   1. Copy this file to your project
//   2. Search for "TODO:" comments and customize
//   3. Build with: go build -o your-batch-job
//   4. Run with: ./your-batch-job
//
// Environment Variables:
//   CANVUS_API_URL    - Required: Canvus API endpoint
//   CANVUS_API_KEY    - Required: API key for authentication
//   BATCH_CHECKPOINT  - Optional: Checkpoint file for resume support

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
	// TODO: Add your imports here
)

// Config holds the batch job configuration
type Config struct {
	APIURL         string
	APIKey         string
	CheckpointFile string
	Concurrency    int
	BatchSize      int
	Timeout        time.Duration
	// TODO: Add your configuration fields here
}

// Checkpoint tracks job progress for resume support
type Checkpoint struct {
	LastProcessedID string    `json:"last_processed_id"`
	ProcessedCount  int       `json:"processed_count"`
	FailedCount     int       `json:"failed_count"`
	StartTime       time.Time `json:"start_time"`
	LastUpdate      time.Time `json:"last_update"`
	// TODO: Add your checkpoint fields here
}

// JobResult holds the result of processing a single item
type JobResult struct {
	ID      string
	Success bool
	Error   error
	Data    interface{}
}

// JobSummary holds the overall job summary
type JobSummary struct {
	TotalItems     int
	ProcessedItems int
	SuccessCount   int
	FailedCount    int
	SkippedCount   int
	Duration       time.Duration
	Errors         []string
}

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Handle shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutdown signal received, finishing current batch...")
		cancel()
	}()

	// Run the batch job
	summary, err := runBatchJob(ctx, cfg)
	if err != nil {
		log.Printf("Batch job failed: %v", err)
		os.Exit(1)
	}

	// Print summary
	printSummary(summary)

	if summary.FailedCount > 0 {
		os.Exit(1)
	}
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	cfg := &Config{
		Concurrency: 5,
		BatchSize:   100,
		Timeout:     30 * time.Minute,
	}

	// Required: API URL
	cfg.APIURL = os.Getenv("CANVUS_API_URL")
	if cfg.APIURL == "" {
		return nil, fmt.Errorf("CANVUS_API_URL environment variable is required")
	}

	// Required: API Key
	cfg.APIKey = os.Getenv("CANVUS_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("CANVUS_API_KEY environment variable is required")
	}

	// Optional: Checkpoint file
	cfg.CheckpointFile = os.Getenv("BATCH_CHECKPOINT")
	if cfg.CheckpointFile == "" {
		cfg.CheckpointFile = "batch_checkpoint.json"
	}

	// TODO: Load your environment variables here
	// Example:
	// if concurrency := os.Getenv("BATCH_CONCURRENCY"); concurrency != "" {
	//     cfg.Concurrency, _ = strconv.Atoi(concurrency)
	// }

	return cfg, nil
}

// runBatchJob executes the batch job
func runBatchJob(ctx context.Context, cfg *Config) (*JobSummary, error) {
	// Create SDK configuration
	sdkCfg := canvus.DefaultSessionConfig()
	sdkCfg.BaseURL = cfg.APIURL
	sdkCfg.RequestTimeout = 60 * time.Second

	// Create session with API key authentication
	session := canvus.NewSession(sdkCfg, canvus.WithAPIKey(cfg.APIKey))

	// Load checkpoint if exists
	checkpoint, err := loadCheckpoint(cfg.CheckpointFile)
	if err != nil {
		log.Printf("No checkpoint found, starting fresh: %v", err)
		checkpoint = &Checkpoint{
			StartTime:  time.Now(),
			LastUpdate: time.Now(),
		}
	} else {
		log.Printf("Resuming from checkpoint: %d items already processed", checkpoint.ProcessedCount)
	}

	// TODO: Fetch items to process
	// Example: Get all canvases
	// items, err := session.ListCanvases(ctx, nil)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to list canvases: %w", err)
	// }

	// Demo: Get all canvases to process
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list canvases: %w", err)
	}

	// Convert to items to process
	items := make([]interface{}, len(canvases))
	for i, c := range canvases {
		items[i] = c
	}

	// Filter items based on checkpoint (resume support)
	items = filterProcessedItems(items, checkpoint)

	summary := &JobSummary{
		TotalItems:     len(items),
		ProcessedItems: checkpoint.ProcessedCount,
	}

	if len(items) == 0 {
		log.Println("No items to process")
		return summary, nil
	}

	log.Printf("Processing %d items with concurrency %d", len(items), cfg.Concurrency)

	// Create channels for work distribution
	itemChan := make(chan interface{}, cfg.BatchSize)
	resultChan := make(chan *JobResult, cfg.BatchSize)

	// Create wait group for workers
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(ctx, session, workerID, itemChan, resultChan)
		}(i)
	}

	// Start result collector
	var collectorWg sync.WaitGroup
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		collectResults(cfg, checkpoint, summary, resultChan)
	}()

	// Send items to workers
	for _, item := range items {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping item distribution")
			break
		case itemChan <- item:
		}
	}
	close(itemChan)

	// Wait for workers to finish
	wg.Wait()
	close(resultChan)

	// Wait for result collector to finish
	collectorWg.Wait()

	// Calculate duration
	summary.Duration = time.Since(checkpoint.StartTime)

	// Save final checkpoint
	if err := saveCheckpoint(cfg.CheckpointFile, checkpoint); err != nil {
		log.Printf("Warning: failed to save checkpoint: %v", err)
	}

	return summary, nil
}

// worker processes items from the item channel
func worker(ctx context.Context, session *canvus.Session, workerID int, items <-chan interface{}, results chan<- *JobResult) {
	for item := range items {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result := processItem(ctx, session, item)
		results <- result
	}
}

// processItem processes a single item
func processItem(ctx context.Context, session *canvus.Session, item interface{}) *JobResult {
	// TODO: Replace with your processing logic
	//
	// Example: Update canvas
	// canvas := item.(*canvus.Canvas)
	// result := &JobResult{ID: canvas.ID}
	//
	// // Do something with the canvas
	// canvas.Description = "Updated by batch job"
	// _, err := session.UpdateCanvas(ctx, canvas.ID, canvas)
	// if err != nil {
	//     result.Success = false
	//     result.Error = err
	//     return result
	// }
	//
	// result.Success = true
	// return result

	// Demo: Just extract the canvas ID
	canvas, ok := item.(*canvus.Canvas)
	if !ok {
		return &JobResult{
			ID:      "unknown",
			Success: false,
			Error:   fmt.Errorf("invalid item type"),
		}
	}

	result := &JobResult{
		ID:      canvas.ID,
		Success: true,
		Data:    canvas,
	}

	// Simulate some work
	// TODO: Replace with actual processing
	time.Sleep(10 * time.Millisecond)

	return result
}

// collectResults collects results and updates checkpoint
func collectResults(cfg *Config, checkpoint *Checkpoint, summary *JobSummary, results <-chan *JobResult) {
	saveInterval := 10 // Save checkpoint every N items

	for result := range results {
		checkpoint.ProcessedCount++
		summary.ProcessedItems++

		if result.Success {
			summary.SuccessCount++
			checkpoint.LastProcessedID = result.ID
		} else {
			summary.FailedCount++
			checkpoint.FailedCount++
			if result.Error != nil {
				errMsg := fmt.Sprintf("%s: %v", result.ID, result.Error)
				summary.Errors = append(summary.Errors, errMsg)
				log.Printf("Error processing %s: %v", result.ID, result.Error)
			}
		}

		// Update checkpoint periodically
		if checkpoint.ProcessedCount%saveInterval == 0 {
			checkpoint.LastUpdate = time.Now()
			if err := saveCheckpoint(cfg.CheckpointFile, checkpoint); err != nil {
				log.Printf("Warning: failed to save checkpoint: %v", err)
			}

			// Log progress
			log.Printf("Progress: %d/%d items processed (%d success, %d failed)",
				summary.ProcessedItems, summary.TotalItems+checkpoint.ProcessedCount-summary.ProcessedItems,
				summary.SuccessCount, summary.FailedCount)
		}
	}
}

// filterProcessedItems filters out already processed items based on checkpoint
func filterProcessedItems(items []interface{}, checkpoint *Checkpoint) []interface{} {
	if checkpoint.LastProcessedID == "" {
		return items
	}

	// TODO: Implement your filtering logic
	// This example skips items until we find the last processed one
	//
	// found := false
	// filtered := make([]interface{}, 0)
	// for _, item := range items {
	//     canvas := item.(*canvus.Canvas)
	//     if found {
	//         filtered = append(filtered, item)
	//     } else if canvas.ID == checkpoint.LastProcessedID {
	//         found = true
	//     }
	// }
	// return filtered

	// For demo, just return all items
	return items
}

// loadCheckpoint loads checkpoint from file
func loadCheckpoint(filename string) (*Checkpoint, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, err
	}

	return &checkpoint, nil
}

// saveCheckpoint saves checkpoint to file
func saveCheckpoint(filename string, checkpoint *Checkpoint) error {
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// printSummary prints the job summary
func printSummary(summary *JobSummary) {
	separator := strings.Repeat("=", 50)
	fmt.Println("\n" + separator)
	fmt.Println("BATCH JOB SUMMARY")
	fmt.Println(separator)
	fmt.Printf("Total Items:     %d\n", summary.TotalItems)
	fmt.Printf("Processed:       %d\n", summary.ProcessedItems)
	fmt.Printf("Successful:      %d\n", summary.SuccessCount)
	fmt.Printf("Failed:          %d\n", summary.FailedCount)
	fmt.Printf("Skipped:         %d\n", summary.SkippedCount)
	fmt.Printf("Duration:        %s\n", summary.Duration.Round(time.Second))

	if len(summary.Errors) > 0 {
		fmt.Println("\nErrors:")
		for i, err := range summary.Errors {
			if i >= 10 {
				fmt.Printf("  ... and %d more errors\n", len(summary.Errors)-10)
				break
			}
			fmt.Printf("  - %s\n", err)
		}
	}
	fmt.Println(separator)
}
