package canvus

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockSession struct {
	canvases         []Canvas
	widgets          map[string][]Widget
	failListCanvases bool
	failListWidgets  map[string]bool
}

func (m *mockSession) ListCanvases(ctx context.Context, filter *Filter) ([]Canvas, error) {
	if m.failListCanvases {
		return nil, errors.New("mock ListCanvases failure")
	}
	return m.canvases, nil
}

func (m *mockSession) ListWidgets(ctx context.Context, canvasID string, filter *Filter) ([]Widget, error) {
	if m.failListWidgets != nil && m.failListWidgets[canvasID] {
		return nil, errors.New("mock ListWidgets failure")
	}
	widgets := m.widgets[canvasID]
	if filter != nil {
		widgets = FilterSlice(widgets, filter)
	}
	return widgets, nil
}

func TestFindWidgetsAcrossCanvases(t *testing.T) {
	ctx := context.Background()
	ms := &mockSession{
		canvases: []Canvas{{ID: "c1"}, {ID: "c2"}},
		widgets: map[string][]Widget{
			"c1": {
				{ID: "w1", WidgetType: "browser", ParentID: "", State: "active"},
				{ID: "w2", WidgetType: "note", ParentID: "", State: "archived"},
			},
			"c2": {
				{ID: "w3", WidgetType: "browser", ParentID: "", State: "active"},
				{ID: "w4", WidgetType: "browser", ParentID: "", State: "inactive"},
			},
		},
	}

	t.Run("ExactMatch", func(t *testing.T) {
		query := map[string]interface{}{"widget_type": "note"}
		matches, err := FindWidgetsAcrossCanvases(ctx, ms, query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []WidgetMatch{{CanvasID: "c1", WidgetID: "w2", Widget: ms.widgets["c1"][1]}}
		if !reflect.DeepEqual(matches, want) {
			t.Errorf("got %+v, want %+v", matches, want)
		}
	})

	t.Run("WildcardMatch", func(t *testing.T) {
		query := map[string]interface{}{"widget_type": "*"}
		matches, err := FindWidgetsAcrossCanvases(ctx, ms, query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matches) != 4 {
			t.Errorf("expected 4 matches, got %d", len(matches))
		}
	})

	t.Run("SuffixMatch", func(t *testing.T) {
		ms.widgets["c1"][0].ParentID = "abc12345"
		query := map[string]interface{}{"parent_id": "*12345"}
		matches, err := FindWidgetsAcrossCanvases(ctx, ms, query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matches) != 1 || matches[0].WidgetID != "w1" {
			t.Errorf("expected w1, got %+v", matches)
		}
	})

	t.Run("ListCanvasesError", func(t *testing.T) {
		ms.failListCanvases = true
		_, err := FindWidgetsAcrossCanvases(ctx, ms, map[string]interface{}{"widget_type": "browser"})
		if err == nil {
			t.Error("expected error, got nil")
		}
		ms.failListCanvases = false
	})

	t.Run("ListWidgetsError", func(t *testing.T) {
		ms.failListWidgets = map[string]bool{"c2": true}
		_, err := FindWidgetsAcrossCanvases(ctx, ms, map[string]interface{}{"widget_type": "browser"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
