// Package main provides the entry point for the Canvus CLI.
//
// This is a placeholder file. Replace with the actual CLI implementation
// migrated from github.com/jaypaulb/Canvus-Go-API/cmd/canvus-cli/
package main

import (
	"fmt"
	"os"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
)

// Version information (set by build flags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// TODO: Replace this placeholder with actual CLI implementation
	// from github.com/jaypaulb/Canvus-Go-API/cmd/canvus-cli/

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("canvus version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
		os.Exit(0)
	}

	// Placeholder: Show that SDK is accessible
	fmt.Println("Canvus CLI")
	fmt.Println("==========")
	fmt.Println()
	fmt.Println("This is a placeholder. To complete the CLI setup:")
	fmt.Println()
	fmt.Println("1. Copy CLI code from Canvus-Go-API/cmd/canvus-cli/")
	fmt.Println("2. Update imports to use github.com/jaypaulb/Canvus-Go-API/canvus")
	fmt.Println("3. Build and test the CLI")
	fmt.Println()
	fmt.Printf("SDK package available: %T\n", canvus.Session{})
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Printf("  CANVUS_URL:     %s\n", getEnvOrDefault("CANVUS_URL", "(not set)"))
	fmt.Printf("  CANVUS_API_KEY: %s\n", maskAPIKey(os.Getenv("CANVUS_API_KEY")))
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
