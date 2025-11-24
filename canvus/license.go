package canvus

import (
	"context"
	"fmt"
)

// LicenseInfo represents the license information for the server.
type LicenseInfo struct {
	Key       string   `json:"key,omitempty"`
	Valid     bool     `json:"valid,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	Type      string   `json:"type,omitempty"`
	Seats     int      `json:"seats,omitempty"`
	IssuedTo  string   `json:"issued_to,omitempty"`
	IssuedBy  string   `json:"issued_by,omitempty"`
	Features  []string `json:"features,omitempty"`
	// Add other fields as needed based on the API response
}

// GetLicenseInfo retrieves the current license information from the Canvus API.
func (s *Session) GetLicenseInfo(ctx context.Context) (*LicenseInfo, error) {
	var info LicenseInfo
	err := s.doRequest(ctx, "GET", "license", nil, &info, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetLicenseInfo: %w", err)
	}
	return &info, nil
}

// GetActivationRequest retrieves the offline activation request token.
func (s *Session) GetActivationRequest(ctx context.Context) (string, error) {
	var resp map[string]interface{}
	err := s.doRequest(ctx, "GET", "license/request", nil, &resp, nil, false)
	if err != nil {
		return "", fmt.Errorf("GetActivationRequest: %w", err)
	}
	if token, ok := resp["request"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("GetActivationRequest: response did not contain 'request' field")
}

// InstallLicense installs a new license key.
func (s *Session) InstallLicense(ctx context.Context, key string) error {
	req := map[string]string{"key": key}
	return s.doRequest(ctx, "POST", "license", req, nil, nil, false)
}

// ActivateLicense activates the license with an offline activation key.
func (s *Session) ActivateLicense(ctx context.Context, activationKey string) error {
	req := map[string]string{"key": activationKey}
	return s.doRequest(ctx, "POST", "license/activate", req, nil, nil, false)
}
