// Package main demonstrates context usage patterns with the Canvus Go SDK.
//
// This example shows:
// - Context cancellation patterns
// - Request timeouts with context.WithTimeout
// - Deadline management with context.WithDeadline
// - Graceful shutdown handling
// - Coordinating multiple operations with contexts
// - Propagating cancellation to child operations
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run cancellation_and_timeouts.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	fmt.Println("Canvus Go SDK - Context, Cancellation, and Timeout Patterns")
	fmt.Println("============================================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Step 3: Demonstrate basic timeout usage
	fmt.Println("Pattern 1: Basic Request Timeout")
	fmt.Println("=================================")

	demonstrateBasicTimeout(session)
	fmt.Println()

	// Step 4: Demonstrate manual cancellation
	fmt.Println("Pattern 2: Manual Cancellation")
	fmt.Println("==============================")

	demonstrateManualCancellation(session)
	fmt.Println()

	// Step 5: Demonstrate deadline management
	fmt.Println("Pattern 3: Deadline Management")
	fmt.Println("==============================")

	demonstrateDeadline(session)
	fmt.Println()

	// Step 6: Demonstrate context propagation
	fmt.Println("Pattern 4: Context Propagation")
	fmt.Println("==============================")

	demonstrateContextPropagation(session)
	fmt.Println()

	// Step 7: Demonstrate graceful shutdown
	fmt.Println("Pattern 5: Graceful Shutdown")
	fmt.Println("============================")

	demonstrateGracefulShutdown(session)
	fmt.Println()

	// Step 8: Demonstrate concurrent operations with shared context
	fmt.Println("Pattern 6: Concurrent Operations")
	fmt.Println("================================")

	demonstrateConcurrentOperations(session)
	fmt.Println()

	// Best practices summary
	fmt.Println("=" + repeatString("=", 59))
	fmt.Println("Context Usage Best Practices")
	fmt.Println("=" + repeatString("=", 59))
	fmt.Println()
	fmt.Println("1. Always pass context as the first parameter")
	fmt.Println("2. Use context.WithTimeout for individual operation limits")
	fmt.Println("3. Use context.WithDeadline for absolute time limits")
	fmt.Println("4. Always call cancel() in defer to release resources")
	fmt.Println("5. Check context.Err() for cancellation cause")
	fmt.Println("6. Propagate context to child goroutines")
	fmt.Println("7. Use signal handling for graceful shutdown")
	fmt.Println("8. Set reasonable timeouts based on operation complexity")

	fmt.Println("\nContext and timeout example completed successfully!")
}

// demonstrateBasicTimeout shows how to set timeouts for requests
func demonstrateBasicTimeout(session *canvus.Session) {
	fmt.Println("Setting a timeout for a single request:")

	// Create a context with a 10-second timeout
	// After 10 seconds, the context will be cancelled automatically
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Always call cancel to release resources

	// Make the request with the timeout context
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Request timed out after 10 seconds")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	fmt.Printf("Retrieved %d canvases within timeout\n", len(canvases))

	// Example: Very short timeout (will likely fail)
	fmt.Println("\nDemonstrating timeout behavior with very short timeout:")

	shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer shortCancel()

	// Small delay to ensure timeout expires
	time.Sleep(2 * time.Millisecond)

	_, err = session.ListCanvases(shortCtx, nil)
	if err != nil {
		if shortCtx.Err() == context.DeadlineExceeded {
			fmt.Println("Request timed out (as expected with 1ms timeout)")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Show the pattern in code
	fmt.Println("\nTimeout pattern:")
	fmt.Println("```go")
	fmt.Println("ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)")
	fmt.Println("defer cancel() // Always call cancel!")
	fmt.Println("")
	fmt.Println("result, err := session.ListCanvases(ctx, nil)")
	fmt.Println("if err != nil {")
	fmt.Println("    if ctx.Err() == context.DeadlineExceeded {")
	fmt.Println("        // Handle timeout specifically")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")
}

// demonstrateManualCancellation shows how to cancel operations manually
func demonstrateManualCancellation(session *canvus.Session) {
	fmt.Println("Creating a cancellable context:")

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine that will cancel after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Cancelling operation from goroutine...")
		cancel()
	}()

	// Try to do a longer operation
	// This demonstrates that cancellation works across goroutines
	fmt.Println("Starting operation that will be cancelled...")

	// Simulate work that checks for cancellation
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Operation cancelled at iteration %d: %v\n", i, ctx.Err())
			return
		default:
			time.Sleep(20 * time.Millisecond)
			fmt.Printf("Working... iteration %d\n", i)
		}
	}

	// Show the pattern in code
	fmt.Println("\nManual cancellation pattern:")
	fmt.Println("```go")
	fmt.Println("ctx, cancel := context.WithCancel(context.Background())")
	fmt.Println("")
	fmt.Println("// In another goroutine or on user action:")
	fmt.Println("go func() {")
	fmt.Println("    <-userCancelChan")
	fmt.Println("    cancel() // Cancel all operations using this context")
	fmt.Println("}()")
	fmt.Println("")
	fmt.Println("// Your operation will receive cancellation")
	fmt.Println("result, err := session.DoOperation(ctx)")
	fmt.Println("```")
}

// demonstrateDeadline shows how to use absolute deadlines
func demonstrateDeadline(session *canvus.Session) {
	fmt.Println("Setting an absolute deadline:")

	// Set a deadline 5 seconds from now
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Check remaining time
	if d, ok := ctx.Deadline(); ok {
		remaining := time.Until(d)
		fmt.Printf("Deadline set for: %s (%.1f seconds remaining)\n",
			d.Format(time.RFC3339), remaining.Seconds())
	}

	// Make a request
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Operation exceeded deadline")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	fmt.Printf("Retrieved %d canvases before deadline\n", len(canvases))

	// Show the pattern in code
	fmt.Println("\nDeadline pattern for batch jobs:")
	fmt.Println("```go")
	fmt.Println("// Must complete by 3 AM")
	fmt.Println("deadline := time.Date(2024, 1, 15, 3, 0, 0, 0, time.Local)")
	fmt.Println("ctx, cancel := context.WithDeadline(context.Background(), deadline)")
	fmt.Println("defer cancel()")
	fmt.Println("")
	fmt.Println("for _, item := range items {")
	fmt.Println("    if err := processItem(ctx, item); err != nil {")
	fmt.Println("        if ctx.Err() == context.DeadlineExceeded {")
	fmt.Println("            log.Printf(\"Batch job hit deadline, %d items remaining\", remaining)")
	fmt.Println("            break")
	fmt.Println("        }")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")
}

// demonstrateContextPropagation shows how to propagate context to child operations
func demonstrateContextPropagation(session *canvus.Session) {
	fmt.Println("Propagating context through operation chain:")

	// Create a parent context with overall timeout
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer parentCancel()

	// processCanvas function receives and uses the parent context
	processCanvas := func(ctx context.Context, canvasID string) error {
		// Create child context with shorter timeout for this specific operation
		childCtx, childCancel := context.WithTimeout(ctx, 5*time.Second)
		defer childCancel()

		// The child context inherits cancellation from parent
		// AND has its own timeout
		_, err := session.GetCanvas(childCtx, canvasID)
		return err
	}

	// Get list of canvases using parent context
	canvases, err := session.ListCanvases(parentCtx, nil)
	if err != nil {
		fmt.Printf("Error listing canvases: %v\n", err)
		return
	}

	fmt.Printf("Processing %d canvases...\n", len(canvases))

	// Process each canvas with propagated context
	for i, canvas := range canvases {
		if i >= 3 { // Limit for demo
			break
		}

		err := processCanvas(parentCtx, canvas.ID)
		if err != nil {
			if parentCtx.Err() != nil {
				fmt.Printf("Parent context cancelled: %v\n", parentCtx.Err())
				break
			}
			fmt.Printf("Error processing canvas %s: %v\n", canvas.ID, err)
			continue
		}
		fmt.Printf("Processed canvas: %s\n", canvas.Name)
	}

	// Show the pattern in code
	fmt.Println("\nContext propagation pattern:")
	fmt.Println("```go")
	fmt.Println("func processAllCanvases(ctx context.Context) error {")
	fmt.Println("    canvases, err := session.ListCanvases(ctx, nil)")
	fmt.Println("    if err != nil {")
	fmt.Println("        return err")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    for _, c := range canvases {")
	fmt.Println("        // Create child context for each canvas")
	fmt.Println("        childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)")
	fmt.Println("        err := processCanvas(childCtx, c.ID)")
	fmt.Println("        cancel() // Always cancel child contexts")
	fmt.Println("        ")
	fmt.Println("        if err != nil {")
	fmt.Println("            // Check if parent was cancelled")
	fmt.Println("            if ctx.Err() != nil {")
	fmt.Println("                return ctx.Err()")
	fmt.Println("            }")
	fmt.Println("        }")
	fmt.Println("    }")
	fmt.Println("    return nil")
	fmt.Println("}")
	fmt.Println("```")
}

// demonstrateGracefulShutdown shows graceful shutdown pattern
func demonstrateGracefulShutdown(session *canvus.Session) {
	fmt.Println("Graceful shutdown with signal handling:")

	// Create a context that will be cancelled on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create a shutdown context with timeout
	// This gives operations time to complete after signal
	shutdownTimeout := 5 * time.Second

	fmt.Println("Application would run until signal received...")
	fmt.Println("(Skipping actual signal wait for demo)")

	// Check context state for demonstration
	select {
	case <-ctx.Done():
		fmt.Printf("Context was cancelled: %v\n", ctx.Err())
	default:
		fmt.Println("Context is still active (no signal received)")
	}

	// Show the pattern in code
	fmt.Println("\nGraceful shutdown pattern:")
	fmt.Println("```go")
	fmt.Println("func main() {")
	fmt.Println("    // Create context that listens for shutdown signals")
	fmt.Println("    ctx, stop := signal.NotifyContext(")
	fmt.Println("        context.Background(),")
	fmt.Println("        syscall.SIGINT, syscall.SIGTERM,")
	fmt.Println("    )")
	fmt.Println("    defer stop()")
	fmt.Println("    ")
	fmt.Println("    // Start your service")
	fmt.Println("    go runService(ctx)")
	fmt.Println("    ")
	fmt.Println("    // Wait for shutdown signal")
	fmt.Println("    <-ctx.Done()")
	fmt.Println("    fmt.Println(\"Shutdown signal received\")")
	fmt.Println("    ")
	fmt.Println("    // Create shutdown context with timeout")
	fmt.Println("    shutdownCtx, cancel := context.WithTimeout(")
	fmt.Println("        context.Background(),")
	fmt.Println("        30*time.Second,")
	fmt.Println("    )")
	fmt.Println("    defer cancel()")
	fmt.Println("    ")
	fmt.Println("    // Graceful shutdown")
	fmt.Println("    if err := shutdown(shutdownCtx); err != nil {")
	fmt.Println("        log.Printf(\"Shutdown error: %v\", err)")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")

	_ = shutdownTimeout // Use variable for documentation purposes
}

// demonstrateConcurrentOperations shows context with concurrent operations
func demonstrateConcurrentOperations(session *canvus.Session) {
	fmt.Println("Running concurrent operations with shared context:")

	// Create a context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get canvases to process
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		fmt.Printf("Error listing canvases: %v\n", err)
		return
	}

	// Limit for demo
	if len(canvases) > 5 {
		canvases = canvases[:5]
	}

	fmt.Printf("Processing %d canvases concurrently...\n", len(canvases))

	// Process canvases concurrently
	var wg sync.WaitGroup
	results := make(chan string, len(canvases))
	errors := make(chan error, len(canvases))

	for _, canvas := range canvases {
		wg.Add(1)
		go func(c canvus.Canvas) {
			defer wg.Done()

			// Check context before starting
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			default:
			}

			// Process the canvas
			_, err := session.GetCanvas(ctx, c.ID)
			if err != nil {
				errors <- err
				return
			}

			results <- c.Name
		}(canvas)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	successCount := 0
	errorCount := 0

	for {
		select {
		case name, ok := <-results:
			if !ok {
				results = nil
			} else {
				successCount++
				fmt.Printf("  Processed: %s\n", name)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else {
				errorCount++
				fmt.Printf("  Error: %v\n", err)
			}
		}

		if results == nil && errors == nil {
			break
		}
	}

	fmt.Printf("Completed: %d success, %d errors\n", successCount, errorCount)

	// Show the pattern in code
	fmt.Println("\nConcurrent operations pattern:")
	fmt.Println("```go")
	fmt.Println("ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)")
	fmt.Println("defer cancel()")
	fmt.Println("")
	fmt.Println("var wg sync.WaitGroup")
	fmt.Println("errChan := make(chan error, len(items))")
	fmt.Println("")
	fmt.Println("for _, item := range items {")
	fmt.Println("    wg.Add(1)")
	fmt.Println("    go func(i Item) {")
	fmt.Println("        defer wg.Done()")
	fmt.Println("        ")
	fmt.Println("        // Each goroutine uses the shared context")
	fmt.Println("        if err := process(ctx, i); err != nil {")
	fmt.Println("            errChan <- err")
	fmt.Println("        }")
	fmt.Println("    }(item)")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("wg.Wait()")
	fmt.Println("close(errChan)")
	fmt.Println("")
	fmt.Println("// If parent context is cancelled, all goroutines will see it")
	fmt.Println("```")
}

// Helper function to repeat a string
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
