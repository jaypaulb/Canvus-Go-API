package canvus

import (
	"context"
	"fmt"
)

// ListWidgets retrieves all widgets for a given canvas. If filter is non-nil, results are filtered client-side.
// This endpoint is read-only: POST, PATCH, DELETE are not supported on /canvases/{id}/widgets.
func (s *Session) ListWidgets(ctx context.Context, canvasID string, filter *Filter) ([]Widget, error) {
	var widgets []Widget
	path := fmt.Sprintf("canvases/%s/widgets", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &widgets, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListWidgets: %w", err)
	}
	if filter != nil {
		widgets = FilterSlice(widgets, filter)
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

// WidgetMatch represents a widget match result across canvases.
type WidgetMatch struct {
	CanvasID string
	WidgetID string
	Widget   Widget
}

// WidgetsLister defines the interface for listing canvases and widgets.
type WidgetsLister interface {
	ListCanvases(ctx context.Context, filter *Filter) ([]Canvas, error)
	ListWidgets(ctx context.Context, canvasID string, filter *Filter) ([]Widget, error)
}

// FindWidgetsAcrossCanvases searches all canvases for widgets matching the given query.
// The query supports exact, wildcard, and partial string matches (see Filter abstraction).
// Returns a slice of WidgetMatch with CanvasID, WidgetID, and the Widget itself.
func FindWidgetsAcrossCanvases(ctx context.Context, lister WidgetsLister, query map[string]interface{}) ([]WidgetMatch, error) {
	canvases, err := lister.ListCanvases(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("FindWidgetsAcrossCanvases: failed to list canvases: %w", err)
	}
	filter := &Filter{Criteria: query}
	var matches []WidgetMatch
	for _, canvas := range canvases {
		widgets, err := lister.ListWidgets(ctx, canvas.ID, filter)
		if err != nil {
			return nil, fmt.Errorf("FindWidgetsAcrossCanvases: failed to list widgets for canvas %s: %w", canvas.ID, err)
		}
		for _, w := range widgets {
			matches = append(matches, WidgetMatch{
				CanvasID: canvas.ID,
				WidgetID: w.ID,
				Widget:   w,
			})
		}
	}
	return matches, nil
}
