package canvus

import (
	"context"
	"fmt"
)

// User represents a user in the Canvus system.
// User contains basic identity and contact information.
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
	Approved  bool   `json:"approved"`
	Blocked   bool   `json:"blocked"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login"`
	State     string `json:"state"`
	// Add other fields as needed
}

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Admin    *bool  `json:"admin,omitempty"`
	Approved *bool  `json:"approved,omitempty"`
	Blocked  *bool  `json:"blocked,omitempty"`
	// Add other fields as needed
}

// UpdateUserRequest is the payload for updating an existing user.
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
	Admin    *bool   `json:"admin,omitempty"`
	Approved *bool   `json:"approved,omitempty"`
	Blocked  *bool   `json:"blocked,omitempty"`
	// Add other fields as needed
}

// ListUsers retrieves all users from the Canvus API.
func (s *Session) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := s.doRequest(ctx, "GET", "users", nil, &users, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListUsers: %w", err)
	}
	return users, nil
}

// GetUser retrieves a user by ID from the Canvus API.
func (s *Session) GetUser(ctx context.Context, id int64) (*User, error) {
	var user User
	endpoint := fmt.Sprintf("users/%d", id)
	err := s.doRequest(ctx, "GET", endpoint, nil, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	return &user, nil
}

// CreateUser creates a new user in the Canvus API.
// req can be CreateUserRequest or map[string]interface{}
func (s *Session) CreateUser(ctx context.Context, req interface{}) (*User, error) {
	var user User
	err := s.doRequest(ctx, "POST", "users", req, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	return &user, nil
}

// UpdateUser updates an existing user by ID in the Canvus API.
// req can be UpdateUserRequest or map[string]interface{}
func (s *Session) UpdateUser(ctx context.Context, id int64, req interface{}) (*User, error) {
	var user User
	endpoint := fmt.Sprintf("users/%d", id)
	err := s.doRequest(ctx, "PATCH", endpoint, req, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser: %w", err)
	}
	return &user, nil
}

// DeleteUser deletes a user by ID in the Canvus API.
func (s *Session) DeleteUser(ctx context.Context, id int64) error {
	endpoint := fmt.Sprintf("users/%d", id)
	err := s.doRequest(ctx, "DELETE", endpoint, nil, nil, nil, false)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}
	return nil
}

// SamlLoginRequest represents the payload for SAML login.
type SamlLoginRequest struct {
	InResponseTo string `json:"inResponseTo"`
	ResponseXML  string `json:"responseXml"`
	Remember     bool   `json:"remember"`
}

// ValidateResetToken checks if a password reset token is valid.
func (s *Session) ValidateResetToken(ctx context.Context, token string) error {
	return s.doRequest(ctx, "GET", "users/password/validate-reset-token", nil, nil, map[string]string{"token": token}, false)
}

// RegisterUser registers a new user.
func (s *Session) RegisterUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	var user User
	err := s.doRequest(ctx, "POST", "users/register", req, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}
	return &user, nil
}

// ConfirmEmail confirms a user's email address using a token.
func (s *Session) ConfirmEmail(ctx context.Context, token string) error {
	req := map[string]string{"token": token}
	return s.doRequest(ctx, "POST", "users/confirm-email", req, nil, nil, false)
}

// CreateResetToken creates a password reset token for a user.
func (s *Session) CreateResetToken(ctx context.Context, email string) error {
	req := map[string]string{"email": email}
	return s.doRequest(ctx, "POST", "users/password/create-reset-token", req, nil, nil, false)
}

// ResetUserPassword resets a user's password using a token.
func (s *Session) ResetUserPassword(ctx context.Context, token, newPassword string) error {
	req := map[string]string{
		"token":    token,
		"password": newPassword,
	}
	return s.doRequest(ctx, "POST", "users/password/reset", req, nil, nil, false)
}

// ChangeUserEmail changes a user's email address.
func (s *Session) ChangeUserEmail(ctx context.Context, userID int64, newEmail string) error {
	req := map[string]string{"email": newEmail}
	return s.doRequest(ctx, "POST", fmt.Sprintf("users/%d/change-email", userID), req, nil, nil, false)
}

// SetUserPassword sets a user's password (admin action).
func (s *Session) SetUserPassword(ctx context.Context, userID int64, newPassword string) error {
	req := map[string]string{"password": newPassword}
	return s.doRequest(ctx, "POST", fmt.Sprintf("users/%d/password", userID), req, nil, nil, false)
}

// BlockUser blocks a user.
func (s *Session) BlockUser(ctx context.Context, userID int64) error {
	return s.doRequest(ctx, "POST", fmt.Sprintf("users/%d/block", userID), nil, nil, nil, false)
}

// UnblockUser unblocks a user.
func (s *Session) UnblockUser(ctx context.Context, userID int64) error {
	return s.doRequest(ctx, "POST", fmt.Sprintf("users/%d/unblock", userID), nil, nil, nil, false)
}

// ApproveUser approves a user.
func (s *Session) ApproveUser(ctx context.Context, userID int64) error {
	return s.doRequest(ctx, "POST", fmt.Sprintf("users/%d/approve", userID), nil, nil, nil, false)
}

// ForcePasswordResetUser forces a password reset for a user.
func (s *Session) ForcePasswordResetUser(ctx context.Context, userID int64) error {
	return s.doRequest(ctx, "POST", fmt.Sprintf("users/%d/reset-password", userID), nil, nil, nil, false)
}

// SamlLogin performs a SAML login.
func (s *Session) SamlLogin(ctx context.Context, req SamlLoginRequest) error {
	return s.doRequest(ctx, "POST", "users/login/saml", req, nil, nil, false)
}
