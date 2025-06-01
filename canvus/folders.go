package canvus

import (
	"context"
	"net/http"
)

// CreateFolderRequest represents the payload for creating a folder.
type CreateFolderRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
}

// ListFolders returns all folders.
func (c *Client) ListFolders(ctx context.Context) ([]Folder, error) {
	var folders []Folder
	endpoint := "canvas-folders"
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &folders, nil, false); err != nil {
		return nil, err
	}
	return folders, nil
}

// CreateFolder creates a new folder.
func (c *Client) CreateFolder(ctx context.Context, req CreateFolderRequest) (Folder, error) {
	var folder Folder
	endpoint := "canvas-folders"
	if err := c.doRequest(ctx, http.MethodPost, endpoint, req, &folder, nil, false); err != nil {
		return Folder{}, err
	}
	return folder, nil
}
