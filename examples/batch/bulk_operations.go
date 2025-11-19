// Package main demonstrates batch operations with the Canvus Go SDK.
//
// This example shows:
// - Using BatchProcessor for bulk operations
// - Automatic retry logic for failed operations
// - Concurrent operation patterns with configurable parallelism
// - Progress tracking with callbacks
// - Analyzing batch results and summaries
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run bulk_operations.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
	// Step 1: Get configuration from environment variables
	apiURL := os.Getenv("CANVUS_API_URL")
	apiKey := os.Getenv("CANVUS_API_KEY")

	if apiURL == "" {
		log.Fatal("CANVUS_API_URL environment variable is required")
	}
	if apiKey == "" {
		log.Fatal("CANVUS_API_KEY environment variable is required")
	}

	fmt.Println("Canvus Go SDK - Batch Operations Example")
	fmt.Println("=========================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Step 3: Create test canvases for batch operations
	fmt.Println("Setting up test environment...")
	fmt.Println("Creating source and target canvases for batch operations...")

	sourceCanvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
		Name: fmt.Sprintf("Batch Source Canvas %d", time.Now().Unix()),
	})
	if err != nil {
		log.Fatalf("Error creating source canvas: %v", err)
	}

	targetCanvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
		Name: fmt.Sprintf("Batch Target Canvas %d", time.Now().Unix()),
	})
	if err != nil {
		log.Fatalf("Error creating target canvas: %v", err)
	}

	fmt.Printf("Source canvas: %s (ID: %s)\n", sourceCanvas.Name, sourceCanvas.ID)
	fmt.Printf("Target canvas: %s (ID: %s)\n\n", targetCanvas.Name, targetCanvas.ID)

	// Ensure cleanup
	defer func() {
		fmt.Println("\nCleaning up: Deleting test canvases...")
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()

		if err := session.DeleteCanvas(cleanupCtx, sourceCanvas.ID); err != nil {
			log.Printf("Warning: Failed to delete source canvas: %v", err)
		}
		if err := session.DeleteCanvas(cleanupCtx, targetCanvas.ID); err != nil {
			log.Printf("Warning: Failed to delete target canvas: %v", err)
		}
		fmt.Println("Test canvases deleted successfully")
	}()

	// Step 4: Create multiple widgets for batch operations
	fmt.Println("Creating widgets for batch operations...")

	var createdWidgetIDs []string
	widgetCount := 10

	for i := 0; i < widgetCount; i++ {
		noteData := map[string]interface{}{
			"widget_type":      "note",
			"text":             fmt.Sprintf("Batch test note #%d", i+1),
			"title":            fmt.Sprintf("Note %d", i+1),
			"background_color": "#4CAF50",
			"location": map[string]interface{}{
				"x": 100.0 + float64(i%5)*250.0,
				"y": 100.0 + float64(i/5)*200.0,
			},
			"size": map[string]interface{}{
				"width":  200.0,
				"height": 150.0,
			},
		}

		note, err := session.CreateNote(ctx, sourceCanvas.ID, noteData)
		if err != nil {
			log.Fatalf("Error creating widget %d: %v", i+1, err)
		}
		createdWidgetIDs = append(createdWidgetIDs, note.ID)
	}

	fmt.Printf("Created %d widgets in source canvas\n\n", len(createdWidgetIDs))

	// Get widgets as Widget type for batch operations
	widgets, err := session.ListWidgets(ctx, sourceCanvas.ID, nil)
	if err != nil {
		log.Fatalf("Error listing widgets: %v", err)
	}

	// Step 5: Configure BatchProcessor with custom settings
	fmt.Println("Configuring BatchProcessor...")

	// Create a custom batch configuration
	batchConfig := &canvus.BatchConfig{
		// MaxConcurrency controls how many operations run in parallel
		// Higher values = faster but more load on server
		MaxConcurrency: 5,

		// Timeout is the overall timeout for the entire batch
		Timeout: 2 * time.Minute,

		// RetryAttempts specifies how many times to retry failed operations
		RetryAttempts: 3,

		// RetryDelay is the base delay between retry attempts
		// Actual delay may include jitter to prevent thundering herd
		RetryDelay: 500 * time.Millisecond,

		// ContinueOnError allows the batch to continue even if some operations fail
		ContinueOnError: true,

		// ProgressCallback is called after each operation completes
		// Use this for progress tracking in UIs or logs
		ProgressCallback: func(completed, total int, results []*canvus.BatchResult) {
			percentage := float64(completed) / float64(total) * 100
			fmt.Printf("Progress: %d/%d (%.0f%%)\n", completed, total, percentage)
		},
	}

	// Create the BatchProcessor with our configuration
	processor := canvus.NewBatchProcessor(session, batchConfig)

	fmt.Printf("Batch configured: %d concurrent operations, %d retries, %.0f second timeout\n\n",
		batchConfig.MaxConcurrency, batchConfig.RetryAttempts, batchConfig.Timeout.Seconds())

	// Step 6: Build batch operations using the fluent builder
	fmt.Println("Building batch operations (copy widgets to target canvas)...")

	builder := canvus.NewBatchOperationBuilder()

	// Add copy operations for each widget
	for i, widget := range widgets {
		// Each operation needs a unique ID for tracking
		operationID := fmt.Sprintf("copy-widget-%d", i+1)

		// Copy the widget to the target canvas
		// We need to pass a pointer to the widget
		w := widget // Create a copy to take address
		builder.Copy(operationID, &w, targetCanvas.ID)
	}

	// Build the operations list
	operations := builder.Build()
	fmt.Printf("Built %d batch operations\n\n", len(operations))

	// Step 7: Execute the batch operations
	fmt.Println("Executing batch operations...")
	fmt.Println("=============================")

	startTime := time.Now()
	results, err := processor.ExecuteBatch(ctx, operations)
	duration := time.Since(startTime)

	if err != nil {
		// Even with errors, we might have partial results
		log.Printf("Batch execution error: %v", err)
	}

	fmt.Printf("\nBatch execution completed in %v\n\n", duration)

	// Step 8: Analyze batch results
	fmt.Println("Analyzing batch results...")

	summary := canvus.Summarize(results)

	fmt.Printf("Summary:\n")
	fmt.Printf("  Total Operations: %d\n", summary.TotalOperations)
	fmt.Printf("  Successful: %d\n", summary.Successful)
	fmt.Printf("  Failed: %d\n", summary.Failed)
	fmt.Printf("  Total Duration: %v\n", summary.TotalDuration)
	fmt.Printf("  Average Duration: %v\n", summary.AverageDuration)
	fmt.Println()

	// Report failed operations
	if len(summary.FailedOperations) > 0 {
		fmt.Println("Failed operations:")
		for _, failed := range summary.FailedOperations {
			fmt.Printf("  - %s: %v (retries: %d)\n",
				failed.OperationID, failed.Error, failed.Retries)
		}
		fmt.Println()
	}

	// Verify copy operations by listing widgets in target canvas
	targetWidgets, err := session.ListWidgets(ctx, targetCanvas.ID, nil)
	if err != nil {
		log.Printf("Warning: Could not verify target canvas: %v", err)
	} else {
		fmt.Printf("Target canvas now has %d widgets (expected: %d)\n\n",
			len(targetWidgets), summary.Successful)
	}

	// Step 9: Demonstrate batch delete operations
	fmt.Println("Demonstrating batch delete operations...")

	// Build delete operations for source canvas widgets
	deleteBuilder := canvus.NewBatchOperationBuilder()

	for i, widget := range widgets {
		operationID := fmt.Sprintf("delete-widget-%d", i+1)

		// Create a delete operation
		// Note: For delete operations, we need to provide metadata
		w := widget // Create a copy to take address
		deleteOp := &canvus.BatchOperation{
			ID:       operationID,
			Type:     canvus.BatchOperationDelete,
			Resource: &w,
			Metadata: map[string]interface{}{
				"canvas_id":   sourceCanvas.ID,
				"widget_type": widget.WidgetType,
			},
		}
		deleteBuilder.Delete(operationID, &w)
		// Replace the last operation with our properly configured one
		ops := deleteBuilder.Build()
		ops[len(ops)-1] = deleteOp
	}

	deleteOps := deleteBuilder.Build()

	// Configure delete operation metadata
	for i := range deleteOps {
		if deleteOps[i].Type == canvus.BatchOperationDelete {
			if deleteOps[i].Metadata == nil {
				deleteOps[i].Metadata = make(map[string]interface{})
			}
			deleteOps[i].Metadata["canvas_id"] = sourceCanvas.ID
			deleteOps[i].Metadata["widget_type"] = "note"
		}
	}

	fmt.Printf("Executing %d delete operations...\n", len(deleteOps))

	// Create a new processor with a progress counter
	var deleteCompleted int32
	deleteConfig := &canvus.BatchConfig{
		MaxConcurrency:  5,
		Timeout:         time.Minute,
		RetryAttempts:   2,
		RetryDelay:      200 * time.Millisecond,
		ContinueOnError: true,
		ProgressCallback: func(completed, total int, _ []*canvus.BatchResult) {
			atomic.StoreInt32(&deleteCompleted, int32(completed))
		},
	}

	deleteProcessor := canvus.NewBatchProcessor(session, deleteConfig)
	deleteResults, err := deleteProcessor.ExecuteBatch(ctx, deleteOps)
	if err != nil {
		log.Printf("Delete batch error: %v", err)
	}

	deleteSummary := canvus.Summarize(deleteResults)
	fmt.Printf("Delete results: %d successful, %d failed\n\n",
		deleteSummary.Successful, deleteSummary.Failed)

	// Step 10: Demonstrate using default configuration
	fmt.Println("Note: You can also use default configuration:")
	fmt.Println("  processor := canvus.NewBatchProcessor(session, nil)")
	fmt.Println("  // Uses DefaultBatchConfig() automatically")
	fmt.Println()

	defaultConfig := canvus.DefaultBatchConfig()
	fmt.Printf("Default configuration:\n")
	fmt.Printf("  MaxConcurrency: %d\n", defaultConfig.MaxConcurrency)
	fmt.Printf("  Timeout: %v\n", defaultConfig.Timeout)
	fmt.Printf("  RetryAttempts: %d\n", defaultConfig.RetryAttempts)
	fmt.Printf("  RetryDelay: %v\n", defaultConfig.RetryDelay)
	fmt.Printf("  ContinueOnError: %v\n", defaultConfig.ContinueOnError)

	// Best practices summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Batch Operations Best Practices:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("1. Set appropriate MaxConcurrency based on server capacity")
	fmt.Println("   - Start with 5-10 and adjust based on performance")
	fmt.Println("   - Higher values increase throughput but also load")
	fmt.Println()
	fmt.Println("2. Use RetryAttempts for handling transient failures")
	fmt.Println("   - 2-3 retries is usually sufficient")
	fmt.Println("   - SDK uses exponential backoff automatically")
	fmt.Println()
	fmt.Println("3. Set ContinueOnError based on your requirements")
	fmt.Println("   - true: process all operations, collect failures")
	fmt.Println("   - false: stop on first failure (faster fail)")
	fmt.Println()
	fmt.Println("4. Use ProgressCallback for user feedback")
	fmt.Println("   - Update progress bars or logs")
	fmt.Println("   - Track intermediate results")
	fmt.Println()
	fmt.Println("5. Always analyze BatchSummary for results")
	fmt.Println("   - Check for partial failures")
	fmt.Println("   - Log failed operations for investigation")

	fmt.Println("\nBatch operations example completed successfully!")
}
