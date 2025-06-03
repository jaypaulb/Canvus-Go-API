// Package canvus contains shared types for the Canvus SDK.
package canvus

import (
	"time"
)

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
	ID              string  `json:"id"`
	Text            string  `json:"text"`
	Title           string  `json:"title"`
	BackgroundColor string  `json:"background_color"`
	Depth           int     `json:"depth"`
	Location        *Point  `json:"location,omitempty"`
	ParentID        string  `json:"parent_id"`
	Pinned          bool    `json:"pinned"`
	Scale           float64 `json:"scale"`
	Size            *Size   `json:"size,omitempty"`
	State           string  `json:"state"`
	WidgetType      string  `json:"widget_type"`
}

type Image struct {
	ID               string  `json:"id"`
	Hash             string  `json:"hash"`
	Title            string  `json:"title"`
	OriginalFilename string  `json:"original_filename"`
	ParentID         string  `json:"parent_id"`
	Pinned           bool    `json:"pinned"`
	Scale            float64 `json:"scale"`
	Size             *Size   `json:"size,omitempty"`
	Location         *Point  `json:"location,omitempty"`
	State            string  `json:"state"`
	WidgetType       string  `json:"widget_type"`
	Depth            int     `json:"depth"`
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
	ID         string  `json:"id"`
	WidgetType string  `json:"widget_type"`
	ParentID   string  `json:"parent_id"`
	Location   *Point  `json:"location,omitempty"`
	Size       *Size   `json:"size,omitempty"`
	Pinned     bool    `json:"pinned"`
	Scale      float64 `json:"scale"`
	State      string  `json:"state"`
	Depth      int     `json:"depth"`
}

type Anchor struct {
	ID          string  `json:"id"`
	AnchorIndex int     `json:"anchor_index"`
	AnchorName  string  `json:"anchor_name"`
	ParentID    string  `json:"parent_id"`
	Pinned      bool    `json:"pinned"`
	Scale       float64 `json:"scale"`
	Size        *Size   `json:"size,omitempty"`
	Location    *Point  `json:"location,omitempty"`
	State       string  `json:"state"`
	WidgetType  string  `json:"widget_type"`
	Depth       int     `json:"depth"`
}

type Browser struct {
	ID string
	// ... other fields
}

type Connector struct {
	ID         string        `json:"id"`
	Src        *ConnectorEnd `json:"src,omitempty"`
	Dst        *ConnectorEnd `json:"dst,omitempty"`
	LineColor  string        `json:"line_color"`
	LineWidth  int           `json:"line_width"`
	State      string        `json:"state"`
	Type       string        `json:"type"`
	WidgetType string        `json:"widget_type"`
}

type ConnectorEnd struct {
	AutoLocation bool   `json:"auto_location"`
	ID           string `json:"id"`
	RelLocation  *Point `json:"rel_location,omitempty"`
	Tip          string `json:"tip"`
}

type Background struct {
	Type string
	// ... other fields
}

type ColorPreset struct {
	Name string
	// ... other fields
}

// ColorPresets represents the color presets for a canvas.
type ColorPresets struct {
	Annotation     []string `json:"annotation"`
	Connector      []string `json:"connector"`
	NoteBackground []string `json:"note_background"`
	NoteText       []string `json:"note_text"`
}

// MipmapInfo represents mipmap information for an asset.
type MipmapInfo struct {
	Resolution struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"resolution"`
	MaxLevel int `json:"max_level"`
	Pages    int `json:"pages"`
}

type VideoInput struct {
	ID string
	// ... other fields
}

// VideoOutput represents a video output channel on a client or canvas.
type VideoOutput struct {
	Index      int    `json:"index,omitempty"`
	Label      string `json:"label,omitempty"`
	Source     string `json:"source,omitempty"`
	Suspended  bool   `json:"suspended,omitempty"`
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Resolution *Size  `json:"resolution,omitempty"`
	State      string `json:"state,omitempty"`
	WidgetType string `json:"widget_type,omitempty"`
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
	CanvasID     string        `json:"canvas_id"`
	ServerID     string        `json:"server_id,omitempty"`
	UserEmail    string        `json:"user_email,omitempty"`
	CenterX      *float64      `json:"center_x,omitempty"`
	CenterY      *float64      `json:"center_y,omitempty"`
	WidgetID     *string       `json:"widget_id,omitempty"`
	PollTimeout  time.Duration `json:"-"`
	PollInterval time.Duration `json:"-"`
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

// CanvasBackground represents the background settings for a canvas.
type CanvasBackground struct {
	Type            string           `json:"type"`
	Haze            *HazeSettings    `json:"haze,omitempty"`
	Grid            *GridSettings    `json:"grid,omitempty"`
	Image           *BackgroundImage `json:"image,omitempty"`
	BackgroundColor string           `json:"background_color,omitempty"`
}

// HazeSettings represents haze background settings.
type HazeSettings struct {
	Color1 string  `json:"color1"`
	Color2 string  `json:"color2"`
	Speed  float64 `json:"speed"`
	Scale  float64 `json:"scale"`
}

// GridSettings represents grid overlay settings.
type GridSettings struct {
	Visible bool   `json:"visible"`
	Color   string `json:"color"`
}

// BackgroundImage represents image background settings.
type BackgroundImage struct {
	Hash string `json:"hash"`
	Fit  string `json:"fit"`
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

type Asset struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url,omitempty"`
	// Add more fields as needed for specific asset types
}

type VideoInputSource struct {
	Name       string `json:"name"`
	Resolution *Size  `json:"resolution,omitempty"`
	Source     string `json:"source"`
}

type Annotation struct {
	ID         string `json:"id"`
	Text       string `json:"text,omitempty"`
	WidgetType string `json:"widget_type"`
	// Add more fields as needed
}
