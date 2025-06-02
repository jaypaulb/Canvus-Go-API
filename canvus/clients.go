package canvus

import (
	"context"
	"fmt"
)

// ClientInfo represents a client in the Canvus system.
type ClientInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	// Add other fields as needed
}

// CreateClientRequest is the payload for creating a new client.
type CreateClientRequest struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	// Add other fields as needed
}

// UpdateClientRequest is the payload for updating an existing client.
type UpdateClientRequest struct {
	Name *string `json:"name,omitempty"`
	// Add other fields as needed
}

// ListClients retrieves all clients from the Canvus API.
func (c *Client) ListClients(ctx context.Context) ([]ClientInfo, error) {
	var clients []ClientInfo
	err := c.doRequest(ctx, "GET", "clients", nil, &clients, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListClients: %w", err)
	}
	return clients, nil
}

// GetClient retrieves a client by ID from the Canvus API.
func (c *Client) GetClient(ctx context.Context, id string) (*ClientInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("GetClient: id is required")
	}
	var client ClientInfo
	endpoint := fmt.Sprintf("clients/%s", id)
	err := c.doRequest(ctx, "GET", endpoint, nil, &client, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetClient: %w", err)
	}
	return &client, nil
}

// CreateClient creates a new client in the Canvus API.
func (c *Client) CreateClient(ctx context.Context, req CreateClientRequest) (*ClientInfo, error) {
	var client ClientInfo
	err := c.doRequest(ctx, "POST", "clients", req, &client, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateClient: %w", err)
	}
	return &client, nil
}

// UpdateClient updates an existing client by ID in the Canvus API.
func (c *Client) UpdateClient(ctx context.Context, id string, req UpdateClientRequest) (*ClientInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("UpdateClient: id is required")
	}
	var client ClientInfo
	endpoint := fmt.Sprintf("clients/%s", id)
	err := c.doRequest(ctx, "PATCH", endpoint, req, &client, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateClient: %w", err)
	}
	return &client, nil
}

// DeleteClient deletes a client by ID in the Canvus API.
func (c *Client) DeleteClient(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("DeleteClient: id is required")
	}
	endpoint := fmt.Sprintf("clients/%s", id)
	err := c.doRequest(ctx, "DELETE", endpoint, nil, nil, nil, false)
	if err != nil {
		return fmt.Errorf("DeleteClient: %w", err)
	}
	return nil
}

// TestClient wraps a Client and manages a temporary test user and token.
type TestClient struct {
	Client      *Client
	userID      int64
	email       string
	password    string
	cleanupFunc func(context.Context) error
}

// UserClient wraps a Client and manages a temporary token for an existing user.
type UserClient struct {
	Client      *Client
	cleanupFunc func(context.Context) error
}

// NewTestClient creates a new test user, logs in as that user, and returns a TestClient.
// The test user and token are deleted on Cleanup.
func NewTestClient(ctx context.Context, adminClient *Client, baseURL, testEmail, testUsername, testPassword string) (*TestClient, error) {
	// 1. Create user
	user, err := adminClient.CreateUser(ctx, CreateUserRequest{
		Name:     testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	if err != nil {
		return nil, err
	}
	// 2. Login as new user
	testClient := NewClient(baseURL)
	err = testClient.Login(ctx, testEmail, testPassword)
	if err != nil {
		// Cleanup user if login fails
		_ = adminClient.DeleteUser(ctx, user.ID)
		return nil, err
	}
	cleanup := func(ctx context.Context) error {
		_ = testClient.Logout(ctx)
		return adminClient.DeleteUser(ctx, user.ID)
	}
	return &TestClient{
		Client:      testClient,
		userID:      user.ID,
		email:       testEmail,
		password:    testPassword,
		cleanupFunc: cleanup,
	}, nil
}

// Cleanup logs out and deletes the test user.
func (tc *TestClient) Cleanup(ctx context.Context) error {
	if tc.cleanupFunc != nil {
		return tc.cleanupFunc(ctx)
	}
	return nil
}

// NewUserClient logs in as an existing user and returns a UserClient with a temporary token.
// The token is invalidated on Cleanup.
func NewUserClient(ctx context.Context, baseURL, email, password string) (*UserClient, error) {
	client := NewClient(baseURL)
	err := client.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	cleanup := func(ctx context.Context) error {
		return client.Logout(ctx)
	}
	return &UserClient{
		Client:      client,
		cleanupFunc: cleanup,
	}, nil
}

// Cleanup logs out and invalidates the token.
func (uc *UserClient) Cleanup(ctx context.Context) error {
	if uc.cleanupFunc != nil {
		return uc.cleanupFunc(ctx)
	}
	return nil
}

// NewClientFromConfig creates a Client using credentials from a config/settings file.
// This is the standard persistent client; no automatic cleanup is performed.
func NewClientFromConfig(baseURL, apiKey string) *Client {
	return NewClient(baseURL, WithAPIKey(apiKey))
}
