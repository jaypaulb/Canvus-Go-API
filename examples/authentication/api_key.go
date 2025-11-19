// Package main demonstrates API key authentication with the Canvus Go SDK.
//
// This example shows:
// - How to configure API key authentication using WithAPIKey
// - How the SDK automatically adds the Private-Token header
// - How to make authenticated requests
// - Proper error handling for authentication failures
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run api_key.go
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

	// Validate required configuration
	if apiURL == "" {
		log.Fatal("CANVUS_API_URL environment variable is required")
	}
	if apiKey == "" {
		log.Fatal("CANVUS_API_KEY environment variable is required")
	}

	fmt.Println("Canvus Go SDK - API Key Authentication Example")
	fmt.Println("===============================================")
	fmt.Printf("Server: %s\n", apiURL)
	fmt.Printf("API Key: %s...%s (masked)\n\n", apiKey[:4], apiKey[len(apiKey)-4:])

	// Step 2: Create session configuration with defaults
	// The default configuration includes:
	// - 3 max retries for transient failures
	// - Exponential backoff between retries
	// - 30 second request timeout
	// - Circuit breaker to prevent overwhelming failed servers
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL

	// Step 3: Create session with API key authentication
	// WithAPIKey() is a functional option that:
	// - Creates a custom HTTP transport
	// - Automatically adds the "Private-Token" header to all requests
	// - Handles TLS configuration (InsecureSkipVerify for self-signed certs)
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Step 4: Create context with timeout
	// Best practice: always use a context with timeout for API calls
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 5: Test the authentication by listing canvases
	fmt.Println("Testing API key authentication...")
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		// Handle different types of errors
		if apiErr, ok := err.(*canvus.APIError); ok {
			switch apiErr.StatusCode {
			case 401:
				log.Fatalf("Authentication failed: Invalid API key\n"+
					"Please verify your CANVUS_API_KEY is correct.\n"+
					"Error: %s", apiErr.Message)
			case 403:
				log.Fatalf("Authorization failed: API key does not have sufficient permissions\n"+
					"Error: %s", apiErr.Message)
			default:
				log.Fatalf("API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			}
		}
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Authentication successful! Found %d canvas(es)\n\n", len(canvases))

	// Step 6: Make additional authenticated requests
	// All subsequent requests will automatically include the API key
	if len(canvases) > 0 {
		// Get details of the first canvas
		canvas := canvases[0]
		fmt.Printf("Fetching canvas details: %s\n", canvas.Name)

		details, err := session.GetCanvas(ctx, canvas.ID)
		if err != nil {
			log.Fatalf("Error getting canvas: %v", err)
		}

		fmt.Printf("  ID: %s\n", details.ID)
		fmt.Printf("  Mode: %s\n", details.Mode)
		fmt.Printf("  State: %s\n", details.State)
		fmt.Println()

		// List widgets on the canvas
		fmt.Printf("Listing widgets on canvas: %s\n", canvas.Name)
		widgets, err := session.ListWidgets(ctx, canvas.ID, nil)
		if err != nil {
			log.Fatalf("Error listing widgets: %v", err)
		}

		fmt.Printf("  Found %d widget(s)\n", len(widgets))
		for i, w := range widgets {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", len(widgets)-5)
				break
			}
			fmt.Printf("  - %s (type: %s)\n", w.ID, w.WidgetType)
		}
	}

	// Step 7: Demonstrate configuring additional session options
	fmt.Println("\nAdvanced: Creating session with custom configuration")

	// Create a session with custom retry and timeout settings
	customCfg := canvus.DefaultSessionConfig()
	customCfg.BaseURL = apiURL
	customCfg.MaxRetries = 5                       // More retries for unreliable networks
	customCfg.RequestTimeout = 60 * time.Second    // Longer timeout for slow operations
	customCfg.RetryWaitMin = 200 * time.Millisecond // Longer wait between retries

	customSession := canvus.NewSession(customCfg, canvus.WithAPIKey(apiKey))

	// Test the custom session
	_, err = customSession.ListCanvases(ctx, nil)
	if err != nil {
		log.Fatalf("Error with custom session: %v", err)
	}
	fmt.Println("Custom session configuration working correctly!")

	fmt.Println("\nAPI key authentication example completed successfully!")
}
