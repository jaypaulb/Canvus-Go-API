package canvus

import (
	"context"
	"fmt"
)

// ListBrowsers retrieves all browsers for a given canvas.
func (s *Session) ListBrowsers(ctx context.Context, canvasID string) ([]Browser, error) {
	var browsers []Browser
	path := fmt.Sprintf("canvases/%s/browsers", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &browsers, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListBrowsers: %w", err)
	}
	return browsers, nil
}

// GetBrowser retrieves a browser by ID for a given canvas.
func (s *Session) GetBrowser(ctx context.Context, canvasID, browserID string) (*Browser, error) {
	var browser Browser
	path := fmt.Sprintf("canvases/%s/browsers/%s", canvasID, browserID)
	err := s.doRequest(ctx, "GET", path, nil, &browser, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetBrowser: %w", err)
	}
	return &browser, nil
}

// CreateBrowser creates a new browser on a canvas.
func (s *Session) CreateBrowser(ctx context.Context, canvasID string, req interface{}) (*Browser, error) {
	var browser Browser
	path := fmt.Sprintf("canvases/%s/browsers", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &browser, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateBrowser: %w", err)
	}
	return &browser, nil
}

// UpdateBrowser updates a browser by ID for a given canvas.
func (s *Session) UpdateBrowser(ctx context.Context, canvasID, browserID string, req interface{}) (*Browser, error) {
	var browser Browser
	path := fmt.Sprintf("canvases/%s/browsers/%s", canvasID, browserID)
	err := s.doRequest(ctx, "PATCH", path, req, &browser, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateBrowser: %w", err)
	}
	return &browser, nil
}

// DeleteBrowser deletes a browser by ID for a given canvas.
func (s *Session) DeleteBrowser(ctx context.Context, canvasID, browserID string) error {
	path := fmt.Sprintf("canvases/%s/browsers/%s", canvasID, browserID)
	err := s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
	if err != nil {
		return fmt.Errorf("DeleteBrowser: %w", err)
	}
	return nil
}
