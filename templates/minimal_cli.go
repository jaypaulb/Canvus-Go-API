// Minimal CLI Template
//
// This template provides a command-line tool structure for building
// utilities that interact with the Canvus API. It includes:
// - Flag parsing for command-line arguments
// - Configuration from environment variables
// - Proper error handling with exit codes
// - Structured logging
//
// Usage:
//   1. Copy this file to your project
//   2. Search for "TODO:" comments and customize
//   3. Build with: go build -o your-tool
//   4. Run with: ./your-tool [flags]
//
// Environment Variables:
//   CANVUS_API_URL - Required: Canvus API endpoint
//   CANVUS_API_KEY - Required: API key for authentication
//
// Exit Codes:
//   0 - Success
//   1 - General error
//   2 - Configuration error

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
	// TODO: Add your imports here
)

// Config holds the application configuration
type Config struct {
	APIURL    string
	APIKey    string
	Timeout   time.Duration
	Verbose   bool
	// TODO: Add your configuration fields here
}

// Exit codes
const (
	exitSuccess = 0
	exitError   = 1
	exitConfig  = 2
)

func main() {
	// Parse command-line flags
	cfg := parseFlags()

	// Load environment configuration
	if err := loadEnvConfig(cfg); err != nil {
		log.Printf("Configuration error: %v", err)
		os.Exit(exitConfig)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Run the main logic
	if err := run(ctx, cfg); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(exitError)
	}

	os.Exit(exitSuccess)
}

// parseFlags parses command-line flags and returns configuration
func parseFlags() *Config {
	cfg := &Config{}

	// Standard flags
	flag.DurationVar(&cfg.Timeout, "timeout", 5*time.Minute, "Operation timeout")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output")

	// TODO: Define your command-line flags here
	// Example:
	// flag.StringVar(&cfg.CanvasID, "canvas", "", "Canvas ID to operate on")
	// flag.BoolVar(&cfg.DryRun, "dry-run", false, "Show what would be done without making changes")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for Canvus operations.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  CANVUS_API_URL    Canvus API endpoint (required)\n")
		fmt.Fprintf(os.Stderr, "  CANVUS_API_KEY    API key for authentication (required)\n")
		// TODO: Document your environment variables here
	}

	flag.Parse()
	return cfg
}

// loadEnvConfig loads configuration from environment variables
func loadEnvConfig(cfg *Config) error {
	// Required: API URL
	cfg.APIURL = os.Getenv("CANVUS_API_URL")
	if cfg.APIURL == "" {
		return fmt.Errorf("CANVUS_API_URL environment variable is required")
	}

	// Required: API Key
	cfg.APIKey = os.Getenv("CANVUS_API_KEY")
	if cfg.APIKey == "" {
		return fmt.Errorf("CANVUS_API_KEY environment variable is required")
	}

	// TODO: Load your environment variables here
	// Example:
	// cfg.DefaultFolderID = os.Getenv("CANVUS_DEFAULT_FOLDER")

	return nil
}

// run executes the main business logic
func run(ctx context.Context, cfg *Config) error {
	// Create SDK configuration
	sdkCfg := canvus.DefaultSessionConfig()
	sdkCfg.BaseURL = cfg.APIURL
	sdkCfg.RequestTimeout = cfg.Timeout

	// Create session with API key authentication
	session := canvus.NewSession(sdkCfg, canvus.WithAPIKey(cfg.APIKey))

	if cfg.Verbose {
		log.Printf("Connected to Canvus API at %s", cfg.APIURL)
	}

	// TODO: Add your business logic here
	//
	// Example 1: List canvases
	// canvases, err := session.ListCanvases(ctx, nil)
	// if err != nil {
	//     return fmt.Errorf("failed to list canvases: %w", err)
	// }
	// for _, canvas := range canvases {
	//     fmt.Printf("Canvas: %s (ID: %s)\n", canvas.Name, canvas.ID)
	// }
	//
	// Example 2: Get a specific canvas
	// canvas, err := session.GetCanvas(ctx, cfg.CanvasID)
	// if err != nil {
	//     return handleAPIError(err, "get canvas")
	// }
	// fmt.Printf("Canvas: %s\n", canvas.Name)
	//
	// Example 3: Create a widget
	// widget := &canvus.Widget{
	//     WidgetType: "Note",
	//     Location:   [2]float64{100, 100},
	//     Size:       [2]float64{200, 150},
	//     Text:       "Hello from CLI",
	// }
	// created, err := session.CreateWidget(ctx, canvasID, widget)
	// if err != nil {
	//     return handleAPIError(err, "create widget")
	// }
	// fmt.Printf("Created widget: %s\n", created.ID)

	// Placeholder: List canvases as a demo
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		return handleAPIError(err, "list canvases")
	}

	if cfg.Verbose {
		log.Printf("Found %d canvases", len(canvases))
	}

	for _, canvas := range canvases {
		fmt.Printf("%s\t%s\n", canvas.ID, canvas.Name)
	}

	return nil
}

// handleAPIError provides detailed error handling for API errors
func handleAPIError(err error, operation string) error {
	if apiErr, ok := err.(*canvus.APIError); ok {
		return fmt.Errorf("%s failed: [%d] %s", operation, apiErr.StatusCode, apiErr.Message)
	}
	return fmt.Errorf("%s failed: %w", operation, err)
}
