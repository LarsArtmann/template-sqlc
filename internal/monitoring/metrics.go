package monitoring

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/promhttp"
)

// Metrics collects and exposes sqlc-related metrics
type Metrics struct {
	// Code generation metrics
	CodeGenDuration prometheus.Histogram
	CodeGenErrors   prometheus.Counter
	CodeGenTotal    prometheus.Counter

	// Database query metrics
	QueryDuration     prometheus.Histogram
	QueryErrors       prometheus.Counter
	QueryTotal        prometheus.Counter
	ActiveConnections prometheus.Gauge

	// User operation metrics
	UserOperations      prometheus.Counter
	UserCreations       prometheus.Counter
	UserAuthentications prometheus.Counter

	// Session metrics
	SessionCreations prometheus.Counter
	SessionActive    prometheus.Gauge

	// Configuration metrics
	ConfigFileSize prometheus.Gauge
	ConfigDatabase prometheus.Counter

	// Build metrics
	BuildDuration prometheus.Histogram
	BuildSuccess  prometheus.Counter
	BuildFailures prometheus.Counter

	registry *prometheus.Registry
	server   *http.Server
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()

	metrics := &Metrics{
		// Code generation metrics
		CodeGenDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:      "sqlc_codegen_duration_seconds",
				Help:      "Duration of sqlc code generation in seconds",
				Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 30},
				Namespace: "sqlc",
				Subsystem: "codegen",
			},
		),
		CodeGenErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_codegen_errors_total",
				Help:      "Total number of sqlc code generation errors",
				Namespace: "sqlc",
				Subsystem: "codegen",
			},
		),
		CodeGenTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_codegen_total",
				Help:      "Total number of sqlc code generation attempts",
				Namespace: "sqlc",
				Subsystem: "codegen",
			},
		),

		// Database query metrics
		QueryDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:      "sqlc_query_duration_seconds",
				Help:      "Duration of database queries in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
				Namespace: "sqlc",
				Subsystem: "query",
			},
		),
		QueryErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_query_errors_total",
				Help:      "Total number of database query errors",
				Namespace: "sqlc",
				Subsystem: "query",
			},
		),
		QueryTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_query_total",
				Help:      "Total number of database queries executed",
				Namespace: "sqlc",
				Subsystem: "query",
			},
		),
		ActiveConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:      "sqlc_database_connections_active",
				Help:      "Number of active database connections",
				Namespace: "sqlc",
				Subsystem: "database",
			},
		),

		// User operation metrics
		UserOperations: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_user_operations_total",
				Help:      "Total number of user operations performed",
				Namespace: "sqlc",
				Subsystem: "user",
			},
		),
		UserCreations: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_user_creations_total",
				Help:      "Total number of user creations performed",
				Namespace: "sqlc",
				Subsystem: "user",
			},
		),
		UserAuthentications: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_user_authentications_total",
				Help:      "Total number of user authentications performed",
				Namespace: "sqlc",
				Subsystem: "user",
			},
		),

		// Session metrics
		SessionCreations: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_session_creations_total",
				Help:      "Total number of session creations performed",
				Namespace: "sqlc",
				Subsystem: "session",
			},
		),
		SessionActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:      "sqlc_sessions_active",
				Help:      "Number of active user sessions",
				Namespace: "sqlc",
				Subsystem: "session",
			},
		),

		// Configuration metrics
		ConfigFileSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:      "sqlc_config_file_size_bytes",
				Help:      "Size of sqlc configuration file in bytes",
				Namespace: "sqlc",
				Subsystem: "config",
			},
		),
		ConfigDatabase: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_config_databases_total",
				Help:      "Total number of databases configured in sqlc.yaml",
				Namespace: "sqlc",
				Subsystem: "config",
			},
		),

		// Build metrics
		BuildDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:      "sqlc_build_duration_seconds",
				Help:      "Duration of build operations in seconds",
				Buckets:   []float64{1, 5, 10, 30, 60, 300, 600},
				Namespace: "sqlc",
				Subsystem: "build",
			},
		),
		BuildSuccess: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_build_success_total",
				Help:      "Total number of successful builds",
				Namespace: "sqlc",
				Subsystem: "build",
			},
		),
		BuildFailures: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:      "sqlc_build_failures_total",
				Help:      "Total number of build failures",
				Namespace: "sqlc",
				Subsystem: "build",
			},
		),

		registry: registry,
	}

	// Register metrics
	registry.MustRegister(
		metrics.CodeGenDuration,
		metrics.CodeGenErrors,
		metrics.CodeGenTotal,
		metrics.QueryDuration,
		metrics.QueryErrors,
		metrics.QueryTotal,
		metrics.ActiveConnections,
		metrics.UserOperations,
		metrics.UserCreations,
		metrics.UserAuthentications,
		metrics.SessionCreations,
		metrics.SessionActive,
		metrics.ConfigFileSize,
		metrics.ConfigDatabase,
		metrics.BuildDuration,
		metrics.BuildSuccess,
		metrics.BuildFailures,
	)

	return metrics
}

// ObserveCodeGen records metrics for code generation
func (m *Metrics) ObserveCodeGen(duration time.Duration, err error) {
	m.CodeGenTotal.Inc()
	m.CodeGenDuration.Observe(duration.Seconds())

	if err != nil {
		m.CodeGenErrors.Inc()
	}
}

// ObserveQuery records metrics for database queries
func (m *Metrics) ObserveQuery(duration time.Duration, err error) {
	m.QueryTotal.Inc()
	m.QueryDuration.Observe(duration.Seconds())

	if err != nil {
		m.QueryErrors.Inc()
	}
}

// RecordUserCreation records a user creation operation
func (m *Metrics) RecordUserCreation() {
	m.UserOperations.Inc()
	m.UserCreations.Inc()
}

// RecordUserAuthentication records a user authentication operation
func (m *Metrics) RecordUserAuthentication(success bool) {
	m.UserOperations.Inc()
	m.UserAuthentications.Inc()
}

// RecordSessionCreation records a session creation operation
func (m *Metrics) RecordSessionCreation() {
	m.SessionCreations.Inc()
}

// SetActiveSessions sets the number of active sessions
func (m *Metrics) SetActiveSessions(count int64) {
	m.SessionActive.Set(float64(count))
}

// SetActiveConnections sets the number of active database connections
func (m *Metrics) SetActiveConnections(count int64) {
	m.ActiveConnections.Set(float64(count))
}

// SetConfigFileSize sets the configuration file size
func (m *Metrics) SetConfigFileSize(size int64) {
	m.ConfigFileSize.Set(float64(size))
}

// SetConfigDatabaseCount sets the number of configured databases
func (m *Metrics) SetConfigDatabaseCount(count int64) {
	m.ConfigDatabase.Set(float64(count))
}

// ObserveBuild records metrics for build operations
func (m *Metrics) ObserveBuild(duration time.Duration, success bool) {
	m.BuildDuration.Observe(duration.Seconds())

	if success {
		m.BuildSuccess.Inc()
	} else {
		m.BuildFailures.Inc()
	}
}

// StartServer starts the metrics HTTP server
func (m *Metrics) StartServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><title>sqlc Metrics</title></head>
<body><h1>sqlc Metrics</h1>
<p><a href="/metrics">Metrics</a></p>
<p><a href="/health">Health Check</a></p>
</body></html>`))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	m.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return m.server.ListenAndServe()
}

// Shutdown gracefully shuts down the metrics server
func (m *Metrics) Shutdown(ctx context.Context) error {
	if m.server != nil {
		return m.server.Shutdown(ctx)
	}
	return nil
}

// Middleware for request tracking
func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start)
		_ = duration // TODO: Record HTTP metrics if needed
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

func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}
