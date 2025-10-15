package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status_code"},
	)

	// Database metrics
	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// Cron job metrics
	cronJobsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cron_jobs_total",
			Help: "Total number of cron job executions",
		},
		[]string{"job_name", "status"},
	)

	cronJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cron_job_duration_seconds",
			Help:    "Cron job execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"job_name"},
	)
)

// InitMetrics initializes Prometheus metrics
func InitMetrics() *prometheus.Registry {
	// Create a new registry
	reg := prometheus.NewRegistry()
	
	// Register default metrics
	reg.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	// Register custom metrics
	reg.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		dbConnectionsActive,
		dbQueriesTotal,
		dbQueryDuration,
		cronJobsTotal,
		cronJobDuration,
	)

	return reg
}

// MetricsMiddleware creates a middleware for HTTP metrics
func MetricsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			// Get route template
			route := mux.CurrentRoute(r)
			routeTemplate := "unknown"
			if route != nil {
				if tpl, err := route.GetPathTemplate(); err == nil {
					routeTemplate = tpl
				}
			}

			next.ServeHTTP(wrapped, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(wrapped.statusCode)

			httpRequestsTotal.WithLabelValues(r.Method, routeTemplate, statusCode).Inc()
			httpRequestDuration.WithLabelValues(r.Method, routeTemplate, statusCode).Observe(duration)
		})
	}
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

// MetricsHandler returns the Prometheus metrics handler
func MetricsHandler(reg *prometheus.Registry) http.Handler {
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}

// RecordDBQuery records database query metrics
func RecordDBQuery(operation, table string, duration time.Duration, err error) {
	dbQueriesTotal.WithLabelValues(operation, table).Inc()
	dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordCronJob records cron job execution metrics
func RecordCronJob(jobName string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	
	cronJobsTotal.WithLabelValues(jobName, status).Inc()
	cronJobDuration.WithLabelValues(jobName).Observe(duration.Seconds())
}

// SetDBConnections sets the active database connections gauge
func SetDBConnections(count int) {
	dbConnectionsActive.Set(float64(count))
}
