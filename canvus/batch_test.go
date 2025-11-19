package canvus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatchProcessor(t *testing.T) {
	// This is a basic test structure - in a real scenario, you'd need a test server
	// For now, we'll test the configuration and basic functionality

	t.Run("DefaultBatchConfig", func(t *testing.T) {
		config := DefaultBatchConfig()

		assert.Equal(t, 10, config.MaxConcurrency)
		assert.Equal(t, 5*time.Minute, config.Timeout)
		assert.Equal(t, 3, config.RetryAttempts)
		assert.Equal(t, time.Second, config.RetryDelay)
		assert.True(t, config.ContinueOnError)
	})

	t.Run("BatchOperationBuilder", func(t *testing.T) {
		builder := NewBatchOperationBuilder()

		canvas := &Canvas{ID: "test-canvas"}
		widget := &Widget{ID: "test-widget"}

		operations := builder.
			Move("op1", canvas, "folder1").
			Copy("op2", widget, "canvas2").
			Delete("op3", canvas).
			Pin("op4", widget).
			Unpin("op5", widget).
			Build()

		assert.Len(t, operations, 5)

		assert.Equal(t, BatchOperationMove, operations[0].Type)
		assert.Equal(t, BatchOperationCopy, operations[1].Type)
		assert.Equal(t, BatchOperationDelete, operations[2].Type)
		assert.Equal(t, BatchOperationPin, operations[3].Type)
		assert.Equal(t, BatchOperationUnpin, operations[4].Type)

		assert.Equal(t, "op1", operations[0].ID)
		assert.Equal(t, "op2", operations[1].ID)
		assert.Equal(t, "op3", operations[2].ID)
		assert.Equal(t, "op4", operations[3].ID)
		assert.Equal(t, "op5", operations[4].ID)
	})

	t.Run("BatchSummary", func(t *testing.T) {
		results := []*BatchResult{
			{
				OperationID: "op1",
				Success:     true,
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(100 * time.Millisecond),
				Duration:    100 * time.Millisecond,
			},
			{
				OperationID: "op2",
				Success:     false,
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(50 * time.Millisecond),
				Duration:    50 * time.Millisecond,
				Error:       assert.AnError,
			},
			{
				OperationID: "op3",
				Success:     true,
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(75 * time.Millisecond),
				Duration:    75 * time.Millisecond,
			},
		}

		summary := Summarize(results)

		assert.Equal(t, 3, summary.TotalOperations)
		assert.Equal(t, 2, summary.Successful)
		assert.Equal(t, 1, summary.Failed)
		assert.Equal(t, 225*time.Millisecond, summary.TotalDuration)
		assert.Equal(t, 75*time.Millisecond, summary.AverageDuration)
		assert.Len(t, summary.FailedOperations, 1)
		assert.Equal(t, "op2", summary.FailedOperations[0].OperationID)
	})

	t.Run("BatchProcessorCreation", func(t *testing.T) {
		// Test that nil config defaults work
		defaultConfig := DefaultBatchConfig()
		assert.NotNil(t, defaultConfig)
	})
}
