// Package canvus contains shared types for the Canvus SDK.
package canvus

type Canvas struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Access      string `json:"access"`
	AssetSize   int64  `json:"asset_size"`
	CreatedAt   string `json:"created_at"`
	FolderID    string `json:"folder_id"`
	InTrash     bool   `json:"in_trash"`
	Mode        string `json:"mode"`
	ModifiedAt  string `json:"modified_at"`
	PreviewHash string `json:"preview_hash"`
	State       string `json:"state"`
}

type Note struct {
	ID   string
	Text string
	// ... other fields
}

type Image struct {
	ID  string
	URL string
	// ... other fields
}

type PDF struct {
	ID   string
	Name string
	// ... other fields
}

type Video struct {
	ID   string
	Name string
	// ... other fields
}

type Widget struct {
	ID   string
	Type string
	// ... other fields
}

type Anchor struct {
	ID string
	// ... other fields
}

type Browser struct {
	ID string
	// ... other fields
}

type Connector struct {
	ID string
	// ... other fields
}

type Background struct {
	Type string
	// ... other fields
}

type ColorPreset struct {
	Name string
	// ... other fields
}

type MipmapInfo struct {
	Hash string
	// ... other fields
}

type VideoInput struct {
	ID string
	// ... other fields
}

type VideoOutput struct {
	ID string
	// ... other fields
}

type Workspace struct {
	ID string
	// ... other fields
}

// Permissions represents access permissions for a resource.
type Permissions struct {
	// ... fields
}

// CreateCanvasRequest is the payload for creating a canvas.
type CreateCanvasRequest struct {
	Name     string `json:"name,omitempty"`
	FolderID string `json:"folder_id,omitempty"`
}

// UpdateCanvasRequest is the payload for updating a canvas (rename, mode change).
type UpdateCanvasRequest struct {
	Name string `json:"name,omitempty"`
	Mode string `json:"mode,omitempty"`
}

// MoveOrCopyCanvasRequest is the payload for moving or copying a canvas.
type MoveOrCopyCanvasRequest struct {
	FolderID  string `json:"folder_id"`
	Conflicts string `json:"conflicts,omitempty"`
}

// CanvasPermissions represents permission overrides on a canvas.
type CanvasPermissions struct {
	EditorsCanShare bool                    `json:"editors_can_share"`
	Users           []CanvasUserPermission  `json:"users"`
	Groups          []CanvasGroupPermission `json:"groups"`
	LinkPermission  string                  `json:"link_permission"`
}

type CanvasUserPermission struct {
	ID         int64  `json:"id"`
	Permission string `json:"permission"`
	Inherited  bool   `json:"inherited"`
}

type CanvasGroupPermission struct {
	ID         int64  `json:"id"`
	Permission string `json:"permission"`
	Inherited  bool   `json:"inherited"`
}
