// Package main demonstrates canvas lifecycle management with the Canvus Go SDK.
//
// This example shows:
// - Creating a new canvas
// - Listing canvases with filtering
// - Getting canvas details
// - Updating canvas properties
// - Copying a canvas
// - Deleting a canvas
// - Handling pagination with ListOptions
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run create_and_manage.go
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

	fmt.Println("Canvus Go SDK - Canvas Management Example")
	fmt.Println("==========================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 3: Create a new canvas
	fmt.Println("Creating a new canvas...")

	createReq := canvus.CreateCanvasRequest{
		Name: fmt.Sprintf("SDK Example Canvas %d", time.Now().Unix()),
		// FolderID can be set to place canvas in a specific folder
		// If empty, it will be placed in the default location
	}

	newCanvas, err := session.CreateCanvas(ctx, createReq)
	if err != nil {
		log.Fatalf("Error creating canvas: %v", err)
	}

	fmt.Printf("Canvas created successfully!\n")
	fmt.Printf("  Name: %s\n", newCanvas.Name)
	fmt.Printf("  ID: %s\n", newCanvas.ID)
	fmt.Printf("  Mode: %s\n", newCanvas.Mode)
	fmt.Printf("  State: %s\n", newCanvas.State)
	fmt.Printf("  Created At: %s\n\n", newCanvas.CreatedAt)

	// Store the canvas ID for cleanup
	canvasID := newCanvas.ID

	// Ensure cleanup happens even if there's an error
	defer func() {
		fmt.Println("\nCleaning up: Deleting created canvas...")
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()

		if err := session.DeleteCanvas(cleanupCtx, canvasID); err != nil {
			log.Printf("Warning: Failed to delete canvas: %v", err)
		} else {
			fmt.Printf("Canvas '%s' deleted successfully\n", canvasID)
		}
	}()

	// Step 4: List all canvases
	fmt.Println("Listing all canvases...")

	allCanvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		log.Fatalf("Error listing canvases: %v", err)
	}

	fmt.Printf("Total canvases: %d\n\n", len(allCanvases))

	// Step 5: List canvases with filtering
	// Use the Filter type for client-side filtering
	fmt.Println("Filtering canvases by name pattern...")

	// Filter for canvases with "SDK Example" in the name
	filter := &canvus.Filter{
		Criteria: map[string]interface{}{
			"name": "*SDK Example*", // Wildcard pattern matching
		},
	}

	filteredCanvases, err := session.ListCanvases(ctx, filter)
	if err != nil {
		log.Fatalf("Error filtering canvases: %v", err)
	}

	fmt.Printf("Found %d canvas(es) matching '*SDK Example*':\n", len(filteredCanvases))
	for _, c := range filteredCanvases {
		fmt.Printf("  - %s (ID: %s)\n", c.Name, c.ID)
	}
	fmt.Println()

	// Step 6: Get canvas details
	fmt.Println("Getting canvas details...")

	canvas, err := session.GetCanvas(ctx, canvasID)
	if err != nil {
		log.Fatalf("Error getting canvas: %v", err)
	}

	fmt.Printf("Canvas Details:\n")
	fmt.Printf("  ID: %s\n", canvas.ID)
	fmt.Printf("  Name: %s\n", canvas.Name)
	fmt.Printf("  Mode: %s\n", canvas.Mode)
	fmt.Printf("  State: %s\n", canvas.State)
	fmt.Printf("  Access: %s\n", canvas.Access)
	fmt.Printf("  Folder ID: %s\n", canvas.FolderID)
	fmt.Printf("  In Trash: %v\n", canvas.InTrash)
	fmt.Printf("  Asset Size: %d bytes\n", canvas.AssetSize)
	fmt.Printf("  Created At: %s\n", canvas.CreatedAt)
	fmt.Printf("  Modified At: %s\n\n", canvas.ModifiedAt)

	// Step 7: Update canvas properties
	fmt.Println("Updating canvas name...")

	updateReq := canvus.UpdateCanvasRequest{
		Name: fmt.Sprintf("Updated SDK Example %d", time.Now().Unix()),
		// Mode can also be updated here
	}

	updatedCanvas, err := session.UpdateCanvas(ctx, canvasID, updateReq)
	if err != nil {
		log.Fatalf("Error updating canvas: %v", err)
	}

	fmt.Printf("Canvas updated successfully!\n")
	fmt.Printf("  New Name: %s\n", updatedCanvas.Name)
	fmt.Printf("  Modified At: %s\n\n", updatedCanvas.ModifiedAt)

	// Step 8: Copy the canvas
	fmt.Println("Copying canvas...")

	copyReq := canvus.MoveOrCopyCanvasRequest{
		FolderID: canvas.FolderID, // Copy to the same folder
		// Conflicts: "rename" // Handle naming conflicts
	}

	copiedCanvas, err := session.CopyCanvas(ctx, canvasID, copyReq)
	if err != nil {
		log.Fatalf("Error copying canvas: %v", err)
	}

	fmt.Printf("Canvas copied successfully!\n")
	fmt.Printf("  Original ID: %s\n", canvasID)
	fmt.Printf("  Copy ID: %s\n", copiedCanvas.ID)
	fmt.Printf("  Copy Name: %s\n\n", copiedCanvas.Name)

	// Clean up the copied canvas as well
	defer func() {
		fmt.Println("Cleaning up: Deleting copied canvas...")
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()

		if err := session.DeleteCanvas(cleanupCtx, copiedCanvas.ID); err != nil {
			log.Printf("Warning: Failed to delete copied canvas: %v", err)
		} else {
			fmt.Printf("Copied canvas '%s' deleted successfully\n", copiedCanvas.ID)
		}
	}()

	// Step 9: Get canvas permissions
	fmt.Println("Getting canvas permissions...")

	permissions, err := session.GetCanvasPermissions(ctx, canvasID)
	if err != nil {
		log.Printf("Note: Could not get permissions (may require admin access): %v\n", err)
	} else {
		fmt.Printf("Canvas Permissions:\n")
		fmt.Printf("  Editors Can Share: %v\n", permissions.EditorsCanShare)
		fmt.Printf("  Link Permission: %s\n", permissions.LinkPermission)
		fmt.Printf("  Users with access: %d\n", len(permissions.Users))
		fmt.Printf("  Groups with access: %d\n\n", len(permissions.Groups))
	}

	// Step 10: Demonstrate filtering by different criteria
	fmt.Println("Additional filtering examples...")

	// Filter by state
	stateFilter := &canvus.Filter{
		Criteria: map[string]interface{}{
			"state": "active",
		},
	}
	activeCanvases, err := session.ListCanvases(ctx, stateFilter)
	if err != nil {
		log.Printf("Error filtering by state: %v", err)
	} else {
		fmt.Printf("Active canvases: %d\n", len(activeCanvases))
	}

	// Filter canvases not in trash
	notInTrashFilter := &canvus.Filter{
		Criteria: map[string]interface{}{
			"in_trash": false,
		},
	}
	notInTrash, err := session.ListCanvases(ctx, notInTrashFilter)
	if err != nil {
		log.Printf("Error filtering by trash status: %v", err)
	} else {
		fmt.Printf("Canvases not in trash: %d\n", len(notInTrash))
	}

	fmt.Println("\nCanvas management example completed successfully!")
	// The deferred cleanup functions will be called here
}
