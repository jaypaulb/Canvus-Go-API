// Package canvus provides batch operations for efficient bulk operations on Canvus resources.
package canvus

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BatchOperationType represents the type of batch operation
type BatchOperationType string

const (
	BatchOperationMove   BatchOperationType = "move"
	BatchOperationCopy   BatchOperationType = "copy"
	BatchOperationDelete BatchOperationType = "delete"
	BatchOperationPin    BatchOperationType = "pin"
	BatchOperationUnpin  BatchOperationType = "unpin"
)

// BatchOperation represents a single operation in a batch
type BatchOperation struct {
	ID       string      // Unique ID for this operation
	Type     BatchOperationType
	Resource interface{} // The resource being operated on (Canvas, Widget, etc.)
	Target   interface{} // Target for move/copy operations (folder ID, canvas ID, etc.)
	Metadata map[string]interface{} // Additional operation-specific data
}

// BatchResult represents the result of a single batch operation
type BatchResult struct {
	OperationID string
	Success     bool
	Error       error
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Retries     int
}

// BatchConfig holds configuration for batch operations
type BatchConfig struct {
	MaxConcurrency    int           // Maximum number of concurrent operations
	Timeout           time.Duration // Overall timeout for the batch
	RetryAttempts     int           // Number of retry attempts for failed operations
	RetryDelay        time.Duration // Delay between retry attempts
	ContinueOnError   bool          // Continue processing if individual operations fail
	ProgressCallback  func(completed, total int, results []*BatchResult) // Optional progress callback
}

// DefaultBatchConfig returns sensible defaults for batch operations
func DefaultBatchConfig() *BatchConfig {
	return &BatchConfig{
		MaxConcurrency:  10,
		Timeout:         5 * time.Minute,
		RetryAttempts:   3,
		RetryDelay:      time.Second,
		ContinueOnError: true,
	}
}

// BatchProcessor handles batch operations with concurrency control and error handling
type BatchProcessor struct {
	session *Session
	config  *BatchConfig
	sem     chan struct{} // Semaphore for concurrency control
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(session *Session, config *BatchConfig) *BatchProcessor {
	if config == nil {
		config = DefaultBatchConfig()
	}

	// Ensure MaxConcurrency is reasonable
	if config.MaxConcurrency <= 0 {
		config.MaxConcurrency = 1
	}
	if config.MaxConcurrency > 100 {
		config.MaxConcurrency = 100
	}

	return &BatchProcessor{
		session: session,
		config:  config,
		sem:     make(chan struct{}, config.MaxConcurrency),
	}
}

// ExecuteBatch executes a batch of operations concurrently
func (bp *BatchProcessor) ExecuteBatch(ctx context.Context, operations []*BatchOperation) ([]*BatchResult, error) {
	if len(operations) == 0 {
		return []*BatchResult{}, nil
	}

	// Create context with timeout
	if bp.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, bp.config.Timeout)
		defer cancel()
	}

	// Channel to collect results
	results := make([]*BatchResult, len(operations))
	resultsChan := make(chan *BatchResult, len(operations))

	// WaitGroup to wait for all operations to complete
	var wg sync.WaitGroup

	// Execute operations concurrently
	for i, op := range operations {
		wg.Add(1)
		go func(idx int, operation *BatchOperation) {
			defer wg.Done()

			// Acquire semaphore
			bp.sem <- struct{}{}
			defer func() { <-bp.sem }()

			result := bp.executeOperation(ctx, operation)
			results[idx] = result
			resultsChan <- result
		}(i, op)
	}

	// Close results channel when all operations complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results and call progress callback if provided
	var completedResults []*BatchResult
	for result := range resultsChan {
		completedResults = append(completedResults, result)

		if bp.config.ProgressCallback != nil {
			bp.config.ProgressCallback(len(completedResults), len(operations), completedResults)
		}
	}

	// Check for overall timeout or cancellation
	if ctx.Err() != nil {
		return results, fmt.Errorf("batch operation cancelled or timed out: %w", ctx.Err())
	}

	return results, nil
}

// executeOperation executes a single operation with retry logic
func (bp *BatchProcessor) executeOperation(ctx context.Context, op *BatchOperation) *BatchResult {
	result := &BatchResult{
		OperationID: op.ID,
		StartTime:   time.Now(),
	}

	for attempt := 0; attempt <= bp.config.RetryAttempts; attempt++ {
		result.Retries = attempt

		var err error
		switch op.Type {
		case BatchOperationMove:
			err = bp.executeMove(ctx, op)
		case BatchOperationCopy:
			err = bp.executeCopy(ctx, op)
		case BatchOperationDelete:
			err = bp.executeDelete(ctx, op)
		case BatchOperationPin:
			err = bp.executePin(ctx, op)
		case BatchOperationUnpin:
			err = bp.executeUnpin(ctx, op)
		default:
			err = fmt.Errorf("unsupported operation type: %s", op.Type)
		}

		if err == nil {
			result.Success = true
			break
		}

		result.Error = err

		// Don't retry on the last attempt or if context is cancelled
		if attempt == bp.config.RetryAttempts || ctx.Err() != nil {
			break
		}

		// Wait before retry (with jitter)
		select {
		case <-ctx.Done():
			return result
		case <-time.After(bp.config.RetryDelay):
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	return result
}

// executeMove executes a move operation
func (bp *BatchProcessor) executeMove(ctx context.Context, op *BatchOperation) error {
	switch resource := op.Resource.(type) {
	case *Canvas:
		targetFolderID, ok := op.Target.(string)
		if !ok {
			return fmt.Errorf("move target must be a folder ID string")
		}
		req := MoveOrCopyCanvasRequest{FolderID: targetFolderID}
		_, err := bp.session.MoveCanvas(ctx, resource.ID, req)
		return err
	case *Widget:
		targetCanvasID, ok := op.Target.(string)
		if !ok {
			return fmt.Errorf("move target must be a canvas ID string")
		}
		return bp.session.MoveWidget(ctx, resource.ID, targetCanvasID)
	default:
		return fmt.Errorf("unsupported resource type for move operation")
	}
}

// executeCopy executes a copy operation
func (bp *BatchProcessor) executeCopy(ctx context.Context, op *BatchOperation) error {
	switch resource := op.Resource.(type) {
	case *Canvas:
		targetFolderID, ok := op.Target.(string)
		if !ok {
			return fmt.Errorf("copy target must be a folder ID string")
		}
		req := MoveOrCopyCanvasRequest{FolderID: targetFolderID}
		_, err := bp.session.CopyCanvas(ctx, resource.ID, req)
		return err
	case *Widget:
		targetCanvasID, ok := op.Target.(string)
		if !ok {
			return fmt.Errorf("copy target must be a canvas ID string")
		}
		return bp.session.CopyWidget(ctx, resource.ID, targetCanvasID)
	default:
		return fmt.Errorf("unsupported resource type for copy operation")
	}
}

// executeDelete executes a delete operation
func (bp *BatchProcessor) executeDelete(ctx context.Context, op *BatchOperation) error {
	switch resource := op.Resource.(type) {
	case *Canvas:
		return bp.session.DeleteCanvas(ctx, resource.ID)
	case *Widget:
		// For delete operations, we need the canvas ID and widget type
		// This information should be stored in the operation metadata
		canvasID, hasCanvas := op.Metadata["canvas_id"].(string)
		widgetType, hasType := op.Metadata["widget_type"].(string)
		if !hasCanvas || !hasType {
			return fmt.Errorf("delete operation requires canvas_id and widget_type in metadata")
		}
		return bp.session.DeleteWidget(ctx, canvasID, resource.ID, widgetType)
	case *User:
		return bp.session.DeleteUser(ctx, resource.ID)
	default:
		return fmt.Errorf("unsupported resource type for delete operation")
	}
}

// executePin executes a pin operation
func (bp *BatchProcessor) executePin(ctx context.Context, op *BatchOperation) error {
	widget, ok := op.Resource.(*Widget)
	if !ok {
		return fmt.Errorf("pin operation only supports widgets")
	}
	return bp.session.PinWidget(ctx, widget.ID)
}

// executeUnpin executes an unpin operation
func (bp *BatchProcessor) executeUnpin(ctx context.Context, op *BatchOperation) error {
	widget, ok := op.Resource.(*Widget)
	if !ok {
		return fmt.Errorf("unpin operation only supports widgets")
	}
	return bp.session.UnpinWidget(ctx, widget.ID)
}

// BatchOperationBuilder helps build batch operations fluently
type BatchOperationBuilder struct {
	operations []*BatchOperation
}

// NewBatchOperationBuilder creates a new batch operation builder
func NewBatchOperationBuilder() *BatchOperationBuilder {
	return &BatchOperationBuilder{
		operations: make([]*BatchOperation, 0),
	}
}

// Move adds a move operation to the batch
func (bob *BatchOperationBuilder) Move(id string, resource interface{}, targetFolderID string) *BatchOperationBuilder {
	bob.operations = append(bob.operations, &BatchOperation{
		ID:       id,
		Type:     BatchOperationMove,
		Resource: resource,
		Target:   targetFolderID,
	})
	return bob
}

// Copy adds a copy operation to the batch
func (bob *BatchOperationBuilder) Copy(id string, resource interface{}, targetCanvasID string) *BatchOperationBuilder {
	bob.operations = append(bob.operations, &BatchOperation{
		ID:       id,
		Type:     BatchOperationCopy,
		Resource: resource,
		Target:   targetCanvasID,
	})
	return bob
}

// Delete adds a delete operation to the batch
func (bob *BatchOperationBuilder) Delete(id string, resource interface{}) *BatchOperationBuilder {
	bob.operations = append(bob.operations, &BatchOperation{
		ID:       id,
		Type:     BatchOperationDelete,
		Resource: resource,
	})
	return bob
}

// Pin adds a pin operation to the batch
func (bob *BatchOperationBuilder) Pin(id string, widget *Widget) *BatchOperationBuilder {
	bob.operations = append(bob.operations, &BatchOperation{
		ID:       id,
		Type:     BatchOperationPin,
		Resource: widget,
	})
	return bob
}

// Unpin adds an unpin operation to the batch
func (bob *BatchOperationBuilder) Unpin(id string, widget *Widget) *BatchOperationBuilder {
	bob.operations = append(bob.operations, &BatchOperation{
		ID:       id,
		Type:     BatchOperationUnpin,
		Resource: widget,
	})
	return bob
}

// Build returns the built batch operations
func (bob *BatchOperationBuilder) Build() []*BatchOperation {
	return bob.operations
}

// BatchSummary provides a summary of batch operation results
type BatchSummary struct {
	TotalOperations   int
	Successful        int
	Failed            int
	TotalDuration     time.Duration
	AverageDuration   time.Duration
	FailedOperations  []*BatchResult
}

// Summarize creates a summary of batch operation results
func Summarize(results []*BatchResult) *BatchSummary {
	summary := &BatchSummary{
		TotalOperations:  len(results),
		FailedOperations: make([]*BatchResult, 0),
	}

	var totalDuration time.Duration
	for _, result := range results {
		if result.Success {
			summary.Successful++
		} else {
			summary.Failed++
			summary.FailedOperations = append(summary.FailedOperations, result)
		}
		totalDuration += result.Duration
	}

	summary.TotalDuration = totalDuration
	if summary.TotalOperations > 0 {
		summary.AverageDuration = totalDuration / time.Duration(summary.TotalOperations)
	}

	return summary
}
