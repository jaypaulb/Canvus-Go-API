package canvus

import (
	"context"
	"fmt"
)

// ListWidgets retrieves all widgets for a given canvas.
// This endpoint is read-only: POST, PATCH, DELETE are not supported on /canvases/{id}/widgets.
func (s *Session) ListWidgets(ctx context.Context, canvasID string) ([]Widget, error) {
	var widgets []Widget
	path := fmt.Sprintf("canvases/%s/widgets", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &widgets, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListWidgets: %w", err)
	}
	return widgets, nil
}

// GetWidget retrieves a widget by ID for a given canvas.
func (s *Session) GetWidget(ctx context.Context, canvasID, widgetID string) (*Widget, error) {
	var widget Widget
	path := fmt.Sprintf("canvases/%s/widgets/%s", canvasID, widgetID)
	err := s.doRequest(ctx, "GET", path, nil, &widget, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetWidget: %w", err)
	}
	return &widget, nil
}

// CreateWidget creates a new widget on a canvas.
func (s *Session) CreateWidget(ctx context.Context, canvasID string, req interface{}) (*Widget, error) {
	var widget Widget
	path := fmt.Sprintf("canvases/%s/widgets", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &widget, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateWidget: %w", err)
	}
	return &widget, nil
}

// UpdateWidget updates a widget by ID.
func (s *Session) UpdateWidget(ctx context.Context, canvasID, widgetID string, req interface{}) (*Widget, error) {
	var widget Widget
	path := fmt.Sprintf("canvases/%s/widgets/%s", canvasID, widgetID)
	err := s.doRequest(ctx, "PATCH", path, req, &widget, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateWidget: %w", err)
	}
	return &widget, nil
}

// DeleteWidget deletes a widget by ID.
func (s *Session) DeleteWidget(ctx context.Context, canvasID, widgetID string) error {
	path := fmt.Sprintf("canvases/%s/widgets/%s", canvasID, widgetID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// PatchParentID updates the parent ID of a widget (parenting).
func (s *Session) PatchParentID(ctx context.Context, canvasID, widgetID, parentID string) (*Widget, error) {
	var widget Widget
	path := fmt.Sprintf("canvases/%s/widgets/%s", canvasID, widgetID)
	req := map[string]interface{}{"parent_id": parentID}
	err := s.doRequest(ctx, "PATCH", path, req, &widget, nil, false)
	if err != nil {
		return nil, fmt.Errorf("PatchParentID: %w", err)
	}
	return &widget, nil
}

// ListAnnotations retrieves all annotations for a given canvas.
// This endpoint is read-only: only GET is supported on /canvases/{id}/annotations.
func (s *Session) ListAnnotations(ctx context.Context, canvasID string) ([]Annotation, error) {
	var annotations []Annotation
	path := fmt.Sprintf("canvases/%s/annotations", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &annotations, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListAnnotations: %w", err)
	}
	return annotations, nil
}
