// Package main demonstrates import/export operations with the Canvus Go SDK.
//
// This example shows:
// - Exporting canvas widgets and assets to a folder
// - Importing widgets from an export folder
// - Round-trip fidelity (widgets maintain properties through export/import)
// - Asset file handling (images, PDFs)
// - Error recovery patterns for import/export operations
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run round_trip.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	fmt.Println("Canvus Go SDK - Import/Export Example")
	fmt.Println("=====================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Step 3: Create source and target canvases
	fmt.Println("Setting up test environment...")

	sourceCanvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
		Name: fmt.Sprintf("Export Source Canvas %d", time.Now().Unix()),
	})
	if err != nil {
		log.Fatalf("Error creating source canvas: %v", err)
	}

	targetCanvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
		Name: fmt.Sprintf("Import Target Canvas %d", time.Now().Unix()),
	})
	if err != nil {
		log.Fatalf("Error creating target canvas: %v", err)
	}

	fmt.Printf("Source canvas: %s (ID: %s)\n", sourceCanvas.Name, sourceCanvas.ID)
	fmt.Printf("Target canvas: %s (ID: %s)\n\n", targetCanvas.Name, targetCanvas.ID)

	// Create temporary export directory
	exportDir := filepath.Join(os.TempDir(), fmt.Sprintf("canvus-export-%d", time.Now().Unix()))
	fmt.Printf("Export directory: %s\n\n", exportDir)

	// Ensure cleanup
	defer func() {
		fmt.Println("\nCleaning up...")
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()

		// Delete canvases
		if err := session.DeleteCanvas(cleanupCtx, sourceCanvas.ID); err != nil {
			log.Printf("Warning: Failed to delete source canvas: %v", err)
		}
		if err := session.DeleteCanvas(cleanupCtx, targetCanvas.ID); err != nil {
			log.Printf("Warning: Failed to delete target canvas: %v", err)
		}

		// Remove export directory
		if err := os.RemoveAll(exportDir); err != nil {
			log.Printf("Warning: Failed to remove export directory: %v", err)
		}

		fmt.Println("Cleanup completed")
	}()

	// Step 4: Create widgets with different types for export
	fmt.Println("Creating widgets in source canvas...")

	// Create several note widgets with different properties
	widgets := []map[string]interface{}{
		{
			"widget_type":      "note",
			"text":             "This is the first note with important information.",
			"title":            "Note 1 - Introduction",
			"background_color": "#FFEB3B", // Yellow
			"location": map[string]interface{}{
				"x": 100.0,
				"y": 100.0,
			},
			"size": map[string]interface{}{
				"width":  300.0,
				"height": 200.0,
			},
		},
		{
			"widget_type":      "note",
			"text":             "Second note with different content and color.",
			"title":            "Note 2 - Details",
			"background_color": "#4CAF50", // Green
			"location": map[string]interface{}{
				"x": 450.0,
				"y": 100.0,
			},
			"size": map[string]interface{}{
				"width":  300.0,
				"height": 200.0,
			},
		},
		{
			"widget_type":      "note",
			"text":             "Third note positioned below the first two.",
			"title":            "Note 3 - Summary",
			"background_color": "#2196F3", // Blue
			"location": map[string]interface{}{
				"x": 100.0,
				"y": 350.0,
			},
			"size": map[string]interface{}{
				"width":  650.0,
				"height": 150.0,
			},
		},
	}

	var createdWidgetIDs []string
	for _, widgetData := range widgets {
		widget, err := session.CreateNote(ctx, sourceCanvas.ID, widgetData)
		if err != nil {
			log.Fatalf("Error creating widget: %v", err)
		}
		createdWidgetIDs = append(createdWidgetIDs, widget.ID)
		fmt.Printf("Created widget: %s (Title: %s)\n", widget.ID, widget.Title)
	}
	fmt.Println()

	// Step 5: Define the export region
	// The region defines which area of the canvas to export
	// Widgets that fall within this region will be included
	fmt.Println("Defining export region...")

	exportRegion := canvus.Rectangle{
		X:      0,      // Left edge
		Y:      0,      // Top edge
		Width:  1000.0, // Region width
		Height: 600.0,  // Region height
	}

	fmt.Printf("Export region: (%.0f, %.0f) - %.0f x %.0f\n\n",
		exportRegion.X, exportRegion.Y, exportRegion.Width, exportRegion.Height)

	// Step 6: Export widgets to folder
	fmt.Println("Exporting widgets to folder...")

	// ExportWidgetsToFolder exports specified widgets along with their assets
	// The sharedCanvasID parameter is used to blank parent_id for widgets
	// whose parent is the shared canvas (pass empty string if not applicable)
	exportPath, err := session.ExportWidgetsToFolder(
		ctx,
		sourceCanvas.ID,  // Canvas containing the widgets
		createdWidgetIDs, // Widget IDs to export
		exportRegion,     // Region for coordinate reference
		"",               // Shared canvas ID (empty = no parent blanking)
		exportDir,        // Base folder for export
	)
	if err != nil {
		log.Fatalf("Error exporting widgets: %v", err)
	}

	fmt.Printf("Widgets exported to: %s\n\n", exportPath)

	// Step 7: Examine the export structure
	fmt.Println("Examining export structure...")

	// List files in export directory
	files, err := os.ReadDir(exportPath)
	if err != nil {
		log.Fatalf("Error reading export directory: %v", err)
	}

	fmt.Printf("Export directory contents:\n")
	for _, file := range files {
		info, _ := file.Info()
		if info != nil {
			fmt.Printf("  - %s (%d bytes)\n", file.Name(), info.Size())
		} else {
			fmt.Printf("  - %s\n", file.Name())
		}
	}
	fmt.Println()

	// Step 8: Read and display export.json
	fmt.Println("Reading export manifest...")

	exportJSONPath := filepath.Join(exportPath, "export.json")
	exportData, err := os.ReadFile(exportJSONPath)
	if err != nil {
		log.Fatalf("Error reading export.json: %v", err)
	}

	var exportManifest struct {
		Widgets []canvus.Widget   `json:"widgets"`
		Assets  map[string]string `json:"assets"`
		Region  *canvus.Rectangle `json:"region"`
	}

	if err := json.Unmarshal(exportData, &exportManifest); err != nil {
		log.Fatalf("Error parsing export.json: %v", err)
	}

	fmt.Printf("Export manifest:\n")
	fmt.Printf("  Widgets: %d\n", len(exportManifest.Widgets))
	fmt.Printf("  Assets: %d\n", len(exportManifest.Assets))
	if exportManifest.Region != nil {
		fmt.Printf("  Region: (%.0f, %.0f) %.0f x %.0f\n",
			exportManifest.Region.X, exportManifest.Region.Y,
			exportManifest.Region.Width, exportManifest.Region.Height)
	}
	fmt.Println()

	// Display widget details from export
	fmt.Println("Exported widget details:")
	for _, w := range exportManifest.Widgets {
		fmt.Printf("  - %s [%s]\n", w.ID, w.WidgetType)
		if w.Location != nil {
			fmt.Printf("    Location: (%.0f, %.0f)\n", w.Location.X, w.Location.Y)
		}
		if w.Size != nil {
			fmt.Printf("    Size: %.0f x %.0f\n", w.Size.Width, w.Size.Height)
		}
	}
	fmt.Println()

	// Step 9: Prepare for import
	fmt.Println("Preparing to import widgets to target canvas...")

	// Create an ExportedWidgetSet from the export data
	exportedSet := &canvus.ExportedWidgetSet{
		Widgets: exportManifest.Widgets,
		Assets:  exportManifest.Assets,
		Region:  exportManifest.Region,
	}

	// Define the target region for import
	// Widgets will be scaled and translated to fit this region
	targetRegion := canvus.Rectangle{
		X:      200.0, // Offset from left
		Y:      200.0, // Offset from top
		Width:  800.0, // Scaled width
		Height: 480.0, // Scaled height
	}

	fmt.Printf("Target region: (%.0f, %.0f) - %.0f x %.0f\n\n",
		targetRegion.X, targetRegion.Y, targetRegion.Width, targetRegion.Height)

	// Step 10: Import widgets to target canvas
	fmt.Println("Importing widgets to target canvas...")

	// ImportWidgetsToRegion imports widgets and scales them to fit the target region
	// It handles both regular widgets and asset widgets (images, PDFs, videos)
	newWidgetIDs, err := session.ImportWidgetsToRegion(
		ctx,
		targetCanvas.ID, // Target canvas
		exportedSet,     // Exported widget set
		targetRegion,    // Target region for scaling
	)
	if err != nil {
		log.Fatalf("Error importing widgets: %v", err)
	}

	fmt.Printf("Imported %d widget(s)\n", len(newWidgetIDs))
	for i, id := range newWidgetIDs {
		fmt.Printf("  %d. New widget ID: %s\n", i+1, id)
	}
	fmt.Println()

	// Step 11: Verify import
	fmt.Println("Verifying import...")

	targetWidgets, err := session.ListWidgets(ctx, targetCanvas.ID, nil)
	if err != nil {
		log.Fatalf("Error listing target canvas widgets: %v", err)
	}

	fmt.Printf("Target canvas now has %d widget(s)\n\n", len(targetWidgets))

	// Compare original and imported widgets
	fmt.Println("Comparing exported and imported widgets:")
	fmt.Println("-----------------------------------------")

	for i, id := range newWidgetIDs {
		importedWidget, err := session.GetWidget(ctx, targetCanvas.ID, id)
		if err != nil {
			log.Printf("Warning: Could not get imported widget %s: %v", id, err)
			continue
		}

		fmt.Printf("\nWidget %d:\n", i+1)
		fmt.Printf("  Original ID: %s\n", exportManifest.Widgets[i].ID)
		fmt.Printf("  Imported ID: %s\n", importedWidget.ID)
		fmt.Printf("  Type: %s\n", importedWidget.WidgetType)

		// Show location transformation
		origLoc := exportManifest.Widgets[i].Location
		if origLoc != nil && importedWidget.Location != nil {
			fmt.Printf("  Location: (%.0f, %.0f) -> (%.0f, %.0f)\n",
				origLoc.X, origLoc.Y,
				importedWidget.Location.X, importedWidget.Location.Y)
		}

		// Show size transformation
		origSize := exportManifest.Widgets[i].Size
		if origSize != nil && importedWidget.Size != nil {
			fmt.Printf("  Size: %.0f x %.0f -> %.0f x %.0f\n",
				origSize.Width, origSize.Height,
				importedWidget.Size.Width, importedWidget.Size.Height)
		}
	}

	// Step 12: Error handling patterns for import/export
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Import/Export Error Handling Patterns:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Pattern 1: Check for missing assets before import
	fmt.Println("1. Validate export before import:")
	fmt.Println("   ```go")
	fmt.Println("   for widgetID, assetFile := range exported.Assets {")
	fmt.Println("       assetPath := filepath.Join(exportDir, assetFile)")
	fmt.Println("       if _, err := os.Stat(assetPath); os.IsNotExist(err) {")
	fmt.Println("           log.Printf(\"Missing asset for widget %s: %s\", widgetID, assetFile)")
	fmt.Println("       }")
	fmt.Println("   }")
	fmt.Println("   ```")
	fmt.Println()

	// Pattern 2: Handle partial import failures
	fmt.Println("2. Handle partial failures gracefully:")
	fmt.Println("   ```go")
	fmt.Println("   newIDs, err := session.ImportWidgetsToRegion(ctx, canvasID, exported, region)")
	fmt.Println("   if err != nil {")
	fmt.Println("       // Some widgets may have been imported before the error")
	fmt.Println("       log.Printf(\"Import error after %d widgets: %v\", len(newIDs), err)")
	fmt.Println("       // Consider cleanup or rollback")
	fmt.Println("   }")
	fmt.Println("   ```")
	fmt.Println()

	// Pattern 3: Verify widget fidelity after import
	fmt.Println("3. Verify widget properties after import:")
	fmt.Println("   ```go")
	fmt.Println("   for i, newID := range newIDs {")
	fmt.Println("       imported, err := session.GetWidget(ctx, targetID, newID)")
	fmt.Println("       if err != nil {")
	fmt.Println("           log.Printf(\"Failed to verify widget %s\", newID)")
	fmt.Println("           continue")
	fmt.Println("       }")
	fmt.Println("       // Compare with original")
	fmt.Println("       if imported.WidgetType != original[i].WidgetType {")
	fmt.Println("           log.Printf(\"Type mismatch for widget %s\", newID)")
	fmt.Println("       }")
	fmt.Println("   }")
	fmt.Println("   ```")

	// Best practices summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Import/Export Best Practices:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("1. Always verify export directory is writable before exporting")
	fmt.Println("2. Check disk space for large exports with many assets")
	fmt.Println("3. Use appropriate target regions to maintain aspect ratios")
	fmt.Println("4. Handle asset type-specific errors (image encoding, PDF parsing)")
	fmt.Println("5. Implement idempotent imports for recovery scenarios")
	fmt.Println("6. Consider batching large imports for better error recovery")
	fmt.Println("7. Preserve export.json for audit and debugging purposes")

	fmt.Println("\nImport/Export example completed successfully!")
}
