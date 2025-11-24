package canvus

import (
	"context"
	"fmt"
)

// ListVideoOutputs retrieves all video outputs for a given client device.
func (s *Session) ListVideoOutputs(ctx context.Context, clientID string) ([]VideoOutput, error) {
	var outputs []VideoOutput
	path := fmt.Sprintf("clients/%s/video-outputs", clientID)
	err := s.doRequest(ctx, "GET", path, nil, &outputs, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListVideoOutputs: %w", err)
	}
	return outputs, nil
}

// GetVideoOutput retrieves a single video output by ID for a specific client.
func (s *Session) GetVideoOutput(ctx context.Context, clientID, outputID string) (*VideoOutput, error) {
	var output VideoOutput
	err := s.doRequest(ctx, "GET", fmt.Sprintf("clients/%s/video-outputs/%s", clientID, outputID), nil, &output, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetVideoOutput: %w", err)
	}
	return &output, nil
}

// SetVideoOutputSource sets the source or suspends a video output for a client device.
func (s *Session) SetVideoOutputSource(ctx context.Context, clientID string, index int, req interface{}) error {
	path := fmt.Sprintf("clients/%s/video-outputs/%d", clientID, index)
	return s.doRequest(ctx, "PATCH", path, req, nil, nil, false)
}

// UpdateVideoOutput updates a video output for a canvas (name, resolution).
func (s *Session) UpdateVideoOutput(ctx context.Context, canvasID, outputID string, req interface{}) (*VideoOutput, error) {
	var output VideoOutput
	path := fmt.Sprintf("canvases/%s/video-outputs/%s", canvasID, outputID)
	err := s.doRequest(ctx, "PATCH", path, req, &output, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateVideoOutput: %w", err)
	}
	return &output, nil
}
