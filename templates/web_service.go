// Web Service Template
//
// This template provides an HTTP service structure for building
// web applications that use the Canvus SDK. It includes:
// - RESTful HTTP endpoints
// - Request/response logging middleware
// - Health check endpoint
// - Graceful shutdown
// - Session management
//
// Usage:
//   1. Copy this file to your project
//   2. Search for "TODO:" comments and customize
//   3. Build with: go build -o your-service
//   4. Run with: ./your-service
//
// Environment Variables:
//   CANVUS_API_URL - Required: Canvus API endpoint
//   CANVUS_API_KEY - Required: API key for authentication
//   HTTP_PORT      - Optional: HTTP port (default: 8080)

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
	// TODO: Add your imports here
)

// Config holds the service configuration
type Config struct {
	APIURL   string
	APIKey   string
	HTTPPort string
	// TODO: Add your configuration fields here
}

// Service holds the service dependencies
type Service struct {
	config  *Config
	session *canvus.Session
	server  *http.Server
}

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create the service
	svc, err := newService(cfg)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	// Start the server
	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
		if err := svc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := svc.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	cfg := &Config{
		HTTPPort: "8080", // Default port
	}

	// Required: API URL
	cfg.APIURL = os.Getenv("CANVUS_API_URL")
	if cfg.APIURL == "" {
		return nil, fmt.Errorf("CANVUS_API_URL environment variable is required")
	}

	// Required: API Key
	cfg.APIKey = os.Getenv("CANVUS_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("CANVUS_API_KEY environment variable is required")
	}

	// Optional: HTTP Port
	if port := os.Getenv("HTTP_PORT"); port != "" {
		cfg.HTTPPort = port
	}

	// TODO: Load your environment variables here
	// Example:
	// cfg.AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")

	return cfg, nil
}

// newService creates a new service instance
func newService(cfg *Config) (*Service, error) {
	// Create SDK configuration
	sdkCfg := canvus.DefaultSessionConfig()
	sdkCfg.BaseURL = cfg.APIURL

	// Create session with API key authentication
	session := canvus.NewSession(sdkCfg, canvus.WithAPIKey(cfg.APIKey))

	svc := &Service{
		config:  cfg,
		session: session,
	}

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", svc.handleHealth)

	// TODO: Add your HTTP handlers here
	// Example:
	// mux.HandleFunc("/api/canvases", svc.handleCanvases)
	// mux.HandleFunc("/api/canvases/", svc.handleCanvas)
	// mux.HandleFunc("/api/widgets", svc.handleWidgets)

	// Demo endpoint: list canvases
	mux.HandleFunc("/api/canvases", svc.handleListCanvases)

	// Apply middleware
	handler := loggingMiddleware(mux)

	// Create HTTP server
	svc.server = &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return svc, nil
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log the request
		log.Printf("%s %s %d %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// handleHealth handles health check requests
func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Add health checks for your service
	// Example: Check database connection, external services, etc.
	//
	// if err := s.db.Ping(); err != nil {
	//     writeError(w, http.StatusServiceUnavailable, "Database unavailable")
	//     return
	// }

	response := map[string]string{
		"status":  "healthy",
		"service": "canvus-web-service",
		"time":    time.Now().UTC().Format(time.RFC3339),
	}

	writeJSON(w, http.StatusOK, response)
}

// handleListCanvases handles requests to list canvases
func (s *Service) handleListCanvases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// List canvases from API
	canvases, err := s.session.ListCanvases(ctx, nil)
	if err != nil {
		handleAPIError(w, err, "list canvases")
		return
	}

	// TODO: Transform response as needed
	// Example: Filter fields, add computed properties

	writeJSON(w, http.StatusOK, canvases)
}

// TODO: Add your HTTP handlers here
//
// Example: Get a specific canvas
// func (s *Service) handleGetCanvas(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodGet {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }
//
//     // Extract canvas ID from URL
//     canvasID := strings.TrimPrefix(r.URL.Path, "/api/canvases/")
//     if canvasID == "" {
//         writeError(w, http.StatusBadRequest, "Canvas ID required")
//         return
//     }
//
//     ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
//     defer cancel()
//
//     canvas, err := s.session.GetCanvas(ctx, canvasID)
//     if err != nil {
//         handleAPIError(w, err, "get canvas")
//         return
//     }
//
//     writeJSON(w, http.StatusOK, canvas)
// }
//
// Example: Create a widget
// func (s *Service) handleCreateWidget(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodPost {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }
//
//     var req struct {
//         CanvasID string          `json:"canvas_id"`
//         Widget   *canvus.Widget  `json:"widget"`
//     }
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         writeError(w, http.StatusBadRequest, "Invalid request body")
//         return
//     }
//
//     ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
//     defer cancel()
//
//     created, err := s.session.CreateWidget(ctx, req.CanvasID, req.Widget)
//     if err != nil {
//         handleAPIError(w, err, "create widget")
//         return
//     }
//
//     writeJSON(w, http.StatusCreated, created)
// }

// handleAPIError handles API errors and writes appropriate HTTP responses
func handleAPIError(w http.ResponseWriter, err error, operation string) {
	if apiErr, ok := err.(*canvus.APIError); ok {
		// Map API status codes to HTTP status codes
		statusCode := apiErr.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}

		writeError(w, statusCode, fmt.Sprintf("%s failed: %s", operation, apiErr.Message))
		return
	}

	// Generic error
	writeError(w, http.StatusInternalServerError, fmt.Sprintf("%s failed: %v", operation, err))
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// writeError writes an error response
func writeError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{
		"error": message,
	}
	writeJSON(w, statusCode, response)
}
