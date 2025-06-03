package canvus

import (
	"context"
	"fmt"
)

// ListConnectors retrieves all connectors for a given canvas.
func (s *Session) ListConnectors(ctx context.Context, canvasID string) ([]Connector, error) {
	var connectors []Connector
	path := fmt.Sprintf("canvases/%s/connectors", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &connectors, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListConnectors: %w", err)
	}
	return connectors, nil
}

// GetConnector retrieves a connector by ID for a given canvas.
func (s *Session) GetConnector(ctx context.Context, canvasID, connectorID string) (*Connector, error) {
	var connector Connector
	path := fmt.Sprintf("canvases/%s/connectors/%s", canvasID, connectorID)
	err := s.doRequest(ctx, "GET", path, nil, &connector, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetConnector: %w", err)
	}
	return &connector, nil
}

// CreateConnector creates a new connector on a canvas.
func (s *Session) CreateConnector(ctx context.Context, canvasID string, req interface{}) (*Connector, error) {
	var connector Connector
	path := fmt.Sprintf("canvases/%s/connectors", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &connector, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateConnector: %w", err)
	}
	return &connector, nil
}

// UpdateConnector updates a connector by ID for a given canvas.
func (s *Session) UpdateConnector(ctx context.Context, canvasID, connectorID string, req interface{}) (*Connector, error) {
	var connector Connector
	path := fmt.Sprintf("canvases/%s/connectors/%s", canvasID, connectorID)
	err := s.doRequest(ctx, "PATCH", path, req, &connector, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateConnector: %w", err)
	}
	return &connector, nil
}

// DeleteConnector deletes a connector by ID for a given canvas.
func (s *Session) DeleteConnector(ctx context.Context, canvasID, connectorID string) error {
	path := fmt.Sprintf("canvases/%s/connectors/%s", canvasID, connectorID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}
