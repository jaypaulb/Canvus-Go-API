package canvus

import (
	"context"
	"fmt"
)

// Folder represents a canvas folder.
type Folder struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"folder_id,omitempty"`
	Access   string `json:"access"`
	InTrash  bool   `json:"in_trash"`
	State    string `json:"state"`
}

// CreateFolderRequest is the payload for creating a folder.
type CreateFolderRequest struct {
	Name     string `json:"name,omitempty"`
	ParentID string `json:"folder_id,omitempty"`
}

// RenameFolderRequest is the payload for renaming a folder.
type RenameFolderRequest struct {
	Name string `json:"name"`
}

// MoveOrCopyFolderRequest is the payload for moving or copying a folder.
type MoveOrCopyFolderRequest struct {
	ParentID  string `json:"folder_id"`
	Conflicts string `json:"conflicts,omitempty"`
}

// FolderPermissions represents permission overrides on a folder.
type FolderPermissions struct {
	EditorsCanShare bool                    `json:"editors_can_share"`
	Users           []FolderUserPermission  `json:"users"`
	Groups          []FolderGroupPermission `json:"groups"`
}

type FolderUserPermission struct {
	ID         int64  `json:"id"`
	Permission string `json:"permission"`
	Inherited  bool   `json:"inherited"`
}

type FolderGroupPermission struct {
	ID         int64  `json:"id"`
	Permission string `json:"permission"`
	Inherited  bool   `json:"inherited"`
}

// ListFolders retrieves all folders.
func (c *Client) ListFolders(ctx context.Context) ([]Folder, error) {
	var folders []Folder
	err := c.doRequest(ctx, "GET", "canvas-folders", nil, &folders, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListFolders: %w", err)
	}
	return folders, nil
}

// GetFolder retrieves a single folder by ID.
func (c *Client) GetFolder(ctx context.Context, id string) (*Folder, error) {
	var folder Folder
	path := fmt.Sprintf("canvas-folders/%s", id)
	err := c.doRequest(ctx, "GET", path, nil, &folder, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetFolder: %w", err)
	}
	return &folder, nil
}

// CreateFolder creates a new folder.
func (c *Client) CreateFolder(ctx context.Context, req CreateFolderRequest) (*Folder, error) {
	var folder Folder
	err := c.doRequest(ctx, "POST", "canvas-folders", req, &folder, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateFolder: %w", err)
	}
	return &folder, nil
}

// RenameFolder renames a folder by ID.
func (c *Client) RenameFolder(ctx context.Context, id string, name string) (*Folder, error) {
	var folder Folder
	path := fmt.Sprintf("canvas-folders/%s", id)
	req := RenameFolderRequest{Name: name}
	err := c.doRequest(ctx, "PATCH", path, req, &folder, nil, false)
	if err != nil {
		return nil, fmt.Errorf("RenameFolder: %w", err)
	}
	return &folder, nil
}

// MoveFolder moves a folder inside another folder.
func (c *Client) MoveFolder(ctx context.Context, id string, parentID string, conflicts string) (*Folder, error) {
	var folder Folder
	path := fmt.Sprintf("canvas-folders/%s/move", id)
	req := MoveOrCopyFolderRequest{ParentID: parentID, Conflicts: conflicts}
	err := c.doRequest(ctx, "POST", path, req, &folder, nil, false)
	if err != nil {
		return nil, fmt.Errorf("MoveFolder: %w", err)
	}
	return &folder, nil
}

// CopyFolder copies a folder inside another folder.
func (c *Client) CopyFolder(ctx context.Context, id string, parentID string, conflicts string) (*Folder, error) {
	var folder Folder
	path := fmt.Sprintf("canvas-folders/%s/copy", id)
	req := MoveOrCopyFolderRequest{ParentID: parentID, Conflicts: conflicts}
	err := c.doRequest(ctx, "POST", path, req, &folder, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CopyFolder: %w", err)
	}
	return &folder, nil
}

// TrashFolder moves a folder to the trash folder.
func (c *Client) TrashFolder(ctx context.Context, id string, trashID string) (*Folder, error) {
	var folder Folder
	path := fmt.Sprintf("canvas-folders/%s/move", id)
	req := MoveOrCopyFolderRequest{ParentID: trashID}
	err := c.doRequest(ctx, "POST", path, req, &folder, nil, false)
	if err != nil {
		return nil, fmt.Errorf("TrashFolder: %w", err)
	}
	return &folder, nil
}

// DeleteFolder permanently deletes a folder by ID.
func (c *Client) DeleteFolder(ctx context.Context, id string) error {
	path := fmt.Sprintf("canvas-folders/%s", id)
	return c.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// DeleteFolderContents deletes all children of a folder.
func (c *Client) DeleteFolderContents(ctx context.Context, id string) error {
	path := fmt.Sprintf("canvas-folders/%s/children", id)
	return c.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// GetFolderPermissions gets the permission overrides on a folder.
func (c *Client) GetFolderPermissions(ctx context.Context, id string) (*FolderPermissions, error) {
	var perms FolderPermissions
	path := fmt.Sprintf("canvas-folders/%s/permissions", id)
	err := c.doRequest(ctx, "GET", path, nil, &perms, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetFolderPermissions: %w", err)
	}
	return &perms, nil
}

// SetFolderPermissions sets permission overrides on a folder.
func (c *Client) SetFolderPermissions(ctx context.Context, id string, perms FolderPermissions) (*FolderPermissions, error) {
	var updated FolderPermissions
	path := fmt.Sprintf("canvas-folders/%s/permissions", id)
	err := c.doRequest(ctx, "POST", path, perms, &updated, nil, false)
	if err != nil {
		return nil, fmt.Errorf("SetFolderPermissions: %w", err)
	}
	return &updated, nil
}
