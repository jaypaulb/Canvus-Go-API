package canvus

import (
	"context"
	"fmt"
)

// VideoInputSource represents a video input source for a client device.
type VideoInputSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListVideoInputs retrieves all video input widgets for a given canvas.
//
// API Limitation: The 'title' field is not exposed by the Canvus API.
// Responses will not include title values. See WarningVideoInputTitleNotExposed.
func (s *Session) ListVideoInputs(ctx context.Context, canvasID string) ([]VideoInput, error) {
	warnOnce(WarningVideoInputTitleNotExposed)
	var inputs []VideoInput
	path := fmt.Sprintf("canvases/%s/video-inputs", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &inputs, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListVideoInputs: %w", err)
	}
	return inputs, nil
}

// GetClientVideoInput retrieves a single video input by ID for a specific client.
func (s *Session) GetClientVideoInput(ctx context.Context, clientID, inputID string) (*VideoInput, error) {
	var input VideoInput
	err := s.doRequest(ctx, "GET", fmt.Sprintf("clients/%s/video-inputs/%s", clientID, inputID), nil, &input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetClientVideoInput: %w", err)
	}
	return &input, nil
}

// GetVideoInput retrieves a single video input widget by ID for a specific canvas.
func (s *Session) GetVideoInput(ctx context.Context, canvasID, inputID string) (*VideoInput, error) {
	var input VideoInput
	err := s.doRequest(ctx, "GET", fmt.Sprintf("canvases/%s/video-inputs/%s", canvasID, inputID), nil, &input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetVideoInput: %w", err)
	}
	return &input, nil
}

// UpdateVideoInput updates a video input widget on a canvas.
func (s *Session) UpdateVideoInput(ctx context.Context, canvasID, inputID string, req map[string]interface{}) (*VideoInput, error) {
	// Ensure widget_type is set for the generic UpdateWidget handler if we were using it,
	// but here we are hitting the specific endpoint or using the generic patchWidgetHandler logic.
	// The MuxDispatch calls patchVideoInput which uses patchWidgetHandler with ELEM_TYPE_VIDEO_INPUT_WIDGET.
	// So we just send the fields to update.
	var input VideoInput
	err := s.doRequest(ctx, "PATCH", fmt.Sprintf("canvases/%s/video-inputs/%s", canvasID, inputID), req, &input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateVideoInput: %w", err)
	}
	return &input, nil
}

// CreateVideoInput creates a new video input widget on a canvas. The payload must include 'source' and 'host-id'.
//
// API Limitation: The 'title' field is not exposed by the Canvus API.
// Any title value in the request will be ignored. See WarningVideoInputTitleNotExposed.
func (s *Session) CreateVideoInput(ctx context.Context, canvasID string, req interface{}) (*VideoInput, error) {
	var input VideoInput
	path := fmt.Sprintf("canvases/%s/video-inputs", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateVideoInput: %w", err)
	}
	return &input, nil
}

// DeleteVideoInput deletes a video input widget by ID for a given canvas.
func (s *Session) DeleteVideoInput(ctx context.Context, canvasID, inputID string) error {
	path := fmt.Sprintf("canvases/%s/video-inputs/%s", canvasID, inputID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// ListClientVideoInputs retrieves all video input sources for a given client device.
func (s *Session) ListClientVideoInputs(ctx context.Context, clientID string) ([]VideoInputSource, error) {
	var sources []VideoInputSource
	path := fmt.Sprintf("clients/%s/video-inputs", clientID)
	err := s.doRequest(ctx, "GET", path, nil, &sources, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListClientVideoInputs: %w", err)
	}
	return sources, nil
}
