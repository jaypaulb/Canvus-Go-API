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

// WidgetsContainId returns all widgets on the given canvas that are fully contained within the bounding box of the source widget (with optional tolerance).
//
// Parameters:
//
//	ctx      - context for cancellation and deadlines
//	s        - the Session (API client)
//	canvasID - the ID of the canvas containing the widgets
//	widgetID - the ID of the source widget (if widget is nil)
//	widget   - the full Widget struct (if available; if nil, widgetID is used to fetch it)
//	tolerance- expands the bounding box by this amount in all directions (use 0 for exact)
//
// Behavior:
//   - If widget is nil, the function fetches the widget using canvasID and widgetID.
//   - If widget is provided, it is used directly.
//   - All widgets on the same canvas are fetched.
//   - Each widget (except the source) is checked to see if it is fully contained within the (optionally tolerance-expanded) bounding box of the source widget.
//   - Returns a slice of all contained widgets.
//
// Usage Example:
//
//	contained, err := canvus.WidgetsContainId(ctx, session, "canvas123", "widget456", nil, 0)
//	// or, if you already have the widget:
//	contained, err := canvus.WidgetsContainId(ctx, session, "canvas123", "", &myWidget, 5)
func WidgetsContainId(ctx context.Context, s *Session, canvasID string, widgetID string, widget *Widget, tolerance float64) ([]Widget, error) {
	var srcWidget Widget
	if widget != nil {
		srcWidget = *widget
	} else {
		if widgetID == "" {
			return nil, fmt.Errorf("WidgetsContainId: widgetID must be provided if widget is nil")
		}
		w, err := s.GetWidget(ctx, canvasID, widgetID)
		if err != nil {
			return nil, fmt.Errorf("WidgetsContainId: failed to fetch widget: %w", err)
		}
		srcWidget = *w
	}

	// Fetch all widgets on the same canvas
	widgets, err := s.ListWidgets(ctx, canvasID, nil)
	if err != nil {
		return nil, fmt.Errorf("WidgetsContainId: failed to list widgets: %w", err)
	}

	srcRect := WidgetBoundingBox(srcWidget)
	// Expand bounding box by tolerance
	srcRect.X -= tolerance
	srcRect.Y -= tolerance
	srcRect.Width += 2 * tolerance
	srcRect.Height += 2 * tolerance

	var contained []Widget
	for _, w := range widgets {
		if w.ID == srcWidget.ID {
			continue // skip self
		}
		if WidgetContainsRect(srcRect, w) {
			contained = append(contained, w)
		}
	}
	return contained, nil
}

// WidgetContainsRect returns true if the given rectangle fully contains the widget's bounding box.
func WidgetContainsRect(rect Rectangle, w Widget) bool {
	return Contains(rect, WidgetBoundingBox(w))
}
