// Package canvus contains shared types for the Canvus SDK.
package canvus

type Canvas struct {
	ID   string
	Name string
	// ... other fields
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

type AuditEvent struct {
	ID string
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

type ServerInfo struct {
	// ... fields
}

// Permissions represents access permissions for a resource.
type Permissions struct {
	// ... fields
}
