package monitoring

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"          // DEPRECATED: prefer go.opentelemetry.io/otel
	"github.com/prometheus/client_golang/prometheus/promhttp" // DEPRECATED: prefer go.opentelemetry.io/otel
)

// Metrics collects and exposes sqlc-related metrics.
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
	ConfigDatabase prometheus.Gauge

	// Build metrics
	BuildDuration prometheus.Histogram
	BuildSuccess  prometheus.Counter
	BuildFailures prometheus.Counter

	registry *prometheus.Registry
	server   *http.Server
}

// NewMetrics creates a new metrics collector.
func NewMetrics() *Metrics {
	return newMetrics(prometheus.NewRegistry())
}

// HistogramConfig holds configuration for a histogram metric.
type HistogramConfig struct {
	Name      string
	Help      string
	Buckets   []float64
	Subsystem string
}

func newHistogram(cfg HistogramConfig) prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      cfg.Name,
			Help:      cfg.Help,
			Buckets:   cfg.Buckets,
			Namespace: "sqlc",
			Subsystem: cfg.Subsystem,
		},
	)
}

func newCounter(name, help, subsystem string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      name,
			Help:      help,
			Namespace: "sqlc",
			Subsystem: subsystem,
		},
	)
}

func newGauge(name, help, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:      name,
			Help:      help,
			Namespace: "sqlc",
			Subsystem: subsystem,
		},
	)
}

func newMetrics(registry *prometheus.Registry) *Metrics {
	metrics := &Metrics{
		// Code generation metrics
		CodeGenDuration: newHistogram(HistogramConfig{
			Name:      "sqlc_codegen_duration_seconds",
			Help:      "Duration of sqlc code generation in seconds",
			Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 30},
			Subsystem: "codegen",
		}),
		CodeGenErrors: newCounter(
			"sqlc_codegen_errors_total",
			"Total number of sqlc code generation errors",
			"codegen",
		),
		CodeGenTotal: newCounter(
			"sqlc_codegen_total",
			"Total number of sqlc code generation attempts",
			"codegen",
		),

		// Database query metrics
		QueryDuration: newHistogram(HistogramConfig{
			Name:      "sqlc_query_duration_seconds",
			Help:      "Duration of database queries in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
			Subsystem: "query",
		}),
		QueryErrors: newCounter(
			"sqlc_query_errors_total",
			"Total number of database query errors",
			"query",
		),
		QueryTotal: newCounter(
			"sqlc_query_total",
			"Total number of database queries executed",
			"query",
		),
		ActiveConnections: newGauge(
			"sqlc_database_connections_active",
			"Number of active database connections",
			"database",
		),

		// User operation metrics
		UserOperations: newCounter(
			"sqlc_user_operations_total",
			"Total number of user operations performed",
			"user",
		),
		UserCreations: newCounter(
			"sqlc_user_creations_total",
			"Total number of user creations performed",
			"user",
		),
		UserAuthentications: newCounter(
			"sqlc_user_authentications_total",
			"Total number of user authentications performed",
			"user",
		),

		// Session metrics
		SessionCreations: newCounter(
			"sqlc_session_creations_total",
			"Total number of session creations performed",
			"session",
		),
		SessionActive: newGauge(
			"sqlc_sessions_active",
			"Number of active user sessions",
			"session",
		),

		// Configuration metrics
		ConfigFileSize: newGauge(
			"sqlc_config_file_size_bytes",
			"Size of sqlc configuration file in bytes",
			"config",
		),
		ConfigDatabase: newGauge(
			"sqlc_config_databases_total",
			"Total number of databases configured in sqlc.yaml",
			"config",
		),

		// Build metrics
		BuildDuration: newHistogram(HistogramConfig{
			Name:      "sqlc_build_duration_seconds",
			Help:      "Duration of build operations in seconds",
			Buckets:   []float64{1, 5, 10, 30, 60, 300, 600},
			Subsystem: "build",
		}),
		BuildSuccess: newCounter(
			"sqlc_build_success_total",
			"Total number of successful builds",
			"build",
		),
		BuildFailures: newCounter(
			"sqlc_build_failures_total",
			"Total number of build failures",
			"build",
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

// observeDurationWithErrors observes duration and increments error counter if error occurs.
func (m *Metrics) observeDurationWithErrors(
	total prometheus.Counter,
	durationHist prometheus.Observer,
	duration time.Duration,
	err error,
	errors prometheus.Counter,
) {
	total.Inc()
	durationHist.Observe(duration.Seconds())

	if err != nil {
		errors.Inc()
	}
}

// ObserveCodeGen records metrics for code generation.
func (m *Metrics) ObserveCodeGen(duration time.Duration, err error) {
	m.observeDurationWithErrors(m.CodeGenTotal, m.CodeGenDuration, duration, err, m.CodeGenErrors)
}

// ObserveQuery records metrics for database queries.
func (m *Metrics) ObserveQuery(duration time.Duration, err error) {
	m.observeDurationWithErrors(m.QueryTotal, m.QueryDuration, duration, err, m.QueryErrors)
}

// RecordUserCreation records a user creation operation.
func (m *Metrics) RecordUserCreation() {
	m.UserOperations.Inc()
	m.UserCreations.Inc()
}

// RecordUserAuthentication records a user authentication operation.
func (m *Metrics) RecordUserAuthentication(success bool) {
	m.UserOperations.Inc()
	m.UserAuthentications.Inc()
}

// RecordSessionCreation records a session creation operation.
func (m *Metrics) RecordSessionCreation() {
	m.SessionCreations.Inc()
}

// SetActiveSessions sets the number of active sessions.
func (m *Metrics) SetActiveSessions(count int64) {
	m.SessionActive.Set(float64(count))
}

// SetActiveConnections sets the number of active database connections.
func (m *Metrics) SetActiveConnections(count int64) {
	m.ActiveConnections.Set(float64(count))
}

// SetConfigFileSize sets the configuration file size.
func (m *Metrics) SetConfigFileSize(size int64) {
	m.ConfigFileSize.Set(float64(size))
}

// SetConfigDatabaseCount sets the number of configured databases.
func (m *Metrics) SetConfigDatabaseCount(count int64) {
	m.ConfigDatabase.Set(float64(count))
}

// ObserveBuild records metrics for build operations.
func (m *Metrics) ObserveBuild(duration time.Duration, success bool) {
	m.BuildDuration.Observe(duration.Seconds())

	if success {
		m.BuildSuccess.Inc()
	} else {
		m.BuildFailures.Inc()
	}
}

// StartServer starts the metrics HTTP server.
func (m *Metrics) StartServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><head><title>sqlc Metrics</title></head>
<body><h1>sqlc Metrics</h1>
<p><a href="/metrics">Metrics</a></p>
<p><a href="/health">Health Check</a></p>
</body></html>`))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	m.server = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return m.server.ListenAndServe()
}

// Shutdown gracefully shuts down the metrics server.
func (m *Metrics) Shutdown(ctx context.Context) error {
	if m.server != nil {
		return m.server.Shutdown(ctx)
	}

	return nil
}

// Middleware for request tracking.
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

// responseWriter wraps http.ResponseWriter to capture status code.
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
