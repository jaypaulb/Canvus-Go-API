package canvus

import (
	"context"
	"fmt"
)

// ListCanvases retrieves all canvases.
func (c *Client) ListCanvases(ctx context.Context) ([]Canvas, error) {
	var canvases []Canvas
	err := c.doRequest(ctx, "GET", "canvases", nil, &canvases, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListCanvases: %w", err)
	}
	return canvases, nil
}

// GetCanvas retrieves a single canvas by ID.
func (c *Client) GetCanvas(ctx context.Context, id string) (*Canvas, error) {
	var canvas Canvas
	path := fmt.Sprintf("canvases/%s", id)
	err := c.doRequest(ctx, "GET", path, nil, &canvas, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetCanvas: %w", err)
	}
	return &canvas, nil
}

// CreateCanvas creates a new canvas.
func (c *Client) CreateCanvas(ctx context.Context, req CreateCanvasRequest) (*Canvas, error) {
	var canvas Canvas
	err := c.doRequest(ctx, "POST", "canvases", req, &canvas, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateCanvas: %w", err)
	}
	return &canvas, nil
}

// UpdateCanvas renames or changes the mode of a canvas by ID.
func (c *Client) UpdateCanvas(ctx context.Context, id string, req UpdateCanvasRequest) (*Canvas, error) {
	var canvas Canvas
	path := fmt.Sprintf("canvases/%s", id)
	err := c.doRequest(ctx, "PATCH", path, req, &canvas, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateCanvas: %w", err)
	}
	return &canvas, nil
}

// DeleteCanvas permanently deletes a canvas by ID.
func (c *Client) DeleteCanvas(ctx context.Context, id string) error {
	path := fmt.Sprintf("canvases/%s", id)
	return c.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// GetCanvasPreview downloads the preview of a canvas, if available.
func (c *Client) GetCanvasPreview(ctx context.Context, id string) ([]byte, error) {
	path := fmt.Sprintf("canvases/%s/preview", id)
	var preview []byte
	err := c.doRequest(ctx, "GET", path, nil, &preview, nil, true)
	if err != nil {
		return nil, fmt.Errorf("GetCanvasPreview: %w", err)
	}
	return preview, nil
}

// RestoreDemoCanvas restores the state of a demo canvas to the last saved state.
func (c *Client) RestoreDemoCanvas(ctx context.Context, id string) error {
	path := fmt.Sprintf("canvases/%s/restore", id)
	return c.doRequest(ctx, "POST", path, nil, nil, nil, false)
}

// SaveDemoState updates the saved demo canvas state with the current changes.
func (c *Client) SaveDemoState(ctx context.Context, id string) error {
	path := fmt.Sprintf("canvases/%s/save", id)
	return c.doRequest(ctx, "POST", path, nil, nil, nil, false)
}

// MoveCanvas moves a canvas to another folder.
func (c *Client) MoveCanvas(ctx context.Context, id string, req MoveOrCopyCanvasRequest) (*Canvas, error) {
	var canvas Canvas
	path := fmt.Sprintf("canvases/%s/move", id)
	err := c.doRequest(ctx, "POST", path, req, &canvas, nil, false)
	if err != nil {
		return nil, fmt.Errorf("MoveCanvas: %w", err)
	}
	return &canvas, nil
}

// CopyCanvas copies a canvas to another folder.
func (c *Client) CopyCanvas(ctx context.Context, id string, req MoveOrCopyCanvasRequest) (*Canvas, error) {
	var canvas Canvas
	path := fmt.Sprintf("canvases/%s/copy", id)
	err := c.doRequest(ctx, "POST", path, req, &canvas, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CopyCanvas: %w", err)
	}
	return &canvas, nil
}

// TrashCanvas moves a canvas to the trash folder.
func (c *Client) TrashCanvas(ctx context.Context, id string, _ string) (*Canvas, error) {
	userID := c.UserID()
	if userID == 0 {
		return nil, fmt.Errorf("TrashCanvas: user ID not set; must login first")
	}
	trashID := fmt.Sprintf("trash.%d", userID)
	var canvas Canvas
	path := fmt.Sprintf("canvases/%s/move", id)
	req := MoveOrCopyCanvasRequest{FolderID: trashID}
	err := c.doRequest(ctx, "POST", path, req, &canvas, nil, false)
	if err != nil {
		return nil, fmt.Errorf("TrashCanvas: %w", err)
	}
	return &canvas, nil
}

// GetCanvasPermissions gets the permission overrides on a canvas.
func (c *Client) GetCanvasPermissions(ctx context.Context, id string) (*CanvasPermissions, error) {
	var perms CanvasPermissions
	path := fmt.Sprintf("canvases/%s/permissions", id)
	err := c.doRequest(ctx, "GET", path, nil, &perms, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetCanvasPermissions: %w", err)
	}
	return &perms, nil
}

// SetCanvasPermissions sets permission overrides on a canvas.
func (c *Client) SetCanvasPermissions(ctx context.Context, id string, perms CanvasPermissions) (*CanvasPermissions, error) {
	var updated CanvasPermissions
	path := fmt.Sprintf("canvases/%s/permissions", id)
	err := c.doRequest(ctx, "POST", path, perms, &updated, nil, false)
	if err != nil {
		return nil, fmt.Errorf("SetCanvasPermissions: %w", err)
	}
	return &updated, nil
}
