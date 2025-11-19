// Package main demonstrates widget operations with the Canvus Go SDK.
//
// This example shows:
// - Creating widgets (notes) on a canvas
// - Listing widgets with filtering
// - Updating widget properties
// - Cross-canvas widget search
// - Using geometry utilities (WidgetsContainId)
// - Deleting widgets
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run create_and_search.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	fmt.Println("Canvus Go SDK - Widget Operations Example")
	fmt.Println("==========================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 3: Create a canvas for our widget examples
	fmt.Println("Creating a test canvas for widget operations...")

	createCanvasReq := canvus.CreateCanvasRequest{
		Name: fmt.Sprintf("Widget Example Canvas %d", time.Now().Unix()),
	}

	canvas, err := session.CreateCanvas(ctx, createCanvasReq)
	if err != nil {
		log.Fatalf("Error creating canvas: %v", err)
	}

	fmt.Printf("Canvas created: %s (ID: %s)\n\n", canvas.Name, canvas.ID)
	canvasID := canvas.ID

	// Ensure cleanup
	defer func() {
		fmt.Println("\nCleaning up: Deleting test canvas...")
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()

		if err := session.DeleteCanvas(cleanupCtx, canvasID); err != nil {
			log.Printf("Warning: Failed to delete canvas: %v", err)
		} else {
			fmt.Printf("Canvas '%s' deleted successfully\n", canvasID)
		}
	}()

	// Step 4: Create multiple note widgets
	fmt.Println("Creating note widgets...")

	// Create a container note (large area)
	containerNote := map[string]interface{}{
		"widget_type":      "note",
		"text":             "This is a container note that will contain other widgets",
		"title":            "Container Note",
		"background_color": "#FFEB3B", // Yellow
		"location": map[string]interface{}{
			"x": 100.0,
			"y": 100.0,
		},
		"size": map[string]interface{}{
			"width":  800.0,
			"height": 600.0,
		},
	}

	createdContainer, err := session.CreateNote(ctx, canvasID, containerNote)
	if err != nil {
		log.Fatalf("Error creating container note: %v", err)
	}
	fmt.Printf("Created container note: %s\n", createdContainer.ID)

	// Create notes inside the container
	notesData := []map[string]interface{}{
		{
			"widget_type":      "note",
			"text":             "First note inside the container",
			"title":            "Note 1",
			"background_color": "#4CAF50", // Green
			"location": map[string]interface{}{
				"x": 150.0,
				"y": 150.0,
			},
			"size": map[string]interface{}{
				"width":  200.0,
				"height": 150.0,
			},
		},
		{
			"widget_type":      "note",
			"text":             "Second note inside the container",
			"title":            "Note 2",
			"background_color": "#2196F3", // Blue
			"location": map[string]interface{}{
				"x": 400.0,
				"y": 150.0,
			},
			"size": map[string]interface{}{
				"width":  200.0,
				"height": 150.0,
			},
		},
		{
			"widget_type":      "note",
			"text":             "Third note OUTSIDE the container",
			"title":            "Note 3",
			"background_color": "#F44336", // Red
			"location": map[string]interface{}{
				"x": 1000.0, // Outside container bounds
				"y": 100.0,
			},
			"size": map[string]interface{}{
				"width":  200.0,
				"height": 150.0,
			},
		},
	}

	var createdNoteIDs []string
	for _, noteData := range notesData {
		note, err := session.CreateNote(ctx, canvasID, noteData)
		if err != nil {
			log.Fatalf("Error creating note: %v", err)
		}
		createdNoteIDs = append(createdNoteIDs, note.ID)
		fmt.Printf("Created note: %s (Title: %s)\n", note.ID, note.Title)
	}
	fmt.Println()

	// Step 5: List all widgets on the canvas
	fmt.Println("Listing all widgets on canvas...")

	widgets, err := session.ListWidgets(ctx, canvasID, nil)
	if err != nil {
		log.Fatalf("Error listing widgets: %v", err)
	}

	fmt.Printf("Found %d widget(s):\n", len(widgets))
	for _, w := range widgets {
		locStr := "unknown"
		if w.Location != nil {
			locStr = fmt.Sprintf("(%.0f, %.0f)", w.Location.X, w.Location.Y)
		}
		sizeStr := "unknown"
		if w.Size != nil {
			sizeStr = fmt.Sprintf("%.0fx%.0f", w.Size.Width, w.Size.Height)
		}
		fmt.Printf("  - %s [%s] at %s size %s\n", w.ID, w.WidgetType, locStr, sizeStr)
	}
	fmt.Println()

	// Step 6: List widgets with filtering
	fmt.Println("Filtering widgets by type...")

	// Filter for note widgets only
	noteFilter := &canvus.Filter{
		Criteria: map[string]interface{}{
			"widget_type": "note",
		},
	}

	noteWidgets, err := session.ListWidgets(ctx, canvasID, noteFilter)
	if err != nil {
		log.Fatalf("Error filtering widgets: %v", err)
	}

	fmt.Printf("Found %d note widget(s)\n\n", len(noteWidgets))

	// Step 7: Get a specific widget
	fmt.Println("Getting specific widget details...")

	if len(createdNoteIDs) > 0 {
		widgetDetails, err := session.GetWidget(ctx, canvasID, createdNoteIDs[0])
		if err != nil {
			log.Fatalf("Error getting widget: %v", err)
		}

		fmt.Printf("Widget Details:\n")
		fmt.Printf("  ID: %s\n", widgetDetails.ID)
		fmt.Printf("  Type: %s\n", widgetDetails.WidgetType)
		fmt.Printf("  State: %s\n", widgetDetails.State)
		fmt.Printf("  Pinned: %v\n", widgetDetails.Pinned)
		fmt.Printf("  Scale: %.2f\n", widgetDetails.Scale)
		if widgetDetails.Location != nil {
			fmt.Printf("  Location: (%.2f, %.2f)\n", widgetDetails.Location.X, widgetDetails.Location.Y)
		}
		if widgetDetails.Size != nil {
			fmt.Printf("  Size: %.2f x %.2f\n", widgetDetails.Size.Width, widgetDetails.Size.Height)
		}
		fmt.Println()
	}

	// Step 8: Update a widget
	fmt.Println("Updating widget properties...")

	if len(createdNoteIDs) > 0 {
		updateReq := map[string]interface{}{
			"widget_type":      "note",
			"text":             "This note has been updated!",
			"background_color": "#9C27B0", // Purple
		}

		updatedWidget, err := session.UpdateNote(ctx, canvasID, createdNoteIDs[0], updateReq)
		if err != nil {
			log.Fatalf("Error updating widget: %v", err)
		}

		fmt.Printf("Widget updated: %s\n", updatedWidget.ID)
		fmt.Printf("  New Text: %s\n", updatedWidget.Text)
		fmt.Printf("  New Color: %s\n\n", updatedWidget.BackgroundColor)
	}

	// Step 9: Use geometry utilities - WidgetsContainId
	fmt.Println("Using geometry utilities: Finding widgets contained in container...")

	// WidgetsContainId finds all widgets that are geometrically contained within
	// the bounding box of the specified widget
	zone, err := canvus.WidgetsContainId(ctx, session, canvasID, createdContainer.ID, nil, 0)
	if err != nil {
		log.Fatalf("Error finding contained widgets: %v", err)
	}

	fmt.Printf("Container widget: %s\n", zone.Container.ID)
	fmt.Printf("Widgets contained inside: %d\n", len(zone.Contents))
	for _, w := range zone.Contents {
		fmt.Printf("  - %s [%s]\n", w.ID, w.WidgetType)
	}
	fmt.Println()

	// Step 10: Cross-canvas widget search
	fmt.Println("Searching for widgets across all canvases...")

	// Find all note widgets across all canvases
	query := map[string]interface{}{
		"widget_type": "note",
	}

	matches, err := canvus.FindWidgetsAcrossCanvases(ctx, session, query)
	if err != nil {
		log.Fatalf("Error searching widgets: %v", err)
	}

	fmt.Printf("Found %d note widget(s) across all canvases:\n", len(matches))

	// Group by canvas
	canvasCount := make(map[string]int)
	for _, match := range matches {
		canvasCount[match.CanvasID]++
	}

	for cID, count := range canvasCount {
		fmt.Printf("  Canvas %s: %d note(s)\n", cID, count)
	}
	fmt.Println()

	// Step 11: Demonstrate widget geometry functions
	fmt.Println("Demonstrating geometry functions...")

	if len(widgets) >= 2 {
		w1 := widgets[0]
		w2 := widgets[1]

		// Check if widgets touch
		touches := canvus.WidgetsTouch(w1, w2)
		fmt.Printf("Widget %s and %s touch: %v\n", w1.ID, w2.ID, touches)

		// Check if one contains the other
		contains := canvus.WidgetContains(w1, w2)
		fmt.Printf("Widget %s contains %s: %v\n", w1.ID, w2.ID, contains)

		// Get bounding boxes
		bbox1 := canvus.WidgetBoundingBox(w1)
		fmt.Printf("Widget %s bounding box: (%.0f, %.0f) %.0fx%.0f\n",
			w1.ID, bbox1.X, bbox1.Y, bbox1.Width, bbox1.Height)
	}
	fmt.Println()

	// Step 12: Delete a widget
	fmt.Println("Deleting a widget...")

	if len(createdNoteIDs) > 0 {
		lastNoteID := createdNoteIDs[len(createdNoteIDs)-1]
		err = session.DeleteNote(ctx, canvasID, lastNoteID)
		if err != nil {
			log.Fatalf("Error deleting widget: %v", err)
		}
		fmt.Printf("Widget %s deleted successfully\n", lastNoteID)
	}

	// Verify deletion by listing again
	remainingWidgets, err := session.ListWidgets(ctx, canvasID, nil)
	if err != nil {
		log.Fatalf("Error listing remaining widgets: %v", err)
	}
	fmt.Printf("Remaining widgets on canvas: %d\n", len(remainingWidgets))

	fmt.Println("\nWidget operations example completed successfully!")
	// The deferred cleanup will delete the canvas and all remaining widgets
}
