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
	ID       string `json:"id"`
	Type     string `json:"type"`
	Location *Point `json:"location,omitempty"`
	Size     *Size  `json:"size,omitempty"`
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
	CanvasID         string     `json:"canvas_id"`
	CanvasSize       *Size      `json:"canvas_size,omitempty"`
	Index            int        `json:"index"`
	InfoPanelVisible bool       `json:"info_panel_visible"`
	Location         *Point     `json:"location,omitempty"`
	Pinned           bool       `json:"pinned"`
	ServerID         string     `json:"server_id"`
	Size             *Size      `json:"size,omitempty"`
	State            string     `json:"state"`
	User             string     `json:"user"`
	ViewRectangle    *Rectangle `json:"view_rectangle,omitempty"`
	WorkspaceName    string     `json:"workspace_name"`
	WorkspaceState   string     `json:"workspace_state"`
}

type Size struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Rectangle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type WorkspaceSelector struct {
	Index *int
	Name  *string
	User  *string
}

type Viewport struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type SetViewportOptions struct {
	WidgetID *string
	X        *float64
	Y        *float64
	Width    *float64
	Height   *float64
	Margin   float64
}

type OpenCanvasOptions struct {
	CanvasID  string
	ServerID  string
	UserEmail string
	CenterX   *float64
	CenterY   *float64
	WidgetID  *string
}

type UpdateWorkspaceRequest struct {
	InfoPanelVisible *bool      `json:"info_panel_visible,omitempty"`
	Pinned           *bool      `json:"pinned,omitempty"`
	ViewRectangle    *Rectangle `json:"view_rectangle,omitempty"`
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
