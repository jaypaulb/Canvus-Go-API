// Package main demonstrates user provisioning and management with the Canvus Go SDK.
//
// This example shows:
// - Creating new users programmatically
// - Updating user properties (admin status, blocked status)
// - Generating API access tokens for users
// - Listing and managing user tokens
// - Deleting users and cleaning up resources
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-admin-api-key-here"
//   go run provision_and_manage.go
//
// Note: This example requires admin privileges to create and manage users.
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

	fmt.Println("Canvus Go SDK - User Provisioning and Management Example")
	fmt.Println("=========================================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	// This API key must belong to an admin user
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Create context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 3: List existing users to see current state
	fmt.Println("Listing existing users...")

	users, err := session.ListUsers(ctx)
	if err != nil {
		// Check for permission errors
		if apiErr, ok := err.(*canvus.APIError); ok {
			if apiErr.StatusCode == 403 {
				log.Fatalf("Permission denied. Make sure you're using an admin API key: %v", err)
			}
		}
		log.Fatalf("Error listing users: %v", err)
	}

	fmt.Printf("Found %d existing user(s)\n\n", len(users))

	// Step 4: Create a new user
	// Generate a unique email using timestamp to avoid conflicts
	timestamp := time.Now().Format("20060102150405")
	newUserEmail := fmt.Sprintf("sdk-test-user-%s@example.com", timestamp)
	newUserName := fmt.Sprintf("SDK Test User %s", timestamp)

	fmt.Printf("Creating new user: %s (%s)\n", newUserName, newUserEmail)

	// Set initial properties for the new user
	// Using pointers for optional boolean fields allows us to explicitly set them
	approved := true
	admin := false
	blocked := false

	createReq := canvus.CreateUserRequest{
		Email:    newUserEmail,
		Name:     newUserName,
		Password: "SecurePassword123!", // Initial password
		Approved: &approved,            // User is approved and can log in
		Admin:    &admin,               // Not an admin initially
		Blocked:  &blocked,             // Not blocked
	}

	newUser, err := session.CreateUser(ctx, createReq)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}

	fmt.Printf("User created successfully!\n")
	fmt.Printf("  ID: %d\n", newUser.ID)
	fmt.Printf("  Name: %s\n", newUser.Name)
	fmt.Printf("  Email: %s\n", newUser.Email)
	fmt.Printf("  Admin: %v\n", newUser.Admin)
	fmt.Printf("  Approved: %v\n", newUser.Approved)
	fmt.Printf("  Blocked: %v\n", newUser.Blocked)
	fmt.Println()

	// Ensure cleanup - delete the user when done
	defer func() {
		fmt.Println("\nCleaning up: Deleting test user...")
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()

		if err := session.DeleteUser(cleanupCtx, newUser.ID); err != nil {
			log.Printf("Warning: Failed to delete user: %v", err)
		} else {
			fmt.Printf("User '%s' (ID: %d) deleted successfully\n", newUser.Email, newUser.ID)
		}
	}()

	// Step 5: Get user details by ID
	fmt.Println("Retrieving user details...")

	retrievedUser, err := session.GetUser(ctx, newUser.ID)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}

	fmt.Printf("Retrieved user: %s (ID: %d)\n", retrievedUser.Name, retrievedUser.ID)
	fmt.Printf("  State: %s\n", retrievedUser.State)
	fmt.Printf("  Created At: %s\n", retrievedUser.CreatedAt)
	fmt.Println()

	// Step 6: Update user properties
	fmt.Println("Updating user properties...")

	// Update the user's name and grant admin privileges
	updatedName := newUserName + " (Updated)"
	adminTrue := true

	updateReq := canvus.UpdateUserRequest{
		Name:  &updatedName,
		Admin: &adminTrue, // Grant admin privileges
	}

	updatedUser, err := session.UpdateUser(ctx, newUser.ID, updateReq)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
	}

	fmt.Printf("User updated successfully!\n")
	fmt.Printf("  New Name: %s\n", updatedUser.Name)
	fmt.Printf("  Admin Status: %v\n", updatedUser.Admin)
	fmt.Println()

	// Step 7: Create an API access token for the user
	// This allows the user to authenticate via API key instead of password
	fmt.Println("Creating API access token for user...")

	tokenReq := canvus.CreateAccessTokenRequest{
		Description: fmt.Sprintf("SDK Example Token %s", timestamp),
	}

	accessToken, err := session.CreateAccessToken(ctx, newUser.ID, tokenReq)
	if err != nil {
		log.Fatalf("Error creating access token: %v", err)
	}

	fmt.Printf("Access token created!\n")
	fmt.Printf("  Token ID: %s\n", accessToken.ID)
	fmt.Printf("  Description: %s\n", accessToken.Description)
	// The plain token is only available immediately after creation
	if accessToken.PlainToken != "" {
		fmt.Printf("  Token (save this - only shown once): %s\n", accessToken.PlainToken)
	}
	fmt.Println()

	// Step 8: List all access tokens for the user
	fmt.Println("Listing user's access tokens...")

	tokens, err := session.ListAccessTokens(ctx, newUser.ID)
	if err != nil {
		log.Fatalf("Error listing access tokens: %v", err)
	}

	fmt.Printf("User has %d access token(s):\n", len(tokens))
	for _, token := range tokens {
		fmt.Printf("  - %s: %s\n", token.ID, token.Description)
	}
	fmt.Println()

	// Step 9: Retrieve a specific access token
	fmt.Println("Retrieving specific access token details...")

	retrievedToken, err := session.GetAccessToken(ctx, newUser.ID, accessToken.ID)
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	fmt.Printf("Token details:\n")
	fmt.Printf("  ID: %s\n", retrievedToken.ID)
	fmt.Printf("  Description: %s\n", retrievedToken.Description)
	fmt.Println()

	// Step 10: Delete the access token
	fmt.Println("Deleting access token...")

	err = session.DeleteAccessToken(ctx, newUser.ID, accessToken.ID)
	if err != nil {
		log.Fatalf("Error deleting access token: %v", err)
	}

	fmt.Printf("Access token %s deleted successfully\n", accessToken.ID)

	// Verify deletion
	remainingTokens, err := session.ListAccessTokens(ctx, newUser.ID)
	if err != nil {
		log.Fatalf("Error listing tokens after deletion: %v", err)
	}

	fmt.Printf("User now has %d access token(s)\n", len(remainingTokens))
	fmt.Println()

	// Step 11: Demonstrate blocking a user
	fmt.Println("Demonstrating user blocking...")

	blockedTrue := true
	blockReq := canvus.UpdateUserRequest{
		Blocked: &blockedTrue,
	}

	blockedUser, err := session.UpdateUser(ctx, newUser.ID, blockReq)
	if err != nil {
		log.Fatalf("Error blocking user: %v", err)
	}

	fmt.Printf("User blocked: %v\n", blockedUser.Blocked)

	// Unblock the user
	blockedFalse := false
	unblockReq := canvus.UpdateUserRequest{
		Blocked: &blockedFalse,
	}

	unblockedUser, err := session.UpdateUser(ctx, newUser.ID, unblockReq)
	if err != nil {
		log.Fatalf("Error unblocking user: %v", err)
	}

	fmt.Printf("User unblocked: %v\n", !unblockedUser.Blocked)
	fmt.Println()

	// Step 12: Change user password
	fmt.Println("Changing user password...")

	newPassword := "NewSecurePassword456!"
	passwordReq := canvus.UpdateUserRequest{
		Password: &newPassword,
	}

	_, err = session.UpdateUser(ctx, newUser.ID, passwordReq)
	if err != nil {
		log.Fatalf("Error changing password: %v", err)
	}

	fmt.Println("Password changed successfully")
	fmt.Println()

	// Final user state
	fmt.Println("Final user state:")
	finalUser, err := session.GetUser(ctx, newUser.ID)
	if err != nil {
		log.Fatalf("Error getting final user state: %v", err)
	}

	fmt.Printf("  ID: %d\n", finalUser.ID)
	fmt.Printf("  Name: %s\n", finalUser.Name)
	fmt.Printf("  Email: %s\n", finalUser.Email)
	fmt.Printf("  Admin: %v\n", finalUser.Admin)
	fmt.Printf("  Approved: %v\n", finalUser.Approved)
	fmt.Printf("  Blocked: %v\n", finalUser.Blocked)

	fmt.Println("\nUser provisioning and management example completed successfully!")
	// The deferred cleanup will delete the test user
}
