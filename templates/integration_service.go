// Integration Service Template
//
// This template provides a production-ready microservice structure for
// building integration services with the Canvus SDK. It includes:
// - Configuration management
// - Prometheus-style metrics
// - Health checks (readiness/liveness)
// - Graceful shutdown
// - Structured logging
//
// Usage:
//   1. Copy this file to your project
//   2. Search for "TODO:" comments and customize
//   3. Build with: go build -o your-service
//   4. Run with: ./your-service
//
// Environment Variables:
//   CANVUS_API_URL    - Required: Canvus API endpoint
//   CANVUS_API_KEY    - Required: API key for authentication
//   HTTP_PORT         - Optional: HTTP port for API (default: 8080)
//   METRICS_PORT      - Optional: HTTP port for metrics (default: 9090)
//   LOG_LEVEL         - Optional: Log level (debug, info, warn, error)

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
	// TODO: Add your imports here
)

// Config holds the service configuration
type Config struct {
	APIURL      string
	APIKey      string
	HTTPPort    string
	MetricsPort string
	LogLevel    string
	// TODO: Add your configuration fields here
}

// Metrics holds service metrics
type Metrics struct {
	RequestsTotal     int64
	RequestsSuccess   int64
	RequestsFailed    int64
	CanvasOperations  int64
	WidgetOperations  int64
	LastRequestTime   time.Time
	ServiceStartTime  time.Time
	mu                sync.RWMutex
}

// Service holds the service dependencies
type Service struct {
	config        *Config
	session       *canvus.Session
	apiServer     *http.Server
	metricsServer *http.Server
	metrics       *Metrics
	ready         atomic.Bool
	healthy       atomic.Bool
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

	// Start servers
	go func() {
		log.Printf("Starting API server on port %s", cfg.HTTPPort)
		if err := svc.apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server error: %v", err)
		}
	}()

	go func() {
		log.Printf("Starting metrics server on port %s", cfg.MetricsPort)
		if err := svc.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Metrics server error: %v", err)
		}
	}()

	// Initialize service (e.g., warm caches, verify connections)
	if err := svc.initialize(); err != nil {
		log.Fatalf("Service initialization failed: %v", err)
	}

	// Mark service as ready
	svc.ready.Store(true)
	log.Println("Service is ready")

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down service...")

	// Mark service as not ready (stops accepting new requests)
	svc.ready.Store(false)

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown servers
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := svc.apiServer.Shutdown(ctx); err != nil {
			log.Printf("API server shutdown error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := svc.metricsServer.Shutdown(ctx); err != nil {
			log.Printf("Metrics server shutdown error: %v", err)
		}
	}()

	wg.Wait()
	log.Println("Service stopped")
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	cfg := &Config{
		HTTPPort:    "8080",
		MetricsPort: "9090",
		LogLevel:    "info",
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

	// Optional settings
	if port := os.Getenv("HTTP_PORT"); port != "" {
		cfg.HTTPPort = port
	}
	if port := os.Getenv("METRICS_PORT"); port != "" {
		cfg.MetricsPort = port
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}

	// TODO: Load your environment variables here
	// Example:
	// cfg.WebhookURL = os.Getenv("WEBHOOK_URL")
	// cfg.RateLimitRPS = parseIntEnv("RATE_LIMIT_RPS", 100)

	return cfg, nil
}

// newService creates a new service instance
func newService(cfg *Config) (*Service, error) {
	// Create SDK configuration
	sdkCfg := canvus.DefaultSessionConfig()
	sdkCfg.BaseURL = cfg.APIURL
	sdkCfg.RequestTimeout = 30 * time.Second

	// Create session with API key authentication
	session := canvus.NewSession(sdkCfg, canvus.WithAPIKey(cfg.APIKey))

	svc := &Service{
		config:  cfg,
		session: session,
		metrics: &Metrics{
			ServiceStartTime: time.Now(),
		},
	}

	// Mark as healthy (can be changed later if dependencies fail)
	svc.healthy.Store(true)

	// Set up API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/health/live", svc.handleLiveness)
	apiMux.HandleFunc("/health/ready", svc.handleReadiness)

	// TODO: Add your API handlers here
	// Example:
	// apiMux.HandleFunc("/api/canvases", svc.handleCanvases)
	// apiMux.HandleFunc("/api/widgets", svc.handleWidgets)
	// apiMux.HandleFunc("/api/sync", svc.handleSync)

	// Demo endpoint
	apiMux.HandleFunc("/api/canvases", svc.handleListCanvases)

	// Apply middleware
	apiHandler := svc.metricsMiddleware(loggingMiddleware(apiMux))

	// Create API server
	svc.apiServer = &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      apiHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Set up metrics routes
	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/metrics", svc.handleMetrics)

	// Create metrics server
	svc.metricsServer = &http.Server{
		Addr:         ":" + cfg.MetricsPort,
		Handler:      metricsMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return svc, nil
}

// initialize performs service initialization tasks
func (s *Service) initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test API connection
	_, err := s.session.ListCanvases(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Canvus API: %w", err)
	}

	// TODO: Add your initialization tasks here
	// Example:
	// - Warm up caches
	// - Verify external service connections
	// - Load initial data

	return nil
}

// metricsMiddleware updates metrics for each request
func (s *Service) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&s.metrics.RequestsTotal, 1)

		// Create a response wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Update metrics based on response
		if wrapped.statusCode >= 200 && wrapped.statusCode < 300 {
			atomic.AddInt64(&s.metrics.RequestsSuccess, 1)
		} else {
			atomic.AddInt64(&s.metrics.RequestsFailed, 1)
		}

		s.metrics.mu.Lock()
		s.metrics.LastRequestTime = time.Now()
		s.metrics.mu.Unlock()
	})
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

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

// handleLiveness handles Kubernetes liveness probe
func (s *Service) handleLiveness(w http.ResponseWriter, r *http.Request) {
	if !s.healthy.Load() {
		http.Error(w, "unhealthy", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// handleReadiness handles Kubernetes readiness probe
func (s *Service) handleReadiness(w http.ResponseWriter, r *http.Request) {
	if !s.ready.Load() {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

// handleMetrics exposes Prometheus-style metrics
func (s *Service) handleMetrics(w http.ResponseWriter, r *http.Request) {
	s.metrics.mu.RLock()
	lastRequest := s.metrics.LastRequestTime
	s.metrics.mu.RUnlock()

	uptime := time.Since(s.metrics.ServiceStartTime).Seconds()

	// Prometheus exposition format
	metrics := fmt.Sprintf(`# HELP canvus_requests_total Total number of requests
# TYPE canvus_requests_total counter
canvus_requests_total %d

# HELP canvus_requests_success_total Total number of successful requests
# TYPE canvus_requests_success_total counter
canvus_requests_success_total %d

# HELP canvus_requests_failed_total Total number of failed requests
# TYPE canvus_requests_failed_total counter
canvus_requests_failed_total %d

# HELP canvus_canvas_operations_total Total number of canvas operations
# TYPE canvus_canvas_operations_total counter
canvus_canvas_operations_total %d

# HELP canvus_widget_operations_total Total number of widget operations
# TYPE canvus_widget_operations_total counter
canvus_widget_operations_total %d

# HELP canvus_uptime_seconds Service uptime in seconds
# TYPE canvus_uptime_seconds gauge
canvus_uptime_seconds %f

# HELP canvus_last_request_timestamp Unix timestamp of last request
# TYPE canvus_last_request_timestamp gauge
canvus_last_request_timestamp %d
`,
		atomic.LoadInt64(&s.metrics.RequestsTotal),
		atomic.LoadInt64(&s.metrics.RequestsSuccess),
		atomic.LoadInt64(&s.metrics.RequestsFailed),
		atomic.LoadInt64(&s.metrics.CanvasOperations),
		atomic.LoadInt64(&s.metrics.WidgetOperations),
		uptime,
		lastRequest.Unix(),
	)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}

// handleListCanvases handles requests to list canvases
func (s *Service) handleListCanvases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	canvases, err := s.session.ListCanvases(ctx, nil)
	if err != nil {
		handleAPIError(w, err, "list canvases")
		return
	}

	// Update metrics
	atomic.AddInt64(&s.metrics.CanvasOperations, 1)

	writeJSON(w, http.StatusOK, canvases)
}

// TODO: Add your API handlers here
//
// Example: Sync handler
// func (s *Service) handleSync(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodPost {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }
//
//     ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
//     defer cancel()
//
//     // Perform sync operation
//     result, err := s.syncData(ctx)
//     if err != nil {
//         handleAPIError(w, err, "sync")
//         return
//     }
//
//     writeJSON(w, http.StatusOK, result)
// }
//
// Example: Webhook handler
// func (s *Service) handleWebhook(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodPost {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }
//
//     var payload WebhookPayload
//     if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
//         writeError(w, http.StatusBadRequest, "Invalid payload")
//         return
//     }
//
//     ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
//     defer cancel()
//
//     // Process webhook
//     if err := s.processWebhook(ctx, &payload); err != nil {
//         handleAPIError(w, err, "process webhook")
//         return
//     }
//
//     w.WriteHeader(http.StatusAccepted)
// }

// handleAPIError handles API errors and writes appropriate HTTP responses
func handleAPIError(w http.ResponseWriter, err error, operation string) {
	if apiErr, ok := err.(*canvus.APIError); ok {
		statusCode := apiErr.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		writeError(w, statusCode, fmt.Sprintf("%s failed: %s", operation, apiErr.Message))
		return
	}
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
