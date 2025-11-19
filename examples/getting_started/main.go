// Package main demonstrates a simple getting started example for the Canvus Go SDK.
//
// This example shows:
// - How to initialize a new session
// - How to authenticate with an API key
// - How to make your first API call (list canvases)
// - How to handle errors properly
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run main.go
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
	// These should be set before running the example
	apiURL := os.Getenv("CANVUS_API_URL")
	apiKey := os.Getenv("CANVUS_API_KEY")

	// Validate required configuration
	if apiURL == "" {
		log.Fatal("CANVUS_API_URL environment variable is required")
	}
	if apiKey == "" {
		log.Fatal("CANVUS_API_KEY environment variable is required")
	}

	fmt.Println("Canvus Go SDK - Getting Started Example")
	fmt.Println("========================================")
	fmt.Printf("Connecting to: %s\n\n", apiURL)

	// Step 2: Create a session configuration
	// DefaultSessionConfig() provides sensible defaults:
	// - 3 max retries
	// - 30 second request timeout
	// - Exponential backoff for retries
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL

	// Step 3: Create a new session with API key authentication
	// WithAPIKey() is a functional option that configures the HTTP client
	// to include the API key in the Private-Token header for all requests
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Step 4: Create a context with timeout
	// Using context allows you to cancel long-running requests
	// and set deadlines for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Always call cancel to release resources

	// Step 5: Make your first API call - List all canvases
	// ListCanvases returns a slice of Canvas objects
	// The second parameter is an optional filter (nil means no filtering)
	fmt.Println("Fetching canvases...")
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		// Check if it's an API error for more details
		if apiErr, ok := err.(*canvus.APIError); ok {
			log.Fatalf("API Error: %s (status: %d, code: %s)",
				apiErr.Message, apiErr.StatusCode, apiErr.Code)
		}
		log.Fatalf("Error listing canvases: %v", err)
	}

	// Step 6: Display the results
	fmt.Printf("Successfully retrieved %d canvas(es)\n\n", len(canvases))

	if len(canvases) == 0 {
		fmt.Println("No canvases found. Create a canvas in Canvus to see it here.")
		return
	}

	// Print details for each canvas
	fmt.Println("Canvas List:")
	fmt.Println("------------")
	for i, canvas := range canvases {
		fmt.Printf("%d. %s\n", i+1, canvas.Name)
		fmt.Printf("   ID: %s\n", canvas.ID)
		fmt.Printf("   Mode: %s\n", canvas.Mode)
		fmt.Printf("   State: %s\n", canvas.State)
		fmt.Printf("   Created: %s\n", canvas.CreatedAt)
		fmt.Println()
	}

	// Step 7: Demonstrate getting a specific canvas
	// Use the first canvas from the list
	firstCanvas := canvases[0]
	fmt.Printf("Fetching details for canvas: %s\n", firstCanvas.Name)

	canvas, err := session.GetCanvas(ctx, firstCanvas.ID)
	if err != nil {
		log.Fatalf("Error getting canvas: %v", err)
	}

	fmt.Printf("Canvas Details:\n")
	fmt.Printf("  Name: %s\n", canvas.Name)
	fmt.Printf("  ID: %s\n", canvas.ID)
	fmt.Printf("  Mode: %s\n", canvas.Mode)
	fmt.Printf("  Access: %s\n", canvas.Access)
	fmt.Printf("  In Trash: %v\n", canvas.InTrash)
	fmt.Printf("  Asset Size: %d bytes\n", canvas.AssetSize)

	fmt.Println("\nGetting started example completed successfully!")
}
