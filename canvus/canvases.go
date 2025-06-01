// Package canvus provides Canvas resource methods for the Canvus SDK.
package canvus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// ListCanvases returns a list of all canvases.
func (c *Client) ListCanvases(ctx context.Context, opts *ListOptions) ([]Canvas, error) {
	var canvases []Canvas
	endpoint := "canvases"
	queryParams := make(map[string]string)
	if opts != nil {
		if opts.Limit > 0 {
			queryParams["limit"] = fmt.Sprintf("%d", opts.Limit)
		}
		if opts.Offset > 0 {
			queryParams["offset"] = fmt.Sprintf("%d", opts.Offset)
		}
		if opts.Filter != "" {
			queryParams["filter"] = opts.Filter
		}
	}
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &canvases, queryParams, false); err != nil {
		return nil, err
	}
	return canvases, nil
}

// GetCanvas returns a single canvas by ID.
func (c *Client) GetCanvas(ctx context.Context, id string, opts *GetOptions) (Canvas, error) {
	var canvas Canvas
	endpoint := fmt.Sprintf("canvases/%s", id)
	queryParams := make(map[string]string)
	if opts != nil {
		if opts.Subscribe {
			queryParams["subscribe"] = "true"
		}
	}
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &canvas, queryParams, false); err != nil {
		return Canvas{}, err
	}
	return canvas, nil
}

// CreateCanvasRequest represents the payload for creating a canvas.
type CreateCanvasRequest struct {
	Name     string
	FolderID string
	// ... other fields
}

// CreateCanvas creates a new canvas.
func (c *Client) CreateCanvas(ctx context.Context, req CreateCanvasRequest) (Canvas, error) {
	var canvas Canvas
	endpoint := "canvases"
	if err := c.doRequest(ctx, http.MethodPost, endpoint, req, &canvas, nil, false); err != nil {
		return Canvas{}, err
	}
	return canvas, nil
}

// UpdateCanvasRequest represents the payload for updating a canvas.
type UpdateCanvasRequest struct {
	Name string
	Mode string
	// ... other fields
}

// UpdateCanvas updates (renames or changes mode of) a canvas.
func (c *Client) UpdateCanvas(ctx context.Context, id string, req UpdateCanvasRequest) (Canvas, error) {
	var canvas Canvas
	endpoint := fmt.Sprintf("canvases/%s", id)
	if err := c.doRequest(ctx, http.MethodPatch, endpoint, req, &canvas, nil, false); err != nil {
		return Canvas{}, err
	}
	return canvas, nil
}

// DeleteCanvas deletes a canvas by ID.
func (c *Client) DeleteCanvas(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("canvases/%s", id)
	return c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil, nil, false)
}

// GetCanvasPreview returns the preview image bytes for a canvas.
func (c *Client) GetCanvasPreview(ctx context.Context, id string) ([]byte, error) {
	endpoint := fmt.Sprintf("canvases/%s/preview", id)
	var data []byte
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &data, nil, true); err != nil {
		return nil, err
	}
	return data, nil
}

// RestoreCanvas restores a demo canvas by ID.
func (c *Client) RestoreCanvas(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("canvases/%s/restore", id)
	return c.doRequest(ctx, http.MethodPost, endpoint, nil, nil, nil, false)
}

// SaveCanvas saves the demo state of a canvas by ID.
func (c *Client) SaveCanvas(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("canvases/%s/save", id)
	return c.doRequest(ctx, http.MethodPost, endpoint, nil, nil, nil, false)
}

// MoveCanvas moves or trashes a canvas to a different folder.
func (c *Client) MoveCanvas(ctx context.Context, id string, folderID string) error {
	endpoint := fmt.Sprintf("canvases/%s/move", id)
	body := map[string]string{"folder_id": folderID}
	return c.doRequest(ctx, http.MethodPost, endpoint, body, nil, nil, false)
}

// CopyCanvas copies a canvas to a different folder.
func (c *Client) CopyCanvas(ctx context.Context, id string, folderID string) (Canvas, error) {
	var canvas Canvas
	endpoint := fmt.Sprintf("canvases/%s/copy", id)
	body := map[string]string{"folder_id": folderID}
	if err := c.doRequest(ctx, http.MethodPost, endpoint, body, &canvas, nil, false); err != nil {
		return Canvas{}, err
	}
	return canvas, nil
}

// GetCanvasPermissions returns the permissions for a canvas.
func (c *Client) GetCanvasPermissions(ctx context.Context, id string) (Permissions, error) {
	var perms Permissions
	endpoint := fmt.Sprintf("canvases/%s/permissions", id)
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &perms, nil, false); err != nil {
		return Permissions{}, err
	}
	return perms, nil
}

// SetCanvasPermissions sets the permissions for a canvas.
func (c *Client) SetCanvasPermissions(ctx context.Context, id string, perms Permissions) error {
	endpoint := fmt.Sprintf("canvases/%s/permissions", id)
	return c.doRequest(ctx, http.MethodPost, endpoint, perms, nil, nil, false)
}

// CreateImageRequest represents the payload for creating an image on a canvas.
type CreateImageRequest struct {
	Title    string                 `json:"title,omitempty"`
	Location map[string]interface{} `json:"location,omitempty"`
	Scale    float64                `json:"scale,omitempty"`
	Pinned   bool                   `json:"pinned,omitempty"`
	// Add other fields as needed
}

// CreateImage uploads an image to a canvas using multipart/form-data.
func (c *Client) CreateImage(ctx context.Context, canvasID string, filePath string, req CreateImageRequest) (Image, error) {
	var img Image
	endpoint := fmt.Sprintf("canvases/%s/images", canvasID)

	// Prepare multipart form
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	// Add json part
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return img, err
	}
	if err := writer.WriteField("json", string(jsonBytes)); err != nil {
		return img, err
	}

	// Add data part (the image file)
	file, err := os.Open(filePath)
	if err != nil {
		return img, err
	}
	defer file.Close()
	part, err := writer.CreateFormFile("data", filePath)
	if err != nil {
		return img, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return img, err
	}

	if err := writer.Close(); err != nil {
		return img, err
	}

	// Build the HTTP request manually (not using doRequest)
	reqURL := c.BaseURL + endpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, buf)
	if err != nil {
		return img, err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	if c.authenticator != nil {
		c.authenticator.Authenticate(httpReq)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return img, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return img, fmt.Errorf("image upload failed: %s", string(b))
	}

	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		return img, err
	}
	return img, nil
}
