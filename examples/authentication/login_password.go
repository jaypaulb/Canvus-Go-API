// Package main demonstrates login/password authentication with the Canvus Go SDK.
//
// This example shows:
// - How to authenticate using username and password
// - How the Login() method obtains and stores a token
// - How to make authenticated requests after login
// - Proper session cleanup with Logout()
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_USERNAME="your-email@example.com"
//   export CANVUS_PASSWORD="your-password"
//   go run login_password.go
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
	username := os.Getenv("CANVUS_USERNAME")
	password := os.Getenv("CANVUS_PASSWORD")

	// Validate required configuration
	if apiURL == "" {
		log.Fatal("CANVUS_API_URL environment variable is required")
	}
	if username == "" {
		log.Fatal("CANVUS_USERNAME environment variable is required")
	}
	if password == "" {
		log.Fatal("CANVUS_PASSWORD environment variable is required")
	}

	fmt.Println("Canvus Go SDK - Login/Password Authentication Example")
	fmt.Println("=====================================================")
	fmt.Printf("Server: %s\n", apiURL)
	fmt.Printf("Username: %s\n\n", username)

	// Step 2: Create session configuration
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL

	// Step 3: Create a new session WITHOUT authentication
	// Unlike WithAPIKey(), we create a bare session first
	// and then call Login() to authenticate
	session := canvus.NewSession(cfg)

	// Step 4: Create context for the login operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 5: Authenticate using Login()
	// Login() sends credentials to POST /users/login and:
	// - Receives a token from the server
	// - Stores the token for future requests
	// - Stores the user ID for reference
	fmt.Println("Logging in...")
	err := session.Login(ctx, username, password)
	if err != nil {
		// Handle different types of login errors
		if apiErr, ok := err.(*canvus.APIError); ok {
			switch apiErr.StatusCode {
			case 401:
				log.Fatalf("Login failed: Invalid username or password\n"+
					"Please verify your credentials.\n"+
					"Error: %s", apiErr.Message)
			case 403:
				log.Fatalf("Login failed: Account may be locked or disabled\n"+
					"Error: %s", apiErr.Message)
			default:
				log.Fatalf("Login API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			}
		}
		log.Fatalf("Login error: %v", err)
	}

	// Get the user ID from the session
	userID := session.UserID()
	fmt.Printf("Login successful! User ID: %d\n\n", userID)

	// Step 6: Ensure we logout when done
	// This is important for security - it invalidates the token
	// Use defer to ensure cleanup happens even if there's an error
	defer func() {
		fmt.Println("\nLogging out...")
		logoutCtx, logoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer logoutCancel()

		if err := session.Logout(logoutCtx); err != nil {
			log.Printf("Warning: Logout failed: %v", err)
		} else {
			fmt.Println("Logged out successfully!")
		}
	}()

	// Step 7: Make authenticated requests
	// The session now uses the token for all requests automatically
	fmt.Println("Making authenticated requests...")

	// List canvases
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		log.Fatalf("Error listing canvases: %v", err)
	}
	fmt.Printf("Found %d canvas(es)\n\n", len(canvases))

	// Display first few canvases
	maxDisplay := 5
	if len(canvases) < maxDisplay {
		maxDisplay = len(canvases)
	}

	for i := 0; i < maxDisplay; i++ {
		canvas := canvases[i]
		fmt.Printf("%d. %s\n", i+1, canvas.Name)
		fmt.Printf("   ID: %s\n", canvas.ID)
		fmt.Printf("   Mode: %s, State: %s\n", canvas.Mode, canvas.State)
	}

	if len(canvases) > 5 {
		fmt.Printf("... and %d more canvas(es)\n", len(canvases)-5)
	}

	// Step 8: Demonstrate additional operations
	if len(canvases) > 0 {
		fmt.Println("\nFetching widgets from first canvas...")
		widgets, err := session.ListWidgets(ctx, canvases[0].ID, nil)
		if err != nil {
			log.Fatalf("Error listing widgets: %v", err)
		}
		fmt.Printf("Found %d widget(s) on canvas '%s'\n", len(widgets), canvases[0].Name)

		// Group widgets by type
		typeCount := make(map[string]int)
		for _, w := range widgets {
			typeCount[w.WidgetType]++
		}

		if len(typeCount) > 0 {
			fmt.Println("Widget types:")
			for wType, count := range typeCount {
				fmt.Printf("  - %s: %d\n", wType, count)
			}
		}
	}

	// Step 9: Demonstrate session information
	fmt.Println("\nSession Information:")
	fmt.Printf("  User ID: %d\n", session.UserID())
	fmt.Printf("  API URL: %s\n", apiURL)

	fmt.Println("\nLogin/password authentication example completed!")
	// The deferred Logout() will be called here
}
